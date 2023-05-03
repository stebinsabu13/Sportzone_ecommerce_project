package usecase

import (
	"context"

	"github.com/stebinsabu13/ecommerce-api/pkg/domain"
	interfaces "github.com/stebinsabu13/ecommerce-api/pkg/repository/interface"
	services "github.com/stebinsabu13/ecommerce-api/pkg/usecase/interface"
)

type adminUseCase struct {
	adminrepo interfaces.AdminRepository
}

func NewAdminUseCase(repo interfaces.AdminRepository) services.AdminUseCase {
	return &adminUseCase{adminrepo: repo}
}
func (c *adminUseCase) FindbyEmail(ctx context.Context, email string) (domain.Admin, error) {
	admin, err := c.adminrepo.FindbyEmail(ctx, email)
	return admin, err
}

// func (c *adminUseCase) SignUpAdmin(ctx context.Context, admin domain.Admin) error {
// 	err := c.adminrepo.SignUpAdmin(ctx, admin)
// 	return err
// }
