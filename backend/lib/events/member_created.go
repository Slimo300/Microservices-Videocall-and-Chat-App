package events

import (
	"github.com/google/uuid"
)

// MemberCreatedEvent holds information about adding user to a group
type MemberCreatedEvent struct {
	ID               uuid.UUID `json:"ID" mapstructure:"ID"`
	GroupID          uuid.UUID `json:"groupID" mapstructure:"groupID"`
	UserID           uuid.UUID `json:"userID" mapstructure:"userID"`
	User             User      `json:"User" mapstructure:"User"`
	Creator          bool      `json:"creator" mapstructure:"creator"`
	Admin            bool      `json:"admin" mapstructure:"admin"`
	Adding           bool      `json:"adding" mapstructure:"adding"`
	DeletingMembers  bool      `json:"deletingMembers" mapstructure:"deletingMembers"`
	DeletingMessages bool      `json:"deletingMessages" mapstructure:"deletingMessages"`
	Muting           bool      `json:"muting" mapstructure:"muting"`
}

// EventName method from Event interface
func (MemberCreatedEvent) EventName() string {
	return "group.membercreated"
}
