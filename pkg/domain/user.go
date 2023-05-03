package domain

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string `json:"firstname" binding:"required" gorm:"not null"`
	LastName  string `json:"lastname" binding:"required" gorm:"not null"`
	Email     string `json:"email" binding:"required" gorm:"uniqueIndex;not null"`
	MobileNum string `json:"mobilenum" binding:"required" gorm:"uniqueIndex;not null"`
	Password  string `json:"password" binding:"required" gorm:"not null"`
}
