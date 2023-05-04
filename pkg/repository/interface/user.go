package interfaces

import (
	"context"

	"github.com/stebinsabu13/ecommerce-api/pkg/domain"
	"github.com/stebinsabu13/ecommerce-api/pkg/utils"
)

type UserRepository interface {
	// FindAll(ctx context.Context) ([]domain.User, error)
	FindbyEmail(ctx context.Context, email string) (domain.User, error)
	SignUpUser(ctx context.Context, user domain.User) error
	FindbyEmailorMobilenum(ctx context.Context, body utils.OtpLogin) (domain.User, error)
}
