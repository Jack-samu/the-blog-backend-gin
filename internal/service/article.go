package service

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Jack-samu/the-blog-backend-gin.git/internal/dtos"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/errs"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/models"
	"gorm.io/gorm"
)

func (s *Service) GetPosts(page, perPage int64) (*dtos.PostListResp, *errs.ErrorResp) {

	posts, total, err := s.r.GetAllPosts(page, perPage)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewError(http.StatusNotFound, "你找的啥啊？", nil)
		}
		log.Printf("404以外的：%s\n", err.Error())
		return nil, errs.NewError(http.StatusInternalServerError, "", err)
	}

	// 转换为响应用的post列表格式
	postList := dtos.ToPostList(posts)

	// 进行实质的响应构造
	postListResp := dtos.NewPostList(postList, total, page)

	return postListResp, nil
}

func (s *Service) GetPost(id uint) (*dtos.PostDetailResp, *errs.ErrorResp) {

	post, err := s.r.GetPost(id, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewError(http.StatusNotFound, "你找的啥啊？", nil)
		}
		log.Printf("404以外的：%s\n", err.Error())
		return nil, errs.NewError(http.StatusInternalServerError, "", err)
	}

	postDetail := dtos.ToPostDetail(post)

	return &dtos.PostDetailResp{
		Post: postDetail,
	}, nil
}

func (s *Service) GetDraft(id uint) (*dtos.DraftDetailResp, *errs.ErrorResp) {

	draft, err := s.r.GetDraft(id, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewError(http.StatusNotFound, "你找的啥啊？", nil)
		}
		log.Printf("404以外的：%s\n", err.Error())
		return nil, errs.NewError(http.StatusInternalServerError, "", err)
	}

	draftDetail := dtos.ToDraftDetail(draft)

	return &dtos.DraftDetailResp{
		Draft: draftDetail,
	}, nil
}

func (s *Service) GetPostsOfUser(page, perPage int64, id string) (*dtos.PostListPersonalResp, *errs.ErrorResp) {

	posts, total, err := s.r.GetAllUsersPosts(page, perPage, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewError(http.StatusNotFound, "你找的啥啊？", nil)
		}
		log.Printf("404以外的：%s\n", err.Error())
		return nil, errs.NewError(http.StatusInternalServerError, "", err)
	}

	// 转换为响应用的post列表格式
	postList := dtos.ToPostList(posts)

	// 进行实质的响应构造
	postListResp := dtos.NewPostListPersonal(postList, total, page)

	return postListResp, nil
}

func (s *Service) GetDraftsOfUser(page, perPage int64, id string) (*dtos.DraftsResp, *errs.ErrorResp) {
	drafts, total, err := s.r.GetAllUsersDrafts(page, perPage, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewError(http.StatusNotFound, "你找的啥啊？", nil)
		}
		log.Printf("404以外的：%s\n", err.Error())
		return nil, errs.NewError(http.StatusInternalServerError, "", err)
	}

	// 转换为响应用的post列表格式
	draftList := dtos.ToDraftList(drafts)

	// 进行实质的响应构造
	draftsResp := dtos.NewDraftsPersonal(draftList, total, page)

	return draftsResp, nil
}

func (s *Service) GetSeries(id string) (*dtos.SeriesResp, *errs.ErrorResp) {

	categories, err := s.r.GetSeries(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewError(http.StatusNotFound, "你找的啥啊？", nil)
		}
		log.Printf("404以外的：%s\n", err.Error())
		return nil, errs.NewError(http.StatusInternalServerError, "", err)
	}

	series := make([]dtos.SeriesItem, len(categories))
	for i := range categories {
		series[i].ID = categories[i].ID
		series[i].Name = categories[i].Name
		// 构造基础post列表
		series[i].Posts = dtos.ToPostsBasic(categories[i].Posts)
	}

	return &dtos.SeriesResp{
		Categories: series,
	}, nil
}

func (s *Service) PublishArticle(req *dtos.ArticleReq, userID string) (int, *errs.ErrorResp) {

	post_id := -1

	if req == nil {
		return post_id, nil
	}

	// 参数的深层次校验，后续看看有么有需要

	err := s.r.Transaction(func(tx *gorm.DB) error {
		var post *models.Post
		var err error

		draft, err := s.r.GetDraft(req.Id, tx)
		if draft == nil && errors.Is(err, gorm.ErrRecordNotFound) {
			// post
			post = &models.Post{
				Title:     req.Title,
				Excerpt:   req.Excerpt,
				Content:   req.Content,
				Cover:     req.Cover,
				UserID:    userID,
				CreatedAt: time.Now(),
			}
			if err = tx.Create(post).Error; err != nil {
				log.Printf("发布过程中创建post实例出错：%s\n", err.Error())
				return err
			}
		} else if err != nil {
			log.Printf("查询对应草稿的其他错误：%s\n", err.Error())
			return err
		} else {
			// draft
			post, err = s.r.Draft2Post(draft, tx)
			if err != nil {
				log.Printf("draft转换post中出错：%s\n", err.Error())
				return err
			}
		}

		if req.Category != "" {
			category_n := strings.ToLower(req.Category)
			if err = s.r.UpdateCategory(tx, post, category_n); err != nil {
				log.Printf("更新category出错：%s\n", err.Error())
				return err
			}
		}

		if len(req.Tags) != 0 {
			toAdd, toRemove := s.compareTags(post.GetTagsName(), req.Tags)
			if err = s.r.UpdateTags(tx, post, toAdd, toRemove); err != nil {
				return err
			}
		}

		post_id = int(post.ID)
		return nil
	})

	if err != nil {
		return post_id, errs.NewError(http.StatusInternalServerError, "", err)
	}

	return post_id, nil
}

func (s *Service) SaveDraft(req *dtos.ArticleReq, userID string) (int, *errs.ErrorResp) {

	draft_id := -1

	if req == nil {
		return draft_id, nil
	}

	// 参数的深层次校验，后续看看有么有需要

	err := s.r.Transaction(func(tx *gorm.DB) error {
		draft, err := s.r.GetDraft(req.Id, tx)
		if draft == nil && errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建draft
			draft = &models.Draft{
				Title:     req.Title,
				Excerpt:   req.Excerpt,
				Content:   req.Content,
				Cover:     req.Cover,
				UserID:    userID,
				CreatedAt: time.Now(),
			}

			if err = tx.Create(draft).Error; err != nil {
				tx.Rollback()
				log.Printf("保存草稿过程中创建draft实例出错：%s\n", err.Error())
				return err
			}
		} else if err != nil {
			log.Printf("查询对应草稿的其他错误：%s\n", err.Error())
			return err
		} else {
			// draft 字段更新
			draft.Title = req.Title
			draft.Excerpt = req.Excerpt
			draft.Content = req.Content
			draft.Cover = req.Cover
		}

		// 更新category
		if req.Category != "" {
			category_n := strings.ToLower(req.Category)
			if err = s.r.UpdateCategory(tx, draft, category_n); err != nil {
				log.Printf("更新category出错：%s\n", err.Error())
				return err
			}
		}

		// 更新tags
		if len(req.Tags) != 0 {
			toAdd, toRemove := s.compareTags(draft.GetTagsName(), req.Tags)
			if err = s.r.UpdateTags(tx, draft, toAdd, toRemove); err != nil {
				return err
			}
		}

		draft_id = int(draft.ID)
		return nil
	})

	if err != nil {
		log.Printf("功亏一篑，提交出错：%s\n", err.Error())
		return draft_id, errs.NewError(http.StatusInternalServerError, "", err)
	}

	return draft_id, nil
}

func (s *Service) DeletePost(post_id uint) *errs.ErrorResp {

	err := s.r.Transaction(func(tx *gorm.DB) error {
		post, err := s.r.GetPost(post_id, tx)
		if err != nil {
			log.Printf("查询要删除的post出错：%s\n", err.Error())
			return err
		}

		err = s.r.DeleteArticle(post, tx)
		if err != nil {
			tx.Rollback()
			log.Printf("删除%s出错：%s\n", post.Title, err.Error())
		}
		return nil
	})

	if err != nil {
		return errs.NewError(http.StatusInternalServerError, "", err)
	}

	return nil
}

func (s *Service) DeleteDraft(draft_id uint) *errs.ErrorResp {

	err := s.r.Transaction(func(tx *gorm.DB) error {
		draft, err := s.r.GetDraft(draft_id, tx)
		if err != nil {
			log.Printf("查询要删除的draft出错：%s\n", err.Error())
			return err
		}

		err = s.r.DeleteArticle(draft, tx)
		if err != nil {
			tx.Rollback()
			log.Printf("删除%s出错：%s\n", draft.Title, err.Error())
		}
		return nil
	})

	if err != nil {
		return errs.NewError(http.StatusInternalServerError, "", err)
	}

	return nil
}

func (s *Service) compareTags(currentTags, newTags []string) (toAdd, toRemove []string) {
	currentNames := make(map[string]bool)
	newNames := make(map[string]bool)

	for _, tag := range currentTags {
		currentNames[tag] = true
	}

	for _, tag := range newTags {
		// 新增参数要小写化
		t := strings.ToLower(tag)
		newNames[t] = true
		if _, ok := currentNames[t]; !ok {
			toAdd = append(toAdd, t)
		}
	}

	for _, tag := range currentTags {
		if _, ok := newNames[tag]; !ok {
			toRemove = append(toRemove, tag)
		}
	}

	return
}
