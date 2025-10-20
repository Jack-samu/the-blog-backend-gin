package service

import (
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/repositories"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/utils"
)

type Service struct {
	r *repositories.Repository
	u *utils.UserReset
}

func NewService(r *repositories.Repository) *Service {
	return &Service{r: r}
}
