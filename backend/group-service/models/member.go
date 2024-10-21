package models

import (
	"github.com/google/uuid"
)

type Member struct {
	id               uuid.UUID
	groupID          uuid.UUID
	userID           uuid.UUID
	user             User
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
	return m.creator
}

func (m Member) CanSendInvite() bool {
	return m.adding || m.creator
}

func (m Member) CanDelete(target Member) bool {
	if target.creator {
		return false
	}
	if m.id == target.id { // user can delete himself if he's not a creator
		return true
	}
	if m.creator { // creator can delete everyone (except himself)
		return true
	}
	if m.deletingMembers && !target.deletingMembers { // user with rights to delete can delete everyone except creator and other user with such rights
		return true
	}
	return false
}

func (m Member) CanAlter(target Member) bool {
	if target.creator {
		return false
	}
	if m.creator {
		return true
	}
	return false
}

type MemberRights struct {
	Adding           bool `json:"adding"`
	DeletingMessages bool `json:"deletingMessages"`
	DeletingMembers  bool `json:"deletingMembers"`
	Muting           bool `json:"muting"`
}

func (m *Member) ApplyRights(rights MemberRights) {
	m.adding = rights.Adding
	m.deletingMembers = rights.DeletingMembers
	m.deletingMessages = rights.DeletingMessages
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

func UnmarshalMemberFromDatabase(memberID, userID, groupID uuid.UUID, user User, adding, delMembers, delMessages, muting, creator bool) Member {
	return Member{
		id:               memberID,
		groupID:          groupID,
		userID:           userID,
		user:             user,
		adding:           adding,
		deletingMembers:  delMembers,
		deletingMessages: delMessages,
		muting:           muting,
		creator:          creator,
	}
}
