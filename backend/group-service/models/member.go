package models

import (
	"github.com/google/uuid"
)

type Member struct {
	ID               uuid.UUID `gorm:"primaryKey" json:"ID"`
	GroupID          uuid.UUID `gorm:"column:group_id;uniqueIndex:idx_first;size:191" json:"groupID"`
	UserID           uuid.UUID `gorm:"column:user_id;uniqueIndex:idx_first;size:191" json:"userID"`
	User             User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Group            Group     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Adding           bool      `gorm:"column:adding" json:"adding"`
	DeletingMembers  bool      `gorm:"column:deleting_members" json:"deletingMembers"`
	DeletingMessages bool      `gorm:"column:deleting_messages" json:"deletingMessages"`
	Muting           bool      `gorm:"column:muting" json:"muting"`
	Admin            bool      `gorm:"column:setting" json:"admin"`
	Creator          bool      `gorm:"column:creator" json:"creator"`
}

func (Member) TableName() string {
	return "members"
}

// Here are methods and constants responsible for resolving users rights in a group when they try to alter
// other members of a group

type role int

const (
	CREATOR role = iota + 1
	ADMIN
	DELETER
	BASIC
)

func (m Member) CanDelete(target Member) bool {
	if m.role(false) < target.role(false) {
		return true
	}
	if m.ID == target.ID && !m.Creator { // user can delete himself
		return true
	}
	return false
}

func (m Member) CanAlter(target Member) bool {
	return m.role(true) < target.role(true)
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

type MemberRights struct {
	Adding           bool `json:"adding"`
	DeletingMessages bool `json:"deletingMessages"`
	DeletingMembers  bool `json:"deletingMembers"`
	Admin            bool `json:"admin"`
	Muting           bool `json:"muting"`
}

func (m *Member) ApplyRights(rights MemberRights) {
	m.Adding = rights.Adding
	m.DeletingMembers = rights.DeletingMembers
	m.DeletingMessages = rights.DeletingMessages
	m.Admin = rights.Admin
	m.Muting = rights.Muting
}
