package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `gorm:"primaryKey"`
	UserName   string    `gorm:"column:username;unique" json:"username"`
	Email      string    `gorm:"column:email;unique" json:"email"`
	Pass       string    `gorm:"column:password" json:"-"`
	PictureURL string    `gorm:"picture_url" json:"pictureUrl"`
	Created    time.Time `gorm:"column:created" json:"created"`
	Updated    time.Time `gorm:"column:updated" json:"updated"`
	Verified   bool      `gorm:"column:verified" json:"verified"`
}

func (User) TableName() string {
	return "users"
}
