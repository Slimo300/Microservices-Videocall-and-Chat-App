package events

import (
	"time"

	"github.com/google/uuid"
)

type InviteRespondedEvent struct {
	ID       uuid.UUID `json:"ID"`
	IssuerID uuid.UUID `json:"issuerID"`
	TargetID uuid.UUID `json:"targetID"`
	Target   User      `json:"target"`
	GroupID  uuid.UUID `json:"groupID"`
	Group    Group     `json:"group"`
	Status   int       `json:"status"`
	Modified time.Time `json:"modified"`
}

func (InviteRespondedEvent) EventName() string {
	return "groups.inviteresponded"
}
