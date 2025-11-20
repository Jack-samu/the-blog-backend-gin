package comment

import (
	"log"
	"os"
	"testing"

	"github.com/Jack-samu/the-blog-backend-gin.git/internal/dtos"
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
	r.GET("/articles/:id/comments", h.GetComments)
	r.GET("/articles/:id/replies", h.GetReplies)

	protected := r.Group("/")
	protected.Use(middleware.Auth())
	{
		// publish文章
		protected.POST("/articles/publish", h.PublishArticle)

		// comment部分
		protected.POST("/articles/comments", h.CreateComment)
		protected.POST("/articles/replies", h.CreateReply)
		protected.POST("/comments/modify", h.ModifyComment)
		protected.POST("/replies/modify", h.ModifyReply)
		protected.DELETE("/comments/:id", h.DeleteComment)
		protected.DELETE("/replies/:id", h.DeleteReply)
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
		&models.Like{},
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
