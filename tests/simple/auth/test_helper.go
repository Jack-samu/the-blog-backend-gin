package auth

import (
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/models"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		t.Fatalf("创建测试用sqlite数据库失败：%s\n", err.Error())
	}
	// 数据库迁移
	err = db.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Draft{},
		&models.Img{},
	)
	if err != nil {
		t.Fatalf("数据库迁移失败：%s\n", err.Error())
	}

	return db
}

func teardownTestDB(db *gorm.DB, t *testing.T) {
	sqlDB, _ := db.DB()
	sqlDB.Close()

	// 后置删除
	err := os.Remove("test.db")
	if err != nil {
		t.Fatalf("后置清除动作失败：%s\n", err.Error())
	}
}

// 辅助工具部分
func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	// 添加环境变量读取
	err := godotenv.Load("../../.env")
	assert.NoError(t, err)

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	cleanup := func() {
		defer db.Close()
	}

	return gormDB, mock, cleanup
}

func verifyMockExpection(t *testing.T, mock sqlmock.Sqlmock) {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("期望的mock未满足：%s\n", err)
	}
}
