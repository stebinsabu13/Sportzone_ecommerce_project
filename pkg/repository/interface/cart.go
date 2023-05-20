package interfaces

import (
	"github.com/stebinsabu13/ecommerce-api/pkg/domain"
	"github.com/stebinsabu13/ecommerce-api/pkg/utils"
)

type CartRepository interface {
	ViewCart(userid uint) ([]utils.ResViewCart, error)
	FindCartById(userid uint) (domain.Cart, error)
	FindProductDetailById(id string) (domain.ProductDetails, error)
	FindProductExsist(id string, cartid uint) (domain.CartItem, error)
	UpdateCartitem(exsistitem domain.CartItem) error
	AddNewitem(item domain.CartItem) error
	DeleteCartitem(item domain.CartItem) error
}
