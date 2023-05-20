package domain

import (
	"gorm.io/gorm"
)

type OrderDetails struct {
	gorm.Model
	UserID           uint           `json:"userid"`
	User             User           `gorm:"foreignkey:UserID"`
	AddressID        uint           `json:"addressid"`
	Address          Address        `gorm:"foreignkey:AddressID"`
	ProductDetailsID uint           `json:"productdetailsid"`
	ProductDetails   ProductDetails `gorm:"foreignkey:ProductDetailsID"`
	Quantity         uint           `json:"quantity" gorm:"not null"`
	OrderStatusID    uint           `json:"orderstatusid"`
	OrderStatus      OrderStatus    `gorm:"foreignkey:StatusID"`
}

type OrderStatus struct {
	ID     uint   `json:"id" gorm:"primarykey;auto_increment"`
	Status string `json:"status" gorm:"not null"`
}
