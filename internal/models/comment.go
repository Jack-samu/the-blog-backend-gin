package models

import "time"

type Comment struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Content   string    `gorm:"type:text"`
	LikeCnt   uint      `gorm:"default:0"`
	CreatedAt time.Time `gorm:"autoUpdateTime:false;index"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli"`

	// 外键外联
	PostID  *uint   `gorm:"index"`
	Replies []Reply `gorm:"foreignKey:CommentID;constraint:OnDelete:CASCADE"`

	UserID string `gorm:"type:varchar(36);"`
	User   *User  `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL"`
}

type Reply struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Content   string    `gorm:"type:text"`
	LikeCnt   uint      `gorm:"default:0"`
	CreatedAt time.Time `gorm:"autoUpdateTime:false;index"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli"`

	// 外键外联，PostID存疑
	PostID    *uint `gorm:"index"`
	CommentID *uint `gorm:"index"`

	UserID string `gorm:"type:varchar(36)"`
	User   *User  `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL"`

	// Reply自我关联
	ParentID *uint   `gorm:"index"`
	Replies  []Reply `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE"`
}

type Like struct {
	UserID     string    `gorm:"type:varchar(36);not null"`
	TargetType string    `gorm:"index"`
	TargetID   *uint     `gorm:"index"`
	CreatedAt  time.Time `gorm:"autoUpdateTime:false;index"`
}
