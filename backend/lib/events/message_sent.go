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
	Files   []File    `json:"files" mapstructure:"files"`
}

type File struct {
	Key       string `json:"key" mapstructure:"key"`
	Extension string `json:"ext" mapstructure:"ext"`
}

func (MessageSentEvent) EventName() string {
	return "wsmessages.created"
}
