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

func (c *orderUseCase) AddtoOrders(addressid, paymentid, userid uint) (utils.RazorpayOrder, error) {
	var body utils.RazorpayOrder
	cart, err := c.cartRepo.FindCartById(userid)
	if err != nil {
		return body, err
	}
	cartitems, err1 := c.orderrepo.Findcartitems(cart.ID)
	if err1 != nil {
		return body, err1
	}

	order := domain.Order{
		UserID:     cart.UserID,
		PlacedDate: time.Now(),
		AddressID:  addressid,
		PaymentID:  paymentid,
		GrandTotal: uint(cart.GrandTotal),
	}
	if paymentid == 1 {
		if err := c.orderrepo.AddtoOrders(cartitems, order); err != nil {
			return body, err
		}
	} else if paymentid == 2 {
		body, err = c.Razorpayment(order.GrandTotal, userid)
		if err != nil {
			return body, err
		}
	}
	return body, nil
}

func (c *orderUseCase) Razorpayment(Amount, userid uint) (razorpayOrder utils.RazorpayOrder, err error) {
	// generate razorpay order
	//razorpay amount is caluculate on pisa for india so make the actual price into paisa
	razorPayAmount := Amount * 100
	razopayOrderId, err := support.GenerateRazorpayOrder(razorPayAmount, "test reciept")
	if err != nil {
		return razorpayOrder, err
	}
	// set all details on razopay order
	razorpayOrder.AmountToPay = Amount
	razorpayOrder.RazorpayAmount = razorPayAmount

	razorpayOrder.RazorpayKey = config.GetCofig().RAZORPAYKEY

	razorpayOrder.RazorpayOrderID = razopayOrderId
	razorpayOrder.UserID = userid

	return razorpayOrder, nil
}

func (c *orderUseCase) CancelOrder(id uint) error {
	orderitem, date, err := c.orderrepo.FindOrderitem(id)
	if err != nil {
		return err
	}
	if orderitem.DeliveredDate != nil {
		return errors.New("already delivered,if not delivered contact customer support")
	}
	if time.Now().After(date.Add(24 * time.Hour)) {
		return errors.New("cancellation time exceeds")
	}
	current := time.Now()
	orderitem.OrderStatusID = 2
	orderitem.CancelledDate = &current
	if err := c.orderrepo.CancelOrder(orderitem); err != nil {
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
	current := time.Now()
	orderitem.OrderStatusID = statusid
	orderitem.DeliveredDate = &current
	if err := c.orderrepo.UpdateStatus(orderitem); err != nil {
		return err
	}
	return nil
}
