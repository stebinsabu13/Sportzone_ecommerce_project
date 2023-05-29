package domain

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	CategoryName string `json:"categoryname" gorm:"not null"`
}

type Product struct {
	gorm.Model
	ModelName  string   `json:"modelname" gorm:"not null"`
	Image      string   `json:"image" gorm:"not null"`
	BrandID    uint     `json:"brandid"`
	Brand      Brand    `gorm:"foreignkey:BrandID"`
	CategoryID uint     `json:"categoryid" gorm:"not null"`
	Category   Category `gorm:"foreignkey:CategoryID"`
}
type ProductDetails struct {
	gorm.Model
	Price             uint            `json:"price" gorm:"not null"`
	Stock             uint            `json:"stock" gorm:"not null"`
	AvailableSizeID   uint            `json:"availablesizeid" gorm:"not null"`
	AvailableSize     AvailableSize   `gorm:"foreignkey:AvailableSizeID"`
	AvailableColourID uint            `json:"availablecolourid" gorm:"not null"`
	AvailableColour   AvailableColour `gorm:"foreignkey:AvailableColourID"`
	ProductID         uint            `json:"productid" gorm:"not null"`
	Product           Product         `gorm:"foreignkey:ProductID"`
	DiscountID        uint            `json:"discountid"`
	Discount          Discount        `gorm:"foreignkey:DiscountID"`
	CouponID          uint            `json:"couponid"`
	Coupon            Coupon          `gorm:"foreignkey:CouponID"`
}

type Brand struct {
	ID        uint   `json:"id" gorm:"primarykey;auto_increment"`
	BrandName string `json:"brandname" gorm:"not null"`
}

type Coupon struct {
	ID         uint `json:"id" gorm:"primarykey;auto_increment"`
	Percentage uint `json:"percentage" gorm:"not null"`
}

type Discount struct {
	ID         uint `json:"id" gorm:"primarykey;auto_increment"`
	Percentage uint `json:"percentage" gorm:"not null"`
}

type AvailableColour struct {
	ID     uint   `json:"id" gorm:"primarykey;auto_increment"`
	Colour string `json:"colour"`
}

type AvailableSize struct {
	ID   uint   `json:"id" gorm:"primarykey;auto_increment"`
	Size string `json:"size" gorm:"not null"`
}
