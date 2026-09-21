package routes

import (
	"github.com/gin-gonic/gin"
	"pajo/controllers"
)

func RegisterAdminRoutes(rg *gin.RouterGroup) {
	admin := controllers.NewAdminController()
	adminRoutes := rg.Group("/admin")
	// adminRoutes.Use(middleware.AuthRequired(), middleware.AdminRequired())
	{
		adminRoutes.GET("/items", admin.ListItems)
		adminRoutes.POST("/items", admin.CreateItem)
		// adminRoutes.PUT("/items/:id", admin.UpdateItem)
		adminRoutes.DELETE("/items/:id", admin.DeleteItem)

		// adminRoutes.GET("/categories", admin.ListCategories)
		// adminRoutes.POST("/categories", admin.CreateCategory)
		// adminRoutes.PUT("/categories/:id", admin.UpdateCategory)
		// adminRoutes.DELETE("/categories/:id", admin.DeleteCategory)

		// adminRoutes.GET("/orders", admin.ListOrders)
		// adminRoutes.GET("/orders/:id", admin.GetOrder)
		// adminRoutes.PUT("/orders/:id/status", admin.UpdateOrderStatus)

		adminRoutes.GET("/carts", admin.ListCarts)
		adminRoutes.GET("/carts/recent", admin.ListRecentCarts)

		// adminRoutes.GET("/users", admin.ListUsers)
		// adminRoutes.GET("/users/:id", admin.GetUser)
		// adminRoutes.PUT("/users/:id/status", admin.UpdateUserStatus)

		// adminRoutes.GET("/coupons", admin.ListCoupons)
		// adminRoutes.POST("/coupons", admin.CreateCoupon)
		// adminRoutes.DELETE("/coupons/:id", admin.DeleteCoupon)
	}
}
