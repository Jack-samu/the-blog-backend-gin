package models

import (
	"time"
)

type Tag struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"type:varchar(36);uniqueIndex;not null"`
	CreatedAt time.Time `gorm:"autoUpdateTime:false;index"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli"`

	Posts  []Post  `gorm:"many2many:post_tags"`
	Drafts []Draft `gorm:"many2many:draft_tags"`
}

type Category struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"size:20;uniqueIndex;not null"`
	UserID    string    `gorm:"type:varchar(36);not null"`
	CreatedAt time.Time `gorm:"autoUpdateTime:false;index"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli"`

	User   *User   `gorm:"foreignKey:UserID;references:ID"`
	Posts  []Post  `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE"`
	Drafts []Draft `gorm:"foreignKey:CategoryID;constraint:OnDelete:SET NULL"`
}

type Post struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Title     string `gorm:"size:100;index;not null"`
	Excerpt   string `gorm:"size:200;not null"`
	Content   string `gorm:"type:longtext"`
	ViewsCnt  int    `gorm:"not null;default:0"`
	LikeCnt   int    `gorm:"not null;default:0"`
	Cover     string
	CreatedAt time.Time `gorm:"autoUpdateTime:false;index"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli"`

	// 外键外联
	UserID     string    `gorm:"type:varchar(36);not null;index"`
	Author     User      `gorm:"foreignKey:UserID"`
	CategoryID *uint     `gorm:"index"`
	Category   *Category `gorm:"foreignKey:CategoryID"`
	Draft      *Draft    `gorm:"foreignKey:PostID"`
	Tags       []Tag     `gorm:"many2many:post_tags;"`

	Comments []Comment `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
}

type Draft struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Title     string    `gorm:"size:100;index"`
	Excerpt   string    `gorm:"size:200"`
	Content   string    `gorm:"type:longtext"`
	CreatedAt time.Time `gorm:"autoUpdateTime:false;index"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli"`
	Cover     string

	UserID string `gorm:"type:varchar(36);not null"`
	Author User   `gorm:"foreignKey:UserID"`

	CategoryID *uint     `gorm:"index"`
	Category   *Category `gorm:"foreignKey:CategoryID"`
	Tags       []Tag     `gorm:"many2many:draft_tags"`
	PostID     *uint     `gorm:"uniqueIndex"`
}

type Img struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"size:100;uniqueIndex;not null"`
	IsAvatar  bool   `gorm:"default:false"`
	CreatedAt time.Time
	UserID    string `gorm:"type:varchar(36);not null"`
}

// 通用接口
type Article interface {
	GetUserID() string
	GetTagsName() []string
	// 后续有需要再添加
}

func (p *Post) GetUserID() string {
	return p.UserID
}

func (d *Draft) GetUserID() string {
	return d.UserID
}

func (p *Post) GetTagsName() []string {
	tags := make([]string, len(p.Tags))

	for i, tag := range p.Tags {
		tags[i] = tag.Name
	}

	return tags
}

func (d *Draft) GetTagsName() []string {
	tags := make([]string, len(d.Tags))

	for i, tag := range d.Tags {
		tags[i] = tag.Name
	}

	return tags
}
