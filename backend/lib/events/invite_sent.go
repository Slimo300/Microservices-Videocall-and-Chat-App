package events

import (
	"time"

	"github.com/google/uuid"
)

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

func (InviteSentEvent) EventName() string {
	return "groups.invitesent"
}

type Group struct {
	Name    string    `json:"name" mapstructure:"name"`
	Picture string    `json:"pictureUrl" mapstructure:"pictureUrl"`
	Created time.Time `json:"created" mapstructure:"created"`
}

type User struct {
	UserName string `json:"username" mapstructure:"username"`
	Picture  string `json:"pictureUrl" mapstructure:"picture"`
}
