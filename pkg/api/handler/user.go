package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	domain "github.com/stebinsabu13/ecommerce-api/pkg/domain"
	"github.com/stebinsabu13/ecommerce-api/pkg/support"
	"github.com/stebinsabu13/ecommerce-api/pkg/utils"

	"github.com/stebinsabu13/ecommerce-api/pkg/auth"
	services "github.com/stebinsabu13/ecommerce-api/pkg/usecase/interface"
)

var signUp_user domain.User

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

func (cr *UserHandler) LoginHandler(c *gin.Context) {
	// implement login logic here

	var body utils.BodyLogin

	if err := c.BindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
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
			"error": "Invalid Password",
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
	c.SetCookie("user-token", tokenString, int(time.Now().Add(60*time.Minute).Unix()), "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{
		"Success": user,
	})
}

func (cr *UserHandler) SignUp(c *gin.Context) {
	if err := c.BindJSON(&signUp_user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if ok := support.Email_validater(signUp_user.Email); !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Email format incorrect",
		})
		return
	}

	if ok := support.MobileNum_validater(signUp_user.MobileNum); !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Not a valid mobile number",
		})
		return
	}
	if _, err := cr.userUseCase.FindbyEmail(c.Request.Context(), signUp_user.Email); err == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "User already Exsists",
		})
		return
	}
	signUp_user.Password, _ = support.HashPassword(signUp_user.Password)
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
		"Success":    "Enter the otp and the responseid",
		"responseid": respSid,
	})
}

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
	err1 := cr.userUseCase.UpdateVerify(c.Request.Context(), session.MobileNum)
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

// func (cr *UserHandler) HomeHandler(c *gin.Context) {
// 	email, ok := c.Get("user-email")
// 	if !ok {
// 		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
// 			"error": "Invalid user",
// 		})
// 	}
// 	user, err := cr.userUseCase.FindbyEmail(c.Request.Context(), email.(string))
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
// 			"error": "Invalid user",
// 		})
// 		return
// 	}
// 	// response := []Response{}
// 	// copier.Copy(&response, &users)

// 	c.JSON(http.StatusOK, user)
// }

func (cr *UserHandler) LogoutHandler(c *gin.Context) {
	c.SetCookie("user-token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{
		"logout": "Success",
	})
}

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
	c.SetCookie("user-token", tokenString, int(time.Now().Add(60*time.Minute).Unix()), "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{
		"Success": user,
	})
}

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

func (cr *UserHandler) ShowOrders(c *gin.Context) {
	id, ok := c.Get("user-id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Not ok",
		})
		return
	}
	orderDetails, err := cr.orderUseCase.OrderDetails(c.Request.Context(), id.(uint))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ORDER DETAILS": orderDetails,
	})
}

func (cr *UserHandler) AddAddress(c *gin.Context) {
	id, ok := c.Get("user-id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Not ok",
		})
		return
	}
	var address domain.Address
	if err := c.BindJSON(&address); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
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
