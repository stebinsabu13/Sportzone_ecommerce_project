package utils

import "time"

type ResponseCategory string

type Address struct {
	ID        uint   `json:"id"`
	HouseName string `json:"housename"`
	Street    string `json:"street"`
	City      string `json:"city"`
	State     string `json:"state"`
	Country   string `json:"country"`
	Pincode   string `json:"pincode"`
}

// struct used to list all products
type ResponseProducts struct {
	Image     string `json:"image"`
	ModelName string `json:"modelname"`
	BrandName string `json:"brandname"`
}

// struct used to view a particular product details
type ResponseProductDetails struct {
	Image      string `json:"image"`
	ModelName  string `json:"modelname"`
	Price      uint   `json:"price"`
	BrandName  string `json:"brandname"`
	Stock      uint   `json:"stock"`
	Percentage uint   `json:"percentage"`
	Colour     string `json:"colour"`
	Size       string `json:"size"`
}

// struct used to list all users
type ResponseUsers struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	MobileNum string `json:"mobilenum"`
	Block     bool   `json:"blocked"`
}

// struct used to view a particular user detail
type ResponseUserDetails struct {
	FirstName string    `json:"firstname"`
	LastName  string    `json:"lastname"`
	Email     string    `json:"email"`
	MobileNum string    `json:"mobilenum"`
	Address   []Address `json:"address"`
}

// struct used to view cart
type ResViewCart struct {
	Image     string `json:"image"`
	ModelName string `json:"modelname"`
	Price     uint   `json:"price"`
	BrandName string `json:"brandname"`
	Size      string `json:"size"`
	Colour    string `json:"colour"`
	Quantity  uint   `json:"quantity"`
	Total     uint   `json:"total"`
}

type ResCartItems struct {
	CartID          uint `json:"cartid"`
	ProductDetailID uint `json:"productdetailid"`
	Quantity        uint `json:"quantity"`
}

// struct used to view orders
type ResOrders struct {
	ID         uint      `json:"id"`
	PlacedDate time.Time `json:"placeddate"`
	HouseName  string    `json:"housename"`
	Street     string    `json:"street"`
	City       string    `json:"city"`
	State      string    `json:"state"`
	Country    string    `json:"country"`
	Pincode    string    `json:"pincode"`
	Mode       string    `json:"mode"`
	GrandTotal uint      `json:"grandtotal"`
}

// struct used to view the order_details
type ResponseOrderDetails struct {
	ID            uint       `json:"id"`
	Image         string     `json:"image"`
	ModelName     string     `json:"modelname"`
	Price         uint       `json:"price"`
	BrandName     string     `json:"brandname"`
	Size          string     `json:"size"`
	Colour        string     `json:"colour"`
	Quantity      uint       `json:"quantity"`
	Status        string     `json:"status"`
	DeliveredDate *time.Time `json:"delivereddate"`
	CancelledDate *time.Time `json:"cancelleddate"`
}

//Admin side order list
type ResAllOrders struct {
	ID         uint      `json:"id"`
	FirstName  string    `json:"firstname"`
	MobileNum  string    `json:"mobilenum"`
	PlacedDate time.Time `json:"placeddate"`
	HouseName  string    `json:"housename"`
	Street     string    `json:"street"`
	Pincode    string    `json:"pincode"`
	Mode       string    `json:"mode"`
	GrandTotal uint      `json:"grandtotal"`
}

// razorpay
type RazorpayOrder struct {
	RazorpayKey     string `json:"razorpay_key"`
	AmountToPay     uint   `json:"amount_to_pay"`
	RazorpayAmount  uint   `json:"razorpay_amount"`
	RazorpayOrderID string `json:"razorpay_order_id"`
	UserID          uint   `json:"userid"`
}
