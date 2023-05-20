package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	services "github.com/stebinsabu13/ecommerce-api/pkg/usecase/interface"
)

type CartHandler struct {
	cartUseCase services.CartUseCase
}

func NewCartHandler(usecase services.CartUseCase) *CartHandler {
	return &CartHandler{
		cartUseCase: usecase,
	}
}

func (cr *CartHandler) ViewCart(c *gin.Context) {
	id, ok := c.Get("user-id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Not ok",
		})
		return
	}
	cartdetail, GrandTotal, err := cr.cartUseCase.ViewCart(id.(uint))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Success":    cartdetail,
		"GrandTotal": GrandTotal,
	})
}

func (cr *CartHandler) AddtoCart(c *gin.Context) {
	prodetid := c.Query("id")
	id, ok := c.Get("user-id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Not ok",
		})
		return
	}
	err := cr.cartUseCase.AdditemtoCart(prodetid, id.(uint))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Success": "Product added",
	})
}

func (cr *CartHandler) RemovefromCart(c *gin.Context) {
	prodetid := c.Query("id")
	id, ok := c.Get("user-id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Not ok",
		})
		return
	}
	err := cr.cartUseCase.RemoveitemFromCart(prodetid, id.(uint))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Success": "Product removed",
	})
}
