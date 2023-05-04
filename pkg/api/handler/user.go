package handler

import (
	"fmt"
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
var loginOtp_user domain.User

type UserHandler struct {
	userUseCase services.UserUseCase
}

// type Response struct {
// 	ID      uint   `copier:"must"`
// 	Name    string `copier:"must"`
// 	Surname string `copier:"must"`
// }

func NewUserHandler(usecase services.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: usecase,
	}
}

func (cr *UserHandler) LoginHandler(c *gin.Context) {
	_, err1 := c.Cookie("user-token")
	if err1 == nil {
		c.Redirect(http.StatusFound, "/user/home")
		return
	}
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
	tokenString, err1 := auth.GenerateJWT(user.Email)
	if err1 != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "error generationg token",
		})
		return
	}
	c.SetCookie("user-token", tokenString, int(time.Now().Add(60*time.Minute).Unix()), "/", "localhost", false, true)
	c.Set("user-email", user.Email)
	c.JSON(http.StatusOK, gin.H{
		"Success": "Login",
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

	_, err1 := utils.TwilioSendOTP("+91" + signUp_user.MobileNum)
	if err1 != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed generating otp",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Success": "Enter the otp",
	})
}

func (cr *UserHandler) SignupOtpverify(c *gin.Context) {
	var otp utils.Otpverify
	if err := c.BindJSON(&otp); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Error binding json",
		})
		return
	}
	if err := utils.TwilioVerifyOTP("+91"+signUp_user.MobileNum, otp.Otp); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Otp",
		})
		return
	}
	signUp_user.Password, _ = support.HashPassword(signUp_user.Password)
	err := cr.userUseCase.SignUpUser(c.Request.Context(), signUp_user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "failed to add",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"User registration": "Success",
	})
}
func (cr *UserHandler) HomeHandler(c *gin.Context) {
	email, ok := c.Get("user-email")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user",
		})
	}
	user, err := cr.userUseCase.FindbyEmail(c.Request.Context(), email.(string))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user",
		})
		return
	}
	// response := []Response{}
	// copier.Copy(&response, &users)

	c.JSON(http.StatusOK, user)
}

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
	loginOtp_user, err := cr.userUseCase.FindbyEmailorMobilenum(c.Request.Context(), body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}
	_, err1 := utils.TwilioSendOTP("+91" + loginOtp_user.MobileNum)
	if err1 != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed generating otp",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Success": "Enter the otp",
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
	fmt.Println(loginOtp_user.MobileNum)
	if err := utils.TwilioVerifyOTP("+91"+loginOtp_user.MobileNum, otp.Otp); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Otp",
		})
		return
	}
	tokenString, err1 := auth.GenerateJWT(loginOtp_user.Email)
	if err1 != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "error generationg token",
		})
		return
	}
	c.SetCookie("user-token", tokenString, int(time.Now().Add(60*time.Minute).Unix()), "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{
		"Success": "Login",
	})
}

func (cr *UserHandler) ShowUserDetails(c *gin.Context) {
	id := c.Param("userid")
	details, err := cr.userUseCase.ShowDetails(c.Request.Context(), id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	address, err := cr.userUseCase.ShowAddress(c.Request.Context(), id)
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
