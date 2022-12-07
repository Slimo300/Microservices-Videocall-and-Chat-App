package events

import (
	"github.com/google/uuid"
)

type MemberCreatedEvent struct {
	ID               uuid.UUID `json:"ID"`
	GroupID          uuid.UUID `json:"groupID"`
	UserID           uuid.UUID `json:"userID"`
	User             User      `json:"User"`
	Adding           bool      `json:"adding"`
	DeletingMembers  bool      `json:"deletingMembers"`
	DeletingMessages bool      `json:"deletingMessages"`
	Admin            bool      `json:"admin"`
	Creator          bool      `json:"creator"`
}

func (MemberCreatedEvent) EventName() string {
	return "groups.membercreated"
}
