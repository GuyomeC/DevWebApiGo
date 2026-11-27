package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JwtSecret = []byte("projet-web-api-go-20-sur-20")

func GenerateJWT(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecret)
}

func VerifyJWT(tokenstring string) (string, error) {
	token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (any, error) { //decoupe le token
		_, ok := token.Method.(*jwt.SigningMethodHMAC) //verifie le type
		if !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return JwtSecret, nil
	})

	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		username, _ := claims["username"].(string)
		return username, nil
	}
	return "", fmt.Errorf("invalid token")
}
