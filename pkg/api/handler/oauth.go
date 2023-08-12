package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/stebinsabu13/ecommerce-api/pkg/auth"
	"github.com/stebinsabu13/ecommerce-api/pkg/config"
	services "github.com/stebinsabu13/ecommerce-api/pkg/usecase/interface"
	"github.com/stebinsabu13/ecommerce-api/pkg/utils"
)

type AuthHandler struct {
	authUseCase services.AuthUseCase
}

func NewAuthHandler(usecase services.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: usecase,
	}
}

func (c *AuthHandler) UserGoogleAuthLoginPage(ctx *gin.Context) {

	ctx.HTML(200, "auth.html", nil)
}

func (c *AuthHandler) UserGoogleAuthInitialize(ctx *gin.Context) {

	// setup the google provider
	goauthClientID := config.GetCofig().GOOGLE_CLIENT
	gouthClientSecret := config.GetCofig().GOOGLE_CLIENT_SECRET
	callbackUrl := "http://localhost:3000/user/login/callback"

	// setup privier
	goth.UseProviders(
		google.New(goauthClientID, gouthClientSecret, callbackUrl, "email", "profile"),
	)

	// start the google login
	gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
}

func (c *AuthHandler) UserGoogleAuthCallBack(ctx *gin.Context) {

	googleUser, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error() + "Failed to get user details from google",
		})
		return
	}

	user := utils.BodySignUpuser{
		FirstName: googleUser.FirstName,
		LastName:  googleUser.LastName,
		Email:     googleUser.Email,
	}

	userID, err := c.authUseCase.GoogleLogin(ctx, user)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error() + "Failed to login with google",
		})
		return
	}
	tokenString, err1 := auth.GenerateJWT(userID)
	if err1 != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "error generationg token",
		})
		return
	}
	ctx.SetCookie("user-token", tokenString, int(time.Now().Add(5*time.Minute).Unix()), "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"Success": user,
	})
}
