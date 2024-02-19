package events

import (
	"github.com/google/uuid"
)

// GroupDeletedEvent holds information about deleting group event
type GroupDeletedEvent struct {
	ID uuid.UUID `json:"groupID" mapstructure:"groupID"`
}

// EventName method from Event interface
func (GroupDeletedEvent) EventName() string {
	return "group.deleted"
}
