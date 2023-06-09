package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
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

// LIST PRODUCTS
//
//	@Summary		API FOR LISTING ALL PRODUCTS
//	@Description	LISTING ALL PRODUCTS FROM ADMINS AND USERS END
//	@Tags			ADMIN USER
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int	false	"Enter the page number to display"
//	@Param			limit	query		int	false	"Number of items to retrieve per page"
//	@Success		200		{object}	utils.Response
//	@Failure		401		{object}	utils.Response
//	@Failure		400		{object}	utils.Response
//	@Failure		500		{object}	utils.Response
//	@Router			/user/products [get]
//	@Router			/admin/product [get]
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

// ADD PRODUCT
//
//	@Summary		API FOR ADDING PRODUCT
//	@ID				ADMIN-ADD-PRODUCT
//	@Description	ADDING PRODUCT FROM ADMINS END
//	@Tags			ADMIN
//	@Accept			json
//	@Produce		json
//	@Param			product_details	body		utils.AddProduct	false	"Enter the product details"
//	@Success		200				{object}	utils.Response
//	@Failure		401				{object}	utils.Response
//	@Failure		400				{object}	utils.Response
//	@Failure		500				{object}	utils.Response
//	@Router			/admin/product/add [post]
func (cr *ProductHandler) AddProduct(c *gin.Context) {
	var body utils.AddProduct
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	var product domain.Product
	copier.Copy(&product, &body)
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

// EDIT PRODUCT
//	@Summary		API FOR EDITING PRODUCT
//	@ID				ADMIN-EDIT-PRODUCT
//	@Description	UPDATING PRODUCT DETAILS FROM ADMINS END
//	@Tags			ADMIN
//	@Accept			json
//	@Produce		json
//	@Param			productid		path		string				true	"Enter the product id that you would like to make the change"
//	@Param			product_details	body		utils.AddProduct	true	"Enter the category details"
//	@Success		200				{object}	utils.Response
//	@Failure		401				{object}	utils.Response
//	@Failure		400				{object}	utils.Response
//	@Failure		500				{object}	utils.Response
//	@Router			/admin/product/update/{productid} [patch]
func (cr ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("productid")
	var body utils.AddProduct
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	var product domain.Product
	copier.Copy(&product, &body)
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

// DELETE PRODUCT
//	@Summary		API FOR DELETING A PRODUCT
//	@ID				ADMIN-DELETE-PRODUCT
//	@Description	DELETING PRODUCT BASED ON PRODUCT ID
//	@Tags			ADMIN
//	@Accept			json
//	@Produce		json
//	@Param			productid	path		string	true	"Enter the product id that you would like to delete"
//	@Success		200			{object}	utils.Response
//	@Failure		401			{object}	utils.Response
//	@Failure		400			{object}	utils.Response
//	@Failure		500			{object}	utils.Response
//	@Router			/admin/product/delete/{productid} [post]
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

// LIST PRODUCTS DETAILS
//
//	@Summary		API FOR LISTING PRODUCTS DETAILS BY ID
//	@Description	LISTING ALL PRODUCTS DETAILS FROM ADMINS AND USERS END
//	@Tags			ADMIN USER
//	@Accept			json
//	@Produce		json
//	@Param			productid	path		string	true	"Enter the product id that you would like to see the details of"
//	@Param			page		query		int		false	"Enter the page number to display"
//	@Param			limit		query		int		false	"Number of items to retrieve"
//	@Success		200			{object}	utils.Response
//	@Failure		401			{object}	utils.Response
//	@Failure		400			{object}	utils.Response
//	@Failure		500			{object}	utils.Response
//	@Router			/user/products/{productid} [get]
//	@Router			/admin/product/detail/{productid} [get]
func (cr *ProductHandler) FindDetailsProductById(c *gin.Context) {
	id := c.Param("productid")
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, err1 := strconv.Atoi(c.DefaultQuery("limit", "5"))
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

// ADD PRODUCT DETAILS
//
//	@Summary		API FOR ADDING PRODUCT DETAILS
//	@ID				ADMIN-ADD-PRODUCT-DETAILS
//	@Description	ADDING PRODUCT DETAILS FROM ADMINS END
//	@Tags			ADMIN
//	@Accept			json
//	@Produce		json
//	@Param			product_details	body		utils.AddProductDetail	false	"Enter the product details"
//	@Success		200				{object}	utils.Response
//	@Failure		401				{object}	utils.Response
//	@Failure		400				{object}	utils.Response
//	@Failure		500				{object}	utils.Response
//	@Router			/admin/product/detail/add [post]
func (cr *ProductHandler) AddProductDetail(c *gin.Context) {
	var body utils.AddProductDetail
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var productdetail domain.ProductDetails
	copier.Copy(&productdetail, &body)
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

// EDIT PRODUCT DETAILS BY ID
//	@Summary		API FOR EDITING PRODUCT DETAILS
//	@ID				ADMIN-EDIT-PRODUCT-DETAILS
//	@Description	UPDATING PRODUCT DETAILS FROM ADMINS END
//	@Tags			ADMIN
//	@Accept			json
//	@Produce		json
//	@Param			productdetailid	path		string					true	"Enter the product details id that you would like to make the change"
//	@Param			product_details	body		utils.AddProductDetail	true	"Enter the product details"
//	@Success		200				{object}	utils.Response
//	@Failure		401				{object}	utils.Response
//	@Failure		400				{object}	utils.Response
//	@Failure		500				{object}	utils.Response
//	@Router			/admin/product/detail/update/{productdetailid} [patch]
func (cr ProductHandler) UpdateProductDetail(c *gin.Context) {
	id := c.Param("productdetailid")
	var body utils.AddProductDetail
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	var productdetail domain.ProductDetails
	copier.Copy(&productdetail, &body)
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

//	@Summary		API FOR DELETE PRODUCT DETAILS
//	@ID				ADMIN-DELETE-PRODUCT-DETAILS
//	@Description	DELETING PRODUCT DETAILS FROM ADMINS END
//	@Tags			ADMIN
//	@Accept			json
//	@Produce		json
//	@Param			productdetailid	path		string	true	"Enter the product details id that you would like to delete"
//	@Success		200				{object}	utils.Response
//	@Failure		401				{object}	utils.Response
//	@Failure		400				{object}	utils.Response
//	@Failure		500				{object}	utils.Response
//	@Router			/admin/product/detail/delete/{productdetailid} [delete]
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

// LIST PRODUCTS BASED ON CATEGORY
//	@Summary		API FOR LISTING ALL PRODUCTS BASED ON CATEGORY
//	@Description	LISTING ALL PRODUCTS FROM ADMINS AND USERS END BASED ON CATEGORY
//	@Tags			ADMIN USER
//	@Accept			json
//	@Produce		json
//	@Param			categoryid	path		string	true	"Enter the category id"
//	@Param			page		query		int		false	"Enter the page number to display"
//	@Param			limit		query		int		false	"Number of items to retrieve per page"
//	@Success		200			{object}	utils.Response
//	@Failure		401			{object}	utils.Response
//	@Failure		400			{object}	utils.Response
//	@Failure		500			{object}	utils.Response
//	@Router			/user/filter/category/{categoryid}/products [get]
//	@Router			/admin/product/bycategory/{categoryid} [get]
func (cr *ProductHandler) ProductsByCategory(c *gin.Context) {
	id := c.Param("categoryid")
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, err1 := strconv.Atoi(c.DefaultQuery("limit", "5"))
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
	products, err := cr.productUseCase.ProductsByCategory(id, pagination)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"products": products,
	})
}

//	@Summary		API FOR LISTING ALL BRANDS
//	@Description	LISTING ALL BRANDS USERS END
//	@Tags			USER
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	utils.Response
//	@Failure		401	{object}	utils.Response
//	@Failure		400	{object}	utils.Response
//	@Failure		500	{object}	utils.Response
//	@Router			/user/filter/brands [get]
func (cr *ProductHandler) ListAllBrands(c *gin.Context) {
	brands, err := cr.productUseCase.ListAllBrands()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"brands": brands,
	})
}

// LIST PRODUCTS BASED ON BRANDS
//	@Summary		API FOR LISTING ALL PRODUCTS BASED ON CATEGBRANDSORY
//	@Description	LISTING ALL PRODUCTS FROM ADMINS AND USERS END BASED ON BRANDS
//	@Tags			ADMIN USER
//	@Accept			json
//	@Produce		json
//	@Param			brandid	path		string	true	"Enter the brand id"
//	@Param			page	query		int		false	"Enter the page number to display"
//	@Param			limit	query		int		false	"Number of items to retrieve per page"
//	@Success		200		{object}	utils.Response
//	@Failure		401		{object}	utils.Response
//	@Failure		400		{object}	utils.Response
//	@Failure		500		{object}	utils.Response
//	@Router			/user/filter/brands/{brandid}/products [get]
//	@Router			/admin/product/bybrands/{brandid} [get]
func (cr *ProductHandler) ProductsByBrands(c *gin.Context) {
	id := c.Param("brandid")
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, err1 := strconv.Atoi(c.DefaultQuery("limit", "5"))
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
	products, err := cr.productUseCase.ProductsByBrands(id, pagination)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"products": products,
	})
}
