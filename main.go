package main

import (
	"log"
	"net/http"

	dao "github.com/Jack-samu/the-blog-backend-gin.git/internal/DAO"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/handler"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/middleware"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/models"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/service"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func customRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// stack := debug.Stack()

				log.Printf("程序异常崩溃：%v\n", r)
				log.Printf("request：%s %s\n", c.Request.Method, c.Request.URL.Path)
				if id := c.GetString("X-Request-ID"); id != "" {
					log.Printf("Request ID:%s\n", id)
				}
				// log.Printf("堆栈：\n%s\n", string(stack))

				c.JSON(http.StatusInternalServerError, gin.H{"err": "服务器错误"})

				c.Abort()
			}
		}()

		c.Next()
	}
}

func main() {
	r := gin.Default()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("读取.env失败")
	}

	// 异常机制
	r.Use(customRecovery())

	// 邮件配置
	utils.InitEmailConfig()
	db := models.InitDB()

	repository := dao.NewRepository(db)
	service := service.NewService(repository)
	handler := handler.NewHandler(service)

	// 路由注册
	r.POST("/upload-img", handler.UploadImage)
	r.GET("/articles/:id/comments", handler.GetComments)
	r.GET("/articles/:id/replies", handler.GetReplies)
	auth := r.Group("/auth")
	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
		auth.GET("/getcode", handler.GetCaptcha)
		auth.POST("/verify", handler.VerifyCaptcha)
		auth.Any("/reset/:token", handler.ResetPassword)
		auth.POST("/refresh", handler.RefreshTheToken)
	}

	article := r.Group("/articles")
	{
		article.GET("", handler.GetArticles)
		article.GET("/:id", handler.GetArticle)
	}

	protected := r.Group("")
	protected.Use(middleware.Auth())
	{
		protected.GET("/auth/:id/profile", handler.Profile)
		protected.GET("/auth/id/photos", handler.GetPhotos)
		protected.POST("/auth/set-avatar", handler.SetAvatar)
		protected.POST("/auth/upload-img", handler.UploadImg)
		protected.POST("logout", handler.Logout)

		// article部分
		protected.GET("/articles/draft/:id", handler.GetDraftEditable)
		protected.GET("/articles/publish", handler.GetPostsOfUser)
		protected.GET("/articles/drafts", handler.GetDraftOfUser)
		protected.POST("/articles/publish", handler.PublishArticle)
		protected.POST("/articles/save", handler.SaveDraft)
		protected.DELETE("/articles/post/:id", handler.DeletePost)
		protected.DELETE("/articles/draft/:id", handler.DeletePost)

		// comment部分
		protected.POST("/articles/comments", handler.CreateComment)
		protected.POST("/articles/replies", handler.CreateReply)
		protected.POST("/comments/modify", handler.ModifyComment)
		protected.POST("/replies/modify", handler.ModifyReply)
		protected.DELETE("/comments/:id", handler.DeleteComment)
		protected.DELETE("/replies/:id", handler.DeleteReply)
	}

	// http://localhost:8080/img1.png
	r.Static("", "./static/images")

	r.Run()
}
