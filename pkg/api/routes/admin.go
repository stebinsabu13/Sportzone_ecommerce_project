package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/stebinsabu13/ecommerce-api/pkg/api/handler"
	"github.com/stebinsabu13/ecommerce-api/pkg/api/middleware"
)

func AdminRoutes(api *gin.RouterGroup, adminHandler *handler.AdminHandler) {
	login := api.Group("/admin")
	{
		login.POST("/login", adminHandler.LoginHandler)
		// login.POST("/signup", adminHandler.SignUp)
	}
	home := api.Group("/admin")
	{
		home.Use(middleware.AuthorizationMiddleware("admin"))
		home.GET("/home", adminHandler.HomeHandler)
		home.POST("/logout", adminHandler.LogoutHandler)
	}
}
