package orm

import (
	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/models"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/events"
	"github.com/google/uuid"
)

func (db *Database) GetUser(userID uuid.UUID) (user models.User, err error) {
	return user, db.First(&user, userID).Error
}

func (db *Database) NewUser(event events.UserRegisteredEvent) error {
	return db.Create(&models.User{
		ID:       event.ID,
		UserName: event.Username,
		Picture:  event.PictureURL,
	}).Error
}

func (db *Database) UpdateUserProfilePictureURL(event events.UserPictureModifiedEvent) error {
	return db.Model(&models.User{ID: event.ID}).Update("picture", event.PictureURL).Error
}
