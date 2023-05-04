package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/stebinsabu13/ecommerce-api/pkg/api/handler"
	"github.com/stebinsabu13/ecommerce-api/pkg/api/middleware"
)

func UserRoutes(api *gin.RouterGroup, userHandler *handler.UserHandler, productHandler *handler.ProductHandler) {

	// login
	login := api.Group("/user")
	{
		// Request JWT
		login.POST("/login", userHandler.LoginHandler)
		login.POST("/login/otp", userHandler.LoginOtp)
		login.POST("/login/otp/verify", userHandler.LoginOtpverify)

	}
	signup := api.Group("/user")
	{
		signup.POST("/signup", userHandler.SignUp)
		signup.POST("/signup/otp/verify", userHandler.SignupOtpverify)
	}
	home := api.Group("/user")
	{
		// Auth middleware
		home.Use(middleware.AuthorizationMiddleware("user"))
		home.GET("/home", userHandler.HomeHandler)
		home.POST("/logout", userHandler.LogoutHandler)
		product := home.Group("/products")
		{
			product.GET("/", productHandler.FindAllProducts)
			product.GET("/:id", productHandler.FindProductById)
		}
		profile := home.Group("/profile")
		{
			profile.GET("/:userid", userHandler.ShowUserDetails)
		}
	}
}
