package routes

import (
	"pajo/controllers"
	"pajo/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(rg *gin.RouterGroup) {
	auth := controllers.NewAuthController()
	authRoutes := rg.Group("/auth")
	{
		// authRoutes.POST("/with-password", auth.AuthWithPassword)
		authRoutes.POST("/request-otp", auth.AuthRequestOtp)
		authRoutes.POST("/with-otp", auth.AuthWithOTP)
		authRoutes.POST("/token/refresh", auth.RefreshToken)
		authRoutes.POST("/logout", middleware.TokenAuth(), auth.AuthLogout)
	}
}
