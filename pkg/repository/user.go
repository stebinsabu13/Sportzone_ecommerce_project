package repository

import (
	"context"
	"errors"

	domain "github.com/stebinsabu13/ecommerce-api/pkg/domain"
	interfaces "github.com/stebinsabu13/ecommerce-api/pkg/repository/interface"
	"gorm.io/gorm"
)

type userDatabase struct {
	DB *gorm.DB
}

func NewUserRepository(Db *gorm.DB) interfaces.UserRepository {
	return &userDatabase{DB: Db}
}

func (c *userDatabase) FindAll(ctx context.Context) ([]domain.User, error) {
	var users []domain.User
	err := c.DB.Find(&users).Error
	return users, err
}

func (c *userDatabase) FindbyEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User
	_ = c.DB.Where("email=?", email).Find(&user)
	if user.ID == 0 {
		return domain.User{}, errors.New("invalid email")
	}
	return user, nil
}

func (c *userDatabase) SignUpUser(ctx context.Context, user domain.User) error {
	err := c.DB.Create(&user).Error
	return err
}
