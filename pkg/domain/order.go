package domain

import (
	"gorm.io/gorm"
)

type OrderDetails struct {
	gorm.Model
	UserID       uint       `json:"userid"`
	User         User       `gorm:"foreignkey:UserID"`
	AddressID    uint       `json:"addressid"`
	Address      Address    `gorm:"foreignkey:AddressID"`
	OrderItemsID uint       `json:"orderitemsid"`
	OrderItems   OrderItems `gorm:"foreignkey:OrderItemsID"`
}
type OrderItems struct {
	gorm.Model
	Quantity         uint            `json:"quantity" gorm:"not null"`
	ProductDetailsID uint            `json:"productdetailsid"`
	ProductDetails   ProductDetails  `gorm:"foreignkey:ProductDetailsID"`
	ColourID         uint            `json:"colourid"`
	AvailableColour  AvailableColour `gorm:"foreignkey:ColourID"`
	SizeID           uint            `json:"sizeid"`
	AvailableSize    AvailableSize   `gorm:"foreignkey:SizeID"`
}
