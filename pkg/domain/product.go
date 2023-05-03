package domain

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	ProductName string `json:"productname" gorm:"not null"`
	CategoryId  uint   `json:"categoryid" gorm:"foreignKey:CategoryID"`
	InventoryId uint   `json:"inventoryid" gorm:"foreignkey:InventoryID"`
	CouponId    uint   `json:"couponid" gorm:"foreignkey:CouponID"`
}

type ProductDetails struct {
	gorm.Model
	ModelName  string `json:"modelname" gorm:"not null"`
	Price      uint   `json:"price" gorm:"not null"`
	Image      string `json:"image" gorm:"not null"`
	ProductId  uint   `json:"productid" gorm:"foreignkey:ProductID"`
	DiscountId uint   `json:"discountid" gorm:"foreignkey:DiscountID"`
	BrandId    uint   `json:"brandid" gorm:"foreignkey:BrandID"`
}

type Brand struct {
	ID        uint   `json:"id" gorm:"primary_key;auto_increment"`
	BrandName string `json:"brandname" gorm:"not null"`
}
type Category struct {
	gorm.Model
	CategoryName string `json:"categoryname" gorm:"not null"`
}

type Inventory struct {
	ID       uint `json:"id" gorm:"primary_key;auto_increment"`
	Quantity uint `json:"quantity" gorm:"not null"`
}

type Coupon struct {
	ID         uint `json:"id" gorm:"primary_key;auto_increment"`
	Percentage uint `json:"percentage" gorm:"not null"`
}

type Discount struct {
	ID         uint `json:"id" gorm:"primary_key;auto_increment"`
	Percentage uint `json:"percentage" gorm:"not null"`
}

type AvailableColour struct {
	ID               uint   `json:"id" gorm:"primary_key;auto_increment"`
	Colour           string `json:"colour"`
	Home             bool   `json:"home"`
	Away             bool   `json:"away"`
	ProductDetailsId uint   `json:"productdetailsid" gorm:"foreignkey:ProductDetailsID"`
}

type AvailableSize struct {
	ID        uint   `json:"id" gorm:"primary_key;auto_increment"`
	Size      string `json:"size" gorm:"not null"`
	ProductId uint   `json:"productid" gorm:"foreignkey:ProductID"`
}
