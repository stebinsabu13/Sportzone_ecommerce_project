package domain

type OtpSession struct {
	ID        uint   `json:"id" gorm:"primarykey;auto_increment"`
	OtpID     string `json:"otpid" gorm:"not null"`
	MobileNum string `json:"mobilenum" gorm:"not null"`
}
