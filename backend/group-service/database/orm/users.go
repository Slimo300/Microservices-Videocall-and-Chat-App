package orm

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/google/uuid"
)

func (r *GroupsGormRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (user *models.User, err error) {
	return user, r.db.WithContext(ctx).First(&user, userID).Error
}

func (r *GroupsGormRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	if err := r.db.WithContext(ctx).Create(&user).Error; err != nil {
		return nil, err
	}
	return user, r.db.WithContext(ctx).First(&user, user.ID).Error
}

func (r *GroupsGormRepository) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	return user, r.db.WithContext(ctx).Model(&user).Updates(*user).Error
}

func (r *GroupsGormRepository) DeleteUser(ctx context.Context, userID uuid.UUID) (user *models.User, err error) {
	if err := r.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return nil, err
	}
	return user, r.db.WithContext(ctx).Delete(&user, userID).Error
}
