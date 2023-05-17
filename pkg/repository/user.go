package repository

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/stebinsabu13/ecommerce-api/pkg/domain"
	interfaces "github.com/stebinsabu13/ecommerce-api/pkg/repository/interface"
	"github.com/stebinsabu13/ecommerce-api/pkg/utils"
	"gorm.io/gorm"
)

type userDatabase struct {
	DB *gorm.DB
}

func NewUserRepository(Db *gorm.DB) interfaces.UserRepository {
	return &userDatabase{DB: Db}
}

// func (c *userDatabase) FindAll(ctx context.Context) ([]domain.User, error) {
// 	var users []domain.User
// 	err := c.DB.Find(&users).Error
// 	return users, err
// }

func (c *userDatabase) FindbyEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User
	_ = c.DB.Where("email=?", email).Find(&user)
	if user.ID == 0 {
		return domain.User{}, errors.New("invalid email")
	}
	if user.Block {
		return user, errors.New("you are blocked")
	}
	if !user.Verified {
		return user, errors.New("you need to verify your mobile number")
	}
	return user, nil
}

func (c *userDatabase) FindbyEmailorMobilenum(ctx context.Context, body utils.OtpLogin) (domain.User, error) {
	var user domain.User
	_ = c.DB.Where("email=? or mobile_num=?", body.Email, body.MobileNum).Find(&user)
	fmt.Println(user)
	if user.ID == 0 {
		return domain.User{}, errors.New("user not exsists")
	}
	if user.Block {
		return user, errors.New("you are blocked")
	}
	if !user.Verified {
		return user, errors.New("you need to verify your mobile number")
	}
	return user, nil
}

func (c *userDatabase) SignUpUser(ctx context.Context, user domain.User) (string, error) {
	err := c.DB.Create(&user).Error
	return user.MobileNum, err
}

func (c *userDatabase) ShowDetails(ctx context.Context, id int) (utils.ResponseUsers, error) {
	var user utils.ResponseUsers
	query := `SELECT first_name,last_name,email,mobile_num from users where id=?`
	if err := c.DB.Raw(query, id).Scan(&user).Error; err != nil {
		return user, err
	}
	return user, nil

}

func (c *userDatabase) ShowAddress(ctx context.Context, id int) ([]utils.Address, error) {
	var address []utils.Address
	query := `select id,house_name,street,city,state,country,pincode from addresses where user_id=?`
	if err := c.DB.Raw(query, id).Scan(&address).Error; err != nil {
		return address, err
	}
	return address, nil
}

func (c *userDatabase) UpdateVerify(ctx context.Context, number string) error {
	err := c.DB.Model(&domain.User{}).Where("mobile_num=?", number).UpdateColumn("verified", true).Error
	if err != nil {
		return err
	}
	return nil
}
