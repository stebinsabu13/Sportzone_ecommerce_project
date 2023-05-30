package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	config "github.com/stebinsabu13/ecommerce-api/pkg/config"
	domain "github.com/stebinsabu13/ecommerce-api/pkg/domain"
)

func ConnectDatabase(cfg config.Config) (*gorm.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)
	db, dbErr := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	db.AutoMigrate(&domain.User{}, &domain.Address{}, &domain.Cart{}, &domain.CartItem{}, &domain.Wallet{},
		&domain.OtpSession{},
		&domain.Admin{},
		&domain.Category{}, &domain.Product{}, &domain.Brand{}, &domain.ProductDetails{}, &domain.Discount{}, &domain.AvailableColour{}, &domain.AvailableSize{},
		&domain.Order{}, &domain.OrderDetails{}, &domain.OrderStatus{},
		&domain.Coupon{}, &domain.CouponType{}, &domain.CouponUsage{},
		&domain.PaymentMode{},
	)

	return db, dbErr
}
