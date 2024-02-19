package events

import (
	"time"

	"github.com/google/uuid"
)

// InviteSentEvent holds information about sending group invite
type InviteSentEvent struct {
	ID       uuid.UUID `json:"ID" mapstructure:"ID"`
	IssuerID uuid.UUID `json:"issuerID" mapstructure:"issuerID"`
	Issuer   User      `json:"issuer" mapstructure:"issuer"`
	TargetID uuid.UUID `json:"targetID" mapstructure:"targetID"`
	GroupID  uuid.UUID `json:"groupID" mapstructure:"groupID"`
	Group    Group     `json:"group" mapstructure:"group"`
	Status   int       `json:"status" mapstructure:"status"`
	Modified time.Time `json:"modified" mapstructure:"modified"`
}

// EventName method from Event interface
func (InviteSentEvent) EventName() string {
	return "invite.created"
}

// Group holds group information to be sent with invite
type Group struct {
	Name    string    `json:"name" mapstructure:"name"`
	Picture string    `json:"pictureUrl" mapstructure:"pictureUrl"`
	Created time.Time `json:"created" mapstructure:"created"`
}

// User holds user information to be sent with invite
type User struct {
	UserName string `json:"username" mapstructure:"username"`
	Picture  string `json:"pictureUrl" mapstructure:"picture"`
}
