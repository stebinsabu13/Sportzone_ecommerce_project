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

// LIST CART_DETAILS
//	@Summary		API FOR DISPLAYING CART TO USER
//	@ID				USER-LIST-CART
//	@Description	LISTING CART AND ITEMS FROM USERS END
//	@Tags			USER
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	utils.Response
//	@Failure		401	{object}	utils.Response
//	@Failure		400	{object}	utils.Response
//	@Failure		500	{object}	utils.Response
//	@Router			/user/cart [get]
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

// ADD TO CART
//	@Summary		API FOR ADDING PRODUCTS TO CART BY USER
//	@ID				USER-ADD-TO-CART
//	@Description	ADDING ITEMS TO CART FROM USERS END
//	@Tags			USER
//	@Accept			json
//	@Produce		json
//	@Param			id	query		string	true	"Enter the product id"
//	@Success		200	{object}	utils.Response
//	@Failure		401	{object}	utils.Response
//	@Failure		400	{object}	utils.Response
//	@Failure		500	{object}	utils.Response
//	@Router			/user/cart/add [put]
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

// REMOVE ITEMS FROM CART
//	@Summary		API FOR REMOVING PRODUCTS TO CART BY USER
//	@ID				USER-REMOVE-FROM-CART
//	@Description	REMOVING ITEMS FROM CART FROM USERS END
//	@Tags			USER
//	@Accept			json
//	@Produce		json
//	@Param			id	query		string	true	"Enter the product id"
//	@Success		200	{object}	utils.Response
//	@Failure		401	{object}	utils.Response
//	@Failure		400	{object}	utils.Response
//	@Failure		500	{object}	utils.Response
//	@Router			/user/cart/remove [put]
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
