// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"github.com/stebinsabu13/ecommerce-api/pkg/api"
	"github.com/stebinsabu13/ecommerce-api/pkg/api/handler"
	"github.com/stebinsabu13/ecommerce-api/pkg/config"
	"github.com/stebinsabu13/ecommerce-api/pkg/db"
	"github.com/stebinsabu13/ecommerce-api/pkg/repository"
	"github.com/stebinsabu13/ecommerce-api/pkg/usecase"
)

// Injectors from wire.go:

func InitializeAPI(cfg config.Config) (*api.ServerHTTP, error) {
	gormDB, err := db.ConnectDatabase(cfg)
	if err != nil {
		return nil, err
	}
	userRepository := repository.NewUserRepository(gormDB)
	userUseCase := usecase.NewUserUseCase(userRepository)
	otpRepository := repository.NewOtpRepository(gormDB)
	otpUseCase := usecase.NewOtpUseCase(otpRepository)
	orderRepository := repository.NewOrderRepository(gormDB)
	cartRepository := repository.NewCartRepository(gormDB)
	orderUseCase := usecase.NewOrderUseCase(orderRepository, cartRepository)
	userHandler := handler.NewUserHandler(userUseCase, otpUseCase, orderUseCase)
	adminRepository := repository.NewAdminRepository(gormDB)
	adminUseCase := usecase.NewAdminUseCase(adminRepository)
	adminHandler := handler.NewAdminHandler(adminUseCase, otpUseCase)
	productRepository := repository.NewProductrepository(gormDB)
	productUseCase := usecase.NewProductUseCase(productRepository)
	productHandler := handler.NewProductHandler(productUseCase)
	cartUseCase := usecase.NewCartUseCase(cartRepository)
	cartHandler := handler.NewCartHandler(cartUseCase)
	orderHandler := handler.NewOrderHandler(orderUseCase)
	serverHTTP := api.NewServerHTTP(userHandler, adminHandler, productHandler, cartHandler, orderHandler)
	return serverHTTP, nil
}
