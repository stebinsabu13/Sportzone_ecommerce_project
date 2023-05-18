package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stebinsabu13/ecommerce-api/pkg/domain"
	"github.com/stebinsabu13/ecommerce-api/pkg/support"
	services "github.com/stebinsabu13/ecommerce-api/pkg/usecase/interface"
	"github.com/stebinsabu13/ecommerce-api/pkg/utils"
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
	page, err := strconv.Atoi(c.Query("page"))
	limit, err1 := strconv.Atoi(c.Query("limit"))
	err = errors.Join(err, err1)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	offset := (page - 1) * limit
	pagination := utils.Pagination{
		Offset: uint(offset),
		Limit:  uint(limit),
	}
	products, err := cr.productUseCase.FindAllProducts(c.Request.Context(), pagination)
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

func (cr *ProductHandler) AddProduct(c *gin.Context) {
	var product domain.ProductDetails
	if err := c.BindJSON(&product); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	err := cr.productUseCase.AddProduct(c.Request.Context(), product)
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

func (cr ProductHandler) EditProduct(c *gin.Context) {
	id := c.Param("productid")
	var product domain.ProductDetails
	if err := c.BindJSON(&product); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	err := cr.productUseCase.EditProduct(c.Request.Context(), product, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Success": "Product edited",
	})
}

func (cr *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("productid")
	err := cr.productUseCase.DeleteProduct(c.Request.Context(), id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Success": "Product deleted",
	})
}
