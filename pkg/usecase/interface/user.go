package interfaces

import (
	"context"

	domain "github.com/stebinsabu13/ecommerce-api/pkg/domain"
	"github.com/stebinsabu13/ecommerce-api/pkg/utils"
)

type UserUseCase interface {
	// FindAll(ctx context.Context) ([]domain.User, error)
	FindbyEmail(ctx context.Context, email string) (domain.User, error)
	SignUpUser(ctx context.Context, user domain.User) error
	FindbyEmailorMobilenum(ctx context.Context, body utils.OtpLogin) (domain.User, error)
	ShowDetails(ctx context.Context, id string) (utils.ResponseUsers, error)
	ShowAddress(ctx context.Context, id string) ([]utils.Address, error)
}
