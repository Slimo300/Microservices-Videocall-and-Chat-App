package events

import (
	"time"

	"github.com/google/uuid"
)

// InviteRespondedEvent holds information about responding to an invite to group
type InviteRespondedEvent struct {
	ID       uuid.UUID `json:"ID" mapstructure:"ID"`
	IssuerID uuid.UUID `json:"issuerID" mapstructure:"issuerID"`
	TargetID uuid.UUID `json:"targetID" mapstructure:"targetID"`
	Target   User      `json:"target" mapstructure:"target"`
	GroupID  uuid.UUID `json:"groupID" mapstructure:"groupID"`
	Group    Group     `json:"group" mapstructure:"group"`
	Status   int       `json:"status" mapstructure:"status"`
	Modified time.Time `json:"modified" mapstructure:"modified"`
}

// EventName method from Event interface
func (InviteRespondedEvent) EventName() string {
	return "invite.updated"
}
