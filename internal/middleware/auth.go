package middleware

import (
	"net/http"
	"strings"

	"github.com/Jack-samu/the-blog-backend-gin.git/internal/utils"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"err": "token缺失"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"err": "token"})
			c.Abort()
			return
		}

		payload, err := utils.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"err": "token无效或者已过期"})
			c.Abort()
			return
		}

		c.Set("user_id", payload.ID)
		c.Next()
	}
}
