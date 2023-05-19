package interfaces

import (
	"context"

	"github.com/stebinsabu13/ecommerce-api/pkg/utils"
)

type OrderUseCase interface {
	OrderDetails(ctx context.Context, id uint) ([]utils.ResponseOrderDetails, error)
}
