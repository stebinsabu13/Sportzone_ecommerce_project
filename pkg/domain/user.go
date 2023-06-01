package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName   string `json:"firstname" binding:"required" gorm:"not null"`
	LastName    string `json:"lastname" binding:"required" gorm:"not null"`
	Email       string `json:"email" binding:"required" gorm:"uniqueIndex;not null"`
	MobileNum   string `json:"mobilenum" binding:"required" gorm:"uniqueIndex;not null"`
	Password    string `json:"password" binding:"required" gorm:"not null"`
	Block       bool   `json:"block" gorm:"default:false"`
	Verified    bool   `json:"verified" gorm:"default:false"`
	ReferalCode string `json:"referalcode" gorm:"uniqueIndex"`
}

type Address struct {
	gorm.Model
	HouseName string `json:"housename" binding:"required" gorm:"not null"`
	Street    string `json:"street" binding:"required" gorm:"not null"`
	City      string `json:"city" binding:"required" gorm:"not null"`
	State     string `json:"state" binding:"required" gorm:"not null"`
	Country   string `json:"country" binding:"required" gorm:"not null"`
	Pincode   string `json:"pincode" binding:"required" gorm:"not null"`
	UserID    uint   `json:"userid"`
	// User      User   `gorm:"foreignkey:UserID"`
}

type PaymentMode struct {
	ID   uint   `json:"id" gorm:"primarykey;auto_increment"`
	Mode string `json:"mode" gorm:"not null"`
}

type Wallet struct {
	ID           uint       `json:"id" gorm:"primarykey;auto_increment"`
	UserID       uint       `json:"userid" gorm:"not null"`
	CreditedDate *time.Time `json:"crediteddate"`
	DebitedDate  *time.Time `json:"debiteddate"`
	Amount       int        `json:"amount" gorm:"not null"`
}
