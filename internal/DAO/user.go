package dao

import (
	"time"

	"github.com/Jack-samu/the-blog-backend-gin.git/internal/models"
)

func (r *DAO) ExistByEmail(email string) (bool, error) {
	var cnt int64
	err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&cnt).Error
	return cnt > 0, err
}

func (r *DAO) ExistByUsername(username string) (bool, error) {
	var cnt int64
	err := r.db.Model(&models.User{}).Where("username = ?", username).Count(&cnt).Error
	return cnt > 0, err
}

func (r *DAO) CreateUser(u *models.User) error {
	return r.db.Create(u).Error
}

func (r *DAO) GetUserByName(username string) (*models.User, error) {
	user := &models.User{}
	err := r.db.Model(&models.User{}).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}

	return user, err
}

func (r *DAO) GetUserById(id string) (*models.User, error) {
	user := &models.User{}
	err := r.db.Model(&models.User{}).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}

	return user, err
}

func (r *DAO) IncreaseFailedLogin(u *models.User) error {
	// u.FailedLogin = u.FailedLogin + 1
	// return r.db.Save(u).Error
	return r.db.Model(&u).Update("failed_login", u.FailedLogin+1).Error
}

func (r *DAO) IncreaseCaptchaCnt(u *models.User) error {
	// u.CaptchaReqCnt = u.CaptchaReqCnt + 1
	// return r.db.Save(u).Error
	return r.db.Model(&u).Update("captcha_req_cnt", u.CaptchaReqCnt+1).Error
}

func (r *DAO) GetUserPosts(userID string) (int64, error) {
	var cnt int64
	err := r.db.Model(&models.Post{}).Where("user_id = ?", userID).Count(&cnt).Error
	return cnt, err
}

func (r *DAO) GetUserDrafts(id string) (int64, error) {
	var drafts int64
	err := r.db.Model(&models.Draft{}).Where("user_id = ?", id).Count(&drafts).Error
	return drafts, err
}

func (r *DAO) GetUserPhotos(id string) ([]*models.Img, error) {
	var imgs []*models.Img
	err := r.db.Model(&models.Img{}).Where("user_id = ?", id).Find(&imgs).Error
	return imgs, err
}

func (r *DAO) SetLastActivity(u *models.User) (string, error) {
	last_activity := time.Now()
	err := r.db.Model(&u).Update("last_activity", last_activity).Error
	return last_activity.GoString(), err
}

func (r *DAO) GetPhoto(id uint) (*models.Img, error) {
	img := &models.Img{}
	err := r.db.Model(&models.Img{}).Where("id = ?", id).First(&img).Error
	return img, err
}

func (r *DAO) DeleteImg(id uint) error {
	return r.db.Where("id = ?", id).Delete(&models.Img{}).Error
}

func (r *DAO) SaveImg(filename, user_id string, is_avatar bool) error {
	img := &models.Img{
		Name:      filename,
		UserID:    user_id,
		IsAvatar:  is_avatar,
		CreatedAt: time.Now(),
	}

	return r.db.Create(img).Error
}
