package database

import (
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/events"
	"github.com/Slimo300/MicroservicesChatApp/backend/search-service/models"
)

type DBLayer interface {
	GetUsers(query string, num int) ([]models.User, error)
	AddUser(user events.UserRegisteredEvent) error
}
