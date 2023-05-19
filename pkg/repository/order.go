package repository

import (
	"context"

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
