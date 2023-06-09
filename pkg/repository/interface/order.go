package interfaces

import (
	"context"
	"time"

	"github.com/stebinsabu13/ecommerce-api/pkg/domain"
	"github.com/stebinsabu13/ecommerce-api/pkg/utils"
)

type OrderRepository interface {
	Orders(ctx context.Context, id uint) ([]utils.ResOrders, error)
	OrderDetail(uint) ([]utils.ResponseOrderDetails, error)
	Findcartitems(id uint) ([]utils.ResCartItems, error)
	AddtoOrders([]utils.ResCartItems, domain.Order) error
	FindOrderitem(uint) (domain.OrderDetails, time.Time, error)
	CancelOrder(uint, domain.OrderDetails) error
	ReturnOrder(domain.OrderDetails) error
	FindCoupon(string) (domain.Coupon, error)
	ValidateCoupon(domain.Coupon, []utils.ResCartItems, *domain.Cart) error

	//Admin Repository

	ListAllOrders() ([]utils.ResAllOrders, error)
	UpdateStatus(domain.OrderDetails) error
}
