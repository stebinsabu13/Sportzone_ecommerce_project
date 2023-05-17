package interfaces

import (
	"context"

	"github.com/stebinsabu13/ecommerce-api/pkg/utils"
)

type OrderRepository interface {
	ShowOrders(ctx context.Context, id string) ([]utils.ResponseOrderDetails, error)
}
