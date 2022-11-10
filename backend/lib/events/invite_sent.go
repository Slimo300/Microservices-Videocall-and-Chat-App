package events

import "github.com/google/uuid"

type InviteSentEvent struct {
	ID       uuid.UUID `json:"id"`
	IssuerID uuid.UUID `json:"issuerID"`
	TargetID uuid.UUID `json:"targetID"`
}

func (InviteSentEvent) EventName() string {
	return "groups.invitesent"
}
