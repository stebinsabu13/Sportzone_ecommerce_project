package usecase

import (
	"context"

	"github.com/stebinsabu13/ecommerce-api/pkg/domain"
	interfaces "github.com/stebinsabu13/ecommerce-api/pkg/repository/interface"
	services "github.com/stebinsabu13/ecommerce-api/pkg/usecase/interface"
	"github.com/stebinsabu13/ecommerce-api/pkg/utils"
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

func (c *adminUseCase) ListAllUsers(ctx context.Context, pagination utils.Pagination) ([]utils.ResponseUsers, error) {
	users, err := c.adminrepo.ListAllUsers(ctx, pagination)
	return users, err
}

func (c *adminUseCase) AccessManage(ctx context.Context, id string, access bool) error {
	err := c.adminrepo.AccessManage(ctx, id, access)
	return err
}

func (c *adminUseCase) ListAllCategories(ctx context.Context) ([]utils.ResponseCategory, error) {
	categories, err := c.adminrepo.ListAllCategories(ctx)
	return categories, err
}

func (c *adminUseCase) AddCategory(ctx context.Context, category domain.Category) error {
	err := c.adminrepo.AddCategory(ctx, category)
	return err
}

func (c *adminUseCase) DeleteCategory(ctx context.Context, id string) error {
	err := c.adminrepo.DeleteCategory(ctx, id)
	return err
}

func (c *adminUseCase) GetFullSalesReport(reqData utils.SalesReport) ([]utils.ResSalesReport, error) {
	return c.adminrepo.GetFullSalesReport(reqData)
}

func (c *adminUseCase) Widgets() (utils.ResWidgets, error) {
	return c.adminrepo.Widgets()
}
