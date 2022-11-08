package events

import (
	"github.com/google/uuid"
)

type UserRegisteredEvent struct {
	ID         uuid.UUID `json:"id"`
	Username   string    `json:"username"`
	PictureURL string    `json:"picture"`
}

func (UserRegisteredEvent) EventName() string {
	return "users.created"
}
