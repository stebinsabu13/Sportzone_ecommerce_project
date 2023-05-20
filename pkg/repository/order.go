package repository

import (
	"context"
	"errors"

	"github.com/stebinsabu13/ecommerce-api/pkg/domain"
	interfaces "github.com/stebinsabu13/ecommerce-api/pkg/repository/interface"
	"github.com/stebinsabu13/ecommerce-api/pkg/utils"
	"gorm.io/gorm"
)

type OrderDatabase struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) interfaces.OrderRepository {
	return &OrderDatabase{
		DB: db,
	}
}

func (c *OrderDatabase) OrderDetails(ctx context.Context, id uint) ([]utils.ResponseOrderDetails, error) {
	var orderDetails []utils.ResponseOrderDetails
	query := `SELECT p.model_name,p.price,p.image,b.brand_name,od.quantity,ad.house_name,ad.street,ad.city,ad.state,ad.country,ad.pincode,os.status FROM product_details p INNER JOIN brands b on b.id=p.brand_id INNER JOIN order_details od on p.id=od.product_details_id inner join addresses ad on ad.id=od.address_id inner join order_statuses os on os.id=od.order_status_id where od.user_id=?`
	err := c.DB.Raw(query, id).Scan(&orderDetails).Error
	if err != nil {
		return orderDetails, err
	}
	return orderDetails, nil
}

func (c *OrderDatabase) AddtoOrders(items []utils.ResCartItems, order domain.Order) error {
	var stock uint
	tx := c.DB.Begin()
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return err
	}
	for _, v := range items {
		orderitem := domain.OrderDetails{
			OrderID:         order.ID,
			OrderStatusID:   3,
			ProductDetailID: v.ProductDetailID,
			Quantity:        v.Quantity,
		}
		if err := tx.Create(&orderitem).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Model(&domain.ProductDetails{}).Where("id=?", v.ProductDetailID).Select("stock").Scan(&stock).Error; err != nil {
			tx.Rollback()
			return err
		}
		if int(stock-v.Quantity) < 0 {
			tx.Rollback()
			return errors.New("can't place orders out of stock product in the cart please remove and come again")
		}
		newstock := stock - v.Quantity
		if err := tx.Model(&domain.ProductDetails{}).Where("id=?", v.ProductDetailID).UpdateColumn("stock", newstock).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	query := `delete from cart_items where cart_id=$1`
	if err := tx.Exec(query, items[0].CartID).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (c *OrderDatabase) Findcartitems(id uint) ([]utils.ResCartItems, error) {
	var cartitems []utils.ResCartItems
	if err := c.DB.Model(&domain.CartItem{}).Where("cart_id=?", id).Select("cart_id,product_detail_id,quantity").Scan(cartitems).Error; err != nil {
		return cartitems, err
	}
	return cartitems, nil
}
