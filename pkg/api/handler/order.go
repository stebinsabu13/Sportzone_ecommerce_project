package handler

import (
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
