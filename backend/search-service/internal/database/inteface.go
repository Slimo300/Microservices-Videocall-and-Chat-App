package database

import (
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/events"
	"github.com/Slimo300/chat-searchservice/internal/models"
)

type DBLayer interface {
	GetUsers(query string, num int) ([]models.User, error)
	AddUser(user events.UserRegisteredEvent) error
	UpdateProfilePicture(ev events.UserPictureModifiedEvent) error
}
