package models

import (
	"github.com/google/uuid"
)

type Member struct {
	id      uuid.UUID
	groupID uuid.UUID
	userID  uuid.UUID
	user    User
	// group            Group
	adding           bool
	deletingMembers  bool
	deletingMessages bool
	muting           bool
	admin            bool
	creator          bool
}

func (m Member) ID() uuid.UUID          { return m.id }
func (m Member) GroupID() uuid.UUID     { return m.groupID }
func (m Member) UserID() uuid.UUID      { return m.userID }
func (m Member) User() User             { return m.user }
func (m Member) Adding() bool           { return m.adding }
func (m Member) DeletingMessages() bool { return m.deletingMessages }
func (m Member) DeletingMembers() bool  { return m.deletingMembers }
func (m Member) Admin() bool            { return m.admin }
func (m Member) Muting() bool           { return m.muting }
func (m Member) Creator() bool          { return m.creator }

func (m Member) CanDeleteGroup() bool {
	return m.creator
}

func (m Member) CanUpdateGroup() bool {
	return m.creator || m.admin
}

func (m Member) CanSendInvite() bool {
	return m.adding || m.admin || m.creator
}

func (m Member) CanDelete(target Member) bool {
	if m.role(false) < target.role(false) {
		return true
	}
	if m.id == target.id && !m.creator { // user can delete himself
		return true
	}
	return false
}

func (m Member) CanAlter(target Member) bool {
	return m.role(true) < target.role(true)
}

type role int

const (
	CREATOR role = iota + 1
	ADMIN
	DELETER
	BASIC
)

func (m Member) role(noDeleter bool) role {
	if m.creator {
		return CREATOR
	}
	if m.admin {
		return ADMIN
	}
	if m.deletingMembers && !noDeleter {
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
	m.adding = rights.Adding
	m.deletingMembers = rights.DeletingMembers
	m.deletingMessages = rights.DeletingMessages
	m.admin = rights.Admin
	m.muting = rights.Muting
}

func newCreatorMember(userID, groupID uuid.UUID) Member {
	return Member{
		id:               uuid.New(),
		userID:           userID,
		groupID:          groupID,
		adding:           true,
		deletingMembers:  true,
		deletingMessages: true,
		muting:           true,
		admin:            true,
		creator:          true,
	}
}

func NewMember(userID, groupID uuid.UUID) Member {
	return Member{
		id:      uuid.New(),
		userID:  userID,
		groupID: groupID,
	}
}

func UnmarshalMemberFromDatabase(memberID, userID, groupID uuid.UUID, user User, adding, delMembers, delMessages, muting, admin, creator bool) Member {
	return Member{
		id:               memberID,
		groupID:          groupID,
		userID:           userID,
		user:             user,
		adding:           adding,
		deletingMembers:  delMembers,
		deletingMessages: delMessages,
		muting:           muting,
		admin:            admin,
		creator:          creator,
	}
}
