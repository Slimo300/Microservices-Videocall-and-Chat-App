package events

import (
	"github.com/google/uuid"
)

type MessageDeletedEvent struct {
	ID      uuid.UUID `json:"messageID" mapstructure:"messageID"`
	GroupID uuid.UUID `json:"groupID" mapstructure:"groupID"`
}

func (MessageDeletedEvent) EventName() string {
	return "messages.deleted"
}
