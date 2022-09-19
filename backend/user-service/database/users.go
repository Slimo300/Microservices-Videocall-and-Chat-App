package database

import (
	"errors"

	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (db *Database) NewVerificationCode(userID uuid.UUID, code string) (verificationCode models.VerificationCode, err error) {
	verificationCode.UserID = userID
	verificationCode.ActivationCode = code
	return verificationCode, db.Create(&verificationCode).Error
}

func (db *Database) VerifyCode(userID uuid.UUID, code string) error {
	var verificationCode models.VerificationCode
	if err := db.First(&verificationCode, userID).Error; err != nil {
		return err
	}
	if verificationCode.ActivationCode != code {
		return errors.New("Invalid code")
	}
	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&verificationCode).Error; err != nil {
			return err
		}
		if err := tx.First(&models.User{}, userID).Update("activated", true).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// func (db *Database) RegisterUser(user models.User) (models.User, error) {
// 	user.ID = uuid.New()
// 	user.Active = time.Now()
// 	user.SignUp = time.Now()
// 	user.LoggedIn = false
// 	user.Picture = ""
// 	return user, db.Create(&user).Error
// }
