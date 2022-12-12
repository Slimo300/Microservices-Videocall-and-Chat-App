package database

import (
	"github.com/Slimo300/MicroservicesChatApp/backend/user-service/models"
	"github.com/google/uuid"
)

type DBLayer interface {
	GetUserById(uid uuid.UUID) (models.User, error)
	SignIn(email, password string) (models.User, error)

	RegisterUser(models.User) (*models.User, *models.VerificationCode, error)
	VerifyCode(code string) (*models.User, error)

	ChangePassword(userID uuid.UUID, oldPassword, newPassword string) error

	GetProfilePictureURL(userID uuid.UUID) (string, error)
	DeleteProfilePicture(userID uuid.UUID) (string, error)

	NewResetPasswordCode(email string) (*models.User, *models.ResetCode, error)
	ResetPassword(code, newPassword string) error
}
