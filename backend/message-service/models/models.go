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
	Deleters []Membership `gorm:"foreignKey:MembershipID;constraint:OnDelete:CASCADE"`
}

func (Message) TableName() string {
	return "messages"
}

type Membership struct {
	MembershipID     uuid.UUID `gorm:"primaryKey"`
	UserID           uuid.UUID `gorm:"uniqueIndex:idx_first;size:191"`
	GroupID          uuid.UUID `gorm:"uniqueIndex:idx_first;size:191"`
	Creator          bool
	DeletingMessages bool
}

func (Membership) TableName() string {
	return "memberships"
}
