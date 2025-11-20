package service

import (
	"errors"
	"log"
	"net/http"

	"github.com/Jack-samu/the-blog-backend-gin.git/internal/dtos"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/errs"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/models"
	"gorm.io/gorm"
)

// 只加载基础评论
func (s *Service) GetComments(id int64) (*dtos.CommentsResp, *errs.ErrorResp) {
	// 后续添加触底刷新
	if id == -1 {
		return nil, errs.NewError(http.StatusBadRequest, "参数错误", nil)
	}

	comments, err := s.r.GetComments(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("没找到对应文章，也就无法查询comments了")
			return nil, errs.NewError(http.StatusNotFound, "", nil)
		}
		log.Printf("具体错误：%s\n", err.Error())
		return nil, errs.NewError(http.StatusInternalServerError, "", err)
	}

	commentsResp := dtos.ToCommentsResp(comments, s.r)
	return commentsResp, nil
}

// 由前台点击触发二级评论加载
func (s *Service) GetReplies(id int64) (*dtos.RepliesResp, *errs.ErrorResp) {
	if id == -1 {
		return nil, errs.NewError(http.StatusBadRequest, "参数错误", nil)
	}

	replies, err := s.r.GetReplies(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("没找到对应comment，也就无法查询comments了")
			return nil, errs.NewError(http.StatusNotFound, "", nil)
		}
		log.Printf("具体错误：%s\n", err.Error())
		return nil, errs.NewError(http.StatusInternalServerError, "", err)
	}

	repliesResp := dtos.ToRepliesResp(replies, s.r)
	return repliesResp, nil
}

func (s *Service) CreateComment(req *dtos.CommentReq, userID string) (*dtos.CommentResp, *errs.ErrorResp) {
	if req.ArticleID == 0 {
		return nil, errs.NewError(http.StatusBadRequest, "无效参数", nil)
	}

	var comment *models.Comment
	var user *models.User

	err := s.r.Transaction(func(tx *gorm.DB) error {

		err := tx.Model(&models.User{}).Where("id = ?", userID).First(&user).Error
		if err != nil {
			log.Printf("%v\n", err.Error())
			return errors.New("用户鉴权失败")
		}

		var cnt int64
		err = tx.Model(&models.Post{}).Where("id = ?", req.ArticleID).Count(&cnt).Error
		if err != nil {
			log.Printf("要进行评论的文章404：%s\n", err.Error())
			return gorm.ErrRecordNotFound
		}

		comment = &models.Comment{
			Content: req.Content,
			UserID:  userID,
			PostID:  uint(req.ArticleID),
		}

		err = tx.Create(comment).Error

		return err
	})

	if err != nil {
		log.Printf("创建comment出错：%s\n", err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewError(http.StatusInternalServerError, "没找到要评论的主体文章", nil)
		}
		return nil, errs.NewError(http.StatusInternalServerError, "", err)
	}

	return &dtos.CommentResp{
		Comment: dtos.ToCommentItem(comment, user, nil),
	}, nil
}

func (s *Service) CreateReply(req *dtos.CommentReq, userID string) (*dtos.ReplyResp, *errs.ErrorResp) {
	if req.CommentID == 0 {
		return nil, errs.NewError(http.StatusBadRequest, "无效参数", nil)
	}

	var reply *models.Reply
	var user *models.User

	err := s.r.Transaction(func(tx *gorm.DB) error {

		err := tx.Model(&models.User{}).Where("id = ?", userID).First(&user).Error
		if err != nil {
			return errors.New("用户鉴权失败")
		}

		var cnt int64
		err = tx.Model(&models.Comment{}).Where("id = ?", req.CommentID).Count(&cnt).Error
		if err != nil {
			log.Printf("要进行评论的基础评论404：%s\n", err.Error())
			return gorm.ErrRecordNotFound
		}

		reply = &models.Reply{
			Content:   req.Content,
			UserID:    userID,
			CommentID: uint(req.CommentID),
		}

		if req.ParentID != 0 {
			reply.ParentID = &req.ParentID
		}

		err = tx.Create(reply).Error

		return err
	})

	if err != nil {
		log.Printf("创建reply出错：%s\n", err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewError(http.StatusInternalServerError, "没找到要评论的主体评论", nil)
		}
		return nil, errs.NewError(http.StatusInternalServerError, "", err)
	}

	return &dtos.ReplyResp{
		Reply: dtos.ToReplyItem(reply, user, nil),
	}, nil
}

func (s *Service) ModifyComment(req *dtos.CommentReq, userID string) (*dtos.CommentResp, *errs.ErrorResp) {
	if req.CommentID == 0 {
		return nil, errs.NewError(http.StatusBadRequest, "参数无效", nil)
	}

	if req.Content == "" {
		return nil, errs.NewError(http.StatusBadRequest, "评论主体内容不能为空", nil)
	}

	var comment *models.Comment

	err := s.r.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&models.Comment{}).Where("id = ?", req.CommentID).First(&comment).Error
		if err != nil {
			log.Printf("查询出错：%s\n", err.Error())
			return gorm.ErrRecordNotFound
		}

		if comment.UserID != userID {
			return errors.New("无权操作")
		}

		if comment.Content == req.Content {
			return errors.New("没有改动")
		}

		err = tx.Model(&comment).Update("content", req.Content).Error

		return err
	})

	if err != nil {
		// 详细错误鉴别并回应
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewError(http.StatusInternalServerError, "没找到要修改的评论", nil)
		}
		return nil, errs.NewError(http.StatusInternalServerError, "", err)
	}

	return &dtos.CommentResp{
		Comment: dtos.ToCommentItem(comment, nil, s.r),
	}, nil
}

func (s *Service) ModifyReply(req *dtos.CommentReq, userID string) (*dtos.ReplyResp, *errs.ErrorResp) {
	if req.ReplyID == 0 {
		return nil, errs.NewError(http.StatusBadRequest, "参数无效", nil)
	}

	if req.Content == "" {
		return nil, errs.NewError(http.StatusBadRequest, "评论主体内容不能为空", nil)
	}

	var reply *models.Reply

	err := s.r.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&models.Reply{}).Where("id = ?", req.ReplyID).First(&reply).Error
		if err != nil {
			log.Printf("查询出错：%s\n", err.Error())
			return gorm.ErrRecordNotFound
		}

		if reply.UserID != userID {
			return errors.New("无权操作")
		}

		if reply.Content == req.Content {
			return errors.New("没有改动")
		}

		err = tx.Model(&reply).Update("content", req.Content).Error

		return err
	})

	if err != nil {
		// 详细错误鉴别并回应
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewError(http.StatusInternalServerError, "没找到要修改的评论", nil)
		}
		return nil, errs.NewError(http.StatusInternalServerError, "", err)
	}

	return &dtos.ReplyResp{
		Reply: dtos.ToReplyItem(reply, nil, s.r),
	}, nil
}

func (s *Service) DeleteComment(id int64, userID string) *errs.ErrorResp {

	err := s.r.Transaction(func(tx *gorm.DB) error {
		var comment models.Comment
		err := tx.Model(&models.Comment{}).
			Preload("User").
			Where("id = ?", id).First(&comment).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}

		if comment.User.ID != userID {
			return errors.New("用户无权删除别人的评论")
		}

		err = tx.Delete(&comment).Error
		return err
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.NewError(http.StatusNotFound, "要删除的comment没找到", nil)
		}
		log.Printf("删除过程中报错：%s\n", err.Error())
		return errs.NewError(http.StatusInternalServerError, "", err)
	}

	return nil
}

func (s *Service) DeleteReply(id int64, userID string) *errs.ErrorResp {

	err := s.r.Transaction(func(tx *gorm.DB) error {
		var reply models.Reply
		err := tx.Model(&models.Reply{}).
			Preload("User").
			Where("id = ?", id).First(&reply).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}

		if reply.User.ID != userID {
			return errors.New("用户无权删除别人的评论")
		}

		err = tx.Delete(&reply).Error
		return err
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.NewError(http.StatusNotFound, "要删除的reply没找到", nil)
		}
		log.Printf("删除过程中报错：%s\n", err.Error())
		return errs.NewError(http.StatusInternalServerError, "", err)
	}

	return nil
}
