package orm

import (
	"fmt"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/apperrors"
	"github.com/Slimo300/MicroservicesChatApp/backend/user-service/database"
	"github.com/Slimo300/MicroservicesChatApp/backend/user-service/models"
	"github.com/google/uuid"
)

func (db *Database) GetUserById(id uuid.UUID) (user models.User, err error) {
	return user, db.First(&user, id).Error
}

// DeleteProfilePicture updates user's property "pictureURL" to empty line
// It returns Authorization error if user requesting operation is not found in database
func (db *Database) DeleteProfilePicture(userID uuid.UUID) (string, error) {
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return "", apperrors.NewAuthorization("User not in database")
	}

	url := user.PictureURL
	if url == "" {
		return "", apperrors.NewForbidden("user has no profile picture")
	}

	if err := db.Model(&user).Update("picture_url", "").Error; err != nil {
		return "", apperrors.NewInternal()
	}
	return url, nil

}

// GetProfilePictureURL returns user's `pictureURL` and if one is not set it generates one and saves to database
func (db *Database) GetProfilePictureURL(userID uuid.UUID) (string, error) {
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return "", apperrors.NewAuthorization("User not in database")
	}
	if user.PictureURL == "" {
		newPictureURL := uuid.NewString()
		if err := db.Model(&user).Update("picture_url", newPictureURL).Error; err != nil {
			return "", apperrors.NewInternal()
		}
		return newPictureURL, nil
	}
	return user.PictureURL, nil
}

func (db *Database) ChangePassword(userID uuid.UUID, oldPassword, newPassword string) error {
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return apperrors.NewAuthorization("User not in database")
	}

	if !database.CheckPassword(user.Pass, oldPassword) {
		return apperrors.NewForbidden("Wrong Password")
	}

	hash, err := database.HashPassword(newPassword)
	if err != nil {
		return apperrors.NewBadRequest(fmt.Sprintf("Invalid password: %v", err))
	}

	if err := db.Model(&user).Update("password", hash).Error; err != nil {
		return apperrors.NewInternal()
	}
	return nil
}

func (db *Database) SignIn(email, password string) (models.User, error) {
	var user models.User
	if err := db.Where(models.User{Email: email}).First(&user).Error; err != nil {
		return models.User{}, apperrors.NewBadRequest("invalid credentials")
	}
	if !database.CheckPassword(user.Pass, password) {
		return models.User{}, apperrors.NewBadRequest("invalid credentials")
	}
	return user, nil
}
