package routes

import (
	"QuickGin/controllers"
	"QuickGin/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(rg *gin.RouterGroup) {
	auth := controllers.NewAuthController()
	authRoutes := rg.Group("/auth")
	authRoutes.Use(middleware.PerRoute(1, 5))
	{
		// authRoutes.POST("/with-password", auth.AuthWithPassword)
		authRoutes.POST("/request-otp", auth.AuthRequestOtp)
		authRoutes.POST("/with-otp", auth.AuthWithOTP)
		authRoutes.POST("/token/refresh", auth.RefreshToken)
		authRoutes.POST("/logout", middleware.TokenAuth(), auth.AuthLogout)
	}
}
