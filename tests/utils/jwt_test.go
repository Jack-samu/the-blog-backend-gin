package utils_test

import (
	"testing"
	"time"

	"github.com/Jack-samu/the-blog-backend-gin.git/internal/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// 测试jwt工具的token生成和token解析
func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name     string
		userID   string
		duration time.Duration
		wantErr  bool
	}{
		{
			name:     "正常",
			userID:   uuid.NewString(),
			duration: time.Hour,
			wantErr:  false,
		},
		{
			name:     "空用户",
			userID:   "",
			duration: time.Hour,
			wantErr:  false,
		},
		{
			name:     "负时间",
			userID:   "bbaa",
			duration: time.Hour,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := utils.GenerateToken(tt.userID, tt.duration)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)

				// token解析
				payload, err := utils.ParseToken(token)
				assert.NoError(t, err)
				assert.Equal(t, tt.userID, payload.ID)
			}
		})
	}
}

func TestParseToken(t *testing.T) {
	userID := "test-123-id"
	token, err := utils.GenerateToken(userID, time.Hour)
	assert.NoError(t, err)

	tests := []struct {
		name        string
		tokenString string
		wantErr     bool
	}{
		{
			name:        "有效token",
			tokenString: token,
			wantErr:     false,
		},
		{
			name:        "空token",
			tokenString: "",
			wantErr:     true,
		},
		{
			name:        "token篡改",
			tokenString: token + "fake",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload, err := utils.ParseToken(tt.tokenString)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, payload)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, payload)
				assert.Equal(t, userID, payload.ID)
			}
		})
	}
}

func TestExpiredToken(t *testing.T) {
	id := uuid.NewString()
	token, err := utils.GenerateToken(id, -time.Hour)
	assert.NoError(t, err)

	payload, err := utils.ParseToken(token)
	assert.Error(t, err)
	assert.Nil(t, payload)
}

// 基准测试
func BenchmarkGenerateToken(b *testing.B) {
	for i := 0; i < b.N; i++ {
		utils.GenerateToken("bench-user", time.Hour)
	}
}

func BenchmarkParseToken(b *testing.B) {
	token, _ := utils.GenerateToken("bench-user", time.Hour)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		utils.ParseToken(token)
	}
}
