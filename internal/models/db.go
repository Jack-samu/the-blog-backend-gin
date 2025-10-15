package models

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	)

	log.Printf("数据库uri：%s\n", dsn)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("数据库句柄创建失败，%v\n", err)
	}

	// 数据库迁移
	err = db.AutoMigrate(
		&User{},
		&Tag{},
		&Category{},
		&Post{},
		&Draft{},
		&PostTags{},
		&DraftTags{},
		&Img{},
		&Comment{},
		&Reply{},
		&Like{},
	)
	if err != nil {
		log.Fatalf("数据库迁移失败：%v\n", err)
	}

	return db
}
