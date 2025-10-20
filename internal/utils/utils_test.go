package utils

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJwt(t *testing.T) {
	fakeUserID := uuid.New().String()
	t.Logf("测试用id：%v\n", fakeUserID)
	token, err := GenerateToken(fakeUserID, 5*time.Minute)
	t.Logf("生成5分钟后过期的token：%v\n", token)
	if err != nil {
		t.Fatalf("生成token出错：%v\n", err.Error())
	}

	payload, err := ParseToken(token)
	if err != nil {
		t.Fatalf("生成token出错：%v\n", err.Error())
	}

	if payload["sub"] != fakeUserID {
		t.Fatalf("解析token后出错，id：%v\n", payload["sub"])
	}
}
