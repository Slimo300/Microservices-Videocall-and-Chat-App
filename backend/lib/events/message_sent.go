package events

import (
	"time"

	"github.com/google/uuid"
)

type MessageSentEvent struct {
	ID      uuid.UUID `json:"messageID" mapstructure:"messageID"`
	GroupID uuid.UUID `json:"groupID" mapstructure:"groupID"`
	UserID  uuid.UUID `json:"userID" mapstructure:"userID"`
	Text    string    `json:"text" mapstructure:"text"`
	Nick    string    `json:"nick" mapstructure:"nick"`
	Posted  time.Time `json:"posted" mapstructure:"posted"`
}

func (MessageSentEvent) EventName() string {
	return "wsmessages.created"
}
