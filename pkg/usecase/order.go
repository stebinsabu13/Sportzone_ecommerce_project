package usecase

import (
	"context"

	interfaces "github.com/stebinsabu13/ecommerce-api/pkg/repository/interface"
	services "github.com/stebinsabu13/ecommerce-api/pkg/usecase/interface"
	"github.com/stebinsabu13/ecommerce-api/pkg/utils"
)

type orderUseCase struct {
	orderrepo interfaces.OrderRepository
}

func NewOrderUseCase(repo interfaces.OrderRepository) services.OrderUseCase {
	return &orderUseCase{
		orderrepo: repo,
	}
}
func (c *orderUseCase) OrderDetails(ctx context.Context, id int) ([]utils.ResponseOrderDetails, error) {
	orderDetails, err := c.orderrepo.OrderDetails(ctx, id)
	return orderDetails, err
}
