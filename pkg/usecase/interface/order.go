package interfaces

import (
	"context"

	"github.com/stebinsabu13/ecommerce-api/pkg/utils"
)

type OrderUseCase interface {
	Orders(ctx context.Context, id uint) ([]utils.ResOrders, error)
	OrderDetail(uint) ([]utils.ResponseOrderDetails, error)
	AddtoOrders(uint, uint, uint) error
	Razorpayment(uint) (utils.RazorpayOrder, error)
	CancelOrder(context.Context, uint, uint) error
	ReturnOrder(uint, uint) error
	ValidateCoupon(uint, string) error

	//Admin UseCase

	ListAllOrders() ([]utils.ResAllOrders, error)
	UpdateStatus(uint, uint) error
}
