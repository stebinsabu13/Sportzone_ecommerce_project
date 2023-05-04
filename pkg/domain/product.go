package domain

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	ProductName string   `json:"productname" gorm:"not null"`
	CategoryID  uint     `json:"categoryid"`
	Category    Category `gorm:"foreignKey:CategoryID"`
	CouponID    uint     `json:"couponid"`
	Coupon      Coupon   `gorm:"foreignkey:CouponID"`
}

type ProductDetails struct {
	gorm.Model
	ModelName   string    `json:"modelname" gorm:"not null"`
	Price       uint      `json:"price" gorm:"not null"`
	Image       string    `json:"image" gorm:"not null"`
	ProductID   uint      `json:"productid"`
	Product     Product   `gorm:"foreignkey:ProductID"`
	DiscountID  uint      `json:"discountid"`
	Discount    Discount  `gorm:"foreignkey:DiscountID"`
	BrandID     uint      `json:"brandid"`
	Brand       Brand     `gorm:"foreignkey:BrandID"`
	InventoryID uint      `json:"inventoryid"`
	Inventory   Inventory `gorm:"foreignkey:InventoryID"`
}

type Brand struct {
	ID        uint   `json:"id" gorm:"primarykey;auto_increment"`
	BrandName string `json:"brandname" gorm:"not null"`
}

type Category struct {
	gorm.Model
	CategoryName string `json:"categoryname" gorm:"not null"`
}

type Inventory struct {
	ID       uint `json:"id" gorm:"primarykey;auto_increment"`
	Quantity uint `json:"quantity" gorm:"not null"`
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
	ID               uint           `json:"id" gorm:"primarykey;auto_increment"`
	Colour           string         `json:"colour"`
	Home             bool           `json:"home"`
	Away             bool           `json:"away"`
	ProductDetailsID uint           `json:"productdetailsid"`
	ProductDetails   ProductDetails `gorm:"foreignkey:ProductDetailsID"`
}

type AvailableSize struct {
	ID        uint    `json:"id" gorm:"primarykey;auto_increment"`
	Size      string  `json:"size" gorm:"not null"`
	ProductID uint    `json:"productid"`
	Product   Product `gorm:"foreignkey:ProductID"`
}
