package models

import (
	"github.com/google/uuid"
)

type Membership struct {
	MembershipID     uuid.UUID `gorm:"primaryKey" json:"memberID"`
	UserID           uuid.UUID `gorm:"uniqueIndex:idx_first;size:191" json:"userID"`
	GroupID          uuid.UUID `gorm:"uniqueIndex:idx_first;size:191" json:"groupID"`
	Username         string    `json:"username"`
	Admin            bool      `json:"admin"`
	Creator          bool      `json:"creator"`
	DeletingMessages bool      `json:"deletingMessages"`
}

func (Membership) TableName() string {
	return "memberships"
}

func (m *Membership) CanDeleteMessage(msg *Message) bool {
	if msg.MemberID == m.MembershipID || m.Creator {
		return true
	}

	if m.Admin && !msg.Member.Creator {
		return true
	}

	if m.DeletingMessages && !msg.Member.Creator && !msg.Member.Admin {
		return true
	}

	return false
}
