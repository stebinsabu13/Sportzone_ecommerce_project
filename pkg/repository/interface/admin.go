package interfaces

import (
	"context"

	"github.com/stebinsabu13/ecommerce-api/pkg/domain"
	"github.com/stebinsabu13/ecommerce-api/pkg/utils"
)

type AdminRepository interface {
	FindbyEmail(ctx context.Context, email string) (domain.Admin, error)
	SignUpAdmin(ctx context.Context, admin domain.Admin) (string, error)
	UpdateVerify(string, string) error
	ListAllUsers(context.Context, utils.Pagination) ([]utils.ResponseUsers, error)
	AccessManage(ctx context.Context, id string, access bool) error
	ListAllCategories(ctx context.Context) ([]utils.ResponseCategory, error)
	AddCategory(ctx context.Context, category domain.Category) error
	DeleteCategory(ctx context.Context, id string) error
	GetFullSalesReport(utils.SalesReport) ([]utils.ResSalesReport, error)
	Widgets() (utils.ResWidgets, error)
	AddCoupon(domain.Coupon) error
	GetAllCoupons(utils.Pagination) ([]domain.Coupon, error)
	UpdateCoupon(domain.Coupon, string) error
	GetCouponByID(string) (domain.Coupon, error)
	DeleteCoupon(string) error
}
