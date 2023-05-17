//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	http "github.com/stebinsabu13/ecommerce-api/pkg/api"
	handler "github.com/stebinsabu13/ecommerce-api/pkg/api/handler"
	config "github.com/stebinsabu13/ecommerce-api/pkg/config"
	db "github.com/stebinsabu13/ecommerce-api/pkg/db"
	repository "github.com/stebinsabu13/ecommerce-api/pkg/repository"
	usecase "github.com/stebinsabu13/ecommerce-api/pkg/usecase"
)

func InitializeAPI(cfg config.Config) (*http.ServerHTTP, error) {
	wire.Build(db.ConnectDatabase,
		repository.NewUserRepository, repository.NewAdminRepository, repository.NewProductrepository, repository.NewOtpRepository,
		usecase.NewUserUseCase, usecase.NewAdminUseCase, usecase.NewProductUseCase, usecase.NewOtpUseCase,
		handler.NewUserHandler, handler.NewAdminHandler, handler.NewProductHandler,
		http.NewServerHTTP)

	return &http.ServerHTTP{}, nil
}
