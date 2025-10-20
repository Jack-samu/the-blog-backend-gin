package utils

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("blog-back")

func GenerateToken(userID string, t time.Duration) (string, error) {
	payload := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(t),
	}

	log.Printf("token生成：%v\n", payload)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString(secretKey)
}

func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	log.Printf("token解析：%v\n", payload["id"])

	if !ok {
		return nil, err
	}

	return payload, nil
}
