package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/stebinsabu13/ecommerce-api/pkg/config"
	"github.com/stebinsabu13/ecommerce-api/pkg/domain"
	interfaces "github.com/stebinsabu13/ecommerce-api/pkg/repository/interface"
	"github.com/stebinsabu13/ecommerce-api/pkg/support"
	services "github.com/stebinsabu13/ecommerce-api/pkg/usecase/interface"
	"github.com/stebinsabu13/ecommerce-api/pkg/utils"
)

type orderUseCase struct {
	orderrepo interfaces.OrderRepository
	cartRepo  interfaces.CartRepository
}

func NewOrderUseCase(repo interfaces.OrderRepository, cartrepo interfaces.CartRepository) services.OrderUseCase {
	return &orderUseCase{
		orderrepo: repo,
		cartRepo:  cartrepo,
	}
}

func (c *orderUseCase) Orders(ctx context.Context, id uint) ([]utils.ResOrders, error) {
	orders, err := c.orderrepo.Orders(ctx, id)
	return orders, err
}

func (c *orderUseCase) OrderDetail(id uint) ([]utils.ResponseOrderDetails, error) {
	return c.orderrepo.OrderDetail(id)
}

func (c *orderUseCase) AddtoOrders(addressid, paymentid, userid uint, couponid *uint) error {
	cart, err := c.cartRepo.FindCartById(userid)
	if err != nil {
		return err
	}
	cartitems, err1 := c.orderrepo.Findcartitems(cart.ID)
	if err1 != nil {
		return err1
	}
	order := domain.Order{
		UserID:     cart.UserID,
		PlacedDate: time.Now(),
		AddressID:  addressid,
		PaymentID:  paymentid,
		GrandTotal: uint(cart.GrandTotal),
		CouponID:   couponid,
	}
	if err := c.orderrepo.AddtoOrders(cartitems, order); err != nil {
		return err
	}
	return nil
}

func (c *orderUseCase) Razorpayment(userid uint, couponid *uint) (utils.RazorpayOrder, error) {
	var razorpayOrder utils.RazorpayOrder
	cart, err := c.cartRepo.FindCartById(userid)
	if err != nil {
		return razorpayOrder, err
	}
	// generate razorpay order
	//razorpay amount is caluculate on pisa for india so make the actual price into paisa
	razorPayAmount := cart.GrandTotal * 100
	razopayOrderId, err := support.GenerateRazorpayOrder(razorPayAmount, "test reciept")
	if err != nil {
		return razorpayOrder, err
	}
	// set all details on razopay order
	razorpayOrder.AmountToPay = uint(cart.GrandTotal)
	razorpayOrder.RazorpayAmount = razorPayAmount

	razorpayOrder.RazorpayKey = config.GetCofig().RAZORPAYKEY

	razorpayOrder.RazorpayOrderID = razopayOrderId
	razorpayOrder.UserID = userid

	return razorpayOrder, nil
}

func (c *orderUseCase) CancelOrder(userid, id, statusid uint) error {
	orderitem, date, err := c.orderrepo.FindOrderitem(id)
	if err != nil {
		return err
	}
	if orderitem.DeliveredDate != nil {
		return errors.New("already delivered,if not delivered contact customer support")
	}
	if orderitem.CancelledDate != nil {
		return errors.New("order already cancelled")
	}
	if time.Now().After(date.Add(24 * time.Hour)) {
		return errors.New("cancellation time exceeds")
	}
	current := time.Now()
	orderitem.OrderStatusID = statusid
	orderitem.CancelledDate = &current
	if err := c.orderrepo.CancelOrder(userid, orderitem); err != nil {
		return err
	}
	return nil
}

func (c *orderUseCase) ReturnOrder(id, statusid uint) error {
	orderitem, _, err := c.orderrepo.FindOrderitem(id)
	if err != nil {
		return err
	}
	if orderitem.CancelledDate != nil {
		return errors.New("order already cancelled")
	}
	if orderitem.DeliveredDate == nil {
		return errors.New("order not delivered yet")
	}
	if orderitem.ReturnSubmitDate != nil {
		return errors.New("request already submitted,contact customer support")
	}
	if time.Now().After(orderitem.DeliveredDate.Add(120 * time.Hour)) {
		return errors.New("returning time exceeds")
	}
	current := time.Now()
	orderitem.OrderStatusID = statusid
	orderitem.ReturnSubmitDate = &current
	if err := c.orderrepo.ReturnOrder(orderitem); err != nil {
		return err
	}
	return nil
}

//Admin usecases

func (c *orderUseCase) ListAllOrders() ([]utils.ResAllOrders, error) {
	return c.orderrepo.ListAllOrders()
}

func (c *orderUseCase) UpdateStatus(id, statusid uint) error {
	orderitem, _, err := c.orderrepo.FindOrderitem(id)
	if err != nil {
		return err
	}
	if orderitem.CancelledDate != nil {
		return errors.New("order already cancelled")
	}
	if statusid == 1 {
		if orderitem.DeliveredDate != nil {
			return errors.New("order already delivered")
		}
		current := time.Now()
		orderitem.OrderStatusID = statusid
		orderitem.DeliveredDate = &current
		if err := c.orderrepo.UpdateStatus(orderitem); err != nil {
			return err
		}
	} else if statusid == 5 {
		if orderitem.ReturnSubmitDate == nil {
			return errors.New("not requested for return")
		}
		orderitem.OrderStatusID = statusid
		if err := c.orderrepo.UpdateStatus(orderitem); err != nil {
			return err
		}
	}
	return nil
}

func (c *orderUseCase) ValidateCoupon(userid uint, code string) (*uint, error) {
	if code != "" {
		cart, err := c.cartRepo.FindCartById(userid)
		if err != nil {
			return nil, err
		}
		coupon, err2 := c.orderrepo.FindCoupon(code)
		if err2 != nil {
			return nil, err2
		}
		cartitems, err1 := c.orderrepo.Findcartitems(cart.ID)
		if err1 != nil {
			return nil, err1
		}
		if err = c.orderrepo.ValidateCoupon(coupon, cartitems, &cart); err != nil {
			return nil, err
		}
		return &coupon.ID, nil
	}
	return nil, nil
}

func (c *orderUseCase) FindCoupon(code string) (*uint, error) {
	if code != "" {
		coupon, err := c.orderrepo.FindCoupon(code)
		if err != nil {
			return nil, err
		}
		return &coupon.ID, err
	}
	return nil, nil
}
