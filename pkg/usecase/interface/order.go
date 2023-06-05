package interfaces

import (
	"context"

	"github.com/stebinsabu13/ecommerce-api/pkg/utils"
)

type OrderUseCase interface {
	Orders(ctx context.Context, id uint) ([]utils.ResOrders, error)
	OrderDetail(uint) ([]utils.ResponseOrderDetails, error)
	AddtoOrders(uint, uint, uint, *uint) error
	Razorpayment(uint, *uint) (utils.RazorpayOrder, error)
	CancelOrder(uint, uint, uint) error
	ReturnOrder(uint, uint) error
	ValidateCoupon(uint, string) (*uint, error)
	FindCoupon(string) (*uint, error)

	//Admin UseCase

	ListAllOrders() ([]utils.ResAllOrders, error)
	UpdateStatus(uint, uint) error
}
