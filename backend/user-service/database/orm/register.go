package orm

import (
	"errors"
	"time"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/apperrors"
	"github.com/Slimo300/MicroservicesChatApp/backend/user-service/database"
	"github.com/Slimo300/MicroservicesChatApp/backend/user-service/models"
	"github.com/google/uuid"
	"github.com/thanhpk/randstr"
	"gorm.io/gorm"
)

func (db *Database) RegisterUser(user models.User) (returnUser models.User, returnCode models.VerificationCode, returnErr error) {
	// checking if username is taken
	if err := db.Where(models.User{UserName: user.UserName}).First(&models.User{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		returnErr = apperrors.NewConflict("username", user.UserName)
		return
	}
	// checking if email is taken
	if err := db.Where(models.User{Email: user.Email}).First(&models.User{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		returnErr = apperrors.NewConflict("email", user.Email)
		return
	}

	hash, err := database.HashPassword(user.Pass)
	if err != nil {
		returnErr = apperrors.NewBadRequest("invalid password")
		return
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		// user creation
		returnUser = models.User{ID: uuid.New(), UserName: user.UserName, Email: user.Email, Pass: hash, Verified: false, Created: now, Updated: now}
		if err := db.Create(&returnUser).Error; err != nil {
			return err
		}
		// verification code creation
		returnCode = models.VerificationCode{UserID: returnUser.ID, ActivationCode: randstr.String(10), Created: now}
		if err := db.Create(&returnCode).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return models.User{}, models.VerificationCode{}, apperrors.NewInternal()
	}

	return
}

func (db *Database) VerifyCode(code string) (models.User, error) {
	// checking if both verification code and the user it is refering to exist
	var verCode models.VerificationCode
	if err := db.Where(models.VerificationCode{ActivationCode: code}).First(&verCode).Error; err != nil {
		return models.User{}, apperrors.NewNotFound("code", code)
	}
	var user models.User
	if err := db.First(&user, verCode.UserID).Error; err != nil {
		return models.User{}, apperrors.NewNotFound("code", code)
	}

	elapsed := time.Now().Sub(verCode.Created)
	if err := db.Transaction(func(tx *gorm.DB) error {

		// we first delete verCode because no matter if code is expired or not
		// it won't be needed outside the scope of this function
		if err := db.Delete(&verCode).Error; err != nil {
			return apperrors.NewInternal()
		}

		// if verification code expired we delete created user and return not found error
		// pretending we don't know what the user wants ¯\_(ツ)_/¯
		if elapsed > db.Config.VerificationCodeDuration {
			if err := db.Delete(&user).Error; err != nil {
				return apperrors.NewInternal()
			}
			return apperrors.NewNotFound("code", code)
		}

		// if verification code is still valid we update the user `verified` property
		if err := db.Model(&user).Update("verified", true).Error; err != nil {
			return apperrors.NewInternal()
		}
		return nil

	}); err != nil {
		return models.User{}, err
	}
	return user, nil

}
