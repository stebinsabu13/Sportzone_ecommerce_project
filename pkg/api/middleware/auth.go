package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stebinsabu13/ecommerce-api/pkg/auth"
)

func AuthorizationMiddleware(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie(role + "-token")
		if err != nil || tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Needs to login",
			})
			return
		}
		claims, err1 := auth.ValidateToken(tokenString)
		if err1 != nil {
			if err1.Error() == "token expired re-login" {
				refreshtokenstring, err2 := auth.GenerateJWT(claims.ID)
				if err2 != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"error": err2.Error(),
					})
					return
				}
				c.SetCookie(role+"-token", refreshtokenstring, int(time.Now().Add(5*time.Minute).Unix()), "/", "localhost", false, true)
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": err1.Error(),
				})
				return
			}
		}
		c.Set(role+"-id", claims.ID)
		c.Next()
	}
}
