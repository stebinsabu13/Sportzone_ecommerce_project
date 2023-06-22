package handler

import (
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
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

// ADMIN LOGIN
//
//	@Summary		API FOR ADMIN LOGIN
//	@ID				ADMIN-LOGIN
//	@Description	VERIFY THE EMAIL,PASSWORD AND GENERATE A JWT TOKEN AND SET IT TO A COOKIE
//	@Tags			ADMIN
//	@Accept			json
//	@Produce		json
//	@Param			login_details	body		utils.BodyLogin	true	"Enter the email and password"
//	@Success		200				{object}	utils.Response
//	@Failure		400				{object}	utils.Response
//	@Failure		401				{object}	utils.Response
//	@Failure		500				{object}	utils.Response
//	@Router			/admin/login [post]
func (cr *AdminHandler) LoginHandler(c *gin.Context) {
	var body utils.BodyLogin
	if err := c.BindJSON(&body); err != nil {
		response := utils.ErrorResponse(400, "Error: Failed to read json body", err.Error(), body)
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	admin, err := cr.AdminUseCase.FindbyEmail(c.Request.Context(), body.Email)
	if err != nil {
		response := utils.ErrorResponse(401, "Error: Admin Doesn't exsist", err.Error(), body)
		c.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}
	if ok := support.CheckPasswordHash(body.Password, admin.Password); !ok {
		response := utils.ErrorResponse(401, "Error: Please check the password", err.Error(), body)
		c.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}
	tokenString, err := auth.GenerateJWT(admin.ID)
	if err != nil {
		response := utils.ErrorResponse(500, "Error: Error generating the jwt token", err.Error(), body)
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}
	c.SetCookie("admin-token", tokenString, int(time.Now().Add(5*time.Minute).Unix()), "/", "sportzone.cloud", true, true)
	c.Set("admin-id", admin.ID)
	response1 := utils.SuccessResponse(200, "Success: Login Successful")
	c.JSON(http.StatusOK, response1)
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

// ADMIN LOGOUT
//
//	@Summary		API FOR ADMIN LOGOUT
//	@ID				ADMIN-LOGOUT
//	@Description	ADMIN LOGOUT
//	@Tags			ADMIN
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	utils.Response
//	@Failure		401	{object}	utils.Response
//	@Failure		400	{object}	utils.Response
//	@Failure		500	{object}	utils.Response
//	@Router			/admin/logout [post]
func (cr *AdminHandler) LogoutHandler(c *gin.Context) {
	c.SetCookie("admin-token", "", -1, "/", "sportzone.cloud", true, true)
	response := utils.SuccessResponse(200, "Success: Logout Successful", nil)
	c.JSON(http.StatusOK, response)
}

// ADMIN SIGN-UP WITH SENDING OTP
//
//	@Summary		API FOR NEW ADMIN SIGN UP
//	@ID				SIGNUP-ADMIN
//	@Description	CREATE A NEW ADMIN WITH REQUIRED DETAILS
//	@Tags			ADMIN
//	@Accept			json
//	@Produce		json
//	@Param			admin_details	body		utils.BodySignUpuser	false	"New Admin Details"
//	@Success		200				{object}	utils.Response
//	@Failure		400				{object}	utils.Response
//	@Failure		400				{object}	utils.Response
//	@Failure		500				{object}	utils.Response
//	@Router			/admin/signup [post]
func (cr *AdminHandler) SignUp(c *gin.Context) {
	var signUp_user utils.BodySignUpuser
	if err := c.BindJSON(&signUp_user); err != nil {
		response := utils.ErrorResponse(400, "Error: Failed to read json body", err.Error(), signUp_user)
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	mobile_num, err := cr.AdminUseCase.SignUpAdmin(c.Request.Context(), signUp_user)
	if err != nil {
		response := utils.ErrorResponse(500, "Error: Check the error", err.Error(), nil)
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}
	respSid, err1 := cr.OtpUseCase.TwilioSendOTP(c.Request.Context(), mobile_num)
	if err1 != nil {
		response := utils.ErrorResponse(500, "Error: Failed to send otp", err.Error(), nil)
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Enter the otp and the responseid", respSid, signUp_user.ReferalCode)
	c.JSON(http.StatusOK, response)
}

// USER SIGN-UP WITH VERIFICATION OF OTP
//
//	@Summary		API FOR NEW ADMIN SIGN UP OTP VERIFICATION
//	@ID				SIGNUP-ADMIN-OTP-VERIFY
//	@Description	VERIFY THE OTP AND UPDATE THE VERIFIED COLUMN
//	@Tags			ADMIN
//	@Accept			json
//	@Produce		json
//	@Param			otp_details	body		utils.Otpverify	false	"otp"
//	@Success		200			{object}	utils.Response
//	@Failure		401			{object}	utils.Response
//	@Failure		400			{object}	utils.Response
//	@Failure		500			{object}	utils.Response
//	@Router			/admin/signup/otp/verify [post]
func (cr *AdminHandler) SignupOtpverify(c *gin.Context) {
	var OTP utils.Otpverify
	if err := c.BindJSON(&OTP); err != nil {
		response := utils.ErrorResponse(400, "Error: Failed to read json body", err.Error(), OTP)
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	session, err := cr.OtpUseCase.TwilioVerifyOTP(c.Request.Context(), OTP)
	if err != nil {
		response := utils.ErrorResponse(401, "Error: Verification failed", err.Error(), OTP)
		c.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}
	err1 := cr.AdminUseCase.UpdateVerify(session.MobileNum, OTP.ReferalCode)
	if err1 != nil {
		response := utils.ErrorResponse(500, "Error: Failed to update the verification status of admin", err.Error(), nil)
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}
	response1 := utils.SuccessResponse(200, "Success: Admin Phone Number Successfully verified", session.MobileNum)
	c.JSON(http.StatusOK, response1)
}

// LIST USERS
//
//	@Summary		API FOR LISTING USERS
//	@ID				ADMIN-LIST-USERS
//	@Description	LISTING ALL EXISTING USERS
//	@Tags			ADMIN
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int	false	"Enter the page number to display"
//	@Param			limit	query		int	false	"Number of items to retrieve per page"
//	@Success		200		{object}	utils.Response
//	@Failure		401		{object}	utils.Response
//	@Failure		400		{object}	utils.Response
//	@Failure		500		{object}	utils.Response
//	@Router			/admin/user [get]
func (cr *AdminHandler) ListAllUsers(c *gin.Context) {
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

// ACCESS HANDLER
//
//	@Summary		API FOR BLOCKING/UNBLOCKING USERS
//	@ID				ADMIN-ACCESS
//	@Description	GRANTING ACCESS FOR INDIVIDUAL USERS.
//	@Tags			ADMIN
//	@Accept			json
//	@Produce		json
//	@Param			userid	path		string	true	"Enter the specific user id"
//	@Param			access	query		string	false	"Enter true/false"
//	@Success		200		{object}	utils.Response
//	@Failure		401		{object}	utils.Response
//	@Failure		400		{object}	utils.Response
//	@Failure		500		{object}	utils.Response
//	@Router			/admin/user/{userid}/make [patch]
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

// LIST CATEGORY
//
//	@Summary		API FOR LISTING ALL CATEGORIES
//	@Description	LISTING ALL CATEGORIES FROM ADMINS AND USERS END
//	@Tags			ADMIN USER
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	utils.Response
//	@Failure		401	{object}	utils.Response
//	@Failure		400	{object}	utils.Response
//	@Failure		500	{object}	utils.Response
//	@Router			/user/filter/category [get]
//	@Router			/admin/category [get]
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

// CATEGORY MANAGEMENT

// ADD CATEGORY
//
//	@Summary		API FOR ADDING CATEGORY
//	@ID				ADMIN-ADD-CATEGORY
//	@Description	ADDING CATEGORY FROM ADMINS END
//	@Tags			ADMIN
//	@Accept			json
//	@Produce		json
//	@Param			category_details	body		utils.AddCategory	true	"Enter the category name"
//	@Success		200					{object}	utils.Response
//	@Failure		401					{object}	utils.Response
//	@Failure		400					{object}	utils.Response
//	@Failure		500					{object}	utils.Response
//	@Router			/admin/category/add [post]
func (cr *AdminHandler) AddCategory(c *gin.Context) {
	var body utils.AddCategory
	if err := c.BindJSON(&body); err != nil {
		response := utils.ErrorResponse(400, "Error: Error binding JSON", err.Error(), body)
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	var category domain.Category
	copier.Copy(&category, &body)
	if err := cr.AdminUseCase.AddCategory(c.Request.Context(), category); err != nil {
		response := utils.ErrorResponse(500, "Error: Faild to add cateogy", err.Error(), category)
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Success: successfully added a new category", category)
	c.JSON(http.StatusOK, response)
}

// DELETE CATEGORY
//
//	@Summary		API FOR DELETING A CATEGORY
//	@ID				ADMIN-DELETE-CATEGORY
//	@Description	DELETING CATEGORY AND ALSO CHECKING WHETHER IT HAS A EXISTING PRODUCT
//	@Tags			ADMIN
//	@Accept			json
//	@Produce		json
//	@Param			categoryid	path		string	true	"Enter the category id that you would like to delete"
//	@Success		200			{object}	utils.Response
//	@Failure		401			{object}	utils.Response
//	@Failure		400			{object}	utils.Response
//	@Failure		500			{object}	utils.Response
//	@Router			/admin/category/delete/{categoryid} [delete]
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

// ADMIN SALES REPORT
//
//	@Summary		API FOR GETTING SALES REPORT
//	@ID				ADMIN-SALES-REPORT
//	@Description	ADMIN SALES REPORT, VIA MONTHLY AND YEARLY
//	@Tags			ADMIN
//	@Accept			json
//	@Produce		json
//	@Param			frequency	query		string	true	"Enter frequency"
//	@Param			month		query		int		false	"Enter the month"
//	@Param			year		query		int		true	"Enter the year"
//	@Param			page_number	query		int		false	"Enter the page number to display"
//	@Param			count		query		int		false	"Number of items to retrieve per page"
//	@Success		200			{object}	utils.Response
//	@Failure		401			{object}	utils.Response
//	@Failure		400			{object}	utils.Response
//	@Failure		500			{object}	utils.Response
//	@Router			/admin/sales [get]
func (cr *AdminHandler) FullSalesReport(c *gin.Context) {
	// time
	monthInt, err1 := strconv.Atoi(c.DefaultQuery("month", "1"))
	month := time.Month(monthInt)
	year, err2 := strconv.Atoi(c.Query("year"))
	frequency := c.Query("frequency")

	// page
	count, err3 := strconv.Atoi(c.DefaultQuery("count", "5"))
	pageNumber, err4 := strconv.Atoi(c.DefaultQuery("page_number", "1"))
	err := errors.Join(err1, err2, err3, err4)
	if err != nil {
		response := utils.ErrorResponse(400, "Failed to parse the requried fields", err1.Error(), nil)
		c.JSON(http.StatusBadRequest, response)
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

// DELETE CATEGORY
//
//	@Summary		API FOR VIEWING THE DASHBOARD
//	@ID				ADMIN-VIEW-DASHBOARD
//	@Description	VIEWING DIFFERENT WIDGETS
//	@Tags			ADMIN
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	utils.ResWidgets
//	@Failure		401	{object}	utils.Response
//	@Failure		400	{object}	utils.Response
//	@Failure		500	{object}	utils.Response
//	@Router			/admin/dashboard [get]
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

// @Summary		API FOR ADDING THE COUPON
// @ID				ADMIN-ADD-COUPON
// @Description	ADDING COUPONS IN THE ADMINS END
// @Tags			ADMIN
// @Accept			json
// @Produce		json
//
// @Param			coupon_details	body		utils.BodyAddCoupon	false	"Enter the details of the coupon"
//
// @Success		200				{object}	utils.ResWidgets
// @Failure		401				{object}	utils.Response
// @Failure		400				{object}	utils.Response
// @Failure		500				{object}	utils.Response
// @Router			/admin/coupon/add [post]
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

// @Summary		API FOR VIEWING ALL COUPON
// @ID				ADMIN-VIEW-COUPON
// @Description	VIEWING COUPONS IN THE ADMINS END
// @Tags			ADMIN
// @Accept			json
// @Produce		json
//
// @Param			page	query		string	false	"Enter the page number"
// @Param			limit	query		string	false	"Enter the number of coupons to retrieve"
// @Success		200		{object}	utils.ResWidgets
// @Failure		401		{object}	utils.Response
// @Failure		400		{object}	utils.Response
// @Failure		500		{object}	utils.Response
// @Router			/admin/coupon [get]
func (cr *AdminHandler) GetAllCoupons(c *gin.Context) {
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

// @Summary		API FOR UPDATING COUPON
// @ID				ADMIN-UPDATING-COUPON
// @Description	UPDATING COUPONS IN THE ADMINS END
// @Tags			ADMIN
// @Accept			json
// @Produce		json
//
// @Param			couponid		query		string				true	"Enter the coupon id to update"
// @Param			coupon_details	body		utils.BodyAddCoupon	false	"Enter the body of the coupon"
// @Success		200				{object}	utils.ResWidgets
// @Failure		401				{object}	utils.Response
// @Failure		400				{object}	utils.Response
// @Failure		500				{object}	utils.Response
// @Router			/admin/coupon/update [patch]
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

// @Summary		API FOR VIEWING COUPON BY ID
// @ID				ADMIN-COUPON-BY-ID
// @Description	VIEWING COUPON BY ID IN THE ADMINS END
// @Tags			ADMIN
// @Accept			json
// @Produce		json
//
// @Param			couponid	path		string	true	"Enter the coupon id to view"
// @Success		200			{object}	utils.ResWidgets
// @Failure		401			{object}	utils.Response
// @Failure		400			{object}	utils.Response
// @Failure		500			{object}	utils.Response
// @Router			/admin/coupon/{couponid} [get]
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

// @Summary		API FOR DELETING COUPON
// @ID				ADMIN-COUPON-DELETION
// @Description	DELETING COUPON IN THE ADMINS END
// @Tags			ADMIN
// @Accept			json
// @Produce		json
//
// @Param			couponid	path		string	true	"Enter the coupon id to delete"
// @Success		200			{object}	utils.ResWidgets
// @Failure		401			{object}	utils.Response
// @Failure		400			{object}	utils.Response
// @Failure		500			{object}	utils.Response
// @Router			/admin/coupon/delete/{couponid} [delete]
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
