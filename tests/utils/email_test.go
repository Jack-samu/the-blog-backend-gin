package utils_test

import (
	"os"
	"testing"

	"github.com/Jack-samu/the-blog-backend-gin.git/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestSendEmail(t *testing.T) {
	// err := godotenv.Load("../../.env")
	// assert.NoError(t, err)
	utils.InitEmailConfig()

	to := os.Getenv("TEST_MAIN")
	t.Logf("测试目标：%s\n", to)

	err := utils.SendCaptcha(to, "测试邮件")
	assert.NoError(t, err)
}
