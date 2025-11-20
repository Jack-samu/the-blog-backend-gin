package models

import (
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID            string `gorm:"type:char(36);primaryKey"`
	Username      string `gorm:"unique;not null"`
	Email         string `gorm:"unique;not null"`
	Pwd           string `gorm:"not null"`
	Bio           string
	Avatar        string
	CreatedAt     time.Time `gorm:"autoUpdateTime:false"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime:milli"`
	LastActivity  time.Time `gorm:"autoUpdateTime:false"`
	FailedLogin   int
	CaptchaReqCnt int

	// 外键外联
	Posts      []Post     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Drafts     []Draft    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Images     []Img      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Categories []Category `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`

	Comments []Comment `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL"`
	Replies  []Reply   `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL"`
}

func (u *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	log.Printf("哈希：%s\n", bytes)

	u.Pwd = string(bytes)
	return nil
}

func (u *User) CheckPassword(pwdString string) bool {
	log.Printf("要校验的密码：%s\n", pwdString)
	err := bcrypt.CompareHashAndPassword([]byte(u.Pwd), []byte(pwdString))
	return err == nil
}
