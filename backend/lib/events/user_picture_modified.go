package events

import (
	"github.com/google/uuid"
)

type UserPictureModifiedEvent struct {
	ID         uuid.UUID `json:"userID" mapstructure:"userID"`
	PictureURL string    `json:"picture" mapstructure:"picture"`
}

func (UserPictureModifiedEvent) EventName() string {
	return "users.picturemodified"
}
