package handler

import (
	"encoding/csv"
	"errors"
	"fmt"
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
	OtpUseCase   services.OtpUseCase
}

func NewAdminHandler(usecase services.AdminUseCase, otpusecase services.OtpUseCase) *AdminHandler {
	return &AdminHandler{
		AdminUseCase: usecase,
		OtpUseCase:   otpusecase,
	}
}

func (cr *AdminHandler) LoginHandler(c *gin.Context) {
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
	c.SetCookie("admin-token", tokenString, int(time.Now().Add(5*time.Minute).Unix()), "/", "sportzone.cloud", true, true)
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

func (cr *AdminHandler) SignUp(c *gin.Context) {
	var signUp_user utils.BodySignUpuser
	if err := c.BindJSON(&signUp_user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	mobile_num, err := cr.AdminUseCase.SignUpAdmin(c.Request.Context(), signUp_user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	respSid, err1 := cr.OtpUseCase.TwilioSendOTP(c.Request.Context(), mobile_num)
	if err1 != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err1.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Success":     "Enter the otp and the responseid",
		"responseid":  respSid,
		"referalcode": signUp_user.ReferalCode,
	})
}

func (cr *AdminHandler) SignupOtpverify(c *gin.Context) {
	var OTP utils.Otpverify
	if err := c.BindJSON(&OTP); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Error binding json",
		})
		return
	}
	session, err := cr.OtpUseCase.TwilioVerifyOTP(c.Request.Context(), OTP)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	err1 := cr.AdminUseCase.UpdateVerify(session.MobileNum, OTP.ReferalCode)
	if err1 != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err1.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"User registration": "Success",
	})
}

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
	monthInt, err1 := strconv.Atoi(c.DefaultQuery("month", "1"))
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
	if salesreport == nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "there is no sales report on this period",
		})
	} else {
		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", "attachment;filename=ecommercesalesreport.csv")

		csvWriter := csv.NewWriter(c.Writer)
		headers := []string{
			"UserID", "FirstName", "Email",
			"ProductDetailID", "ProductName", "Price",
			"DiscountPercentage", "Quantity", "OrderID",
			"PlacedDate", "PaymentMode", "OrderStatus", "Total",
		}

		if err := csvWriter.Write(headers); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		grandtotal := 0
		for _, sales := range salesreport {
			discount := (sales.DiscountPercentage * sales.Price) / 100
			total := sales.Quantity * (sales.Price - discount)
			row := []string{
				fmt.Sprintf("%v", sales.UserID),
				sales.FirstName,
				sales.Email,
				fmt.Sprintf("%v", sales.ProductDetailID),
				sales.ProductName,
				fmt.Sprintf("%v", sales.Price),
				fmt.Sprintf("%v", sales.DiscountPercentage),
				fmt.Sprintf("%v", sales.Quantity),
				fmt.Sprintf("%v", sales.OrderID),
				sales.PlacedDate.Format("2006-01-02 15:04:05"),
				sales.PaymentMode,
				sales.OrderStatus,
				fmt.Sprintf("%v", total),
			}

			if err := csvWriter.Write(row); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}
			if sales.PaymentMode == "Razorpay" || sales.PaymentMode == "Wallet" || (sales.PaymentMode == "Cash on Delivery" && sales.OrderStatus == "Delivered") {
				grandtotal += int(total)
			}
		}
		rowtotal := []string{
			fmt.Sprintf("Grand Total=%v", grandtotal),
		}
		if err := csvWriter.Write(rowtotal); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		csvWriter.Flush()
	}
}

func (cr *AdminHandler) Dashboard(c *gin.Context) {
	reswidgets, err := cr.AdminUseCase.Widgets()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Widgets": reswidgets,
	})
}

func (cr *AdminHandler) AddCoupon(c *gin.Context) {
	var couponBody utils.BodyAddCoupon
	if err := c.BindJSON(&couponBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	date, err := time.Parse("2006-01-02", couponBody.ExpirationDate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	coupon := domain.Coupon{
		CouponCode:         couponBody.Code,
		CouponType:         couponBody.Type,
		Discount:           couponBody.Discount,
		UsageLimit:         couponBody.UsageLimit,
		ExpirationDate:     date,
		MinimumOrderAmount: couponBody.MinOrderAmount,
		ProductID:          couponBody.ProductID,
	}
	if err := cr.AdminUseCase.AddCoupon(coupon); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Success": "coupon added",
	})
}

func (cr *AdminHandler) GetAllCoupons(c *gin.Context) {
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
	coupons, err := cr.AdminUseCase.GetAllCoupons(pagination)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Coupons": coupons,
	})
}

func (cr *AdminHandler) UpdateCoupon(c *gin.Context) {
	couponid := c.Query("couponid")
	var couponBody utils.BodyAddCoupon
	if err := c.BindJSON(&couponBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	date, err := time.Parse("2006-01-02", couponBody.ExpirationDate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	coupon := domain.Coupon{
		CouponCode:         couponBody.Code,
		CouponType:         couponBody.Type,
		Discount:           couponBody.Discount,
		UsageLimit:         couponBody.UsageLimit,
		ExpirationDate:     date,
		MinimumOrderAmount: couponBody.MinOrderAmount,
		ProductID:          couponBody.ProductID,
	}
	if err := cr.AdminUseCase.UpdateCoupon(coupon, couponid); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Success": "coupon updated",
	})
}

func (cr *AdminHandler) GetCouponByID(c *gin.Context) {
	couponid := c.Param("couponid")
	coupon, err := cr.AdminUseCase.GetCouponByID(couponid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Success": coupon,
	})
}

func (cr *AdminHandler) DeleteCoupon(c *gin.Context) {
	couponid := c.Param("couponid")
	if err := cr.AdminUseCase.DeleteCoupon(couponid); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Success": "Deleted",
	})
}
