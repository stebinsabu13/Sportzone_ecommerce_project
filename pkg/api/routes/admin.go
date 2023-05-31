package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/stebinsabu13/ecommerce-api/pkg/api/handler"
	"github.com/stebinsabu13/ecommerce-api/pkg/api/middleware"
)

func AdminRoutes(api *gin.RouterGroup, adminHandler *handler.AdminHandler, productHandler *handler.ProductHandler, orderHandler *handler.OrderHandler) {
	login := api.Group("/admin")
	{
		login.POST("/login", adminHandler.LoginHandler)
		// login.POST("/signup", adminHandler.SignUp)
	}
	home := api.Group("/admin")
	{
		home.Use(middleware.AuthorizationMiddleware("admin"))
		// home.GET("/home", adminHandler.HomeHandler)
		home.POST("/logout", adminHandler.LogoutHandler)
		// // sales report
		sales := home.Group("/sales")
		{
			sales.GET("/", adminHandler.FullSalesReport)
		}
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
			product.PATCH("/update/:productid", productHandler.UpdateProduct)
			product.DELETE("/delete/:productid", productHandler.DeleteProduct)
			productdetail := product.Group("/detail")
			{
				productdetail.GET("/:productid", productHandler.FindDetailsProductById)
				productdetail.POST("/add", productHandler.AddProductDetail)
				productdetail.PATCH("/update/:productdetailid", productHandler.UpdateProductDetail)
				productdetail.DELETE("/delete/:productdetailid", productHandler.DeleteProductDetail)
			}
		}
		orders := home.Group("/order")
		{
			orders.GET("/", orderHandler.ListAllOrders)
			orders.GET("/detail", orderHandler.ShowOrderDetail)
			orders.PATCH("/update/status", orderHandler.UpdateStatus)
		}
		dashboard := home.Group("/dashboard")
		{
			dashboard.GET("/", adminHandler.Dashboard)
		}
		coupon := home.Group("/coupon")
		{
			coupon.POST("/add", adminHandler.AddCoupon)
			coupon.GET("/", adminHandler.GetAllCoupons)
			coupon.PATCH("/update", adminHandler.UpdateCoupon)
			coupon.GET("/:couponid", adminHandler.GetCouponByID)
			coupon.DELETE("/delete/:couponid", adminHandler.DeleteCoupon)
		}
	}
}
