package models

import (
	"time"

	"github.com/google/uuid"
)

type invite_status int

const (
	INVITE_AWAITING invite_status = iota + 1
	INVITE_ACCEPT
	INVITE_DECLINE
)

type Invite struct {
	ID       uuid.UUID     `gorm:"primaryKey"`
	IssId    uuid.UUID     `gorm:"column:iss_id;size:191" json:"issID"`
	Iss      User          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"issuer"`
	TargetID uuid.UUID     `gorm:"column:target_id;size:191" json:"targetID"`
	Target   User          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	GroupID  uuid.UUID     `gorm:"column:group_id;size:191" json:"groupID"`
	Group    Group         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"group"`
	Status   invite_status `gorm:"column:status" json:"status"`
	Created  time.Time     `gorm:"column:created" json:"created"`
	Modified time.Time     `gorm:"column:modified" json:"modified"`
}

func (Invite) TableName() string {
	return "invites"
}
