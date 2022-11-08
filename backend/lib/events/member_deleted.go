package events

import (
	"github.com/google/uuid"
)

type MemberDeletedEvent struct {
	ID      uuid.UUID `json:"id"`
	GroupID uuid.UUID `json:"groupID"`
	UserID  uuid.UUID `json:"userID"`
}

func (MemberDeletedEvent) EventName() string {
	return "groups.memberdeleted"
}
