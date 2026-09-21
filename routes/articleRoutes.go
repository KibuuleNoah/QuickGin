package routes

import (
	"github.com/gin-gonic/gin"
	"QuickGin/controllers"
	"QuickGin/middleware"
)

func RegisterArticleRoutes(rg *gin.RouterGroup) {
	articles := controllers.NewArticleController()
	articleRoutes := rg.Group("/articles")
	{
		articleRoutes.GET("/", articles.ListArticles)
		articleRoutes.GET("/:id", articles.GetArticle)
		articleRoutes.POST("/", middleware.AuthRequired(), articles.CreateArticle)
		articleRoutes.PUT("/:id", middleware.AuthRequired(), articles.UpdateArticle)
		articleRoutes.DELETE("/:id", middleware.AuthRequired(), articles.DeleteArticle)
	}
}
