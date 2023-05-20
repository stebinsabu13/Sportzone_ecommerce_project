package domain

type Cart struct {
	ID         uint `json:"id" gorm:"primarykey;auto_increment"`
	UserID     uint `json:"userid" gorm:"not null"`
	GrandTotal int  `json:"grandtotal" gorm:"default:0"`
}

type CartItem struct {
	ID              uint `json:"id" gorm:"primarykey;auto_increment"`
	CartID          uint `json:"cartid" gorm:"not null"`
	ProductDetailID uint `json:"productdetailid" gorm:"not null"`
	Quantity        uint `json:"quantity" gorm:"not null"`
	Total           uint `json:"total" gorm:"not null"`
}
