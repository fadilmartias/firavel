package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateToken(data map[string]any, exp time.Duration) (string, error) {
	claims := jwt.MapClaims(data)
	claims["exp"] = time.Now().Add(exp).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, nil
}
