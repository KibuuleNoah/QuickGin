package routes

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup) {
	RegisterAuthRoutes(rg)
	RegisterAdminRoutes(rg)
}
