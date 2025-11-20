package auth

import (
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	dao "github.com/Jack-samu/the-blog-backend-gin.git/internal/DAO"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func createTestUser() *models.User {
	user := &models.User{
		ID:            uuid.NewString(),
		Username:      "test-user",
		Email:         os.Getenv("TEST_MAIN"),
		Bio:           "say something",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		LastActivity:  time.Now(),
		FailedLogin:   0,
		CaptchaReqCnt: 0,
	}
	user.SetPassword("Guess123")
	return user
}

func createTestImg(is_avatar bool) *models.Img {
	img := &models.Img{
		Name:      uuid.NewString(),
		IsAvatar:  is_avatar,
		UserID:    uuid.NewString(),
		CreatedAt: time.Now(),
	}

	return img
}

// ExistByEmail方法的record not found版测试
func TestUserExistByEmail(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	defer verifyMockExpection(t, mock)

	repo := dao.NewRepository(db)
	user := createTestUser()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `users` WHERE email = ?")).
		WithArgs(user.Email).
		WillReturnRows(rows)

	exist, err := repo.ExistByEmail(user.Email)
	assert.NoError(t, err)
	assert.False(t, exist)
}

// ExistByUsername方法的record not found版测试
func TestUserExistByName(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	defer verifyMockExpection(t, mock)

	repo := dao.NewRepository(db)
	user := createTestUser()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `users` WHERE username = ?")).
		WithArgs(user.Username).
		WillReturnRows(rows)

	exist, err := repo.ExistByUsername(user.Username)
	assert.NoError(t, err)
	assert.False(t, exist)
}

// GetUserPosts的测试
func TestGetUserPosts(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	defer verifyMockExpection(t, mock)

	repo := dao.NewRepository(db)
	user := createTestUser()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `posts` WHERE user_id = ?")).
		WithArgs(user.ID).
		WillReturnRows(rows)

	cnt, err := repo.GetUserPosts(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), cnt)
}

// GetUserDrafts的测试
func TestGetUserDrafts(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	defer verifyMockExpection(t, mock)

	repo := dao.NewRepository(db)
	user := createTestUser()

	// 因为预设了结果为0
	rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `drafts` WHERE user_id = ?")).
		WithArgs(user.ID).
		WillReturnRows(rows)

	cnt, err := repo.GetUserDrafts(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), cnt)
}

// IncreaseCaptchaCnt的测试，改
func TestIncreaseCaptchaCnt(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	defer verifyMockExpection(t, mock)

	repo := dao.NewRepository(db)
	user := createTestUser()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `captcha_req_cnt`=?,`updated_at`=? WHERE `id` = ?")).
		WithArgs(1, sqlmock.AnyArg(), user.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.IncreaseCaptchaCnt(user)
	assert.NoError(t, err)
	assert.Equal(t, 1, user.CaptchaReqCnt)
}

// IncreaseFailedLogin的测试，改
func TestIncreaseFailedLogin(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	defer verifyMockExpection(t, mock)

	repo := dao.NewRepository(db)
	user := createTestUser()

	mock.ExpectBegin()
	// 不会连空格都要一致吧
	// 还真tm的空格差
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `failed_login`=?,`updated_at`=? WHERE `id` = ?")).
		WithArgs(1, sqlmock.AnyArg(), user.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.IncreaseFailedLogin(user)
	assert.NoError(t, err)
	assert.Equal(t, 1, user.FailedLogin)
}

// GetUserPhotos的测试
func TestGetPhotos(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	defer verifyMockExpection(t, mock)

	repo := dao.NewRepository(db)
	user := createTestUser()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `imgs` WHERE user_id = ?")).
		WithArgs(user.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "user_id", "is_avatar", "created_at"}))

	imgs, err := repo.GetUserPhotos(user.ID)
	assert.NoError(t, err)
	assert.Empty(t, imgs)
}

// CreateUser的测试，增
func TestCreateUser(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	defer verifyMockExpection(t, mock)

	repo := dao.NewRepository(db)
	user := createTestUser()

	mock.ExpectBegin()
	// 这他妈居然还要在时间精度上进行匹配，真的恶
	mock.ExpectExec(regexp.QuoteMeta(
		"INSERT INTO `users` (`id`,`username`,`email`,`pwd`,`bio`,`created_at`,`updated_at`,`last_activity`,`failed_login`,`captcha_req_cnt`) "+"VALUES (?,?,?,?,?,?,?,?,?,?)")).
		WithArgs(
			user.ID, user.Username, user.Email, user.Pwd, user.Bio, sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), user.FailedLogin, user.CaptchaReqCnt,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.CreateUser(user)
	assert.NoError(t, err)
}

// GetUserById的测试
func TestGetUserById(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	defer verifyMockExpection(t, mock)

	repo := dao.NewRepository(db)
	user := createTestUser()

	rows := sqlmock.NewRows([]string{
		`id`, `username`, `email`, `pwd`, `bio`,
		"created_at", "updated_at", "last_activity",
		"failed_login", "captcha_req_cnt",
	}).AddRow(
		user.ID, user.Username, user.Email, user.Pwd, user.Bio,
		user.CreatedAt, user.UpdatedAt, user.LastActivity,
		user.FailedLogin, user.CaptchaReqCnt,
	)

	// 用户查询
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE id = ? ORDER BY `users`.`id` LIMIT ?")).
		WithArgs(user.ID, 1).
		WillReturnRows(rows)

	u, err := repo.GetUserById(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.Username, u.Username)
}

// GetUserByName的测试
func TestGetUserByName(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	defer verifyMockExpection(t, mock)

	repo := dao.NewRepository(db)
	user := createTestUser()

	rows := sqlmock.NewRows([]string{
		`id`, `username`, `email`, `pwd`, `bio`,
		"created_at", "updated_at", "last_activity",
		"failed_login", "captcha_req_cnt",
	}).AddRow(
		user.ID, user.Username, user.Email, user.Pwd, user.Bio,
		user.CreatedAt, user.UpdatedAt, user.LastActivity,
		user.FailedLogin, user.CaptchaReqCnt,
	)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE username = ? ORDER BY `users`.`id` LIMIT ?")).
		WithArgs(user.Username, 1).
		WillReturnRows(rows)

	u, err := repo.GetUserByName(user.Username)
	assert.NoError(t, err)
	assert.Equal(t, user.Username, u.Username)
}

// SetLastActivity的测试，改
func TestSetLastAvtivity(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	defer verifyMockExpection(t, mock)

	repo := dao.NewRepository(db)
	user := createTestUser()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `last_activity`=?,`updated_at`=? WHERE `id` = ?")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), user.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	_, err := repo.SetLastActivity(user)
	assert.NoError(t, err)
}

// SaveImg的普通测试，增
func TestCreateImg(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	defer verifyMockExpection(t, mock)

	repo := dao.NewRepository(db)
	img := createTestImg(false)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		"INSERT INTO `imgs` (`name`,`is_avatar`,`created_at`,`user_id`) VALUES (?,?,?,?)")).
		WithArgs(img.Name, img.IsAvatar, sqlmock.AnyArg(), img.UserID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.SaveImg(img.Name, img.UserID, false)
	assert.NoError(t, err)
}

// SaveImg的头像添加测试，增
func TestCreateAvatar(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	defer verifyMockExpection(t, mock)

	repo := dao.NewRepository(db)
	img := createTestImg(true)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		"INSERT INTO `imgs` (`name`,`is_avatar`,`created_at`,`user_id`) VALUES (?,?,?,?)")).
		WithArgs(img.Name, img.IsAvatar, sqlmock.AnyArg(), img.UserID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.SaveImg(img.Name, img.UserID, true)
	assert.NoError(t, err)
}

// GetPhoto的测试
func TestGetImg(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	defer verifyMockExpection(t, mock)

	repo := dao.NewRepository(db)
	img := createTestImg(false)

	rows := sqlmock.NewRows([]string{`id`, `name`, `is_avatar`, `user_id`, "created_at"}).
		AddRow(img.ID, img.Name, img.IsAvatar, img.UserID, img.CreatedAt)

	// 用户查询
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `imgs` WHERE id = ? ORDER BY `imgs`.`id` LIMIT ?")).
		WithArgs(img.ID, 1).
		WillReturnRows(rows)

	img, err := repo.GetPhoto(img.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, img)
}

// DeleteImg的测试，删
func TestDeleteImg(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	defer verifyMockExpection(t, mock)

	repo := dao.NewRepository(db)
	img := createTestImg(false)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `imgs` WHERE id = ?")).
		WithArgs(img.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.DeleteImg(img.ID)
	assert.NoError(t, err)
}
