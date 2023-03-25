package models

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `gorm:"primaryKey" json:"ID"`
	UserName string    `gorm:"column:username;unique" json:"username"`
	Picture  string    `gorm:"picture_url" json:"pictureUrl"`
}

func (User) TableName() string {
	return "users"
}
