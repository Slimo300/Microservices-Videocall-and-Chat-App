package database

import (
	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/models"
	"github.com/google/uuid"
)

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
