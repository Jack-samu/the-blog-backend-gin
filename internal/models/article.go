package models

import "time"

type PostTags struct {
	PostID uint `gorm:"primaryKey"`
	TagID  uint `gorm:"primaryKey"`
}

type DraftTags struct {
	DraftID uint `gorm:"primaryKey"`
	TagID   uint `gorm:"primaryKey"`
}

type Tag struct {
	ID     uint    `gorm:"primaryKey;autoIncrement"`
	Name   string  `gorm:"size:20;uniqueIndex;not null"`
	Posts  []Post  `gorm:"many2many:post_tags"`
	Drafts []Draft `gorm:"many2many:draft_tags"`
}

type Category struct {
	ID     uint   `gorm:"primaryKey;autoIncrement"`
	Name   string `gorm:"size:20;uniqueIndex;not null"`
	UserID string `gorm:"type:varchar(36);not null"`
	Posts  []Post `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE"`
	Drafts []Post `gorm:"foreignKey:CategoryID;constraint:OnDelete:SET NULL"`
}

type Post struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Title     string    `gorm:"size:100;index"`
	Excerpt   string    `gorm:"size:200"`
	Content   string    `gorm:"type:text"`
	ViewsCnt  int       `gorm:"not null;default:0"`
	LikeCnt   int       `gorm:"default:0"`
	Cover     *string   `gorm:"size:100"`
	CreatedAt time.Time `gorm:"autoUpdateTime:false;index"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli"`

	// 外键外联
	UserID     string    `gorm:"type:varchar(36);not null"`
	CategoryID *uint     `gorm:"index"`
	Category   *Category `gorm:"foreignKey:CategoryID"`
	Draft      *Draft    `gorm:"foreignKey:PostID"`
	Tags       []Tag     `gorm:"many2many:post_tags;"`

	Comments []Comment `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
	Replies  []Reply   `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
}

type Draft struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	Title      string    `gorm:"size:100;index"`
	Excerpt    string    `gorm:"size:200"`
	Content    string    `gorm:"type:text"`
	Cover      *string   `gorm:"size:100"`
	CreatedAt  time.Time `gorm:"autoUpdateTime:false;index"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime:milli"`
	UserID     string    `gorm:"type:varchar(36);not null"`
	CategoryID *uint     `gorm:"index"`
	Category   *Category `gorm:"foreignKey:CategoryID"`
	PostID     *uint     `gorm:"uniqueIndex"`
	Post       *Post     `gorm:"foreignKey:PostID;constraint:OnDelete:SET NULL"`
	Tags       []Tag     `gorm:"many2many:draft_tags;"`
}

type Img struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"size:100;uniqueIndex;not null"`
	IsAvatar  bool   `gorm:"default:false"`
	CreatedAt time.Time
	UserID    string `gorm:"type:varchar(36);not null"`
}
