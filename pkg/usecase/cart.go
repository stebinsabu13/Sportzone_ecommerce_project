package usecase

import (
	"errors"

	"github.com/stebinsabu13/ecommerce-api/pkg/domain"
	interfaces "github.com/stebinsabu13/ecommerce-api/pkg/repository/interface"
	services "github.com/stebinsabu13/ecommerce-api/pkg/usecase/interface"
	"github.com/stebinsabu13/ecommerce-api/pkg/utils"
)

type cartUseCase struct {
	cartRepo interfaces.CartRepository
}

func NewCartUseCase(repo interfaces.CartRepository) services.CartUseCase {
	return &cartUseCase{
		cartRepo: repo,
	}
}
func (c *cartUseCase) ViewCart(userid uint) ([]utils.ResViewCart, int, error) {
	cartdetail, err1 := c.cartRepo.ViewCart(userid)
	cart, err2 := c.cartRepo.FindCartById(userid)
	err := errors.Join(err1, err2)
	if err != nil {
		return cartdetail, cart.GrandTotal, err
	}
	return cartdetail, cart.GrandTotal, nil
}

func (c *cartUseCase) AdditemtoCart(id string, userid uint) error {
	cart, err := c.cartRepo.FindCartById(userid)
	if err != nil {
		return err
	}
	productdetail, err1 := c.cartRepo.FindProductDetailById(id)
	if err1 != nil {
		return err1
	}
	if productdetail.Stock <= 0 {
		return errors.New("out of stock")
	} else {
		exsistitem, err2 := c.cartRepo.FindProductExsist(id, cart.ID)
		if err2 != nil {
			return err2
		}
		if exsistitem.ID != 0 {
			exsistitem.Quantity += 1
			exsistitem.Total = exsistitem.Quantity * productdetail.Price
			if err := c.cartRepo.UpdateCartitem(exsistitem); err != nil {
				return err
			}
		} else {
			item := domain.CartItem{
				CartID:          cart.ID,
				ProductDetailID: productdetail.ID,
				Quantity:        1,
				Total:           productdetail.Price,
			}
			if err := c.cartRepo.AddNewitem(item); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *cartUseCase) RemoveitemFromCart(id string, userid uint) error {
	cart, err := c.cartRepo.FindCartById(userid)
	if err != nil {
		return err
	}
	exsistitem, err1 := c.cartRepo.FindProductExsist(id, cart.ID)
	if err1 != nil {
		return err1
	}
	productdetail, err2 := c.cartRepo.FindProductDetailById(id)
	if err2 != nil {
		return err2
	}
	if exsistitem.Quantity > 1 {
		exsistitem.Quantity -= 1
		exsistitem.Total = exsistitem.Quantity * productdetail.Price
		if err := c.cartRepo.UpdateCartitem(exsistitem); err != nil {
			return err
		}
	} else {
		if err := c.cartRepo.DeleteCartitem(exsistitem); err != nil {
			return err
		}
	}
	return nil
}
