package events

import (
	"github.com/google/uuid"
)

type UserRegisteredEvent struct {
	ID         uuid.UUID `json:"userID" mapstructure:"userID"`
	Username   string    `json:"username" mapstructure:"username"`
	PictureURL string    `json:"picture" mapstructure:"picture"`
}

func (UserRegisteredEvent) EventName() string {
	return "users.created"
}
