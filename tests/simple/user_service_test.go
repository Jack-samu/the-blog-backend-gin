package simple

import (
	"testing"

	"github.com/Jack-samu/the-blog-backend-gin.git/internal/repositories"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/service"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/utils"
	"github.com/stretchr/testify/assert"
)

// 还缺了忘记密码部分和照片部分
func TestUserServiceWithoutAvatar(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db, t)

	repo := repositories.NewRepository(db)
	s := service.NewService(repo)

	err := s.Register("test-user", "test@test.com", "test123", "guest what", "")
	assert.Empty(t, err)

	// 重复注册测试
	err = s.Register("test-user", "test@test.com", "test123", "guest what", "")
	assert.Empty(t, err.Err)
	assert.Contains(t, err.Msg, "用户邮箱已注册")

	err = s.Register("test-user", "test1@test.com", "test123", "guest what", "")
	assert.Empty(t, err.Err)
	assert.Contains(t, err.Msg, "昵称重复")

	// 登录测试
	userInfo, err := s.Login("test-user", "test123")
	assert.Empty(t, err)
	assert.Equal(t, "test-user", userInfo.UserInfo.Username)

	// 登录后的token校验
	payload, err1 := utils.ParseToken(userInfo.Token)
	assert.NoError(t, err1)
	_, err1 = utils.ParseToken(userInfo.RefreshToken)
	assert.NoError(t, err1)

	// token刷新
	refreshResp, err := s.RefreshTheToken(payload.ID)
	assert.Empty(t, err)
	assert.Equal(t, userInfo.UserInfo.ID, refreshResp.UserInfo.ID)
	_, err1 = utils.ParseToken(refreshResp.Token)
	assert.NoError(t, err1)

	// 用户信息获取
	profileResp, err := s.Profile(userInfo.UserInfo.ID)
	assert.Empty(t, err)
	assert.Equal(t, userInfo.UserInfo.Username, profileResp.Username)
	assert.Equal(t, userInfo.UserInfo.Email, profileResp.Email)
	assert.Empty(t, profileResp.Avatar)

	// 用户图片获取
	photosResp, err := s.GetPhotos(userInfo.UserInfo.ID)
	assert.Empty(t, err)
	assert.Empty(t, photosResp.Photos)

	last_activity, err := s.Logout(userInfo.UserInfo.ID)
	assert.Empty(t, err)
	t.Logf("最后活动时间：%s\n", last_activity)
}
