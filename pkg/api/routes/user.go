package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/stebinsabu13/ecommerce-api/pkg/api/handler"
	"github.com/stebinsabu13/ecommerce-api/pkg/api/middleware"
)

func UserRoutes(api *gin.RouterGroup, userHandler *handler.UserHandler, productHandler *handler.ProductHandler, cartHandler *handler.CartHandler, orderHandler *handler.OrderHandler) {

	// login
	login := api.Group("/user")
	{
		// Request JWT
		login.POST("/login", userHandler.LoginHandler)
		login.POST("/login/otp", userHandler.LoginOtp)
		login.POST("/login/otp/verify", userHandler.LoginOtpverify)
	}
	forgotpassword := api.Group("/user/forgot/password")
	{
		forgotpassword.POST("/", userHandler.ForgotPassword)
		forgotpassword.PATCH("/otp/verify", userHandler.ForgotPasswordOtpverify)
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
		// home.GET("/home", userHandler.HomeHandler)
		home.POST("/logout", userHandler.LogoutHandler)
		product := home.Group("/products")
		{
			product.GET("/", productHandler.FindAllProducts)
			product.GET("/:id", productHandler.FindDetailsProductById)
		}
		profile := home.Group("/profile")
		{
			profile.GET("/", userHandler.ShowUserDetails)
			profile.POST("/add/address", userHandler.AddAddress)
			profile.PATCH("/edit/details", userHandler.UpdateProfile)
		}
		orders := home.Group("/orders")
		{
			orders.GET("/", orderHandler.ShowOrders)
			orders.GET("/detail", orderHandler.ShowOrderDetail)
			orders.PATCH("/cancel", orderHandler.CancelOrder)
		}
		cart := home.Group("/cart")
		{
			cart.GET("/", cartHandler.ViewCart)
			cart.PUT("/add", cartHandler.AddtoCart)
			cart.PUT("/remove", cartHandler.RemovefromCart)
		}
		checkout := home.Group("/checkout")
		{
			checkout.GET("/add", orderHandler.AddtoOrders)
			checkout.GET("/success", orderHandler.RazorpaymentSuccess)
		}
	}
}
