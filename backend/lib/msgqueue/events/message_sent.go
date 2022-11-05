package events

import (
	"time"

	"github.com/google/uuid"
)

type MessageSentEvent struct {
	ID      uuid.UUID `json:"id"`
	GroupID uuid.UUID `json:"groupID"`
	UserID  uuid.UUID `json:"userID"`
	Text    string    `json:"text"`
	Nick    string    `json:"nick"`
	Posted  time.Time `json:"posted"`
}

func (MessageSentEvent) EventName() string {
	return "message.created"
}
