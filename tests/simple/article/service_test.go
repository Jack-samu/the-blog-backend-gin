package article

import (
	"net/http"
	"testing"

	dao "github.com/Jack-samu/the-blog-backend-gin.git/internal/DAO"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/dtos"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/service"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type tFixture struct {
	db     *gorm.DB
	repo   *dao.DAO
	serv   *service.Service
	userID string
}

func setupTestFixture(t *testing.T) *tFixture {
	db := setupTestDB(t)
	repo := dao.NewRepository(db)
	serv := service.NewService(repo)

	return &tFixture{
		db:     db,
		repo:   repo,
		serv:   serv,
		userID: createTestUser(serv, t),
	}
}

func TestGetPosts(t *testing.T) {
	fixture := setupTestFixture(t)
	defer teardownTestDB(fixture.db)

	resp, err := fixture.serv.GetPosts(1, 10)
	assert.Empty(t, err)
	assert.Empty(t, resp.Posts)
}

func TestGetPost(t *testing.T) {
	fixture := setupTestFixture(t)
	defer teardownTestDB(fixture.db)

	resp, err := fixture.serv.GetPost(1)
	assert.NotEmpty(t, err)
	assert.Empty(t, resp)
}

func TestGetDraft(t *testing.T) {
	fixture := setupTestFixture(t)
	defer teardownTestDB(fixture.db)

	resp, err := fixture.serv.GetDraft(1)
	assert.NotEmpty(t, err)
	assert.Empty(t, resp)
}

func TestGetPerson(t *testing.T) {
	fixture := setupTestFixture(t)
	defer teardownTestDB(fixture.db)

	// 空的列表且查询中并没有系统错误
	postsResp, err := fixture.serv.GetPostsOfUser(1, 10, fixture.userID)
	assert.Nil(t, err)
	assert.Empty(t, postsResp.Posts)

	draftsResp, err := fixture.serv.GetDraftsOfUser(1, 10, fixture.userID)
	assert.Nil(t, err)
	assert.Empty(t, draftsResp.Drafts)
}

func TestGetSeries(t *testing.T) {
	fixture := setupTestFixture(t)
	defer teardownTestDB(fixture.db)

	// 对不存在的用户进行查询，因为用户鉴权在handler层进行了
	seriesResp, err := fixture.serv.GetSeries("abab")
	assert.Nil(t, seriesResp)
	assert.Equal(t, http.StatusNotFound, err.Code)

	// 用户存在，但没有分类信息
	seriesResp, err = fixture.serv.GetSeries(fixture.userID)
	assert.Empty(t, seriesResp.Categories)
	assert.Nil(t, err)
}

func TestPublish(t *testing.T) {
	fixture := setupTestFixture(t)
	defer teardownTestDB(fixture.db)

	req := &dtos.ArticleReq{
		Id:      0,
		Title:   "测试",
		Excerpt: "摘要abstract",
		Content: "正文",
		Cover:   "",

		// 附加项
		Category: "技术",
		Tags:     []string{"t1", "t2", "t3"},
	}

	// 不存在的userID得出了系统漏洞
	// 对于真实存在的用户进行publish
	postId, err := fixture.serv.PublishArticle(req, fixture.userID)
	assert.Empty(t, err)
	assert.Equal(t, int(1), postId)

	// 针对id获取post
	postResp, err := fixture.serv.GetPost(uint(postId))
	assert.Nil(t, err)
	assert.Equal(t, req.Title, postResp.Post.Title)

	// 删除post
	err = fixture.serv.DeletePost(uint(postId))
	assert.Nil(t, err)

	// 删除文章后，查询个人发布文章
	postsResp, err := fixture.serv.GetPostsOfUser(1, 10, fixture.userID)
	assert.Nil(t, err)
	assert.Equal(t, uint(0), postsResp.Cnt)
	assert.Empty(t, postsResp.Posts)
}

func TestSaveDraft(t *testing.T) {
	fixture := setupTestFixture(t)
	defer teardownTestDB(fixture.db)

	req := &dtos.ArticleReq{
		Id:      0,
		Title:   "测试",
		Excerpt: "摘要abstract",
		Content: "正文",
		Cover:   "",

		// 附加项
		Category: "技术",
		Tags:     []string{"t1", "t2", "t3"},
	}

	// 创建draft，创建tags因为之前的post会在日志输出错误
	draftId, err := fixture.serv.SaveDraft(req, fixture.userID)
	assert.Equal(t, int(1), draftId)
	assert.Empty(t, err)

	// 针对id查询draft
	draftResp, err := fixture.serv.GetDraft(uint(draftId))
	assert.Nil(t, err)
	assert.Equal(t, req.Title, draftResp.Draft.Title)

	// 删除draft
	err = fixture.serv.DeleteDraft(uint(draftId))
	assert.Nil(t, err)

	// 删除草稿后查询个人草稿
	draftsResp, err := fixture.serv.GetDraftsOfUser(1, 10, fixture.userID)
	assert.Nil(t, err)
	assert.Equal(t, uint(0), draftsResp.Cnt)
	assert.Empty(t, draftsResp.Drafts)
}
