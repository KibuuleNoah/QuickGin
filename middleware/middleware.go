package middleware

import "github.com/gin-gonic/gin"

func NoCacheMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only apply this in development mode
		if gin.Mode() == gin.DebugMode {
			c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
		}
		c.Next()
	}
}
