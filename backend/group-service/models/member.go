package models

import (
	"fmt"
	"reflect"

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

// Here are methods and constants responsible for changing rights of a member

type operation int

const (
	REVOKE operation = iota - 1
	IGNORE
	GRANT
)

type MemberRights struct {
	Adding           operation `json:"adding,omitempty"`
	DeletingMessages operation `json:"deletingMessages,omitempty"`
	DeletingMembers  operation `json:"deletingMembers,omitempty"`
	Admin            operation `json:"admin,omitempty"`
}

func (m *Member) ApplyRights(rights MemberRights) error {
	val := reflect.ValueOf(rights)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		switch val.Field(i).Interface() {
		case IGNORE:
			continue
		case GRANT:
			m.grant(typ.Field(i).Name)
		case REVOKE:
			m.revoke(typ.Field(i).Name)
		default:
			return fmt.Errorf("Unsupported action code: %v", val.Field(i).Interface())
		}
	}
	return nil
}

func (m *Member) grant(field string) {
	reflect.ValueOf(m).Elem().FieldByName(field).SetBool(true)
}

func (m *Member) revoke(field string) {
	reflect.ValueOf(m).Elem().FieldByName(field).SetBool(false)
}
