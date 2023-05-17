package repository

import (
	"gorm.io/gorm"
)

type OrderDatabase struct {
	DB *gorm.DB
}

// func (c *OrderDatabase) ShowOrders(ctx context.Context, id string) ([]utils.ResponseOrderDetails, error) {

// }
