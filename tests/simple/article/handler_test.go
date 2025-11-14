package article

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	dao "github.com/Jack-samu/the-blog-backend-gin.git/internal/DAO"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/dtos"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/handler"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type hFixture struct {
	db *gorm.DB
	h  *handler.Handler
}

func setupFixture(t *testing.T) *hFixture {
	db := setupTestDB(t)
	repo := dao.NewRepository(db)
	service := service.NewService(repo)
	handler := handler.NewHandler(service)

	// 准备测试用户
	err := service.Register("test-user", "test@test.com", "test1234", "guest what", "")
	assert.Nil(t, err)

	return &hFixture{
		db: db,
		h:  handler,
	}
}

func TestGetArticleHandler(t *testing.T) {
	r := gin.Default()
	fixture := setupFixture(t)
	setupRoutes(fixture.h, r)
	defer teardownTestDB(fixture.db)

	gin.SetMode(gin.TestMode)
	var token string

	t.Run("首页文章列表，公共路由", func(t *testing.T) {
		// 不带参数
		recorder := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/articles", nil)
		assert.NoError(t, err)
		r.ServeHTTP(recorder, req)

		var resp *dtos.PostListResp
		json.Unmarshal(recorder.Body.Bytes(), &resp)
		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, uint(0), resp.Cnt)

		// 带查询参数
		recorder = httptest.NewRecorder()
		req, err = http.NewRequest(http.MethodGet, "/articles?page=1&per_page=10", nil)
		assert.NoError(t, err)
		r.ServeHTTP(recorder, req)

		json.Unmarshal(recorder.Body.Bytes(), &resp)
		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, uint(0), resp.Cnt)
	})

	t.Run("登录准备", func(t *testing.T) {
		data := map[string]interface{}{
			"username": "test-user",
			"password": "test1234",
		}
		jsonData, err := json.Marshal(&data)
		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(jsonData))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(recorder, req)

		var loginResp *dtos.LoginResp
		json.Unmarshal(recorder.Body.Bytes(), &loginResp)
		assert.Equal(t, http.StatusOK, recorder.Code)
		token = loginResp.Token
	})

	t.Run("获取个人post", func(t *testing.T) {

		recorder := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/articles/publish", nil)
		assert.NoError(t, err)
		r.ServeHTTP(recorder, req)
		// 得到认证不通过
		assert.Equal(t, http.StatusUnauthorized, recorder.Code)

		// 添加认证
		recorder = httptest.NewRecorder()
		req.Header.Set("Authorization", "Bearer "+token)
		r.ServeHTTP(recorder, req)
		var postsResp *dtos.PostListResp
		json.Unmarshal(recorder.Body.Bytes(), &postsResp)
		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, uint(0), postsResp.Cnt)

		// 通过id查找不存在的post
		req, err = http.NewRequest(http.MethodGet, "/articles/publish/1", nil)
		assert.NoError(t, err)
		recorder = httptest.NewRecorder()
		req.Header.Set("Authorization", "Bearer "+token)
		r.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusNotFound, recorder.Code)
	})

	t.Run("获取个人draft", func(t *testing.T) {

		recorder := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/articles/drafts", nil)
		assert.NoError(t, err)
		r.ServeHTTP(recorder, req)
		// 得到认证不通过
		assert.Equal(t, http.StatusUnauthorized, recorder.Code)

		// 添加鉴权
		recorder = httptest.NewRecorder()
		req.Header.Set("Authorization", "Bearer "+token)
		r.ServeHTTP(recorder, req)
		var draftsResp *dtos.DraftsResp
		json.Unmarshal(recorder.Body.Bytes(), &draftsResp)
		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, uint(0), draftsResp.Cnt)

		// 通过id查找不存在的draft
		req, err = http.NewRequest(http.MethodGet, "/articles/draft/1", nil)
		assert.NoError(t, err)
		recorder = httptest.NewRecorder()
		req.Header.Set("Authorization", "Bearer "+token)
		r.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusNotFound, recorder.Code)
	})

	t.Run("保存草稿后查看对应信息", func(t *testing.T) {
		// 保存草稿
		reqData := dtos.ArticleReq{
			Title:   "测试",
			Excerpt: "摘要",
			Content: "正文",
			Cover:   "",
		}
		jsonData, err := json.Marshal(&reqData)
		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/articles/save", bytes.NewBuffer(jsonData))
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)
		r.ServeHTTP(recorder, req)
		// 草稿已保存，可以通过id和列表访问
		var resp map[string]interface{}
		json.Unmarshal(recorder.Body.Bytes(), &resp)
		assert.Equal(t, http.StatusCreated, recorder.Code)

		recorder = httptest.NewRecorder()
		req, err = http.NewRequest(http.MethodGet, "/articles/draft/1", nil)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)
		r.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusOK, recorder.Code)

		recorder = httptest.NewRecorder()
		req, err = http.NewRequest(http.MethodGet, "/articles/drafts", nil)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)
		r.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("发布", func(t *testing.T) {
		// 发布前面保存的draft
		reqData := map[string]interface{}{
			"id":      1,
			"title":   "测试1",
			"excerpt": "摘要1",
			"content": "正文1",
			"cover":   "bbbb",
		}
		jsonData, err := json.Marshal(&reqData)
		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/articles/publish", bytes.NewBuffer(jsonData))
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)
		r.ServeHTTP(recorder, req)

		// post已发布，可以通过id和列表访问
		var resp map[string]interface{}
		json.Unmarshal(recorder.Body.Bytes(), &resp)
		assert.Equal(t, http.StatusCreated, recorder.Code)

		recorder = httptest.NewRecorder()
		req, err = http.NewRequest(http.MethodGet, "/articles/1", nil)
		assert.NoError(t, err)
		r.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusOK, recorder.Code)

		recorder = httptest.NewRecorder()
		req, err = http.NewRequest(http.MethodGet, "/articles/drafts", nil)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)
		r.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusOK, recorder.Code)
	})
}
