package events

import (
	"github.com/google/uuid"
)

// MemberDeletedEvent holds information about member being deleted from group
type MemberDeletedEvent struct {
	ID      uuid.UUID `json:"ID" mapstructure:"ID"`
	GroupID uuid.UUID `json:"groupID" mapstructure:"groupID"`
	UserID  uuid.UUID `json:"userID" mapstructure:"userID"`
}

// EventName method from Event interface
func (MemberDeletedEvent) EventName() string {
	return "group.memberdeleted"
}
