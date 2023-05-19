package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stebinsabu13/ecommerce-api/pkg/domain"
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

func (cr *ProductHandler) AddProduct(c *gin.Context) {
	var product domain.Product
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

func (cr ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("productid")
	var product domain.Product
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

func (cr *ProductHandler) FindDetailsProductById(c *gin.Context) {
	id := c.Param("productid")
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
	productdetails, err := cr.productUseCase.FindProductById(c.Request.Context(), id, pagination)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Details": productdetails,
	})
}

func (cr *ProductHandler) AddProductDetail(c *gin.Context) {
	var productdetail domain.ProductDetails
	if err := c.BindJSON(&productdetail); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := cr.productUseCase.AddProductDetail(c.Request.Context(), productdetail); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Success": "Product detail added",
	})
}

func (cr ProductHandler) UpdateProductDetail(c *gin.Context) {
	id := c.Param("productdetailid")
	var productdetail domain.ProductDetails
	if err := c.BindJSON(&productdetail); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	err := cr.productUseCase.EditProductDetail(c.Request.Context(), productdetail, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Success": "Product detail edited",
	})
}

func (cr *ProductHandler) DeleteProductDetail(c *gin.Context) {
	id := c.Param("productdetailid")
	err := cr.productUseCase.DeleteProductDetail(c.Request.Context(), id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Success": "Product detail deleted",
	})
}
