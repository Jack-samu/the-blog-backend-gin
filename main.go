package main

import (
	"log"

	"github.com/Jack-samu/the-blog-backend-gin.git/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	r := gin.Default()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("读取.env失败")
	}

	_ = models.InitDB()
	r.Run()
}
