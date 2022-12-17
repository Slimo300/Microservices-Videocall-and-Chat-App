package events

import (
	"github.com/google/uuid"
)

type MemberDeletedEvent struct {
	ID      uuid.UUID `json:"ID" mapstructure:"ID"`
	GroupID uuid.UUID `json:"groupID" mapstructure:"groupID"`
	UserID  uuid.UUID `json:"userID" mapstructure:"userID"`
}

func (MemberDeletedEvent) EventName() string {
	return "groups.memberdeleted"
}
