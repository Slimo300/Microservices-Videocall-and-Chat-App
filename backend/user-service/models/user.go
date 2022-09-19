package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	UserName  string    `gorm:"column:username;unique" json:"username"`
	Email     string    `gorm:"column:email;unique" json:"email"`
	Pass      string    `gorm:"column:password" json:"-"`
	Picture   string    `gorm:"picture_url" json:"pictureUrl"`
	Active    time.Time `gorm:"column:activity" json:"activity"`
	SignUp    time.Time `gorm:"column:signup" json:"signup"`
	LoggedIn  bool      `gorm:"column:logged" json:"logged"`
	Activated bool      `gorm:"column:activated" json:"activated"`
}

func (User) TableName() string {
	return "users"
}
