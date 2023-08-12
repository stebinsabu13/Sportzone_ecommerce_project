package api

import (
	"net/http"

	_ "github.com/stebinsabu13/ecommerce-api/docs"

	"github.com/gin-gonic/gin"
	"github.com/stebinsabu13/ecommerce-api/pkg/api/handler"
	"github.com/stebinsabu13/ecommerce-api/pkg/api/routes"
	"github.com/stebinsabu13/ecommerce-api/pkg/config"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type ServerHTTP struct {
	engine *gin.Engine
}

func NewServerHTTP(authHandler *handler.AuthHandler, userHandler *handler.UserHandler, adminHandler *handler.AdminHandler, productHandler *handler.ProductHandler, cartHandler *handler.CartHandler, orderHandler *handler.OrderHandler) *ServerHTTP {

	engine := gin.New()

	engine.Use(gin.Logger())

	// swagger docs
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// set up routes
	routes.UserRoutes(engine.Group("/"), authHandler, userHandler, productHandler, cartHandler, orderHandler)
	routes.AdminRoutes(engine.Group("/"), adminHandler, productHandler, orderHandler)

	// no handler
	engine.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"StatusCode": 404,
			"msg":        "invalid url",
		})
	})

	return &ServerHTTP{engine: engine}
}

func (s *ServerHTTP) Start() {
	port := config.GetLocalPort()
	s.engine.LoadHTMLGlob("static/*.html")
	if err := s.engine.Run(":" + port); err != nil {
		return
	}
}
