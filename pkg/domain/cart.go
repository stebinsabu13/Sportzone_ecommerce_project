package domain

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	UserID uint `json:"userid"`
	User   User `gorm:"foreignkey:UserID"`
}
