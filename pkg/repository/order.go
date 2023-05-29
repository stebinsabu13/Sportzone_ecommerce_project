package repository

import (
	"context"
	"errors"
	"time"

	"github.com/stebinsabu13/ecommerce-api/pkg/domain"
	interfaces "github.com/stebinsabu13/ecommerce-api/pkg/repository/interface"
	"github.com/stebinsabu13/ecommerce-api/pkg/utils"
	"gorm.io/gorm"
)

type OrderDatabase struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) interfaces.OrderRepository {
	return &OrderDatabase{
		DB: db,
	}
}

func (c *OrderDatabase) Orders(ctx context.Context, id uint) ([]utils.ResOrders, error) {
	var orders []utils.ResOrders
	query := `SELECT o.id,o.placed_date,o.grand_total,ad.house_name,ad.street,ad.city,ad.state,ad.country,ad.pincode,pm.mode from orders o
	inner join addresses ad on ad.id=o.address_id
	inner join payment_modes pm on pm.id=o.payment_id where o.user_id=?`
	err := c.DB.Raw(query, id).Scan(&orders).Error
	if err != nil {
		return orders, err
	}
	return orders, nil
}

func (c *OrderDatabase) OrderDetail(id uint) ([]utils.ResponseOrderDetails, error) {
	var orderdetail []utils.ResponseOrderDetails
	query := `select od.id,p.image,p.model_name,pd.price,b.brand_name,osize.size,ocolour.colour,od.quantity,os.status,od.delivered_date,od.cancelled_date,d.percentage from products p
	inner join product_details pd on p.id=pd.product_id
	left join brands b on b.id=p.brand_id
	inner join order_details od on od.product_detail_id=pd.id
	inner join available_sizes osize on pd.available_size_id=osize.id
	inner join available_colours ocolour on pd.available_colour_id=ocolour.id
	inner join order_statuses os on od.order_status_id=os.id
	left join discounts d on pd.discount_id=d.id where od.order_id=?`
	if err := c.DB.Raw(query, id).Scan(&orderdetail).Error; err != nil {
		return orderdetail, err
	}
	return orderdetail, nil
}
func (c *OrderDatabase) AddtoOrders(items []utils.ResCartItems, order domain.Order) error {
	var stock uint
	tx := c.DB.Begin()
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return err
	}
	for _, v := range items {
		orderitem := domain.OrderDetails{
			OrderID:         order.ID,
			OrderStatusID:   3,
			DeliveredDate:   nil,
			CancelledDate:   nil,
			ProductDetailID: v.ProductDetailID,
			Quantity:        v.Quantity,
		}
		if err := tx.Create(&orderitem).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Model(&domain.ProductDetails{}).Where("id=?", v.ProductDetailID).Select("stock").Scan(&stock).Error; err != nil {
			tx.Rollback()
			return err
		}
		if int(stock-v.Quantity) < 0 {
			tx.Rollback()
			return errors.New("can't place orders out of stock product in the cart please remove and come again")
		}
		newstock := stock - v.Quantity
		if err := tx.Model(&domain.ProductDetails{}).Where("id=?", v.ProductDetailID).UpdateColumn("stock", newstock).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	query := `delete from cart_items where cart_id=$1`
	if err := tx.Exec(query, items[0].CartID).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&domain.Cart{}).Where("id=?", items[0].CartID).UpdateColumn("grand_total", 0).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (c *OrderDatabase) Findcartitems(id uint) ([]utils.ResCartItems, error) {
	var cartitems []utils.ResCartItems
	if err := c.DB.Model(&domain.CartItem{}).Where("cart_id=?", id).Select("cart_id,product_detail_id,quantity").Scan(&cartitems).Error; err != nil {
		return cartitems, err
	}
	return cartitems, nil
}

func (c *OrderDatabase) FindOrderitem(id uint) (domain.OrderDetails, time.Time, error) {
	var item domain.OrderDetails
	var date time.Time
	if err := c.DB.Where("id=?", id).Find(&item).Error; err != nil {
		return item, date, err
	}
	if err := c.DB.Model(&domain.Order{}).Select("placed_date").Where("id=?", item.OrderID).Scan(&date).Error; err != nil {
		return item, date, err
	}
	return item, date, nil
}

func (c *OrderDatabase) CancelOrder(item domain.OrderDetails) error {
	prodetail := struct {
		Price      uint
		Stock      uint
		percentage int
	}{
		Price:      0,
		Stock:      0,
		percentage: 0,
	}
	var grandtotal int
	tx := c.DB.Begin()
	if err := tx.Model(&domain.OrderDetails{}).Where("id=?", item.ID).UpdateColumns(&item).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&domain.Discount{}).Joins("join product_details on discounts.id=product_details.discount_id").Where("product_details.id=?", item.ProductDetailID).Select("discounts.percentage").Scan(&prodetail.percentage).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&domain.ProductDetails{}).Where("id=?", item.ProductDetailID).Select("price,stock").Scan(&prodetail).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&domain.ProductDetails{}).Where("id=?", item.ProductDetailID).UpdateColumn("stock", (prodetail.Stock + item.Quantity)).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&domain.Order{}).Where("id=?", item.OrderID).Select("grand_total").Scan(&grandtotal).Error; err != nil {
		tx.Rollback()
		return err
	}
	discount := (prodetail.percentage * int(prodetail.Price)) / 100
	total := grandtotal - (int(item.Quantity) * (int(prodetail.Price) - discount))
	if err := tx.Model(&domain.Order{}).Where("id=?", item.OrderID).UpdateColumn("grand_total", total).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (c *OrderDatabase) ReturnOrder(item domain.OrderDetails) error {
	if err := c.DB.Model(&domain.OrderDetails{}).Where("id=?", item.ID).UpdateColumns(&item).Error; err != nil {
		return err
	}
	return nil
}

//Admin repository

func (c *OrderDatabase) ListAllOrders() ([]utils.ResAllOrders, error) {
	var allorders []utils.ResAllOrders
	query := `select o.id,u.first_name,u.mobile_num,o.placed_date,ad.house_name,ad.street,ad.pincode,pm.mode,o.grand_total from users u
	inner join orders o on u.id=o.user_id
	inner join addresses ad on u.id=ad.user_id
	inner join payment_modes pm on o.payment_id=pm.id`
	if err := c.DB.Raw(query).Scan(&allorders).Error; err != nil {
		return allorders, err
	}
	return allorders, nil
}

func (c *OrderDatabase) UpdateStatus(item domain.OrderDetails) error {
	var userid uint
	prodetail := struct {
		Price      uint
		Stock      uint
		percentage int
	}{
		Price:      0,
		Stock:      0,
		percentage: 0,
	}
	tx := c.DB.Begin()
	if err := tx.Model(&domain.OrderDetails{}).Where("id=?", item.ID).UpdateColumns(&item).Error; err != nil {
		tx.Rollback()
		return err
	}
	if item.OrderStatusID == 5 {
		if err := tx.Model(&domain.Order{}).Where("id=?", item.OrderID).Select("user_id").Scan(&userid).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Model(&domain.ProductDetails{}).Where("id=?", item.ProductDetailID).Select("price,stock").Scan(&prodetail).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Model(&domain.Discount{}).Joins("join product_details on discounts.id=product_details.discount_id").Where("product_details.id=?", item.ProductDetailID).Select("discounts.percentage").Scan(&prodetail.percentage).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Model(&domain.ProductDetails{}).Where("id=?", item.ProductDetailID).UpdateColumn("stock", (prodetail.Stock + item.Quantity)).Error; err != nil {
			tx.Rollback()
			return err
		}
		discount := (prodetail.percentage * int(prodetail.Price)) / 100
		current := time.Now()
		wallet := domain.Wallet{
			UserID:       userid,
			CreditedDate: &current,
			Amount:       int(item.Quantity) * (int(prodetail.Price) - discount),
		}
		if err := tx.Create(&wallet).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
