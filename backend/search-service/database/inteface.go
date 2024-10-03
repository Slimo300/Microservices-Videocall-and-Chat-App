package database

import (
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/search-service/models"
	"github.com/google/uuid"
)

type DBLayer interface {
	GetUsers(query string, num int) ([]models.User, error)
	AddUser(userID uuid.UUID, username string) error
	UpdateProfilePicture(userID uuid.UUID, hasPicture bool) error
}
