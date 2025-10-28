package formal

import (
	"testing"

	"gorm.io/gorm"
)

func teardownTestDB(db *gorm.DB) {
	sqlDB, _ := db.DB()
	sqlDB.Close()

	// 后面要添加数据库表清理操作
}

func TestUserFlow(t *testing.T) {}
