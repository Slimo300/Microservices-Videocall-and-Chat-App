package models

import "github.com/google/uuid"

type role int

const (
	CREATOR role = iota + 1
	ADMIN
	DELETER
	BASIC
)

type MemberRights struct {
	Adding           *bool `json:"adding" binding:"required"`
	DeletingMessages *bool `json:"deleting" binding:"required"`
	DeletingMembers  *bool `json:"deletingMembers" binding:"required"`
	Admin            *bool `json:"setting" binding:"required"`
}

type Member struct {
	ID               uuid.UUID `gorm:"primaryKey"`
	GroupID          uuid.UUID `gorm:"column:group_id;uniqueIndex:idx_first;size:191" json:"group_id"`
	UserID           uuid.UUID `gorm:"column:user_id;uniqueIndex:idx_first;size:191" json:"user_id"`
	User             User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Group            Group     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Adding           bool      `gorm:"column:adding" json:"adding"`
	DeletingMembers  bool      `gorm:"column:deleting_members" json:"deletingMembers"`
	DeletingMessages bool      `gorm:"column:deleting_messages" json:"deletingMessages"`
	Admin            bool      `gorm:"column:setting" json:"setting"`
	Creator          bool      `gorm:"column:creator" json:"creator"`
}

func (Member) TableName() string {
	return "members"
}

func (m Member) CanDelete(target Member) bool {
	if m.role(false) < target.role(false) {
		return true
	}
	return false
}

func (m Member) CanAlter(target Member) bool {
	if m.role(true) < target.role(true) {
		return true
	}
	return false
}

func (m Member) role(noDeleter bool) role {
	if m.Creator {
		return CREATOR
	}
	if m.Admin {
		return ADMIN
	}
	if m.DeletingMembers && !noDeleter {
		return DELETER
	}
	return BASIC
}
