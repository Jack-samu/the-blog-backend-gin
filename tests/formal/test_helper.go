package formal

import (
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/handler"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupRoutes(h *handler.Handler, r *gin.Engine) {
	r.POST("/auth/register", h.Register)
	r.POST("/auth/login", h.Login)
	// 基础列表
	r.GET("/articles", h.GetArticles)
	// 详情页面
	r.GET("/articles/:id", h.GetArticle)
	// 获取分类的所有文章
	r.GET("/articles/series/:id", h.GetSeries)

	protected := r.Group("")
	protected.Use(middleware.Auth())
	{
		protected.POST("/auth/:id/profile", h.Profile)
		protected.POST("/auth/id/photos", h.GetPhotos)
		protected.POST("/auth/refresh", h.RefreshTheToken)
		protected.POST("/auth/logout", h.Logout)

		// 获取草稿内容
		protected.GET("/articles/draft/:id", h.GetDraftEditable)
		// 获取post内容
		protected.GET("/articles/publish", h.GetPostsOfUser)
		// draft列表
		protected.GET("/articles/drafts", h.GetDraftOfUser)
		// 发布post
		protected.POST("/articles/publish", h.PublishArticle)
		// 保存draft
		protected.POST("/articles/save", h.SaveDraft)
		// 删除post
		protected.DELETE("/articles/post/:id", h.DeletePost)
		// 删除draft
		protected.DELETE("/articles/draft/:id", h.DeletePost)
	}
}

func teardownTestDB(db *gorm.DB) {
	db.Exec("DELETE FROM draft_tags")
	db.Exec("DELETE FROM post_tags")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM posts")
	db.Exec("DELETE FROM drafts")
	db.Exec("DELETE FROM imgs")
	db.Exec("DELETE FROM categories")
	db.Exec("DELETE FROM tags")
	db.Exec("DELETE FROM likes")
	// db.Exec("DELETE FROM alembic_version")

	sqlDB, _ := db.DB()
	sqlDB.Close()
}
