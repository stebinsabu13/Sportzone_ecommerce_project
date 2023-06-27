package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

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

func (c *userDatabase) FindbyEmail(ctx context.Context, email string) (utils.ResponseUsers, error) {
	var user utils.ResponseUsers
	query := `SELECT * from users where email=$1`
	c.DB.Raw(query, email).Scan(&user)
	// _ = c.DB.Where("email=?", email).Find(&user)
	if user.ID == 0 {
		return user, errors.New("invalid email")
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

func (c *userDatabase) SignUpUser(ctx context.Context, user utils.BodySignUpuser) (string, error) {
	var userid uint
	tx := c.DB.Begin()
	query1 := `insert into users(created_at,updated_at,first_name,last_name,email,mobile_num,password,referal_code)values($1,$2,$3,$4,$5,$6,$7,$8) returning id`
	if err := tx.Raw(query1, time.Now(), time.Now(), user.FirstName, user.LastName, user.Email, user.MobileNum, user.Password, user.ReferalCode).Scan(&userid).Error; err != nil {
		tx.Rollback()
		return user.MobileNum, err
	}
	query := `insert into carts(user_id)values($1)`
	if err := tx.Exec(query, userid).Error; err != nil {
		tx.Rollback()
		return user.MobileNum, err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return user.MobileNum, err
	}
	return user.MobileNum, nil
}

func (c *userDatabase) ShowDetails(ctx context.Context, id uint) (utils.ResponseUsers, error) {
	var user utils.ResponseUsers
	query := `SELECT id,first_name,last_name,email,mobile_num,referal_code from users where id=?`
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

func (c *userDatabase) UpdateVerify(number, referalcode string) error {
	var userid1, userid2 uint
	tx := c.DB.Begin()
	if err := tx.Model(&domain.User{}).Where("mobile_num=?", number).UpdateColumn("verified", true).Error; err != nil {
		tx.Rollback()
		return err
	}
	if referalcode != "" {
		if err := tx.Model(&domain.User{}).Where("referal_code=? AND block=?", referalcode, false).Select("id").Scan(&userid2); err.Error != nil {
			tx.Rollback()
			return err.Error
		} else if err.RowsAffected != 0 {
			if err1 := tx.Model(&domain.User{}).Where("mobile_num=?", number).Select("id").Scan(&userid1).Error; err1 != nil {
				tx.Rollback()
				return err1
			}
			current := time.Now()
			wallet1 := domain.Wallet{
				UserID:       userid1,
				CreditedDate: &current,
				Amount:       10,
			}
			if err1 := tx.Create(&wallet1).Error; err1 != nil {
				tx.Rollback()
				return err1
			}
			wallet2 := domain.Wallet{
				UserID:       userid2,
				CreditedDate: &current,
				Amount:       50,
			}
			if err1 := tx.Create(&wallet2).Error; err1 != nil {
				tx.Rollback()
				return err1
			}
		}
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
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

func (c *userDatabase) ListAllCategories(ctx context.Context) ([]utils.ResponseCategory, error) {
	var categories []utils.ResponseCategory
	query := `select category_name from categories where deleted_at is null`
	result := c.DB.Raw(query).Scan(&categories).Error
	if result != nil {
		return categories, errors.New("failed to get categories")
	}
	return categories, nil
}

func (c *userDatabase) ViewWallet(userid uint) (wallet []utils.ResWallet, balance int, err error) {
	var result *int
	query := `SELECT id,credited_date,debited_date,amount from wallets where user_id=?`
	if err = c.DB.Raw(query, userid).Scan(&wallet).Error; err != nil {
		return
	}
	if err = c.DB.Model(&domain.Wallet{}).Select("sum(amount) as balance").Where("user_id=?", userid).Scan(&result).Error; err != nil {
		return
	} else {
		if result != nil {
			balance = *result
		} else {
			balance = 0
		}
	}
	return
}
