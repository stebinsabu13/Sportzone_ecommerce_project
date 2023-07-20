package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"

	domain "github.com/stebinsabu13/ecommerce-api/pkg/domain"
	"github.com/stebinsabu13/ecommerce-api/pkg/support"
	"github.com/stebinsabu13/ecommerce-api/pkg/utils"

	"github.com/stebinsabu13/ecommerce-api/pkg/auth"
	services "github.com/stebinsabu13/ecommerce-api/pkg/usecase/interface"
)

type UserHandler struct {
	userUseCase  services.UserUseCase
	otpUseCase   services.OtpUseCase
	orderUseCase services.OrderUseCase
}

func NewUserHandler(usecase services.UserUseCase, otpusecase services.OtpUseCase, orderusecase services.OrderUseCase) *UserHandler {
	return &UserHandler{
		userUseCase:  usecase,
		otpUseCase:   otpusecase,
		orderUseCase: orderusecase,
	}
}

// USER USERLOGIN
//
//	@Summary		API FOR USER LOGIN
//	@ID				USER-LOGIN
//	@Description	VERIFY THE EMAIL,PASSWORD AND GENERATE A JWT TOKEN AND SET IT TO A COOKIE
//	@Tags			USER
//	@Accept			json
//	@Produce		json
//	@Param			login_details	body		utils.BodyLogin	true	"Enter the email and password"
//	@Success		200				{object}	utils.Response
//	@Failure		401				{object}	utils.Response
//	@Failure		400				{object}	utils.Response
//	@Failure		500				{object}	utils.Response
//	@Router			/user/login [post]
func (cr *UserHandler) LoginHandler(c *gin.Context) {
	// implement login logic here

	var body utils.BodyLogin

	if err := c.BindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": errors.New("failed to bind the required fields"),
		})
		return
	}
	user, err := cr.userUseCase.FindbyEmail(c.Request.Context(), body.Email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}
	ok := support.CheckPasswordHash(body.Password, user.Password)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "invalid password",
		})
		return
	}
	tokenString, err1 := auth.GenerateJWT(user.ID)
	if err1 != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "error generationg token",
		})
		return
	}
	c.SetCookie("user-token", tokenString, int(time.Now().Add(5*time.Minute).Unix()), "/", "sportzone.cloud", false, true)
	c.JSON(http.StatusOK, gin.H{
		"Success": user,
	})
}

// USER SIGN-UP WITH SENDING OTP
//
//	@Summary		API FOR NEW USER SIGN UP
//	@ID				SIGNUP-USER
//	@Description	CREATE A NEW USER WITH REQUIRED DETAILS
//	@Tags			USER
//	@Accept			json
//	@Produce		json
//	@Param			user_details	body		utils.BodySignUpuser	false	"New user Details"
//	@Success		200				{object}	utils.Response
//	@Failure		400				{object}	utils.Response
//	@Failure		500				{object}	utils.Response
//	@Router			/user/signup [post]
func (cr *UserHandler) SignUp(c *gin.Context) {
	var signUp_user utils.BodySignUpuser
	if err := c.BindJSON(&signUp_user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	mobile_num, err := cr.userUseCase.SignUpUser(c.Request.Context(), signUp_user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	respSid, err1 := cr.otpUseCase.TwilioSendOTP(c.Request.Context(), mobile_num)
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

// USER SIGN-UP WITH VERIFICATION OF OTP
//
//	@Summary		API FOR NEW USER SIGN UP OTP VERIFICATION
//	@ID				SIGNUP-USER-OTP-VERIFY
//	@Description	VERIFY THE OTP AND UPDATE THE VERIFIED COLUMN
//	@Tags			USER
//	@Accept			json
//	@Produce		json
//	@Param			otp_details	body		utils.Otpverify	false	"otp"
//	@Success		200			{object}	utils.Response
//	@Failure		401			{object}	utils.Response
//	@Failure		400			{object}	utils.Response
//	@Failure		500			{object}	utils.Response
//	@Router			/user/signup/otp/verify [post]
func (cr *UserHandler) SignupOtpverify(c *gin.Context) {
	var OTP utils.Otpverify
	if err := c.BindJSON(&OTP); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Error binding json",
		})
		return
	}
	session, err := cr.otpUseCase.TwilioVerifyOTP(c.Request.Context(), OTP)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	err1 := cr.userUseCase.UpdateVerify(session.MobileNum, OTP.ReferalCode)
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

// USERLOGOUT
//
//	@Summary		API FOR USER LOGOUT
//	@ID				USER-LOGOUT
//	@Description	LOGOUT USER AND ALSO CLEAR COOKIES
//	@Tags			USER
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	utils.Response
//	@Failure		401	{object}	utils.Response
//	@Failure		400	{object}	utils.Response
//	@Failure		500	{object}	utils.Response
//	@Router			/user/logout [post]
func (cr *UserHandler) LogoutHandler(c *gin.Context) {
	c.SetCookie("user-token", "", -1, "/", "sportzone.cloud", false, true)
	c.JSON(http.StatusOK, gin.H{
		"logout": "Success",
	})
}

// USER LOGIN WITH SENDING OTP
//
//	@Summary		API FOR USER LOGIN USING OTP
//	@ID				LOGIN-USER-OTP
//	@Description	LOGIN A USER USING OTP
//	@Tags			USER
//	@Accept			json
//	@Produce		json
//	@Param			user_details	body		utils.OtpLogin	true	"Enter email and mobile number"
//	@Success		200				{object}	utils.Response
//	@Failure		400				{object}	utils.Response
//	@Failure		500				{object}	utils.Response
//	@Router			/user/login/otp [post]
func (cr *UserHandler) LoginOtp(c *gin.Context) {
	var body utils.OtpLogin
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	user, err := cr.userUseCase.FindbyEmailorMobilenum(c.Request.Context(), body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}
	respSid, err1 := cr.otpUseCase.TwilioSendOTP(c.Request.Context(), "+91"+user.MobileNum)
	if err1 != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err1.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Success":     "Enter the otp and response id",
		"response id": respSid,
	})
}

// USER LOGIN WITH VERIFICATION OF OTP
//
//	@Summary		API FOR USER LOGIN OTP VERIFICATION
//	@ID				LOGIN-USER-OTP-VERIFY
//	@Description	VERIFY THE OTP AND MAKE THE USER LOGGED IN
//	@Tags			USER
//	@Accept			json
//	@Produce		json
//	@Param			otp_details	body		utils.Otpverify	false	"otp"
//	@Success		200			{object}	utils.Response
//	@Failure		401			{object}	utils.Response
//	@Failure		400			{object}	utils.Response
//	@Failure		500			{object}	utils.Response
//	@Router			/user/login/otp/verify [post]
func (cr *UserHandler) LoginOtpverify(c *gin.Context) {
	var otp utils.Otpverify
	if err := c.BindJSON(&otp); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Error binding json",
		})
		return
	}
	session, err := cr.otpUseCase.TwilioVerifyOTP(c.Request.Context(), otp)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	user, err2 := cr.userUseCase.FindbyEmailorMobilenum(c.Request.Context(), utils.OtpLogin{Email: "", MobileNum: session.MobileNum})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err2.Error(),
		})
		return
	}
	tokenString, err1 := auth.GenerateJWT(user.ID)
	if err1 != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err1.Error(),
		})
		return
	}
	c.SetCookie("user-token", tokenString, int(time.Now().Add(60*time.Minute).Unix()), "/", "sportzone.cloud", true, true)
	c.JSON(http.StatusOK, gin.H{
		"Success": user,
	})
}

// VIEW PROFILE
//
//	@Summary		API FOR VIEW PROFILE
//	@ID				USER-PROFILE VIEEW
//	@Description	VIEW USER PROFILE
//	@Tags			USER
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	utils.Response
//	@Failure		401	{object}	utils.Response
//	@Failure		400	{object}	utils.Response
//	@Failure		500	{object}	utils.Response
//	@Router			/user/profile [get]
func (cr *UserHandler) ShowUserDetails(c *gin.Context) {
	id, ok := c.Get("user-id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Not ok",
		})
		return
	}
	details, err := cr.userUseCase.ShowDetails(c.Request.Context(), id.(uint))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	address, err := cr.userUseCase.ShowAddress(c.Request.Context(), id.(uint))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	profile := support.BuildProfile(details, address)
	c.JSON(http.StatusOK, gin.H{
		"Profile": profile,
	})
}

// LIST ADDRESS
//
//	@Summary		API FOR LISTING ADDRESSES
//	@ID				USER-LIST-ADDRESS
//	@Description	LISTING ALL ADDRESSES FOR THE PARTICULAR USER
//	@Tags			USER
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	utils.Response
//	@Failure		401	{object}	utils.Response
//	@Failure		400	{object}	utils.Response
//	@Failure		500	{object}	utils.Response
//	@Router			/user/profile/address [get]
func (cr *UserHandler) ShowAllAddress(c *gin.Context) {
	id, ok := c.Get("user-id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Not ok",
		})
		return
	}
	addresses, err := cr.userUseCase.ShowAddress(c.Request.Context(), id.(uint))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"addresses": addresses,
	})
}

// ADD ADDRESS
//
//	@Summary		API FOR ADDING ADDRESS
//	@ID				USER-ADD-ADDRESS
//	@Description	ADDING NEW ADDRESS TO USER PROFILE
//	@Tags			USER
//	@Accept			json
//	@Produce		json
//	@Param			address_details	body		utils.AddAddress	true	"Add the address details"
//	@Success		200				{object}	utils.Response
//	@Failure		401				{object}	utils.Response
//	@Failure		400				{object}	utils.Response
//	@Failure		500				{object}	utils.Response
//	@Router			/user/profile/address/add [post]
func (cr *UserHandler) AddAddress(c *gin.Context) {
	var body utils.AddAddress
	id, ok := c.Get("user-id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Not ok",
		})
		return
	}
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var address domain.Address
	copier.Copy(&address, &body)
	address.UserID = id.(uint)
	if err := cr.userUseCase.AddAddress(c.Request.Context(), address); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Success": "Address added",
	})
}

// EDIT PROFILE
//
//	@Summary		API FOR EDIT PROFILE
//	@ID				USER-PROFILE EDIT
//	@Description	EDIT/UPDATE USER PROFILE
//	@Tags			USER
//	@Accept			json
//	@Produce		json
//	@Param			update_details	body		utils.EditProfileReq	true	"Edit the details as per wish"
//	@Success		200				{object}	utils.Response
//	@Failure		401				{object}	utils.Response
//	@Failure		400				{object}	utils.Response
//	@Failure		500				{object}	utils.Response
//	@Router			/user/profile/edit/profile [patch]
func (cr *UserHandler) UpdateProfile(c *gin.Context) {
	var profile utils.EditProfileReq
	if err := c.BindJSON(&profile); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	id, ok := c.Get("user-id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Not ok",
		})
		return
	}
	if err := cr.userUseCase.EditProfile(c.Request.Context(), profile, id.(uint)); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Success": "Profile updated",
	})
}

// USER FORGOT PASSWORD
//
//	@Summary		API FOR USER FORGOT PASSWORD OPTION
//	@ID				USER-FORGOT-PASSWORD
//	@Description	VERIFY THE EMAIL AND NUMBER AND FIND THE DATA. SEND THE OTP AND VERIFY WITH NEW PASSWORD AND OTP.
//	@Tags			USER
//	@Accept			json
//	@Produce		json
//	@Param			login_details	body		utils.OtpLogin	true	"Enter the email and phoneNumber"
//	@Success		200				{object}	utils.Response
//	@Failure		401				{object}	utils.Response
//	@Failure		400				{object}	utils.Response
//	@Failure		500				{object}	utils.Response
//	@Router			/user/forgot/password [post]
func (cr *UserHandler) ForgotPassword(c *gin.Context) {
	var bodydetail utils.OtpLogin
	if err := c.BindJSON(&bodydetail); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	user, err := cr.userUseCase.FindbyEmailorMobilenum(c.Request.Context(), bodydetail)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	respSid, err := cr.otpUseCase.TwilioSendOTP(c.Request.Context(), user.MobileNum)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Success":    "Enter the otp,responseid and the new password",
		"responseid": respSid,
	})
}

// USER FORGOT PASSWORD OTP VERIFY
//
//	@Summary		API FOR USER FORGOT PASSWORD OTP VERIFICATION
//	@ID				USER-FORGOT-PASSWORD-OTP-VERIFY
//	@Description	VERIFY THE OTP AND ENTER A NEW PASSWORD
//	@Tags			USER
//	@Accept			json
//	@Produce		json
//	@Param			verify_details	body		utils.Otpverify	false	"Enter the Otp and New Password"
//	@Success		200				{object}	utils.Response
//	@Failure		401				{object}	utils.Response
//	@Failure		400				{object}	utils.Response
//	@Failure		500				{object}	utils.Response
//	@Router			/user/forgot/password/otp/verify [patch]
func (cr *UserHandler) ForgotPasswordOtpverify(c *gin.Context) {
	var changepassbody utils.Otpverify
	if err := c.BindJSON(&changepassbody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	session, err := cr.otpUseCase.TwilioVerifyOTP(c.Request.Context(), changepassbody)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	newpassword, err1 := support.HashPassword(changepassbody.NewPassword)
	if err1 != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err1.Error(),
		})
		return
	}
	if err := cr.userUseCase.ChangePassword(c.Request.Context(), newpassword, session.MobileNum); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err1.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Success": "Password updated",
	})
}

// LIST CATEGORY
//
//	@Summary		API FOR LISTING ALL CATEGORIES
//	@Description	LISTING ALL CATEGORIES FROM USERS END
//	@Tags			USER
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	utils.Response
//	@Failure		401	{object}	utils.Response
//	@Failure		400	{object}	utils.Response
//	@Failure		500	{object}	utils.Response
//	@Router			/user/filter/category [get]
func (cr *UserHandler) ListAllCategories(c *gin.Context) {
	categories, err := cr.userUseCase.ListAllCategories(c.Request.Context())
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

// @Summary		API FOR VIEWING THE WALLER
// @Description	VIEWING THE WALLET FROM USERS END
// @Tags			USER
// @Accept			json
// @Produce		json
// @Success		200	{object}	utils.Response
// @Failure		401	{object}	utils.Response
// @Failure		400	{object}	utils.Response
// @Failure		500	{object}	utils.Response
// @Router			/user/wallet [get]
func (cr *UserHandler) ViewWallet(c *gin.Context) {
	userid, ok := c.Get("user-id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Not ok",
		})
		return
	}
	wallet, balance, err := cr.userUseCase.ViewWallet(userid.(uint))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Wallet":            wallet,
		"Balance in wallet": balance,
	})
}
