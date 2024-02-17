package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `gorm:"primaryKey"`
	UserName   string    `gorm:"column:username;unique" json:"username"`
	Email      string    `gorm:"column:email;unique" json:"email"`
	Pass       string    `gorm:"column:password" json:"-"`
	PictureURL string    `gorm:"picture_url" json:"pictureUrl"`
	Verified   bool      `gorm:"column:verified" json:"verified"`
}

func (User) TableName() string {
	return "users"
}
