package models

import (
	"github.com/google/uuid"
)

type Membership struct {
	MembershipID     uuid.UUID `gorm:"primaryKey"`
	UserID           uuid.UUID `gorm:"uniqueIndex:idx_first;size:191"`
	GroupID          uuid.UUID `gorm:"uniqueIndex:idx_first;size:191"`
	Admin            bool
	Creator          bool
	DeletingMessages bool
}

func (Membership) TableName() string {
	return "memberships"
}
