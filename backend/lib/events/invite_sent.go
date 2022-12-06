package events

import (
	"time"

	"github.com/google/uuid"
)

type InviteSentEvent struct {
	ID       uuid.UUID `json:"ID"`
	IssuerID uuid.UUID `json:"issuerID"`
	Issuer   User      `json:"issuer"`
	TargetID uuid.UUID `json:"targetID"`
	GroupID  uuid.UUID `json:"groupID"`
	Group    Group     `json:"group"`
	Status   int       `json:"status"`
	Modified time.Time `json:"modified"`
}

func (InviteSentEvent) EventName() string {
	return "groups.invitesent"
}

type Group struct {
	ID      uuid.UUID `json:"ID"`
	Name    string    `json:"name"`
	Picture string    `json:"pictureUrl"`
	Created time.Time `json:"created"`
}

type User struct {
	ID       uuid.UUID `json:"ID"`
	UserName string    `json:"username"`
	Picture  string    `json:"pictureUrl"`
}
