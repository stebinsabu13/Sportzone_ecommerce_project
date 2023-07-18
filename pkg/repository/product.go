package repository

import (
	"context"
	"errors"

	"github.com/stebinsabu13/ecommerce-api/pkg/domain"
	interfaces "github.com/stebinsabu13/ecommerce-api/pkg/repository/interface"
	"github.com/stebinsabu13/ecommerce-api/pkg/utils"
	"gorm.io/gorm"
)

type ProductDatabase struct {
	DB *gorm.DB
}

func NewProductrepository(Db *gorm.DB) interfaces.ProductRepository {
	return &ProductDatabase{
		DB: Db,
	}
}
func (c *ProductDatabase) FindAllProducts(ctx context.Context, pagination utils.Pagination) ([]utils.ResponseProducts, error) {
	var products []utils.ResponseProducts
	offset := pagination.Offset
	limit := pagination.Limit
	query := `SELECT p.id,p.model_name,p.image,b.brand_name,c.category_name FROM products p
	LEFT JOIN brands b on p.brand_id=b.id
	INNER JOIN categories c on p.category_id=c.id where p.deleted_at is null LIMIT $1 OFFSET $2`
	result := c.DB.Raw(query, limit, offset).Scan(&products)
	if result.Error != nil {
		return products, errors.New("failed to load products")
	}
	return products, nil
}

func (c *ProductDatabase) AddProduct(ctx context.Context, product domain.Product) error {
	result := c.DB.Create(&product).Error
	if result != nil {
		return errors.New("failed to add product")
	}
	return nil
}

func (c *ProductDatabase) EditProduct(ctx context.Context, product domain.Product, id string) error {
	result := c.DB.Where("id=?", id).UpdateColumns(&product).Error
	if result != nil {
		return errors.New("failed to update product")
	}
	return nil
}

func (c *ProductDatabase) DeleteProduct(ctx context.Context, id string) error {
	result := c.DB.Where("id=?", id).Delete(&domain.Product{}).Error
	if result != nil {
		return errors.New("failed to delete")
	}
	return nil
}

func (c *ProductDatabase) FindProductById(ctx context.Context, id string, pagination utils.Pagination) ([]utils.ResponseProductDetails, error) {
	var Product []utils.ResponseProductDetails
	offset := pagination.Offset
	limit := pagination.Limit
	//Finding the Product
	query := `select p.model_name,p.image,b.brand_name,pd.stock,pd.price,c.colour,s.size,d.percentage from products p
	left join brands b on b.id=p.brand_id
	inner join product_details pd on pd.product_id=p.id
	inner join available_colours c on c.id=pd.available_colour_id
	inner join available_sizes s on s.id=pd.available_size_id
	left join discounts d on d.id=pd.discount_id where p.id=$1 and pd.deleted_at is null LIMIT $2 OFFSET $3`
	result := c.DB.Raw(query, id, limit, offset).Scan(&Product)
	if result.Error != nil {
		return Product, errors.New("failed to get product")
	}
	return Product, nil
}

func (c *ProductDatabase) AddProductDetail(ctx context.Context, productdetail domain.ProductDetails) error {
	result := c.DB.Create(&productdetail).Error
	if result != nil {
		return result
	}
	return nil
}

func (c *ProductDatabase) EditProductDetail(ctx context.Context, productdetail domain.ProductDetails, id string) error {
	result := c.DB.Where("id=?", id).UpdateColumns(&productdetail).Error
	if result != nil {
		return result
	}
	return nil
}

func (c *ProductDatabase) DeleteProductDetail(ctx context.Context, id string) error {
	result := c.DB.Where("id=?", id).Delete(&domain.ProductDetails{}).Error
	if result != nil {
		return result
	}
	return nil
}

func (c *ProductDatabase) ProductsByCategory(id string, pagination utils.Pagination) ([]utils.ResponseProducts, error) {
	var products []utils.ResponseProducts
	offset := pagination.Offset
	limit := pagination.Limit
	query := `SELECT p.id,p.model_name,p.image,b.brand_name,c.category_name FROM products p
	LEFT JOIN brands b on p.brand_id=b.id
	INNER JOIN categories c on p.category_id=c.id where p.deleted_at is null and p.category_id=$1 LIMIT $2 OFFSET $3`
	result := c.DB.Raw(query, id, limit, offset).Scan(&products)
	if result.Error != nil {
		return products, errors.New("failed to load products")
	}
	return products, nil
}

func (c *ProductDatabase) ListAllBrands() ([]utils.ResBrands, error) {
	var brands []utils.ResBrands
	query := `Select brand_name from brands`
	result := c.DB.Raw(query).Scan(&brands).Error
	if result != nil {
		return brands, errors.New("failed to get brands")
	}
	return brands, nil
}

func (c *ProductDatabase) ProductsByBrands(id string, pagination utils.Pagination) ([]utils.ResponseProducts, error) {
	var products []utils.ResponseProducts
	offset := pagination.Offset
	limit := pagination.Limit
	query := `SELECT p.id,p.model_name,p.image,b.brand_name,c.category_name FROM products p
		LEFT JOIN brands b on p.brand_id=b.id
		INNER JOIN categories c on p.category_id=c.id where p.deleted_at is null and p.brand_id=$1 LIMIT $2 OFFSET $3`
	result := c.DB.Raw(query, id, limit, offset).Scan(&products)
	if result.Error != nil {
		return products, errors.New("failed to load products")
	}
	return products, nil
}
