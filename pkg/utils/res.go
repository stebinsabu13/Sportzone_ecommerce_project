package utils

type Colours string

type Size string

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

// type Product string

// struct used to list all products
type ResponseProducts struct {
	Image     string `json:"image"`
	ModelName string `json:"modelname"`
	Price     uint   `json:"price"`
	BrandName string `json:"brandname"`
}

// struct used to view a particular product details
type ResponseProductDetails struct {
	Image              string    `json:"image"`
	ModelName          string    `json:"modelname"`
	Price              uint      `json:"price"`
	BrandName          string    `json:"brandname"`
	DiscountPercentage string    `json:"discountpercentage"`
	AvailableColours   []Colours `json:"availablecolours"`
	AvailableSizes     []Size    `json:"availablesizes"`
	PayableAmount      uint      `json:"payableamount"`
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
