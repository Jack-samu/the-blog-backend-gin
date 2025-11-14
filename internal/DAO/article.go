package dao

import (
	"errors"
	"log"
	"time"

	"github.com/Jack-samu/the-blog-backend-gin.git/internal/models"
	"gorm.io/gorm"
)

func (r *DAO) GetAllPosts(page, perPage int64) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	// 暂时不添加条件查询，后续功能完善可以添加
	err := r.db.Model(&models.Post{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage

	// 通过偏移量确定开始，每页数量确定获取的数据数量
	err = r.db.Model(&models.Post{}).
		Preload("Author").
		Preload("Category", "id IS NOT NULL").
		Preload("Tags", "id IS NOT NULL").
		Order("created_at DESC").
		Offset(int(offset)).
		Limit(int(perPage)).
		Find(&posts).Error

	return posts, total, err
}

func (r *DAO) GetPost(id uint, tx *gorm.DB) (*models.Post, error) {
	if tx == nil {
		tx = r.db
	}

	var post models.Post
	err := tx.Model(&models.Post{}).
		Preload("Author").
		Preload("Category", "id IS NOT NULL").
		Preload("Tags", "id IS NOT NULL").
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Preload("User").Order("created_at DESC") // 评论及评论用户
		}).
		First(&post, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return &post, err
}

func (r *DAO) GetDraft(id uint, tx *gorm.DB) (*models.Draft, error) {
	if tx == nil {
		tx = r.db
	}

	var draft models.Draft

	err := tx.Model(&models.Draft{}).
		Preload("Author").
		Preload("Category", "id IS NOT NULL").
		Preload("Tags", "id IS NOT NULL").
		First(&draft, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return &draft, err
}

func (r *DAO) GetAllUsersPosts(page, perPage int64, id string) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	// 暂时不添加条件查询，后续功能完善可以添加
	err := r.db.Model(&models.Post{}).
		Where("user_id = ?", id).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage

	// 通过偏移量确定开始，每页数量确定获取的数据数量
	err = r.db.Model(&models.Post{}).
		Preload("Author").
		Preload("Category", "id IS NOT NULL").
		Preload("Tags", "id IS NOT NULL").
		Where("user_id = ?", id).
		Order("created_at DESC").
		Offset(int(offset)).
		Limit(int(perPage)).
		Find(&posts).Error

	return posts, total, err
}

func (r *DAO) GetAllUsersDrafts(page, perPage int64, id string) ([]models.Draft, int64, error) {
	var drafts []models.Draft
	var total int64

	// 暂时不添加条件查询，后续功能完善可以添加
	err := r.db.Model(&models.Draft{}).
		Where("user_id = ?", id).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage

	// 通过偏移量确定开始，每页数量确定获取的数据数量
	err = r.db.Model(&models.Draft{}).
		Preload("Author").
		Preload("Category", "id IS NOT NULL").
		Preload("Tags", "id IS NOT NULL").
		Where("user_id = ?", id).
		Order("created_at DESC").
		Offset(int(offset)).
		Limit(int(perPage)).
		Find(&drafts).Error

	return drafts, total, err
}

func (r *DAO) GetSeries(id string) ([]models.Category, error) {
	var categories []models.Category
	var user models.User

	err := r.db.Model(&models.User{}).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}

	err = r.db.Model(&models.Category{}).
		Preload("Posts").
		Where("user_id = ?", id).
		Find(&categories).Error

	return categories, err
}

// 创建post数据库记录并返回操作对象
func (r *DAO) CreatePost(title, excerpt, content, cover, user_id string, tx *gorm.DB) (*models.Post, error) {
	if tx == nil {
		tx = r.db
	}

	post := &models.Post{
		Title:     title,
		Excerpt:   excerpt,
		Content:   content,
		Cover:     cover,
		UserID:    user_id,
		CreatedAt: time.Now(),
	}

	err := tx.Create(&post).Error

	return post, err
}

func (r *DAO) CreateDraft(title, excerpt, content, cover, user_id string, tx *gorm.DB) (*models.Draft, error) {
	if tx == nil {
		tx = r.db
	}

	draft := &models.Draft{
		Title:     title,
		Excerpt:   excerpt,
		Content:   content,
		Cover:     cover,
		UserID:    user_id,
		CreatedAt: time.Now(),
	}

	err := tx.Create(&draft).Error

	return draft, err
}

// 针对Post和Draft的category的更新
func (r *DAO) UpdateCategory(tx *gorm.DB, article models.Article, c string) error {
	var category models.Category

	err := tx.Where("name = ? AND user_id = ?", c, article.GetUserID()).
		First(&category).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		category.Name = c
		category.UserID = article.GetUserID()

		if err = tx.Create(&category).Error; err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	switch a := article.(type) {
	case *models.Post:
		if err = tx.Model(&category).Association("Posts").Append(a); err != nil {
			return err
		}
	case *models.Draft:
		if err = tx.Model(&category).Association("Drafts").Append(a); err != nil {
			return err
		}
	}

	return nil
}

// 针对draft和post的tags更新
func (r *DAO) UpdateTags(tx *gorm.DB, article models.Article, toAdd, toRemove []string) error {
	// 删旧
	for _, tag := range toRemove {
		var t models.Tag
		if err := tx.Model(&models.Tag{}).Where("name = ?", tag).First(&t).Error; err != nil {
			log.Printf("查询旧标签%s出错：%s\n", tag, err.Error())
			continue
		} else {
			switch a := article.(type) {
			case (*models.Post):
				if err = tx.Model(a).Association("Tags").Delete(&t); err != nil {
					log.Printf("取消旧标签%s关联出错：%s\n", tag, err.Error())
					continue
				}
			case (*models.Draft):
				if err = tx.Model(a).Association("Tags").Delete(&t); err != nil {
					log.Printf("取消旧标签%s关联出错：%s\n", tag, err.Error())
					continue
				}
			}
		}
	}

	// 添新
	for _, t := range toAdd {
		tag := models.Tag{
			Name: t,
		}
		if err := tx.Create(&tag).Error; err != nil {
			log.Printf("创建新标签%s出错：%s\n", t, err.Error())
			continue
		}
		switch a := article.(type) {
		case (*models.Post):
			if err := tx.Model(a).Association("Tags").Append(&tag); err != nil {
				log.Printf("添加新标签%s关联出错：%s\n", t, err.Error())
				continue
			}
		case (*models.Draft):
			if err := tx.Model(a).Association("Tags").Append(&tag); err != nil {
				log.Printf("添加新标签%s关联出错：%s\n", t, err.Error())
				continue
			}
		}
	}

	return nil
}

func (r *DAO) Draft2Post(draft *models.Draft, tx *gorm.DB) (*models.Post, error) {

	post, err := r.CreatePost(draft.Title, draft.Excerpt, draft.Content, draft.Cover, draft.UserID, tx)
	if err != nil {
		return nil, err
	} else {
		err = tx.Delete(draft).Error
		return post, err
	}
}

func (r *DAO) DeleteArticle(article models.Article, tx *gorm.DB) error {
	if tx == nil {
		tx = r.db
	}

	switch a := article.(type) {
	case (*models.Post):
		if err := tx.Model(&models.Post{}).Delete(a).Error; err != nil {
			return err
		}
	case (*models.Draft):
		if err := tx.Model(&models.Draft{}).Delete(a).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *DAO) Transaction(f func(*gorm.DB) error) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		log.Printf("事务开启失败：%s\n", tx.Error.Error())
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("事务执行出现panic，已回滚：%v\n", r)
			// 后续可以添加监控
		}
	}()

	if err := f(tx); err != nil {
		tx.Rollback()
		log.Printf("事务执行失败，已回滚：%s\n", err.Error())
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Printf("功亏一篑，提交出错：%s\n", err.Error())
		return err
	}

	return nil
}
