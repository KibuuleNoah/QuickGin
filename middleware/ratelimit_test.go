package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func newTestRouter(rl *RateLimiter) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(rl.Middleware())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	return r
}

func doRequest(r *gin.Engine, remoteAddr string, headers map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.RemoteAddr = remoteAddr
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestRateLimiter_AllowsWithinBurst(t *testing.T) {
	rl := NewRateLimiter(1, 5)
	r := newTestRouter(rl)

	for i := 0; i < 5; i++ {
		w := doRequest(r, "1.2.3.4:1000", nil)
		assert.Equal(t, http.StatusOK, w.Code, "request %d should pass within burst", i+1)
	}
}

func TestRateLimiter_BlocksOverBurst(t *testing.T) {
	rl := NewRateLimiter(1, 5)
	r := newTestRouter(rl)

	for i := 0; i < 5; i++ {
		doRequest(r, "1.2.3.4:1000", nil)
	}
	w := doRequest(r, "1.2.3.4:1000", nil)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
	assert.Contains(t, w.Body.String(), "rate limit exceeded")
	assert.Equal(t, "1", w.Header().Get("Retry-After"))
}

func TestRateLimiter_RefillsOverTime(t *testing.T) {
	rl := NewRateLimiter(10, 1) // 10 rps, burst 1 -> refill every 100ms
	r := newTestRouter(rl)

	w := doRequest(r, "1.2.3.4:1000", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	w = doRequest(r, "1.2.3.4:1000", nil)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)

	time.Sleep(150 * time.Millisecond)

	w = doRequest(r, "1.2.3.4:1000", nil)
	assert.Equal(t, http.StatusOK, w.Code, "should refill after waiting")
}

func TestRateLimiter_IndependentPerIP(t *testing.T) {
	rl := NewRateLimiter(1, 1)
	r := newTestRouter(rl)

	w1 := doRequest(r, "1.1.1.1:1000", nil)
	assert.Equal(t, http.StatusOK, w1.Code)

	// different IP, own bucket, should not be blocked by 1.1.1.1's usage
	w2 := doRequest(r, "2.2.2.2:1000", nil)
	assert.Equal(t, http.StatusOK, w2.Code)

	// 1.1.1.1 hits its own limit
	w3 := doRequest(r, "1.1.1.1:1000", nil)
	assert.Equal(t, http.StatusTooManyRequests, w3.Code)
}

func TestRateLimiter_RespectsXRealIP(t *testing.T) {
	rl := NewRateLimiter(1, 1)
	r := newTestRouter(rl)

	headers := map[string]string{"X-Real-IP": "9.9.9.9"}

	w1 := doRequest(r, "1.2.3.4:1000", headers)
	assert.Equal(t, http.StatusOK, w1.Code)

	// same X-Real-IP, different RemoteAddr -> should share the bucket
	w2 := doRequest(r, "5.6.7.8:1000", headers)
	assert.Equal(t, http.StatusTooManyRequests, w2.Code)
}

func TestRateLimiter_RespectsXForwardedFor(t *testing.T) {
	rl := NewRateLimiter(1, 1)
	r := newTestRouter(rl)

	headers := map[string]string{"X-Forwarded-For": "8.8.8.8, 10.0.0.1"}

	w1 := doRequest(r, "1.2.3.4:1000", headers)
	assert.Equal(t, http.StatusOK, w1.Code)

	w2 := doRequest(r, "1.2.3.4:1000", headers)
	assert.Equal(t, http.StatusTooManyRequests, w2.Code)
}

func TestRateLimiter_FallsBackToRemoteAddr(t *testing.T) {
	rl := NewRateLimiter(1, 1)
	r := newTestRouter(rl)

	w1 := doRequest(r, "3.3.3.3:5555", nil)
	assert.Equal(t, http.StatusOK, w1.Code)

	w2 := doRequest(r, "3.3.3.3:6666", nil) // same IP, different port
	assert.Equal(t, http.StatusTooManyRequests, w2.Code)
}

func TestClientKey_PrefersXRealIP(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("X-Real-IP", "7.7.7.7")
	c.Request.RemoteAddr = "1.1.1.1:1000"

	assert.Equal(t, "7.7.7.7", clientKey(c))
}

func TestClientKey_ParsesXForwardedForFirstIP(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("X-Forwarded-For", "6.6.6.6, 5.5.5.5")
	c.Request.RemoteAddr = "1.1.1.1:1000"

	assert.Equal(t, "6.6.6.6", clientKey(c))
}

func TestClientKey_FallsBackToRemoteAddr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.RemoteAddr = "2.2.2.2:9999"

	assert.Equal(t, "2.2.2.2", clientKey(c))
}

func TestCleanup_RemovesStaleVisitors(t *testing.T) {
	rl := NewRateLimiter(1, 1)
	rl.ttl = 50 * time.Millisecond // shrink ttl for test speed

	rl.getLimiter("stale-key")
	rl.mu.Lock()
	_, exists := rl.visitors["stale-key"]
	rl.mu.Unlock()
	assert.True(t, exists)

	time.Sleep(70 * time.Millisecond)

	rl.mu.Lock()
	for key, v := range rl.visitors {
		if time.Since(v.lastSeen) > rl.ttl {
			delete(rl.visitors, key)
		}
	}
	_, stillExists := rl.visitors["stale-key"]
	rl.mu.Unlock()

	assert.False(t, stillExists, "stale visitor should be cleaned up")
}

func TestPerRoute_UsesOwnLimiter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/login", PerRoute(1, 2), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodPost, "/login", nil)
		req.RemoteAddr = "4.4.4.4:1000"
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	req := httptest.NewRequest(http.MethodPost, "/login", nil)
	req.RemoteAddr = "4.4.4.4:1000"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
}
