package utils

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var secretKey = []byte("blog-back")

type Payload struct {
	ID string
	jwt.RegisteredClaims
}

func GenerateToken(userID string, t time.Duration) (string, error) {
	payload := &Payload{
		ID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(), // 加个唯一标识
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(t)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "babalababa",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenString, err := token.SignedString(secretKey)
	log.Printf("用户'%s'token生成：%s\n", userID, tokenString)

	return tokenString, err
}

func ParseToken(tokenString string) (*Payload, error) {
	if tokenString == "" {
		return nil, jwt.ErrInvalidKey
	}

	token, err := jwt.ParseWithClaims(tokenString, &Payload{}, func(token *jwt.Token) (interface{}, error) {
		// 签名验证
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrInvalidKey
	}

	payload, ok := token.Claims.(*Payload)
	if !ok {
		return nil, jwt.ErrInvalidKey
	}

	log.Printf("token解析, 用户id：%s\n", payload.ID)

	return payload, nil
}
