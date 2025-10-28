package simple

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Jack-samu/the-blog-backend-gin.git/internal/dtos"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/handler"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/middleware"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/repositories"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRoutes(h *handler.Handler, r *gin.Engine) {
	r.POST("/auth/register", h.Register)
	r.POST("/auth/login", h.Login)

	protected := r.Group("/auth")
	protected.Use(middleware.Auth())
	{
		protected.POST("/:id/profile", h.Profile)
		protected.POST("/photos", h.GetPhotos)
		protected.POST("/refresh", h.RefreshTheToken)
		protected.POST("/logout", h.Logout)
	}
}

func TestUserFlowSimple(t *testing.T) {
	r := gin.Default()
	db := setupTestDB(t)
	defer teardownTestDB(db, t)

	repo := repositories.NewRepository(db)
	s := service.NewService(repo)
	h := handler.NewHandler(s)
	setupRoutes(h, r)
	gin.SetMode(gin.TestMode)

	testReq := dtos.RegisterReq{
		Username: "阿巴阿巴",
		Email:    "ababa@test.com",
		Password: "bbbbbaby",
		Bio:      "Guest what",
		Avatar:   "",
	}

	var token string
	var refreshToken string
	var user_id string

	t.Run("注册", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		form, _ := json.Marshal(testReq)
		req, err := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(form))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(recorder, req)
		var resp map[string]interface{}

		json.Unmarshal(recorder.Body.Bytes(), &resp)
		assert.Equal(t, http.StatusCreated, recorder.Code)
		assert.Contains(t, resp["msg"], "注册成功")
	})

	t.Run("登录", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		form, _ := json.Marshal(testReq)
		req, err := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(form))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(recorder, req)
		var resp dtos.LoginResp

		json.Unmarshal(recorder.Body.Bytes(), &resp)
		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, testReq.Username, resp.UserInfo.Username)
		token = resp.Token
		refreshToken = resp.RefreshToken
		user_id = resp.UserInfo.ID
	})

	t.Run("个人信息获取", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/auth/"+user_id+"/profile", nil)
		assert.NoError(t, err)

		req.Header.Set("Authorization", "Bearer "+refreshToken)
		r.ServeHTTP(recorder, req)
		var resp dtos.ProfileResp

		json.Unmarshal(recorder.Body.Bytes(), &resp)
		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, testReq.Username, resp.Username)
	})

	t.Run("token刷新", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/auth/refresh", nil)
		assert.NoError(t, err)

		req.Header.Set("Authorization", "Bearer "+refreshToken)
		r.ServeHTTP(recorder, req)
		var refreshResp dtos.RefreshResp

		json.Unmarshal(recorder.Body.Bytes(), &refreshResp)
		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.NotEqual(t, token, refreshResp.Token)
		token = refreshResp.Token
	})

	t.Run("退出登录", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/auth/logout", nil)
		assert.NoError(t, err)

		req.Header.Set("Authorization", "Bearer "+token)
		r.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})
}
