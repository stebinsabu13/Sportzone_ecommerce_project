package repository

import (
	"context"
	"errors"

	domain "github.com/stebinsabu13/ecommerce-api/pkg/domain"
	interfaces "github.com/stebinsabu13/ecommerce-api/pkg/repository/interface"
	"github.com/stebinsabu13/ecommerce-api/pkg/utils"
	"gorm.io/gorm"
)

type adminDatabase struct {
	DB *gorm.DB
}

func NewAdminRepository(Db *gorm.DB) interfaces.AdminRepository {
	return &adminDatabase{DB: Db}
}

// func (c *adminDatabase) FindAll(ctx context.Context) ([]domain.Admin, error) {
// 	var admin []domain.Admin
// 	err := c.DB.Find(&admin).Error
// 	return admin, err
// }

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

func (c *adminDatabase) ListAllUsers(ctx context.Context, pagination utils.Pagination) ([]utils.ResponseUsers, error) {
	var users []utils.ResponseUsers
	limit := pagination.Limit
	offset := pagination.Offset
	query := `SELECT first_name,last_name,email,mobile_num,block from users LIMIT $1 OFFSET $2`
	result := c.DB.Raw(query, limit, offset).Scan(&users).Error
	if result != nil {
		return users, errors.New("failed to get all users")
	}
	return users, nil
}

func (c *adminDatabase) AccessManage(ctx context.Context, id string, access bool) error {
	result := c.DB.Model(&domain.User{}).Where("id=?", id).UpdateColumn("block", access).Error
	if result != nil {
		return errors.New("failed to update")
	}
	return nil
}

func (c *adminDatabase) ListAllCategories(ctx context.Context) ([]utils.ResponseCategory, error) {
	var categories []utils.ResponseCategory
	query := `select category_name from categories where deleted_at is null`
	result := c.DB.Raw(query).Scan(&categories).Error
	if result != nil {
		return categories, errors.New("failed to get categories")
	}
	return categories, nil
}

func (c *adminDatabase) AddCategory(ctx context.Context, category domain.Category) error {
	result := c.DB.Create(&category).Error
	if result != nil {
		return errors.New("failed to add the category")
	}
	return nil
}

func (c *adminDatabase) DeleteCategory(ctx context.Context, id string) error {
	result := c.DB.Where("id=?", id).Delete(&domain.Category{}).Error
	if result != nil {
		return errors.New("failed to delete")
	}
	return nil
}

func (c *adminDatabase) GetFullSalesReport(reqData utils.SalesReport) ([]utils.ResSalesReport, error) {
	var salesreport []utils.ResSalesReport
	if reqData.Frequency == "MONTHLY" {
		result := c.DB.Where("YEAR(created_at) = ? AND MONTH(created_at) = ?", reqData.Year, reqData.Month).Select("")
		if result.Error != nil {
			return salesreport, result.Error
		}
	}
	return salesreport, nil
}
