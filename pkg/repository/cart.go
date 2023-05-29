package repository

import (
	"github.com/stebinsabu13/ecommerce-api/pkg/domain"
	interfaces "github.com/stebinsabu13/ecommerce-api/pkg/repository/interface"
	"github.com/stebinsabu13/ecommerce-api/pkg/utils"
	"gorm.io/gorm"
)

type CartDatabase struct {
	DB *gorm.DB
}

func NewCartRepository(db *gorm.DB) interfaces.CartRepository {
	return &CartDatabase{
		DB: db,
	}
}

func (c *CartDatabase) ViewCart(userid uint) ([]utils.ResViewCart, error) {
	var cartdetail []utils.ResViewCart
	query := `select p.model_name,p.image,pd.price,citm.quantity,citm.total,b.brand_name,osize.size,ocolour.colour,d.percentage from products p
	inner join product_details pd on p.id=pd.product_id
	inner join cart_items citm on pd.id=citm.product_detail_id
	left join brands b on p.brand_id=b.id
	inner join carts cs on citm.cart_id=cs.id
	inner join available_sizes osize on osize.id=pd.available_size_id
	inner join available_colours ocolour on ocolour.id=pd.available_colour_id
	left join discounts d on pd.discount_id=d.id where cs.user_id=?`
	err := c.DB.Raw(query, userid).Scan(&cartdetail).Error
	if err != nil {
		return cartdetail, err
	}
	return cartdetail, nil
}

func (c *CartDatabase) FindCartById(userid uint) (domain.Cart, error) {
	var cart domain.Cart
	if err := c.DB.Where("user_id=?", userid).Find(&cart).Error; err != nil {
		return cart, err
	}
	return cart, nil
}

func (c *CartDatabase) FindProductDetailById(id string) (domain.ProductDetails, int, error) {
	var productdetail domain.ProductDetails
	var discount int
	if err := c.DB.Where("id=?", id).Find(&productdetail).Error; err != nil {
		return productdetail, discount, err
	}
	if err := c.DB.Model(&domain.Discount{}).Where("id=?", productdetail.DiscountID).Select("percentage").Scan(&discount).Error; err != nil {
		return productdetail, discount, err
	}
	return productdetail, discount, nil
}

func (c *CartDatabase) FindProductExsist(id string, cartid uint) (domain.CartItem, error) {
	var exsistitem domain.CartItem
	result := c.DB.Where("product_detail_id=$1 and cart_id=$2", id, cartid).Find(&exsistitem)
	if result.Error != nil {
		return exsistitem, result.Error
	}
	return exsistitem, nil
}

func (c *CartDatabase) UpdateCartitem(exsistitem domain.CartItem) error {
	var grandtotal int
	tx := c.DB.Begin()
	if err := tx.Model(&domain.CartItem{}).Where("id=?", exsistitem.ID).UpdateColumns(&exsistitem).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&domain.CartItem{}).Where("cart_id=?", exsistitem.CartID).Select("SUM(total)").Scan(&grandtotal).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&domain.Cart{}).Where("id=?", exsistitem.CartID).UpdateColumn("grand_total", grandtotal).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (c *CartDatabase) AddNewitem(item domain.CartItem) error {
	var grandtotal int
	tx := c.DB.Begin()
	if err := tx.Create(&item).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&domain.CartItem{}).Where("cart_id=?", item.CartID).Select("SUM(total)").Scan(&grandtotal).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&domain.Cart{}).Where("id=?", item.CartID).UpdateColumn("grand_total", grandtotal).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (c *CartDatabase) DeleteCartitem(item domain.CartItem) error {
	var grandtotal int
	tx := c.DB.Begin()
	if err := tx.Model(&domain.CartItem{}).Where("id=?", item.ID).Delete(&item).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&domain.CartItem{}).Where("cart_id=?", item.CartID).Select("SUM(total)").Scan(&grandtotal).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&domain.Cart{}).Where("id=?", item.CartID).UpdateColumn("grand_total", grandtotal).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
