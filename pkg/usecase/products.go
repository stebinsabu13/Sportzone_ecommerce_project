package usecase

import (
	"context"

	"github.com/stebinsabu13/ecommerce-api/pkg/domain"
	interfaces "github.com/stebinsabu13/ecommerce-api/pkg/repository/interface"
	services "github.com/stebinsabu13/ecommerce-api/pkg/usecase/interface"
	"github.com/stebinsabu13/ecommerce-api/pkg/utils"
)

type productUseCase struct {
	Productrepo interfaces.ProductRepository
}

func NewProductUseCase(repo interfaces.ProductRepository) services.ProductUseCase {
	return &productUseCase{
		Productrepo: repo,
	}
}

func (c *productUseCase) FindAllProducts(ctx context.Context, pagination utils.Pagination) ([]utils.ResponseProducts, error) {
	products, err := c.Productrepo.FindAllProducts(ctx, pagination)
	return products, err
}

func (c *productUseCase) AddProduct(ctx context.Context, product domain.Product) error {
	err := c.Productrepo.AddProduct(ctx, product)
	return err
}

func (c *productUseCase) EditProduct(ctx context.Context, product domain.Product, id string) error {
	err := c.Productrepo.EditProduct(ctx, product, id)
	return err
}

func (c *productUseCase) DeleteProduct(ctx context.Context, id string) error {
	err := c.Productrepo.DeleteProduct(ctx, id)
	return err
}

func (c *productUseCase) FindProductById(ctx context.Context, id string, pagination utils.Pagination) ([]utils.ResponseProductDetails, error) {
	productDetails, err := c.Productrepo.FindProductById(ctx, id, pagination)
	return productDetails, err
}

func (c *productUseCase) AddProductDetail(ctx context.Context, productdetail domain.ProductDetails) error {
	return c.Productrepo.AddProductDetail(ctx, productdetail)
}

func (c *productUseCase) EditProductDetail(ctx context.Context, productdetail domain.ProductDetails, id string) error {
	err := c.Productrepo.EditProductDetail(ctx, productdetail, id)
	return err
}

func (c *productUseCase) DeleteProductDetail(ctx context.Context, id string) error {
	err := c.Productrepo.DeleteProductDetail(ctx, id)
	return err
}

func (c *productUseCase) ProductsByCategory(id string, pagination utils.Pagination) ([]utils.ResponseProducts, error) {
	return c.Productrepo.ProductsByCategory(id, pagination)
}

func (c *productUseCase) ListAllBrands() ([]utils.ResBrands, error) {
	return c.Productrepo.ListAllBrands()
}

func (c *productUseCase) ProductsByBrands(id string, pagination utils.Pagination) ([]utils.ResponseProducts, error) {
	return c.Productrepo.ProductsByBrands(id, pagination)
}
