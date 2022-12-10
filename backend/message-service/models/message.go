package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID       uuid.UUID    `gorm:"primaryKey" json:"messageID"`
	Posted   time.Time    `json:"created"`
	Text     string       `json:"text"`
	UserID   uuid.UUID    `json:"userID"`
	GroupID  uuid.UUID    `json:"groupID"`
	Nick     string       `json:"nick"`
	Deleters []Membership `gorm:"many2many:users_who_deleted;constraint:OnDelete:CASCADE;"`
}

func (Message) TableName() string {
	return "messages"
}
