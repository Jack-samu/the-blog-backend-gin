package utils

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestSendEmail(t *testing.T) {
	err := godotenv.Load("../../.env")
	assert.NoError(t, err)
	InitEmailConfig()

	to := os.Getenv("TEST_MAIN")
	t.Logf("测试目标：%s\n", to)

	err = sendEmail(to, "test测试测测测", "测试邮件")
	assert.NoError(t, err)
}
