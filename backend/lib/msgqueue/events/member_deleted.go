package events

import "github.com/google/uuid"

type MemberDeletedEvent struct {
	ID uuid.UUID `json:"id"`
}

func (MemberDeletedEvent) EventName() string {
	return "member.deleted"
}
