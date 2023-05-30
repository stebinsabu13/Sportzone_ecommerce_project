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
		result := c.DB.Model(&domain.Order{}).Where("EXTRACT(YEAR FROM orders.placed_date) = ? AND EXTRACT(MONTH FROM orders.placed_date) = ?", reqData.Year, reqData.Month).
			Joins("JOIN order_details od on orders.id=od.order_id").
			Joins("JOIN product_details pd on pd.id=od.product_detail_id").
			Joins("JOIN products p on p.id=pd.product_id").
			Joins("JOIN payment_modes pm on pm.id=orders.payment_id").
			Joins("JOIN users u on orders.user_id=u.id").
			Joins("JOIN order_statuses os on os.id=od.order_status_id").
			Joins("JOIN discounts d on d.id=pd.discount_id").
			Select("u.id as userid,u.first_name,u.email,od.product_detail_id as productdetailid,p.model_name as productname,od.quantity,orders.id as orderid,orders.placed_date,pm.mode as paymentmode,pd.price,d.percentage as discountpercentage,os.status as orderstatus").
			Order("orders.placed_date DESC").Scan(&salesreport)
		if result.Error != nil {
			return salesreport, result.Error
		}
	}
	if reqData.Frequency == "YEARLY" {
		result := c.DB.Model(&domain.Order{}).Where("EXTRACT(YEAR FROM orders.placed_date) = ?", reqData.Year).
			Joins("JOIN order_details od on orders.id=od.order_id").
			Joins("JOIN product_details pd on pd.id=od.product_detail_id").
			Joins("JOIN products p on p.id=pd.product_id").
			Joins("JOIN payment_modes pm on pm.id=orders.payment_id").
			Joins("JOIN users u on orders.user_id=u.id").
			Joins("JOIN order_statuses os on os.id=od.order_status_id").
			Joins("JOIN discounts d on d.id=pd.discount_id").
			Select("u.id as userid,u.first_name,u.email,od.product_detail_id as productdetailid,p.model_name as productname,od.quantity,orders.id as orderid,orders.placed_date,pm.mode as paymentmode,pd.price,d.percentage as discountpercentage,os.status as orderstatus").
			Order("orders.placed_date DESC").Scan(&salesreport)
		if result.Error != nil {
			return salesreport, result.Error
		}
	}
	return salesreport, nil
}

func (c *adminDatabase) Widgets() (utils.ResWidgets, error) {
	var widgets utils.ResWidgets
	if err := c.DB.Model(&domain.User{}).Select("count(users)").Where("block='t'").Scan(&widgets.Numberofblockedusers).Error; err != nil {
		return widgets, err
	}
	if err := c.DB.Model(&domain.OrderDetails{}).Select("count(order_details)").Where("delivered_date is null and cancelled_date is null").Scan(&widgets.Numberofpendingorders).Error; err != nil {
		return widgets, err
	}
	if err := c.DB.Model(&domain.ProductDetails{}).Select("count(product_details)").Where("deleted_at is null").Scan(&widgets.Numberofproducts).Error; err != nil {
		return widgets, err
	}
	return widgets, nil
}

func (c *adminDatabase) AddCoupon(coupon domain.Coupon) error {
	if err := c.DB.Create(&coupon).Error; err != nil {
		return err
	}
	return nil
}
