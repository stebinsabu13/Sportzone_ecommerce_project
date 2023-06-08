package usecase

import (
	"context"
	"errors"

	"github.com/stebinsabu13/ecommerce-api/pkg/domain"
	interfaces "github.com/stebinsabu13/ecommerce-api/pkg/repository/interface"
	"github.com/stebinsabu13/ecommerce-api/pkg/support"
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

func (c *adminUseCase) SignUpAdmin(ctx context.Context, admin utils.BodySignUpuser) (string, error) {
	if _, err := c.adminrepo.FindbyEmail(ctx, admin.Email); err == nil {
		return "", errors.New("user already exsists")
	}
	hash, err := support.HashPassword(admin.Password)
	if err != nil {
		return "", errors.New("error while hashing password")
	}
	ADMIN := domain.Admin{
		Name:      admin.FirstName,
		Email:     admin.Email,
		MobileNum: admin.MobileNum,
		Password:  hash,
	}
	mobile_num, err := c.adminrepo.SignUpAdmin(ctx, ADMIN)
	return mobile_num, err
}

func (c *adminUseCase) UpdateVerify(number, refercode string) error {
	err := c.adminrepo.UpdateVerify(number, refercode)
	return err
}

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

func (c *adminUseCase) AddCoupon(coupon domain.Coupon) error {
	return c.adminrepo.AddCoupon(coupon)
}

func (c *adminUseCase) GetAllCoupons(pagination utils.Pagination) ([]domain.Coupon, error) {
	return c.adminrepo.GetAllCoupons(pagination)
}

func (c *adminUseCase) UpdateCoupon(coupon domain.Coupon, couponid string) error {
	return c.adminrepo.UpdateCoupon(coupon, couponid)
}

func (c *adminUseCase) GetCouponByID(id string) (domain.Coupon, error) {
	return c.adminrepo.GetCouponByID(id)
}

func (c *adminUseCase) DeleteCoupon(id string) error {
	return c.adminrepo.DeleteCoupon(id)
}
