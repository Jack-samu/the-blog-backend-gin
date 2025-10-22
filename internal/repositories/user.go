package repositories

import (
	"time"

	"github.com/Jack-samu/the-blog-backend-gin.git/internal/models"
)

func (r *Repository) ExistByEmail(email string) (bool, error) {
	var cnt int64
	err := r.db.Where("email = ?", email).Count(&cnt).Error
	return cnt > 0, err
}

func (r *Repository) ExistByUsername(username string) (bool, error) {
	var cnt int64
	err := r.db.Where("username = ?", username).Count(&cnt).Error
	return cnt > 0, err
}

func (r *Repository) CreateUser(u *models.User) error {
	return r.db.Create(u).Error
}

func (r *Repository) AddAvatar(u *models.User, filename string) error {
	img := &models.Img{
		Name:      filename,
		IsAvatar:  true,
		UserID:    u.ID,
		CreatedAt: time.Now(),
	}
	return r.db.Create(img).Error
}

func (r *Repository) SavePic(img *models.Img) error {
	return r.db.Create(img).Error
}

func (r *Repository) GetUserByNameWithAvatar(username string) (*models.User, string, error) {
	user := &models.User{}
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, "", err
	}

	var avatar string
	var img models.Img
	err = r.db.Model(&models.Img{}).
		Where("user_id = ? AND is_avatar = ?", user.ID, true).
		First(&img).Error
	if err == nil {
		avatar = img.Name
	}

	return user, avatar, err
}

func (r *Repository) IncreaseFailedLogin(u *models.User) error {
	u.FailedLogin = u.FailedLogin + 1
	return r.db.Save(u).Error
}

func (r *Repository) GetUserPosts(userID string) (int64, error) {
	var cnt int64
	err := r.db.Model(&models.Post{}).Where("user_id = ?", userID).Count(&cnt).Error
	return cnt, err
}

func (r *Repository) GetUserById(userID string) (*models.User, error) {
	var user *models.User
	err := r.db.Where("id = ?", userID).First(&user).Error
	return user, err
}

func (r *Repository) GetUserByName(username string) (*models.User, error) {
	var user *models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return user, err
}

func (r *Repository) SetLastActivity(u *models.User) (string, error) {
	u.LastActivity = time.Now()
	err := r.db.Save(u).Error
	return u.LastActivity.GoString(), err
}

func (r *Repository) GetUserByIdWithAvatar(id string) (*models.User, string, error) {
	user := &models.User{}
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, "", err
	}

	var avatar string
	var img models.Img
	err = r.db.Model(&models.Img{}).
		Where("user_id = ? AND is_avatar = ?", user.ID, true).
		First(&img).Error
	if err == nil {
		avatar = img.Name
	}

	return user, avatar, err
}

func (r *Repository) IncreaseCaptchaCnt(u *models.User) error {
	u.CaptchaReqCnt = u.CaptchaReqCnt + 1
	return r.db.Save(u).Error
}

func (r *Repository) GetUserDrafts(id string) (int64, error) {
	var drafts int64
	err := r.db.Model(&models.Draft{}).Where("user_id = ?", id).Count(&drafts).Error
	return drafts, err
}

func (r *Repository) GetUserPhotos(id string) ([]models.Img, error) {
	var imgs []models.Img
	err := r.db.Model(&models.Img{}).Where("user_id = ?", id).Find(&imgs).Error
	return imgs, err
}

func (r *Repository) GetPhoto(id uint) (*models.Img, error) {
	img := &models.Img{}
	err := r.db.Where("id = ?", id).First(&img).Error
	return img, err
}

func (r *Repository) SaveImgWithUser(userID, filename string) error {

	img := &models.Img{
		Name:   filename,
		UserID: userID,
	}
	err := r.db.Create(img).Error

	return err
}

func (r *Repository) DeleteImg(id uint) error {
	return r.db.Delete("id = ?", id).Error
}
