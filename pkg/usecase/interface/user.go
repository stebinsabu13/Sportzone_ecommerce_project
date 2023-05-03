package interfaces

import (
	"context"

	domain "github.com/stebinsabu13/ecommerce-api/pkg/domain"
)

type UserUseCase interface {
	FindAll(ctx context.Context) ([]domain.User, error)
	FindbyEmail(ctx context.Context, email string) (domain.User, error)
	SignUpUser(ctx context.Context, user domain.User) error
}
