package interfaces

import (
	"context"

	"github.com/stebinsabu13/ecommerce-api/pkg/domain"
)

type AdminUseCase interface {
	FindbyEmail(ctx context.Context, email string) (domain.Admin, error)
	// SignUpAdmin(ctx context.Context, admin domain.Admin) error
}
