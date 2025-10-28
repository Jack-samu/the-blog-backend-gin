package main

import (
	"log"

	"github.com/Jack-samu/the-blog-backend-gin.git/internal/handler"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/middleware"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/models"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/repositories"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/service"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	r := gin.Default()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("读取.env失败")
	}

	// 邮件配置
	utils.InitEmailConfig()
	db := models.InitDB()

	repository := repositories.NewRepository(db)
	service := service.NewService(repository)
	handler := handler.NewHandler(service)

	// 路由注册
	r.POST("/upload-img", handler.UploadImage)
	auth := r.Group("/auth")
	{
		auth.POST("/register", handler.Register)
	}

	protected := r.Group("/")
	protected.Use(middleware.Auth())
	{
		protected.POST("logout", handler.Logout)
	}

	// http://localhost:8080/img1.png
	r.Static("", "./static/images")

	r.Run()
}
