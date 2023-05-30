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

// type BodyCoupon struct {
// 	Code string `json:"code"`
// }
