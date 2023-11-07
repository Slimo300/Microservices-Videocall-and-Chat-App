package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID       uuid.UUID     `gorm:"primaryKey" json:"messageID"`
	Posted   time.Time     `json:"created"`
	Text     string        `json:"text"`
	MemberID uuid.UUID     `gorm:"column:member_id;size:191" json:"memberID"`
	Member   Membership    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Deleters []Membership  `gorm:"many2many:users_who_deleted;constraint:OnDelete:CASCADE;"`
	Files    []MessageFile `gorm:"foreignKey:MessageID" json:"files"`
}

func (Message) TableName() string {
	return "messages"
}

type MessageFile struct {
	MessageID string `gorm:"size:191"`
	Key       string `gorm:"primaryKey" json:"key"`
	Extention string `json:"ext"`
}

func (MessageFile) TableName() string {
	return "messagefiles"
}
