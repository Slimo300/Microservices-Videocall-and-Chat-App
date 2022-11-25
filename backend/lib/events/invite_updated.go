package events

import "github.com/google/uuid"

type InviteRespondedEvent struct {
	ID       uuid.UUID `json:"id"`
	IssuerID uuid.UUID `json:"issID"`
	Answer   bool      `json:"status"`
}

func (InviteRespondedEvent) EventName() string {
	return "groups.inviteresponded"
}
