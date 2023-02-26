package events

import (
	"github.com/google/uuid"
)

// MemberUpdatedEvent holds information about member's privileges being changed
type MemberUpdatedEvent struct {
	ID               uuid.UUID `json:"ID" mapstructure:"ID"`
	GroupID          uuid.UUID `json:"groupID" mapstructure:"groupID"`
	UserID           uuid.UUID `json:"userID" mapstructure:"userID"`
	User             User      `json:"User" mapstructure:"User"`
	Adding           bool      `json:"adding" mapstructure:"adding"`
	DeletingMessages bool      `json:"deletingMessages" mapstructure:"deletingMessages"`
	DeletingMembers  bool      `json:"deletingMembers" mapstructure:"deletingMembers"`
	Admin            bool      `json:"admin" mapstructure:"admin"`
}

// EventName method from Event interface
func (MemberUpdatedEvent) EventName() string {
	return "groups.memberupdated"
}
