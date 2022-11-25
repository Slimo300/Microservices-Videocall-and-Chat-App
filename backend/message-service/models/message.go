package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID       uuid.UUID `gorm:"primaryKey"`
	Posted   time.Time
	Text     string
	UserID   uuid.UUID
	GroupID  uuid.UUID
	Nick     string
	Deleters []Membership `gorm:"many2many:users_who_deleted;constraint:OnDelete:CASCADE;"`
}

func (Message) TableName() string {
	return "messages"
}
