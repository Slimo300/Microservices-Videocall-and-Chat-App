package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID       uuid.UUID `gorm:"primaryKey"`
	Posted   time.Time `gorm:"column:posted" json:"posted"`
	Text     string    `gorm:"column:text" json:"text"`
	GroupID  uuid.UUID `gorm:"column:id_group;size:191" json:"groupID"`
	MemberID uuid.UUID `gorm:"column:id_member;size:191" json:"memberID"`
	Nick     string    `gorm:"column:nick" json:"nick"`
}

func (Message) TableName() string {
	return "messages"
}

type Membership struct {
	UserID  uuid.UUID `gorm:"primaryKey;size:191"`
	GroupID uuid.UUID `gorm:"primaryKey;size:191"`
	Deleted bool      `gorm:"column:deleted" json:"deleted"`
}

func (Membership) TableName() string {
	return "memberships"
}
