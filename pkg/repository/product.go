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
	query := `SELECT p.model_name,p.image,b.brand_name FROM products p LEFT JOIN brands b on p.brand_id=b.id where p.deleted_at is null LIMIT $1 OFFSET $2`
	result := c.DB.Raw(query, limit, offset).Scan(&products)
	if result.Error != nil {
		return products, errors.New("failed to load products")
	}
	return products, nil
}

// func (c *ProductDatabase) FindProductById(ctx context.Context, id string) (utils.ResponseProductDetails, error) {
// 	var (
// 		Product        utils.Product
// 		ProductDetails utils.ResponseProductDetails
// 		PartDetail     utils.ResponseProducts
// 		colours        []utils.Colours
// 		//Type           utils.JersyType
// 		Size []utils.Size
// 	)
// 	//Finding the Product
// 	query := `select p.product_name from products p inner join product_details pd on p.id=pd.product_id where pd.id=?`
// 	result := c.DB.Raw(query, id).Pluck("product_name", &Product)
// 	if result.Error != nil {
// 		return ProductDetails, errors.New("failed to get data")
// 	}
// 	if Product == "Shoe" {
// 		query := `SELECT colour from available_colours where product_details_id=?`
// 		if err := c.DB.Raw(query, id).Pluck("colour", &colours).Error; err != nil {
// 			return ProductDetails, errors.New("failed to get data")
// 		}
// 		if err := c.DB.Raw("select size from available_sizes where product_id=?", id).Scan(&Size).Error; err != nil {
// 			return ProductDetails, errors.New("failed to get data")
// 		}
// 		query1 := `SELECT p.model_name,p.price,p.image,b.brand_name FROM product_details p INNER JOIN brands b on p.brand_id=b.id`
// 		result := c.DB.Raw(query1).Scan(&PartDetail)
// 		if result.Error != nil {
// 			return ProductDetails, errors.New("failed to load products")
// 		}
// 		ProductDetails.Image = PartDetail.Image
// 		ProductDetails.BrandName = PartDetail.BrandName
// 		ProductDetails.Price = PartDetail.Price
// 		ProductDetails.ModelName = PartDetail.ModelName
// 		ProductDetails.AvailableColours = colours
// 		ProductDetails.AvailableSizes = Size
// 	}
// 	// if Product == "Jersy" {
// 	// 	query := `SELECT home,away from available_colours where product_details_id=?`
// 	// 	if err := c.DB.Raw(query, id).Scan(&Type).Error; err != nil {
// 	// 		return ProductDetails, errors.New("failed to get data")
// 	// 	}
// 	// }
// 	return ProductDetails, nil
// }

// func (c *ProductDatabase) FindAvailableColours(ctx context.Context, id string) ([]utils.Colours, error) {
// 	var colours []utils.Colours
// 	//Finding the Product
// 	query := `SELECT colour from available_colours where product_details_id=?`
// 	if err := c.DB.Raw(query, id).Scan(&colours).Error; err != nil {
// 		return colours, errors.New("failed to get available colours")
// 	}
// 	return colours, nil
// }

// func (c *ProductDatabase) FindAvailableSize(ctx context.Context, id string) ([]utils.Size, error) {
// 	var Size []utils.Size
// 	query := `select size from available_sizes where product_id=?`
// 	if err := c.DB.Raw(query, id).Scan(&Size).Error; err != nil {
// 		return Size, errors.New("failed to get available sizes")
// 	}
// 	return Size, nil
// }

// func (c *ProductDatabase) FindProductDesc(ctx context.Context, id string) (utils.ResponseProducts, error) {
// 	var PartDetail utils.ResponseProducts
// 	query := `SELECT p.model_name,p.price,p.image,b.brand_name FROM product_details p INNER JOIN brands b on p.brand_id=b.id where p.id=?`
// 	result := c.DB.Raw(query, id).Scan(&PartDetail)
// 	if result.Error != nil {
// 		return PartDetail, errors.New("failed to load products description")
// 	}
// 	return PartDetail, nil
// }

// func (c *ProductDatabase) FindProductDiscount(ctx context.Context, id string) (uint, error) {
// 	var Discount uint
// 	query := `SELECT d.percentage from discounts d inner join product_details pd on pd.discount_id=d.id where pd.id=?`
// 	result := c.DB.Raw(query, id).Scan(&Discount)
// 	if result.Error != nil {
// 		return Discount, errors.New("failed to load products discount")
// 	}
// 	return Discount, nil
// }

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
