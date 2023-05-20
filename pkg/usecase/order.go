package usecase

import (
	"context"
	"time"

	"github.com/stebinsabu13/ecommerce-api/pkg/domain"
	interfaces "github.com/stebinsabu13/ecommerce-api/pkg/repository/interface"
	services "github.com/stebinsabu13/ecommerce-api/pkg/usecase/interface"
	"github.com/stebinsabu13/ecommerce-api/pkg/utils"
)

type orderUseCase struct {
	orderrepo interfaces.OrderRepository
	cartRepo  interfaces.CartRepository
}

func NewOrderUseCase(repo interfaces.OrderRepository, cartrepo interfaces.CartRepository) services.OrderUseCase {
	return &orderUseCase{
		orderrepo: repo,
		cartRepo:  cartrepo,
	}
}
func (c *orderUseCase) OrderDetails(ctx context.Context, id uint) ([]utils.ResponseOrderDetails, error) {
	orderDetails, err := c.orderrepo.OrderDetails(ctx, id)
	return orderDetails, err
}

func (c *orderUseCase) AddtoOrders(addressid, paymentid, userid uint) error {
	// var order domain.Order
	cart, err := c.cartRepo.FindCartById(userid)
	if err != nil {
		return err
	}
	cartitems, err1 := c.orderrepo.Findcartitems(cart.ID)
	if err1 != nil {
		return err1
	}
	order := domain.Order{
		UserID:     cart.UserID,
		PlacedDate: time.Now(),
		AddressID:  addressid,
		PaymentID:  paymentid,
		GrandTotal: uint(cart.GrandTotal),
	}
	if err := c.orderrepo.AddtoOrders(cartitems, order); err != nil {
		return err
	}
	return nil
}
