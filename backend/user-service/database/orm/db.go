package orm

import (
	"fmt"
	"time"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/apperrors"
	"github.com/Slimo300/MicroservicesChatApp/backend/user-service/database"
	"github.com/Slimo300/MicroservicesChatApp/backend/user-service/models"
	"github.com/google/uuid"
	"github.com/thanhpk/randstr"
	"gorm.io/gorm"
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
		return "", apperrors.NewBadRequest("user has no profile picture")
	}

	if err := db.Model(&user).Update("picture_url", "").Error; err != nil {
		return "", apperrors.NewInternal()
	}
	return url, nil

}

// GetProfilePictureURL fetches user's `pictureURL` and if one is not set it generates one and saves to database
// Return values are pictureURL, boolean informing if new url was generated and error
func (db *Database) GetProfilePictureURL(userID uuid.UUID) (string, bool, error) {
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return "", false, apperrors.NewAuthorization("User not in database")
	}
	if user.PictureURL == "" {
		newPictureURL := uuid.NewString()
		if err := db.Model(&user).Update("picture_url", newPictureURL).Error; err != nil {
			return "", true, apperrors.NewInternal()
		}
		return newPictureURL, true, nil
	}
	return user.PictureURL, false, nil
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
	if err := db.Where(models.User{Email: email, Verified: true}).First(&user).Error; err != nil {
		return models.User{}, apperrors.NewBadRequest("wrong email or password")
	}
	if !database.CheckPassword(user.Pass, password) {
		return models.User{}, apperrors.NewBadRequest("wrong email or password")
	}
	return user, nil
}

func (db *Database) NewResetPasswordCode(email string) (*models.User, *models.ResetCode, error) {
	var user models.User
	if err := db.Where(models.User{Email: email}).First(&user).Error; err != nil {
		return nil, nil, nil
	}

	// here we check if user has any existing reset code and if any exist we delete it
	var resetCode models.ResetCode
	if err := db.First(&resetCode, user.ID).Error; err != gorm.ErrRecordNotFound {
		if err := db.Delete(&resetCode).Error; err != nil {
			return nil, nil, apperrors.NewInternal()
		}
	}

	resetCode = models.ResetCode{UserID: user.ID, Created: time.Now(), ResetCode: randstr.String(10)}
	if err := db.Create(&resetCode).Error; err != nil {
		return nil, nil, apperrors.NewInternal()
	}

	return &user, &resetCode, nil
}

func (db *Database) ResetPassword(code, newPassword string) error {
	var resetCode models.ResetCode
	if err := db.Where(models.ResetCode{ResetCode: code}).First(&resetCode).Error; err != nil {
		return apperrors.NewNotFound("reset code", code)
	}

	currentTime := time.Now()
	if currentTime.Sub(resetCode.Created) > db.Config.ResetCodeDuration {
		if err := db.Delete(&resetCode).Error; err != nil {
			return apperrors.NewInternal()
		}
		return apperrors.NewNotFound("reset code", code)
	}

	var user models.User
	if err := db.First(&user, resetCode.UserID).Error; err != nil {
		return apperrors.NewBadRequest("User does not exist")
	}

	hash, err := database.HashPassword(newPassword)
	if err != nil {
		return apperrors.NewBadRequest("Invalid password")
	}

	if err := db.Model(&user).Update("password", hash).Error; err != nil {
		return apperrors.NewInternal()
	}

	return nil
}
