package service

import (
	dao "github.com/Jack-samu/the-blog-backend-gin.git/internal/DAO"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/utils"
)

type Service struct {
	r *dao.DAO
	u *utils.UserReset
}

func NewService(r *dao.DAO) *Service {
	return &Service{r: r}
}
