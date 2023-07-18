package interfaces

import (
	"context"

	"github.com/stebinsabu13/ecommerce-api/pkg/domain"
	"github.com/stebinsabu13/ecommerce-api/pkg/utils"
)

type ProductUseCase interface {
	FindAllProducts(context.Context, utils.Pagination) ([]utils.ResponseProducts, error)
	FindProductById(context.Context, string, utils.Pagination) ([]utils.ResponseProductDetails, error)
	AddProduct(ctx context.Context, product domain.Product) error
	EditProduct(ctx context.Context, product domain.Product, id string) error
	DeleteProduct(ctx context.Context, id string) error
	AddProductDetail(context.Context, domain.ProductDetails) error
	EditProductDetail(context.Context, domain.ProductDetails, string) error
	DeleteProductDetail(ctx context.Context, id string) error
	ProductsByCategory(string, utils.Pagination) ([]utils.ResponseProducts, error)
	ListAllBrands() ([]utils.ResBrands, error)
	ProductsByBrands(string, utils.Pagination) ([]utils.ResponseProducts, error)
}
