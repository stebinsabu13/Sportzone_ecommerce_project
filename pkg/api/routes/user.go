package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/stebinsabu13/ecommerce-api/pkg/api/handler"
	"github.com/stebinsabu13/ecommerce-api/pkg/api/middleware"
)

func UserRoutes(api *gin.RouterGroup, authHandler *handler.AuthHandler, userHandler *handler.UserHandler, productHandler *handler.ProductHandler, cartHandler *handler.CartHandler, orderHandler *handler.OrderHandler) {

	// login
	login := api.Group("/user")
	{
		// Request JWT
		login.POST("/login", userHandler.LoginHandler)
		login.POST("/login/otp", userHandler.LoginOtp)
		login.POST("/login/otp/verify", userHandler.LoginOtpverify)
		login.GET("/login/google", authHandler.UserGoogleAuthLoginPage)
		login.GET("/login/initialize", authHandler.UserGoogleAuthInitialize)
		login.GET("/login/callback", authHandler.UserGoogleAuthCallBack)
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
			product.GET("/:productid", productHandler.FindDetailsProductById)
		}
		filter := home.Group("/filter")
		{
			filter.GET("/category", userHandler.ListAllCategories)
			filter.GET("/category/:categoryid/products", productHandler.ProductsByCategory)
			filter.GET("/brands", productHandler.ListAllBrands)
			filter.GET("/brands/:brandid/products", productHandler.ProductsByBrands)
		}
		profile := home.Group("/profile")
		{
			profile.GET("/", userHandler.ShowUserDetails)
			profile.PATCH("/edit/details", userHandler.UpdateProfile)
			address := profile.Group("/address")
			{
				address.GET("/", userHandler.ShowAllAddress)
				address.POST("/add", userHandler.AddAddress)
			}
		}
		orders := home.Group("/orders")
		{
			orders.GET("/", orderHandler.ShowOrders)
			orders.GET("/detail", orderHandler.ShowOrderDetail)
			orders.PATCH("/cancel", orderHandler.CancelOrder)
			orders.PATCH("/return", orderHandler.ReturnOrder)
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
		wallet := home.Group("/wallet")
		{
			wallet.GET("/", userHandler.ViewWallet)
		}
	}
}
