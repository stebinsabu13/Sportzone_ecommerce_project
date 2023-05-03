package utils

type BodyLogin struct {
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

type Otpverify struct {
	Otp string `binding:"required"`
}
