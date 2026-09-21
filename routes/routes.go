package routes

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup) {
	RegisterUserRoutes(rg)
	RegisterAuthRoutes(rg)
	RegisterAdminRoutes(rg)
	RegisterArticleRoutes(rg)
}
