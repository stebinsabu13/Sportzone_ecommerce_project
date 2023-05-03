package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stebinsabu13/ecommerce-api/pkg/support"
	services "github.com/stebinsabu13/ecommerce-api/pkg/usecase/interface"
)

type ProductHandler struct {
	productUseCase services.ProductUseCase
}

func NewProductHandler(service services.ProductUseCase) *ProductHandler {
	return &ProductHandler{
		productUseCase: service,
	}
}

func (cr *ProductHandler) FindAllProducts(c *gin.Context) {
	products, err := cr.productUseCase.FindAllProducts(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Products_list": products,
	})
}

func (cr *ProductHandler) FindProductById(c *gin.Context) {
	id := c.Param("id")
	descProd, er := cr.productUseCase.FindProductDesc(c.Request.Context(), id)
	if er != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": er.Error(),
		})
		return
	}
	availablecolours, err := cr.productUseCase.FindAvailableColours(c.Request.Context(), id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	availablesizes, err1 := cr.productUseCase.FindAvailableSize(c.Request.Context(), id)
	if err1 != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err1.Error(),
		})
		return
	}
	discount, err2 := cr.productUseCase.FindProductDiscount(c.Request.Context(), id)
	if err2 != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err2.Error(),
		})
		return
	}
	productDetails := support.CalculateTotalPrice(descProd, availablecolours, availablesizes, discount)
	c.JSON(http.StatusOK, gin.H{
		"Product_Details": productDetails,
	})
}
