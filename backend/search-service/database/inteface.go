package database

import (
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/events"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/search-service/models"
)

type DBLayer interface {
	GetUsers(query string, num int) ([]models.User, error)
	AddUser(user events.UserRegisteredEvent) error
	UpdateProfilePicture(ev events.UserPictureModifiedEvent) error
}
