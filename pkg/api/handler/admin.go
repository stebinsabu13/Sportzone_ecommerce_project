package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stebinsabu13/ecommerce-api/pkg/auth"
	"github.com/stebinsabu13/ecommerce-api/pkg/domain"
	"github.com/stebinsabu13/ecommerce-api/pkg/support"
	services "github.com/stebinsabu13/ecommerce-api/pkg/usecase/interface"
	"github.com/stebinsabu13/ecommerce-api/pkg/utils"
)

type AdminHandler struct {
	AdminUseCase services.AdminUseCase
}

func NewAdminHandler(usecase services.AdminUseCase) *AdminHandler {
	return &AdminHandler{
		AdminUseCase: usecase,
	}
}

func (cr *AdminHandler) LoginHandler(c *gin.Context) {
	_, err := c.Cookie("admin-token")
	if err == nil {
		c.Redirect(http.StatusFound, "/admin/home")
		return
	}
	var body utils.BodyLogin
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	admin, err := cr.AdminUseCase.FindbyEmail(c.Request.Context(), body.Email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if ok := support.CheckPasswordHash(body.Password, admin.Password); !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Password",
		})
		return
	}
	tokenString, err := auth.GenerateJWT(admin.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Not able to generate token, login again",
		})
		return
	}
	c.SetCookie("admin-token", tokenString, int(time.Now().Add(60*time.Minute).Unix()), "/", "localhost", false, true)
	c.Set("admin-id", admin.ID)
	c.JSON(http.StatusOK, gin.H{
		"Success": "Admin Login",
	})
}

// func (cr AdminHandler) HomeHandler(c *gin.Context) {
// 	email, ok := c.Get("admin-email")
// 	if !ok {
// 		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
// 			"error": "Unauthorized admin",
// 		})
// 		return
// 	}
// 	admin, err := cr.AdminUseCase.FindbyEmail(c.Request.Context(), email.(string))
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
// 			"error": "Unauthorized admin",
// 		})
// 		return
// 	}
// 	c.JSON(http.StatusOK, admin)
// }

func (cr *AdminHandler) LogoutHandler(c *gin.Context) {
	c.SetCookie("admin-token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{
		"logout": "Success",
	})
}

// func (cr *AdminHandler) SignUp(c *gin.Context) {
// 	var admin domain.Admin
// 	if err := c.BindJSON(&admin); err != nil {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	if ok := support.Email_validater(admin.Email); !ok {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
// 			"error": "Email format incorrect",
// 		})
// 		return
// 	}

// 	if ok := support.MobileNum_validater(admin.MobileNum); !ok {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
// 			"error": "Not a valid mobile number",
// 		})
// 		return
// 	}
// 	if _, err := cr.AdminUseCase.FindbyEmail(c.Request.Context(), admin.Email); err == nil {
// 		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
// 			"error": "User already Exsists",
// 		})
// 		return
// 	}

// 	admin.Password, _ = support.HashPassword(admin.Password)
// 	err := cr.AdminUseCase.SignUpAdmin(c.Request.Context(), admin)
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{
// 		"User registration": "Success",
// 	})
// }

func (cr *AdminHandler) ListAllUsers(c *gin.Context) {
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
	users, err := cr.AdminUseCase.ListAllUsers(c.Request.Context(), pagination)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Users": users,
	})
}

func (cr *AdminHandler) AccessManage(c *gin.Context) {
	id := c.Param("userid")
	str := c.Query("access")
	access, _ := strconv.ParseBool(str)
	err := cr.AdminUseCase.AccessManage(c.Request.Context(), id, access)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Access": "Updated",
	})
}

func (cr *AdminHandler) ListAllCategories(c *gin.Context) {
	categories, err := cr.AdminUseCase.ListAllCategories(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Categories": categories,
	})
}

func (cr *AdminHandler) AddCategory(c *gin.Context) {
	var category domain.Category
	if err := c.BindJSON(&category); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "error while binding json",
		})
		return
	}
	if err := cr.AdminUseCase.AddCategory(c.Request.Context(), category); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Success": "Category added",
	})
}

func (cr *AdminHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("categoryid")
	if err := cr.AdminUseCase.DeleteCategory(c.Request.Context(), id); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Success": "Category deleted",
	})
}

func (cr *AdminHandler) FullSalesReport(c *gin.Context) {
	// time
	monthInt, err1 := strconv.Atoi(c.Query("month"))
	month := time.Month(monthInt)
	year, err2 := strconv.Atoi(c.Query("year"))
	frequency := c.Query("frequency")

	// page
	count, err3 := strconv.Atoi(c.Query("count"))
	pageNumber, err4 := strconv.Atoi(c.Query("page_number"))
	err := errors.Join(err1, err2, err3, err4)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	offset := (pageNumber - 1) * count
	reqData := utils.SalesReport{
		Month:     month,
		Year:      year,
		Frequency: frequency,
		Pagination: utils.Pagination{
			Offset: uint(offset),
			Limit:  uint(count),
		},
	}
	salesreport, err5 := cr.AdminUseCase.GetFullSalesReport(reqData)
	if err5 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err5.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Report": salesreport,
	})
}
