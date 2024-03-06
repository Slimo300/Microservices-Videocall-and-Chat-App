package orm

import (
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/google/uuid"
)

func (db *Database) GetUserByID(userID uuid.UUID) (user *models.User, err error) {
	return user, db.First(&user, userID).Error
}

func (db *Database) CreateUser(user *models.User) (*models.User, error) {
	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}
	return user, db.First(&user, user.ID).Error
}

func (db *Database) UpdateUser(user *models.User) (*models.User, error) {
	return user, db.Model(&user).Updates(*user).Error
}

func (db *Database) DeleteUser(userID uuid.UUID) (user *models.User, err error) {
	if err := db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return user, db.Delete(&user, userID).Error
}
