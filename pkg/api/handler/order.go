package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

func (cr *OrderHandler) AddtoOrders(c *gin.Context) {
	addressid, _ := strconv.Atoi(c.Query("addressid"))
	paymentid, _ := strconv.Atoi(c.Query("paymentid"))
	id, ok := c.Get("user-id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Not ok",
		})
		return
	}
	if err := cr.orderUseCase.AddtoOrders(uint(addressid), uint(paymentid), id.(uint)); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Success": "Order placed",
	})
}

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

func (cr *OrderHandler) CancelOrder(c *gin.Context) {
	id, err1 := strconv.Atoi(c.Query("orderdetailid"))
	if err1 != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err1.Error(),
		})
		return
	}
	if err := cr.orderUseCase.CancelOrder(uint(id)); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Success": "Order cancelled",
	})
}

//Admin Handlers

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
	}
	c.JSON(http.StatusOK, gin.H{
		"Success": "Status updated",
	})
}
