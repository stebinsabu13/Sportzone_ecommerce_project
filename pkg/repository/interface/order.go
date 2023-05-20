package interfaces

import (
	"context"

	"github.com/stebinsabu13/ecommerce-api/pkg/domain"
	"github.com/stebinsabu13/ecommerce-api/pkg/utils"
)

type OrderRepository interface {
	OrderDetails(ctx context.Context, id uint) ([]utils.ResponseOrderDetails, error)
	Findcartitems(id uint) ([]utils.ResCartItems, error)
	AddtoOrders([]utils.ResCartItems, domain.Order) error
}
