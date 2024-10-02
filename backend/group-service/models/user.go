package models

import "github.com/google/uuid"

type User struct {
	ID         uuid.UUID `gorm:"primaryKey" json:"ID"`
	UserName   string    `gorm:"column:username;unique" json:"username"`
	HasPicture bool      `gorm:"has_picture" json:"hasPicture"`
}

func (User) TableName() string {
	return "users"
}
