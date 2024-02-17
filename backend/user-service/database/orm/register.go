package orm

import (
	"errors"
	"fmt"
	"time"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/apperrors"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/models"
	"github.com/google/uuid"
	"github.com/thanhpk/randstr"
	"gorm.io/gorm"
)

func (db *Database) RegisterUser(user models.User) (*models.User, *models.VerificationCode, error) {
	// checking if username is taken
	if err := db.Where(models.User{UserName: user.UserName}).First(&models.User{}).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, apperrors.NewConflict(fmt.Sprintf("Username %s already taken", user.UserName))
	}
	// checking if email is taken
	if err := db.Where(models.User{Email: user.Email}).First(&models.User{}).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, apperrors.NewConflict(fmt.Sprintf("Email %s already taken", user.Email))
	}

	hash, err := database.HashPassword(user.Pass)
	if err != nil {
		return nil, nil, apperrors.NewBadRequest("invalid password")
	}

	var newUser models.User
	var code models.VerificationCode

	if err := db.Transaction(func(tx *gorm.DB) error {
		// user creation
		newUser = models.User{ID: uuid.New(), UserName: user.UserName, Email: user.Email, Pass: hash, Verified: false}
		if err := tx.Create(&newUser).Error; err != nil {
			return err
		}
		// verification code creation
		code = models.VerificationCode{UserID: newUser.ID, ActivationCode: randstr.String(10), Created: time.Now()}
		if err = tx.Create(&code).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, nil, err
	}

	return &newUser, &code, nil
}

func (db *Database) VerifyCode(code string) (*models.User, error) {
	// checking if both verification code and the user it is refering to exist
	var verCode models.VerificationCode
	if err := db.Where(models.VerificationCode{ActivationCode: code}).First(&verCode).Error; err != nil {
		return nil, apperrors.NewNotFound(fmt.Sprintf("No activation code %s", code))
	}
	var user models.User
	if err := db.First(&user, verCode.UserID).Error; err != nil {
		return nil, apperrors.NewNotFound(fmt.Sprintf("No activation code %s", code))
	}

	if err := db.Transaction(func(tx *gorm.DB) error {

		// we first delete verCode because no matter if code is expired or not
		// it won't be needed outside the scope of this function
		if err := tx.Delete(&verCode).Error; err != nil {
			return err
		}

		// if verification code expired we delete created user and return not found error
		// pretending we don't know what the user wants ¯\_(ツ)_/¯
		if time.Since(verCode.Created) > db.Config.VerificationCodeDuration {
			if err := tx.Delete(&user).Error; err != nil {
				return err
			}
			return apperrors.NewNotFound(fmt.Sprintf("No activation code %s", code))
		}

		// if verification code is still valid we update the user `verified` property
		if err := tx.Model(&user).Update("verified", true).Error; err != nil {
			return err
		}
		return nil

	}); err != nil {
		return nil, err
	}
	return &user, nil

}
