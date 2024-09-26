package models

import (
	"time"

	"github.com/google/uuid"
)

type Group struct {
	ID         uuid.UUID `gorm:"primaryKey" json:"ID"`
	Name       string    `gorm:"column:name" json:"name"`
	HasPicture bool      `gorm:"column:hasPicture" json:"hasPicture"`
	Created    time.Time `gorm:"column:created" json:"created"`
	Members    []Member  `gorm:"foreignKey:GroupID"`
}

func (Group) TableName() string {
	return "groups"
}
