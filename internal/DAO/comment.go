package dao

import (
	"errors"
	"log"

	"github.com/Jack-samu/the-blog-backend-gin.git/internal/models"
	"gorm.io/gorm"
)

func (r DAO) GetComments(id int64) ([]*models.Comment, error) {
	var comments []*models.Comment

	var post models.Post
	if err := r.db.Take(&post, "id = ?", id).Error; err != nil {
		log.Println("文章主体不存在")
		// 评论的主体post都不存在
		return nil, gorm.ErrRecordNotFound
	}

	err := r.db.Model(&models.Comment{}).Where("post_id = ?", id).
		Preload("User").
		Preload("Replies").
		Find(&comments).
		Error

	return comments, err
}

func (r *DAO) GetReplies(id int64) ([]*models.Reply, error) {
	var replies []*models.Reply

	var comment models.Comment
	err := r.db.Take(&comment, "id = ?", id).Error
	if err != nil {
		// 二级评论的主体comment都不存在
		return nil, gorm.ErrRecordNotFound
	}

	err = r.db.Model(&models.Reply{}).
		Where("comment_id = ?", id).
		Preload("User").
		Preload("Replies", "parent_id IS NOT NULL").
		Find(&replies).Error

	return replies, err
}

func (r *DAO) IsLiked(userID, targetType string, targetID uint) (bool, error) {
	var record int64

	err := r.db.Model(&models.Like{}).
		Where("target_type = ?", targetType).
		Where("target_id = ?", targetID).
		Where("user_id = ?", userID).
		Count(&record).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}

	return true, nil
}
