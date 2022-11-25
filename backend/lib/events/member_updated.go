package events

import (
	"github.com/google/uuid"
)

type MemberUpdatedEvent struct {
	ID               uuid.UUID `json:"id"`
	GroupID          uuid.UUID `json:"groupID"`
	UserID           uuid.UUID `json:"userID"`
	Adding           int       `json:"adding"`
	DeletingMessages int       `json:"deletingMessages"`
	DeletingMembers  int       `json:"deletingMembers"`
	Admin            int       `json:"admin"`
}

func (MemberUpdatedEvent) EventName() string {
	return "groups.memberupdated"
}
