package events

import "github.com/google/uuid"

type MessageDeletedEvent struct {
	ID      uuid.UUID `json:"id"`
	GroupID uuid.UUID `json:"groupID"`
}

func (MessageDeletedEvent) EventName() string {
	return "messages.deleted"
}
