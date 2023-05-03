package interfaces

import (
	"context"

	"github.com/stebinsabu13/ecommerce-api/pkg/utils"
)

type ProductRepository interface {
	FindAllProducts(ctx context.Context) ([]utils.ResponseProducts, error)
	// FindProductById(ctx context.Context, id string) (utils.ResponseProductDetails, error)
	// FindProduct(ctx context.Context, id string) (utils.Product, error)
	FindAvailableColours(ctx context.Context, id string) ([]utils.Colours, error)
	FindAvailableSize(ctx context.Context, id string) ([]utils.Size, error)
	FindProductDesc(ctx context.Context, id string) (utils.ResponseProducts, error)
	FindProductDiscount(ctx context.Context, id string) (uint, error)
}
