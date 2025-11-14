package dao

import (
	"gorm.io/gorm"
)

type DAO struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *DAO {
	return &DAO{db: db}
}
