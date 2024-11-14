package models

import (
	"github.com/google/uuid"
)

type Member struct {
	id               uuid.UUID
	userID           uuid.UUID
	groupID          uuid.UUID
	username         string
	creator          bool
	admin            bool
	deletingMessages bool
}

func (m Member) ID() uuid.UUID          { return m.id }
func (m Member) UserID() uuid.UUID      { return m.userID }
func (m Member) GroupID() uuid.UUID     { return m.groupID }
func (m Member) Username() string       { return m.username }
func (m Member) Creator() bool          { return m.creator }
func (m Member) Admin() bool            { return m.admin }
func (m Member) DeletingMessages() bool { return m.deletingMessages }

func NewMember(memberID, userID, groupID uuid.UUID, username string, isCreator bool) Member {
	return Member{
		id:               memberID,
		userID:           userID,
		groupID:          groupID,
		username:         username,
		creator:          isCreator,
		admin:            isCreator,
		deletingMessages: isCreator,
	}
}

func UnmarshalMemberFromDatabase(memberID, userID, groupID uuid.UUID, username string, isCreator, isAdmin, canDeleteMsgs bool) Member {
	return Member{
		id:               memberID,
		userID:           userID,
		groupID:          groupID,
		username:         username,
		creator:          isCreator,
		admin:            isAdmin,
		deletingMessages: canDeleteMsgs,
	}
}

func (m *Member) CanDeleteMessage(msg *Message) bool {
	if msg.member.groupID != m.groupID {
		return false
	}
	if msg.memberID == m.id || m.creator {
		return true
	}
	if m.admin && !msg.member.creator {
		return true
	}
	if m.deletingMessages && !msg.member.creator && !msg.member.admin {
		return true
	}
	return false
}

func (m *Member) UpdateRights(admin, deletingMessages bool) (updated bool) {
	if m.admin == admin && m.deletingMessages == deletingMessages {
		return false
	}
	m.admin = admin
	m.deletingMessages = deletingMessages
	return true
}
