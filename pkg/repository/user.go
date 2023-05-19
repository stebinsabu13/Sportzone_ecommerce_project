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
func (c *userDatabase) FindbyUserID(ctx context.Context, id uint) (domain.User, error) {
	var user domain.User
	_ = c.DB.Where("id=?", id).Find(&user)
	if user.ID == 0 {
		return domain.User{}, errors.New("invalid user id")
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
	if err != nil {
		return user.MobileNum, err
	}
	return user.MobileNum, nil
}

func (c *userDatabase) ShowDetails(ctx context.Context, id uint) (utils.ResponseUsers, error) {
	var user utils.ResponseUsers
	query := `SELECT first_name,last_name,email,mobile_num from users where id=?`
	if err := c.DB.Raw(query, id).Scan(&user).Error; err != nil {
		return user, err
	}
	return user, nil

}

func (c *userDatabase) ShowAddress(ctx context.Context, id uint) ([]utils.Address, error) {
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

func (c *userDatabase) AddAddress(ctx context.Context, address domain.Address) error {
	err := c.DB.Create(&address).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *userDatabase) EditProfile(ctx context.Context, profile utils.EditProfileReq, id uint) error {
	result := c.DB.Model(&domain.User{}).Where("id=?", id).Updates(domain.User{FirstName: profile.FirstName, LastName: profile.LastName, Email: profile.Email})
	if result.RowsAffected == 0 {
		return errors.New("no row updated")
	} else if result.Error != nil {
		return result.Error
	}
	return nil
}

func (c *userDatabase) ChangePassword(ctx context.Context, newpassword string, mobile_num string) error {
	result := c.DB.Model(&domain.User{}).Where("mobile_num=?", mobile_num).UpdateColumn("password", newpassword)
	if result.RowsAffected == 0 {
		return errors.New("no row updated")
	} else if result.Error != nil {
		return result.Error
	}
	return nil
}
