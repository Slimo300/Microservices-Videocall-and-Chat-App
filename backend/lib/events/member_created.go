package events

import (
	"github.com/google/uuid"
)

type MemberCreatedEvent struct {
	ID               uuid.UUID `json:"memberID"`
	GroupID          uuid.UUID `json:"groupID"`
	UserID           uuid.UUID `json:"userID"`
	Adding           bool      `json:"adding"`
	DeletingMembers  bool      `json:"deletingMembers"`
	DeletingMessages bool      `json:"deletingMessages"`
	Admin            bool      `json:"setting"`
	Creator          bool      `json:"creator"`
}

func (MemberCreatedEvent) EventName() string {
	return "groups.membercreated"
}
