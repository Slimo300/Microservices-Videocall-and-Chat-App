package database

import (
	"errors"
	"time"

	"github.com/Slimo300/MicroservicesChatApp/backend/user-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (db *Database) GetUserById(id uuid.UUID) (user models.User, err error) {
	return user, db.First(&user, id).Error
}

func (db *Database) GetUserByEmail(email string) (user models.User, err error) {
	return user, db.Where(models.User{Email: email}).First(&user).Error
}

func (db *Database) GetUserByUsername(username string) (user models.User, err error) {
	return user, db.Where(models.User{UserName: username}).First(&user).Error
}

func (db *Database) RegisterUser(user models.User) (models.User, error) {
	user.ID = uuid.New()
	user.Active = time.Now()
	user.SignUp = time.Now()
	user.LoggedIn = false
	user.Picture = ""
	return user, db.Create(&user).Error
}

func (db *Database) IsEmailInDatabase(email string) bool {
	if err := db.Where(models.User{Email: email}).First(&models.User{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}

func (db *Database) IsUsernameInDatabase(username string) bool {
	if err := db.Where(models.User{UserName: username}).First(&models.User{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}

func (db *Database) SignInUser(id uuid.UUID) (err error) {
	return db.First(&models.User{}, id).Updates(models.User{LoggedIn: true, Active: time.Now()}).Error

}

func (db *Database) SignOutUser(id uuid.UUID) error {
	return db.First(&models.User{ID: id}).Updates(models.User{LoggedIn: false, Active: time.Now()}).Error
}

func (db *Database) SetPassword(userID uuid.UUID, password string) error {
	return db.First(&models.User{}, userID).Update("password", password).Error
}

func (db *Database) DeleteProfilePicture(userID uuid.UUID) error {
	return db.First(&models.User{}, userID).Update("picture", "").Error
}

func (db *Database) SetProfilePicture(userID uuid.UUID, newURI string) error {
	return db.First(&models.User{}, userID).Update("picture", newURI).Error
}

func (db *Database) GetProfilePictureURL(userID uuid.UUID) (string, error) {
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return "", err
	}
	return user.Picture, nil
}

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
