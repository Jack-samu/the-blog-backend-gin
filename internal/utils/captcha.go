package utils

import (
	"crypto/rand"
	"math/big"
	"sync"
	"time"
)

const (
	lowerLetters   = "abcdefghijklmnopqrstuvwxyz"
	upperLetters   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberLetters  = "0123456789"
	specialLetters = "!@#$%&*?"
)

type Captcha struct {
	Code  string
	Since time.Time
}

type UserReset struct {
	sync.Map
}

func GenerateCaptcha() (*Captcha, error) {
	// 6位验证码生成
	letters := []byte(lowerLetters + upperLetters + numberLetters + specialLetters)
	ret := make([]byte, 6)
	for i := 0; i < 6; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return nil, err
		}
		ret[i] = letters[num.Int64()]
	}

	return &Captcha{Code: string(ret), Since: time.Now()}, nil
}
