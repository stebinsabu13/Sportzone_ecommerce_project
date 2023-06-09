package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stebinsabu13/ecommerce-api/pkg/support"
	services "github.com/stebinsabu13/ecommerce-api/pkg/usecase/interface"
)

type OrderHandler struct {
	orderUseCase services.OrderUseCase
}

func NewOrderHandler(usecase services.OrderUseCase) *OrderHandler {
	return &OrderHandler{
		orderUseCase: usecase,
	}
}

// PLACE A NEW ORDER
//	@Summary		API FOR PLACING A NEW ORDER
//	@ID				USER-PROCEED-ORDER
//	@Description	Users can place a new order with the cart items.
//	@Tags			USER
//	@Accept			json
//	@Produce		json
//	@Param			paymentid	query		string	true	"Enter the payment id"
//	@Param			addressid	query		string	true	"Enter the address id"
//	@Param			code		query		string	false	"If you have a coupon,Enter the coupon code"
//	@Success		200			{object}	utils.Response
//	@Failure		401			{object}	utils.Response
//	@Failure		400			{object}	utils.Response
//	@Failure		500			{object}	utils.Response
//	@Router			/user/checkout/add [get]
func (cr *OrderHandler) AddtoOrders(c *gin.Context) {
	code := c.DefaultQuery("code", "")
	addressid, _ := strconv.Atoi(c.Query("addressid"))
	paymentid, _ := strconv.Atoi(c.Query("paymentid"))
	id, ok := c.Get("user-id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Not ok",
		})
		return
	}
	couponid, err := cr.orderUseCase.ValidateCoupon(id.(uint), code)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if paymentid == 2 {
		body, err := cr.orderUseCase.Razorpayment(id.(uint), couponid)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.HTML(200, "app.html", gin.H{
			"UserID":      body.UserID,
			"Orderid":     body.RazorpayOrderID,
			"Total_price": body.AmountToPay,
		})
	} else {
		if err := cr.orderUseCase.AddtoOrders(uint(addressid), uint(paymentid), id.(uint), couponid); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"Success": "Order placed",
		})
	}
}

// VIEW ORDERS
//	@Summary		API FOR VIEWING ORDERS
//	@Description	Users can view all orders.
//	@Tags			USER
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	utils.Response
//	@Failure		401	{object}	utils.Response
//	@Failure		400	{object}	utils.Response
//	@Failure		500	{object}	utils.Response
//	@Router			/user/orders [get]
func (cr *OrderHandler) ShowOrders(c *gin.Context) {
	id, ok := c.Get("user-id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Not ok",
		})
		return
	}
	orders, err := cr.orderUseCase.Orders(c.Request.Context(), id.(uint))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ORDERS": orders,
	})
}

// VIEW ORDERS DETAILS
//	@Summary		API FOR VIEWING ORDERS DETAILS
//	@Description	Users can the selected order details.
//	@Tags			ADMIN USER
//	@Accept			json
//	@Produce		json
//	@Param			orderid	query		uint	true	"Enter the order id"
//	@Success		200		{object}	utils.Response
//	@Failure		401		{object}	utils.Response
//	@Failure		400		{object}	utils.Response
//	@Failure		500		{object}	utils.Response
//	@Router			/user/orders/detail [get]
//	@Router			/admin/order/detail [get]
func (cr *OrderHandler) ShowOrderDetail(c *gin.Context) {
	id, err1 := strconv.Atoi(c.Query("orderid"))
	if err1 != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err1.Error(),
		})
		return
	}
	orderDetail, err := cr.orderUseCase.OrderDetail(uint(id))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ORDER DETAILS": orderDetail,
	})
}

// CANCEL ORDER
//	@Summary		API FOR CANCELLING A ORDER
//	@Description	Users can cancel orders
//	@Tags			USER
//	@Accept			json
//	@Produce		json
//	@Param			orderdetailid	query		uint	true	"Enter the order details id"
//	@Param			statusid		query		uint	true	"Enter the status id"
//	@Success		200				{object}	utils.Response
//	@Failure		401				{object}	utils.Response
//	@Failure		400				{object}	utils.Response
//	@Failure		500				{object}	utils.Response
//	@Router			/user/orders/cancel [patch]
func (cr *OrderHandler) CancelOrder(c *gin.Context) {
	userid, ok := c.Get("user-id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized access",
		})
		return
	}
	id, err1 := strconv.Atoi(c.Query("orderdetailid"))
	statusid, err2 := strconv.Atoi(c.Query("statusid"))
	err := errors.Join(err1, err2)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := cr.orderUseCase.CancelOrder(userid.(uint), uint(id), uint(statusid)); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Success": "Order cancelled",
	})
}

func (cr *OrderHandler) RazorpaymentSuccess(c *gin.Context) {
	orderID := c.Query("order_id")
	userID, err1 := strconv.Atoi(c.Query("user_id"))
	addressid, err2 := strconv.Atoi(c.Query("addressid"))
	paymentid, err3 := strconv.Atoi(c.Query("paymentid"))
	code := c.Query("code")
	payment_refID := c.Query("payment_ref")
	signatureID := c.Query("signature")
	response := gin.H{
		"data":    false,
		"message": "Payment failed",
	}
	err := errors.Join(err3, err2, err1)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	if err := support.VeifyRazorpayPayment(orderID, payment_refID, signatureID); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	couponid, er := cr.orderUseCase.FindCoupon(code)
	if er != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": er.Error(),
		})
		return
	}
	if err := cr.orderUseCase.AddtoOrders(uint(addressid), uint(paymentid), uint(userID), couponid); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}
	response["data"] = true
	response["message"] = "Payment Success."
	c.JSON(http.StatusOK, response)
}

// RETURN ORDER
//	@Summary		API FOR RETURNING A ORDER
//	@Description	Users can return orders
//	@Tags			USER
//	@Accept			json
//	@Produce		json
//	@Param			orderdetailid	query		uint	true	"Enter the order details id"
//	@Param			statusid		query		uint	true	"Enter the status id"
//	@Success		200				{object}	utils.Response
//	@Failure		401				{object}	utils.Response
//	@Failure		400				{object}	utils.Response
//	@Failure		500				{object}	utils.Response
//	@Router			/user/orders/return [patch]
func (cr *OrderHandler) ReturnOrder(c *gin.Context) {
	id, err1 := strconv.Atoi(c.Query("orderdetailid"))
	statusid, err2 := strconv.Atoi(c.Query("statusid"))
	err := errors.Join(err1, err2)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := cr.orderUseCase.ReturnOrder(uint(id), uint(statusid)); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Success": "Order returning request submitted",
	})
}

//Admin Handlers

// VIEW ORDERS
//	@Summary		API FOR VIEWING ALL ORDERS
//	@Description	Admin can view all orders.
//	@Tags			ADMIN
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	utils.Response
//	@Failure		401	{object}	utils.Response
//	@Failure		400	{object}	utils.Response
//	@Failure		500	{object}	utils.Response
//	@Router			/admin/order [get]
func (cr *OrderHandler) ListAllOrders(c *gin.Context) {
	Allorders, err := cr.orderUseCase.ListAllOrders()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ORDERS": Allorders,
	})
}

// CHANGE STATUS OF ORDER
//	@Summary		API FOR CHANGING THE STATUS OF A ORDER
//	@Description	Admin can change the ststus of orders
//	@Tags			ADMIN
//	@Accept			json
//	@Produce		json
//	@Param			orderdetailid	query		uint	true	"Enter the order details id"
//	@Param			statusid		query		uint	true	"Enter the status id"
//	@Success		200				{object}	utils.Response
//	@Failure		401				{object}	utils.Response
//	@Failure		400				{object}	utils.Response
//	@Failure		500				{object}	utils.Response
//	@Router			/admin/order/update/status [post]
func (cr *OrderHandler) UpdateStatus(c *gin.Context) {
	id, err1 := strconv.Atoi(c.Query("orderdetailid"))
	statusid, err2 := strconv.Atoi(c.Query("statusid"))
	err := errors.Join(err1, err2)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := cr.orderUseCase.UpdateStatus(uint(id), uint(statusid)); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Success": "Status updated",
	})
}
