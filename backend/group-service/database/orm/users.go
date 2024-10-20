package orm

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *GroupsGormRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (models.User, error) {
	var u User
	if err := r.db.WithContext(ctx).First(&u, userID).Error; err != nil {
		return models.User{}, err
	}
	return models.UnmarshalUserFromDatabase(u.ID, u.UserName, u.HasPicture), nil
}

func (r *GroupsGormRepository) CreateUser(ctx context.Context, user models.User) error {
	u := marshalUser(user)
	return r.db.WithContext(ctx).Create(&u).Error
}

func (r GroupsGormRepository) UpdateUser(ctx context.Context, userID uuid.UUID, updateFn func(u *models.User) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var u User
		if err := tx.First(&u, userID).Error; err != nil {
			return err
		}
		user := models.UnmarshalUserFromDatabase(u.ID, u.UserName, u.HasPicture)
		if err := updateFn(&user); err != nil {
			return err
		}
		u = marshalUser(user)
		return tx.Save(&u).Error
	})
}
