package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"QuickGin/config"
	"QuickGin/db"
	_ "QuickGin/docs"
	"QuickGin/forms"
	"QuickGin/middleware"
	"QuickGin/routes"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	cfg := config.Load()

	if cfg.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	if !cfg.IsProd() {
		r.Use(gin.Logger())
	}

	if err := r.SetTrustedProxies(nil); err != nil {
		log.Fatalf("failed to set trusted proxies: %v", err)
	}

	binding.Validator = new(forms.DefaultValidator)

	limiter := middleware.NewRateLimiter(float64(cfg.RateLimitRPS), cfg.RateLimitBurst)
	r.Use(limiter.Middleware())

	r.Use(middleware.CORS())
	r.Use(middleware.RequestID())
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	if err := db.InitAppDB(); err != nil {
		log.Fatalf("failed to init db: %v", err)
	}
	if err := db.InitAppCache(db.PostgresCache); err != nil {
		log.Fatalf("failed to init cache: %v", err)
	}

	// controllers.NewWebController(r)
	apiV1 := r.Group("/api/v1")
	routes.RegisterRoutes(apiV1)

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	if !cfg.IsProd() {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3002"
	}
	if _, err := strconv.Atoi(port); err != nil {
		log.Fatalf("invalid PORT value: %q", port)
	}

	log.Printf("\n\n PORT: %s \n ENV: %s \n SSL: %s \n Version: %s \n\n", port, cfg.Env, cfg.DBSSLMode, os.Getenv("API_VERSION"))

	srv := &http.Server{
		Addr:              "0.0.0.0:" + cfg.Port,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	if err := db.AppDB().Close(); err != nil {
		log.Printf("error closing db: %v", err)
	}

	fmt.Println("server exited")
}
