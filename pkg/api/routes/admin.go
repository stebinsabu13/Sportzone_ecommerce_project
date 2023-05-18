package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/stebinsabu13/ecommerce-api/pkg/api/handler"
	"github.com/stebinsabu13/ecommerce-api/pkg/api/middleware"
)

func AdminRoutes(api *gin.RouterGroup, adminHandler *handler.AdminHandler, productHandler *handler.ProductHandler) {
	login := api.Group("/admin")
	{
		login.POST("/login", adminHandler.LoginHandler)
		// login.POST("/signup", adminHandler.SignUp)
	}
	home := login.Group("/")
	{
		home.Use(middleware.AuthorizationMiddleware("admin"))
		// home.GET("/home", adminHandler.HomeHandler)
		home.POST("/logout", adminHandler.LogoutHandler)
		user := home.Group("/user")
		{
			user.GET("/", adminHandler.ListAllUsers)
			user.PATCH("/:userid/make", adminHandler.AccessManage)
		}
		category := home.Group("/category")
		{
			category.GET("/", adminHandler.ListAllCategories)
			category.POST("/add", adminHandler.AddCategory)
			category.DELETE("/delete/:categoryid", adminHandler.DeleteCategory)
		}
		product := home.Group("/product")
		{
			product.GET("/", productHandler.FindAllProducts)
			product.POST("/add", productHandler.AddProduct)
			product.PATCH("/update/:productid", productHandler.EditProduct)
			product.DELETE("/delete/:productid", productHandler.DeleteProduct)
		}
	}
}
