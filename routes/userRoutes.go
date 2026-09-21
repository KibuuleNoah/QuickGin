package routes

import (
	"github.com/gin-gonic/gin"
	"QuickGin/controllers"
	"QuickGin/middleware"
)

func RegisterUserRoutes(rg *gin.RouterGroup) {

	user := controllers.NewUserController()
	userRoutes := rg.Group("/user")
	{
		userRoutes.POST("/", user.CreateUser)
		userRoutes.GET("/me", middleware.AuthRequired(), user.GetProfile)
		userRoutes.PUT("/me", middleware.AuthRequired(), user.UpdateUser)
		// userRoutes.DELETE("/me", middleware.AuthRequired(), user.DeleteAccount)
	}

}
