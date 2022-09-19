package database

import (
	"errors"
	"time"

	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/models"
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
