package interfaces

import (
	"context"

	"github.com/stebinsabu13/ecommerce-api/pkg/domain"
	"github.com/stebinsabu13/ecommerce-api/pkg/utils"
)

type ProductUseCase interface {
	FindAllProducts(context.Context, utils.Pagination) ([]utils.ResponseProducts, error)
	FindProductById(context.Context, string, utils.Pagination) ([]utils.ResponseProductDetails, error)
	// FindProduct(ctx context.Context, id string) (utils.Product, error)
	// FindAvailableColours(ctx context.Context, id string) ([]utils.Colours, error)
	// FindAvailableSize(ctx context.Context, id string) ([]utils.Size, error)
	// FindProductDesc(ctx context.Context, id string) (utils.ResponseProducts, error)
	// FindProductDiscount(ctx context.Context, id string) (uint, error)
	AddProduct(ctx context.Context, product domain.Product) error
	EditProduct(ctx context.Context, product domain.Product, id string) error
	DeleteProduct(ctx context.Context, id string) error
	AddProductDetail(context.Context, domain.ProductDetails) error
	EditProductDetail(context.Context, domain.ProductDetails, string) error
	DeleteProductDetail(ctx context.Context, id string) error
}
