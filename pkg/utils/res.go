package utils

import "time"

type ResponseCategory string
type ResBrands string

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
	ID           uint   `json:"id"`
	Image        string `json:"image"`
	ModelName    string `json:"modelname"`
	BrandName    string `json:"brandname"`
	CategoryName string `json:"categoryname"`
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
	Image      string `json:"image"`
	ModelName  string `json:"modelname"`
	Price      uint   `json:"price"`
	BrandName  string `json:"brandname"`
	Size       string `json:"size"`
	Colour     string `json:"colour"`
	Quantity   uint   `json:"quantity"`
	Percentage int    `json:"discountpercentage"`
	Total      uint   `json:"total"`
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
	Percentage    int        `json:"discountpercentage"`
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
	RazorpayKey     string `json:"razorpaykey"`
	AmountToPay     uint   `json:"amounttopay"`
	RazorpayAmount  int    `json:"razorpayamount"`
	RazorpayOrderID string `json:"razorpayorderid"`
	UserID          uint   `json:"userid"`
}

//salesreport
type ResSalesReport struct {
	UserID          uint      `json:"userid"`
	FirstName       string    `json:"firstname"`
	Email           string    `json:"email"`
	ProductDetailID uint      `json:"productdetailid"`
	ProductName     string    `json:"productname"`
	Quantity        uint      `json:"quantity"`
	OrderID         uint      `json:"orderid"`
	PlacedDate      time.Time `json:"placeddate"`
	PaymentMode     string    `json:"paymentmode"`
}

type ResWidgets struct {
	Numberofblockedusers  int `json:"numberofblockedusers"`
	Numberofproducts      int `json:"numberofproducts"`
	Numberofpendingorders int `json:"numberofpendingorders"`
}
