package utils

type BodyLogin struct {
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

type Otpverify struct {
	Otp   string `json:"otp" binding:"required"`
	OtpID string `json:"otpid" binding:"required"`
}

type OtpLogin struct {
	Email     string
	MobileNum string
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
