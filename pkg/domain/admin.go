package domain

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	Name      string `json:"name" gorm:"not null"`
	Email     string `json:"email" gorm:"uniqueIndex;not null"`
	MobileNum string `json:"mobilenum" gorm:"uniqueIndex;not null"`
	Password  string `json:"password" gorm:"not null"`
	Verified  bool   `json:"verified" gorm:"default:false"`
}
