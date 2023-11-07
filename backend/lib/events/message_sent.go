package events

import (
	"time"

	"github.com/google/uuid"
)

// MessageSentEvent holds information about message being sent
type MessageSentEvent struct {
	ID        uuid.UUID `json:"messageID" mapstructure:"messageID"`
	MemberID  uuid.UUID `json:"memberID" mapstructure:"memberID"`
	GroupID   uuid.UUID `json:"groupID" mapstructure:"groupID"`
	UserID    uuid.UUID `json:"userID" mapstructure:"userID"`
	Text      string    `json:"text" mapstructure:"text"`
	Nick      string    `json:"nick" mapstructure:"nick"`
	Posted    time.Time `json:"posted" mapstructure:"posted"`
	Files     []File    `json:"files" mapstructure:"files"`
	ServiceID uuid.UUID `json:"serviceID" mapstructure:"serviceID"`
}

// File holds information about file send alongside message
type File struct {
	Key       string `json:"key" mapstructure:"key"`
	Extension string `json:"ext" mapstructure:"ext"`
}

// EventName method from Event interface
func (MessageSentEvent) EventName() string {
	return "wsmessages.created"
}
