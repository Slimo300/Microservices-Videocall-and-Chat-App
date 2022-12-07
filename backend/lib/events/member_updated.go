package events

import (
	"github.com/google/uuid"
)

type MemberUpdatedEvent struct {
	ID               uuid.UUID `json:"ID"`
	GroupID          uuid.UUID `json:"groupID"`
	UserID           uuid.UUID `json:"userID"`
	User             User      `json:"User"`
	Adding           bool      `json:"adding"`
	DeletingMessages bool      `json:"deletingMessages"`
	DeletingMembers  bool      `json:"deletingMembers"`
	Admin            bool      `json:"admin"`
}

func (MemberUpdatedEvent) EventName() string {
	return "groups.memberupdated"
}
