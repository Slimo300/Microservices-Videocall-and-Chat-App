package events

import (
	"github.com/google/uuid"
)

// UserPictureModifiedEvent holds information about user changing his profile picture
type UserPictureModifiedEvent struct {
	ID         uuid.UUID `json:"userID" mapstructure:"userID"`
	HasPicture bool      `json:"hasPicture" mapstructure:"hasPicture"`
}

// EventName method from Event interface
func (UserPictureModifiedEvent) EventName() string {
	return "user.picturemodified"
}
