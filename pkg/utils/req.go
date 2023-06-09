package utils

import "time"

type BodyLogin struct {
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

type Otpverify struct {
	Otp         string `json:"otp" binding:"required"`
	OtpID       string `json:"otpid"`
	NewPassword string `json:"newpassword"`
	ReferalCode string `json:"referalcode"`
}

type OtpLogin struct {
	Email     string `json:"email"`
	MobileNum string `json:"mobilenum"`
}

type Pagination struct {
	Offset uint `json:"offset"`
	Limit  uint `json:"limit"`
}

type EditProfileReq struct {
	FirstName string `json:"firstname" binding:"required"`
	LastName  string `json:"lastname" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
}

type SalesReport struct {
	Month     time.Month `json:"startdate"`
	Year      int        `json:"year"`
	Frequency string     `json:"frequency"`
	EndDate   time.Time  `json:"enddate"`
	Pagination
}

type BodyAddCoupon struct {
	Code           string `json:"code" binding:"required"`
	Type           uint   `json:"type" binding:"required"`
	Discount       uint   `json:"discount" binding:"required"`
	UsageLimit     uint   `json:"usagelimit" binding:"required"`
	ExpirationDate string `json:"expdate" binding:"required"`
	MinOrderAmount *uint  `json:"minorderamount"`
	ProductID      *int   `json:"productid"`
}

type BodySignUpuser struct {
	FirstName       string `json:"firstname" binding:"required"`
	LastName        string `json:"lastname" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	MobileNum       string `json:"mobilenum" binding:"required,min=10,max=10"`
	Password        string `json:"password" binding:"required,eqfield=ConfirmPassword"`
	ConfirmPassword string `json:"confirmpassword" binding:"required"`
	ReferalCode     string `json:"referalcode"`
}

type AddCategory struct {
	CategoryName string `json:"categoryname" binding:"required"`
}
type AddAddress struct {
	HouseName string `json:"housename" binding:"required"`
	Street    string `json:"street" binding:"required"`
	City      string `json:"city" binding:"required" `
	State     string `json:"state" binding:"required" `
	Country   string `json:"country" binding:"required"`
	Pincode   string `json:"pincode" binding:"required" `
}

type AddProduct struct {
	ModelName  string `json:"modelname" binding:"required"`
	Image      string `json:"image" binding:"required"`
	BrandID    uint   `json:"brandid"`
	CategoryID uint   `json:"categoryid" binding:"required"`
}

type AddProductDetail struct {
	Price             uint `json:"price" binding:"required"`
	Stock             uint `json:"stock" binding:"required"`
	AvailableSizeID   uint `json:"availablesizeid" binding:"required"`
	AvailableColourID uint `json:"availablecolourid" binding:"required"`
	ProductID         uint `json:"productid" binding:"required"`
	DiscountID        uint `json:"discountid"`
}
