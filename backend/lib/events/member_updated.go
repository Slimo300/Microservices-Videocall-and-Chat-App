package events

import (
	"github.com/google/uuid"
)

type MemberUpdatedEvent struct {
	ID               uuid.UUID `json:"id"`
	GroupID          uuid.UUID `json:"groupID"`
	UserID           uuid.UUID `json:"userID"`
	DeletingMessages bool      `json:"deletingMessages,omitempty"`
	DeletingMembers  bool      `json:"deletingMembers,omitempty"`
	Adding           bool      `json:"adding,omitempty"`
	Setting          bool      `json:"setting,omitempty"`
}

func (MemberUpdatedEvent) EventName() string {
	return "groups.memberupdated"
}
