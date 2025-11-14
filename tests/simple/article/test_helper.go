package article

import (
	"log"
	"os"
	"testing"

	"github.com/Jack-samu/the-blog-backend-gin.git/internal/handler"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/middleware"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/models"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupRoutes(h *handler.Handler, r *gin.Engine) {
	r.POST("/auth/login", h.Login)
	// 基础列表
	r.GET("/articles", h.GetArticles)
	// 详情页面
	r.GET("/articles/:id", h.GetArticle)
	// 获取分类的所有文章
	r.GET("/articles/series/:id", h.GetSeries)

	protected := r.Group("/")
	protected.Use(middleware.Auth())
	{
		// 编辑草稿
		protected.GET("/articles/draft/:id", h.GetDraftEditable)
		// 编辑post
		protected.GET("/articles/publish", h.GetPostsOfUser)
		// draft列表
		protected.GET("/articles/drafts", h.GetDraftOfUser)
		// post列表
		protected.POST("/articles/publish", h.PublishArticle)
		// 保存draft
		protected.POST("/articles/save", h.SaveDraft)
		// 删除post
		protected.DELETE("/articles/post/:id", h.DeletePost)
		// 删除draft
		protected.DELETE("/articles/draft/:id", h.DeletePost)
	}
}

func setupTestDB(t *testing.T) *gorm.DB {

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		t.Fatalf("创建测试用sqlite数据库失败：%s\n", err.Error())
	}
	// 数据库迁移
	err = db.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Draft{},
		&models.Img{},
		&models.Comment{},
		&models.Reply{},
	)
	if err != nil {
		t.Fatalf("数据库迁移失败：%s\n", err.Error())
	}

	return db
}

func teardownTestDB(db *gorm.DB) {
	sqlDB, _ := db.DB()
	sqlDB.Close()

	// 后置删除
	err := os.Remove("test.db")
	if err != nil {
		log.Fatalf("后置清除动作失败：%s\n", err.Error())
	}
}

func createTestUser(s *service.Service, t *testing.T) string {
	err := s.Register("test-user", "test@test.com", "test123", "guest what", "")
	assert.Empty(t, err)
	resp, err := s.Login("test-user", "test123")
	assert.Empty(t, err)
	return resp.UserInfo.ID
}
