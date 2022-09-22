package database

import (
	"github.com/Slimo300/MicroservicesChatApp/backend/user-service/models"
	"github.com/google/uuid"
)

type DBLayer interface {
	IsEmailInDatabase(email string) bool
	IsUsernameInDatabase(username string) bool

	GetUserById(uid uuid.UUID) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	GetUserByUsername(username string) (models.User, error)

	RegisterUser(models.User) (models.User, error)
	SignInUser(id uuid.UUID) error
	SignOutUser(id uuid.UUID) error

	SetPassword(userID uuid.UUID, password string) error
	GetProfilePictureURL(userID uuid.UUID) (string, error)
	SetProfilePicture(userID uuid.UUID, newURI string) error
	DeleteProfilePicture(userID uuid.UUID) error

	NewVerificationCode(userID uuid.UUID, code string) (models.VerificationCode, error)
	VerifyCode(userID uuid.UUID, code string) error
}
