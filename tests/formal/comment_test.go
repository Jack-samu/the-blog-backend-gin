package formal

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	dao "github.com/Jack-samu/the-blog-backend-gin.git/internal/DAO"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/dtos"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/handler"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/models"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func preparation(s *service.Service, t *testing.T) (string, int) {
	err := s.Register("test-user", "test@test.com", "test1234", "guest what", "")
	assert.Empty(t, err)
	resp, err := s.Login("test-user", "test1234")
	assert.Empty(t, err)

	req := &dtos.ArticleReq{
		Title:   "测试",
		Content: "test-content",
		Excerpt: "test-excerpt",
		Cover:   "",
	}

	postID, err := s.PublishArticle(req, resp.UserInfo.ID)
	assert.Nil(t, err)

	return resp.UserInfo.ID, postID
}

func TestComment(t *testing.T) {

	// err := godotenv.Load("../../.env")
	// if err != nil {
	// 	log.Fatal("读取.env失败")
	// }

	r := gin.Default()
	db := models.InitDB()
	defer teardownTestDB(db)

	repo := dao.NewRepository(db)
	serv := service.NewService(repo)
	h := handler.NewHandler(serv)

	setupRoutes(h, r)

	t.Run("查询不存在的comment和reply", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/articles/1/comments", nil)
		assert.NoError(t, err)
		r.ServeHTTP(recorder, req)
		var resp map[string]interface{}
		json.Unmarshal(recorder.Body.Bytes(), &resp)
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.Equal(t, "没有找到对应的文章，无法刷新评论", resp["err"])

		recorder = httptest.NewRecorder()
		req, err = http.NewRequest(http.MethodGet, "/articles/1/replies", nil)
		assert.NoError(t, err)
		r.ServeHTTP(recorder, req)
		json.Unmarshal(recorder.Body.Bytes(), &resp)
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.Equal(t, "没有找到对应的评论，无法加载", resp["err"])

		recorder = httptest.NewRecorder()
		req, err = http.NewRequest(http.MethodGet, "/articles/-1/comments", nil)
		assert.NoError(t, err)
		r.ServeHTTP(recorder, req)
		json.Unmarshal(recorder.Body.Bytes(), &resp)
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Equal(t, "参数错误", resp["err"])
	})

	//准备文章主体和用户主体
	var token string
	var commentID int64
	_, post_id := preparation(serv, t)

	t.Run("创建comment而后修改comment", func(t *testing.T) {
		// 登录
		data := map[string]interface{}{
			"username": "test-user",
			"password": "test1234",
		}
		loginReq, err := json.Marshal(&data)
		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(loginReq))
		assert.NoError(t, err)
		r.ServeHTTP(recorder, req)
		var resp map[string]interface{}
		json.Unmarshal(recorder.Body.Bytes(), &resp)
		assert.Equal(t, http.StatusOK, recorder.Code)
		token, _ = resp["token"].(string)
		t.Logf("鉴权用：%s\n", token)

		// 创建评论
		commentReq := &dtos.CommentReq{
			ArticleID: int64(post_id),
			Content:   "comment测试",
		}
		reqData, err := json.Marshal(commentReq)
		assert.NoError(t, err)
		recorder = httptest.NewRecorder()
		req, err = http.NewRequest(http.MethodPost, "/articles/comments", bytes.NewBuffer(reqData))
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)
		r.ServeHTTP(recorder, req)
		json.Unmarshal(recorder.Body.Bytes(), &resp)
		t.Logf("响应：%v\n", recorder.Body)
		assert.Equal(t, http.StatusCreated, recorder.Code)
		comment, ok := resp["comment"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "comment测试", comment["content"])
		comment_id, ok := comment["id"].(float64)
		assert.True(t, ok)
		commentID = int64(comment_id)
		assert.NoError(t, err)

		// 查看结果
		recorder = httptest.NewRecorder()
		req, err = http.NewRequest(http.MethodGet, "/articles/2/comments", nil)
		assert.NoError(t, err)
		r.ServeHTTP(recorder, req)
		json.Unmarshal(recorder.Body.Bytes(), &resp)
		assert.Equal(t, http.StatusOK, recorder.Code)
		comments, ok := resp["comments"].([]interface{})
		assert.True(t, ok)
		comment, ok = comments[0].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "comment测试", comment["content"])

		recorder = httptest.NewRecorder()
		req, err = http.NewRequest(http.MethodGet, "/articles/1/replies", nil)
		assert.NoError(t, err)
		r.ServeHTTP(recorder, req)
		json.Unmarshal(recorder.Body.Bytes(), &resp)
		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Empty(t, resp["replies"])
		// 修改评论
		commentReq = &dtos.CommentReq{
			ArticleID: int64(post_id),
			Content:   "comment测试，改",
			CommentID: commentID,
		}
		reqData, err = json.Marshal(commentReq)
		assert.NoError(t, err)
		recorder = httptest.NewRecorder()
		req, err = http.NewRequest(http.MethodPost, "/comments/modify", bytes.NewBuffer(reqData))
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)
		r.ServeHTTP(recorder, req)
		json.Unmarshal(recorder.Body.Bytes(), &resp)
		t.Logf("响应：%v\n", recorder.Body)
		assert.Equal(t, http.StatusCreated, recorder.Code)
		comment, ok = resp["comment"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "comment测试，改", comment["content"])
		// 查看结果
		recorder = httptest.NewRecorder()
		req, err = http.NewRequest(http.MethodGet, "/articles/2/comments", nil)
		assert.NoError(t, err)
		r.ServeHTTP(recorder, req)
		json.Unmarshal(recorder.Body.Bytes(), &resp)
		assert.Equal(t, http.StatusOK, recorder.Code)
		comments, ok = resp["comments"].([]interface{})
		assert.True(t, ok)
		comment, ok = comments[0].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "comment测试，改", comment["content"])
	})

	t.Run("创建reply且修改reply", func(t *testing.T) {
		// 创建reply
		replyReq := &dtos.CommentReq{
			ArticleID: 0,
			Content:   "reply测试",
			CommentID: commentID,
		}
		t.Logf("reply参数：%v\n", replyReq)
		reqData, err := json.Marshal(replyReq)
		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/articles/replies", bytes.NewBuffer(reqData))
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)
		r.ServeHTTP(recorder, req)
		var resp map[string]interface{}
		json.Unmarshal(recorder.Body.Bytes(), &resp)
		assert.Equal(t, http.StatusCreated, recorder.Code)
		reply, ok := resp["reply"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "reply测试", reply["content"])
		// 查结果
		recorder = httptest.NewRecorder()
		req, err = http.NewRequest(http.MethodGet, "/articles/1/replies", nil)
		assert.NoError(t, err)
		r.ServeHTTP(recorder, req)
		json.Unmarshal(recorder.Body.Bytes(), &resp)
		assert.Equal(t, http.StatusOK, recorder.Code)
		replies, ok := resp["replies"].([]interface{})
		assert.True(t, ok)
		reply, ok = replies[0].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "reply测试", reply["content"])
		reply_id, ok := reply["id"].(float64)
		assert.True(t, ok)

		// 修改reply
		replyReq = &dtos.CommentReq{
			ReplyID:   uint(reply_id),
			Content:   "reply测试，改",
			CommentID: int64(commentID),
		}
		t.Logf("reply修改参数：%v\n", replyReq)
		reqData, err = json.Marshal(replyReq)
		assert.NoError(t, err)
		recorder = httptest.NewRecorder()
		req, err = http.NewRequest(http.MethodPost, "/replies/modify", bytes.NewBuffer(reqData))
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)
		r.ServeHTTP(recorder, req)
		json.Unmarshal(recorder.Body.Bytes(), &resp)
		assert.Equal(t, http.StatusCreated, recorder.Code)
		reply, ok = resp["reply"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "reply测试，改", reply["content"])

		// 查结果
		recorder = httptest.NewRecorder()
		req, err = http.NewRequest(http.MethodGet, "/articles/1/replies", nil)
		assert.NoError(t, err)
		r.ServeHTTP(recorder, req)
		json.Unmarshal(recorder.Body.Bytes(), &resp)
		assert.Equal(t, http.StatusOK, recorder.Code)
		replies, ok = resp["replies"].([]interface{})
		assert.True(t, ok)
		reply, ok = replies[0].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "reply测试，改", reply["content"])
	})

	t.Run("删除comment", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodDelete, "/comments/1", nil)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)
		r.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusCreated, recorder.Code)

		// 查看结果
		recorder = httptest.NewRecorder()
		req, err = http.NewRequest(http.MethodGet, "/articles/2/comments", nil)
		assert.NoError(t, err)
		r.ServeHTTP(recorder, req)
		var resp map[string]interface{}
		json.Unmarshal(recorder.Body.Bytes(), &resp)
		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Empty(t, resp["comments"])

		recorder = httptest.NewRecorder()
		req, err = http.NewRequest(http.MethodGet, "/articles/1/replies", nil)
		assert.NoError(t, err)
		r.ServeHTTP(recorder, req)
		json.Unmarshal(recorder.Body.Bytes(), &resp)
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.Equal(t, "没有找到对应的评论，无法加载", resp["err"])
	})
}
