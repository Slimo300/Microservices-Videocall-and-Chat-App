package events

import (
	"github.com/google/uuid"
)

// MessageDeletedEvent holds information about message being deleted
type MessageDeletedEvent struct {
	ID      uuid.UUID `json:"messageID" mapstructure:"messageID"`
	GroupID uuid.UUID `json:"groupID" mapstructure:"groupID"`
}

// EventName method from Event interface
func (MessageDeletedEvent) EventName() string {
	return "messages.deleted"
}
