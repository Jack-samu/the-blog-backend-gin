package handler

import (
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/service"
)

type Handler struct {
	s *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{s: s}
}
