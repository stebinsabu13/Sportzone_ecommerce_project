package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stebinsabu13/ecommerce-api/pkg/config"
)

type Claims struct {
	ID uint
	jwt.RegisteredClaims
}

func GenerateJWT(id uint) (string, error) {

	expireTime := time.Now().Add(60 * time.Minute)

	// create token with expire time and claims id as user id
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
		},
	})

	// conver the token into signed string
	tokenString, err := token.SignedString([]byte(config.GetJWTCofig()))

	if err != nil {
		return "", err
	}
	// refresh token add next time
	return tokenString, nil
}

func ValidateToken(tokenString string) (Claims, error) {
	claims := Claims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(config.GetJWTCofig()), nil
		},
	)
	if err != nil || !token.Valid {
		return claims, errors.New("not valid token")
	}
	//checking the expiry of the token
	if time.Now().Unix() > claims.ExpiresAt.Unix() {
		return claims, errors.New("token expired re-login")
	}
	return claims, nil
}
