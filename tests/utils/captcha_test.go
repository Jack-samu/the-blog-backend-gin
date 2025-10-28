package utils_test

import (
	"testing"

	"github.com/Jack-samu/the-blog-backend-gin.git/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestGenerateCap(t *testing.T) {
	captcha, err := utils.GenerateCaptcha()
	assert.NoError(t, err)
	t.Logf("验证码：%s\n", captcha.Code)
}

func BenchmarkGenerateCaptcha(b *testing.B) {
	for i := 0; i < b.N; i++ {
		utils.GenerateCaptcha()
	}
}
