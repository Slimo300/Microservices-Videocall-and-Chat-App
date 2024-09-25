package events

import (
	"github.com/google/uuid"
)

// UserRegisteredEvent holds information about new users
type UserRegisteredEvent struct {
	ID       uuid.UUID `json:"userID" mapstructure:"userID"`
	Username string    `json:"username" mapstructure:"username"`
}

// EventName method from Event interface
func (UserRegisteredEvent) EventName() string {
	return "user.created"
}
