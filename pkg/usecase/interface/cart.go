package interfaces

import "github.com/stebinsabu13/ecommerce-api/pkg/utils"

type CartUseCase interface {
	ViewCart(uint) ([]utils.ResViewCart, int, error)
	AdditemtoCart(string, uint) error
	RemoveitemFromCart(string, uint) error
}
