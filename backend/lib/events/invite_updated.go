package events

import (
	"time"

	"github.com/google/uuid"
)

type InviteRespondedEvent struct {
	ID         uuid.UUID `json:"inviteID"`
	IssuerID   uuid.UUID `json:"issuerID"`
	IssuerName string    `json:"issuerName"`
	TargetID   uuid.UUID `json:"targetID"`
	TargetName string    `json:"targetName"`
	GroupID    uuid.UUID `json:"groupID"`
	GroupName  string    `json:"groupName"`
	Status     int       `json:"status"`
	Modified   time.Time `json:"modified"`
}

func (InviteRespondedEvent) EventName() string {
	return "groups.inviteresponded"
}
