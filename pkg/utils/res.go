package utils

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

// struct used to view the order_details
type ResponseOrderDetails struct {
	Image     string `json:"image"`
	ModelName string `json:"modelname"`
	Price     uint   `json:"price"`
	BrandName string `json:"brandname"`
	Quantity  uint   `json:"quantity"`
	HouseName string `json:"housename"`
	Street    string `json:"street"`
	City      string `json:"city"`
	State     string `json:"state"`
	Country   string `json:"country"`
	Pincode   string `json:"pincode"`
	Status    string `json:"status"`
}

type ResViewCart struct {
	Image     string `json:"image"`
	ModelName string `json:"modelname"`
	Price     uint   `json:"price"`
	BrandName string `json:"brandname"`
	Quantity  uint   `json:"quantity"`
	Total     uint   `json:"total"`
}

type ResCartItems struct {
	CartID          uint `json:"cartid"`
	ProductDetailID uint `json:"productdetailid"`
	Quantity        uint `json:"quantity"`
}
