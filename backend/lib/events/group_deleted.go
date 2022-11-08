package events

import (
	"github.com/google/uuid"
)

type GroupDeletedEvent struct {
	ID uuid.UUID `json:"id"`
}

func (GroupDeletedEvent) EventName() string {
	return "groups.deleted"
}
