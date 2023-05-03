package repository

import (
	"context"
	"errors"

	domain "github.com/stebinsabu13/ecommerce-api/pkg/domain"
	interfaces "github.com/stebinsabu13/ecommerce-api/pkg/repository/interface"
	"gorm.io/gorm"
)

type adminDatabase struct {
	DB *gorm.DB
}

func NewAdminRepository(Db *gorm.DB) interfaces.AdminRepository {
	return &adminDatabase{DB: Db}
}

func (c *adminDatabase) FindAll(ctx context.Context) ([]domain.Admin, error) {
	var admin []domain.Admin
	err := c.DB.Find(&admin).Error
	return admin, err
}

func (c *adminDatabase) FindbyEmail(ctx context.Context, email string) (domain.Admin, error) {
	var admin domain.Admin
	_ = c.DB.Where("email=?", email).Find(&admin)
	if admin.ID == 0 {
		return domain.Admin{}, errors.New("invalid email")
	}
	return admin, nil
}

// func (c *adminDatabase) SignUpAdmin(ctx context.Context, admin domain.Admin) error {
// 	err := c.DB.Create(&admin).Error
// 	return err
// }
