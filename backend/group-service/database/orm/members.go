package orm

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/google/uuid"
)

func (r *GroupsGormRepository) GetMemberByID(ctx context.Context, memberID uuid.UUID) (member *models.Member, err error) {
	return member, r.db.WithContext(ctx).Preload("User").First(&member, memberID).Error
}

func (r *GroupsGormRepository) GetMemberByUserGroupID(ctx context.Context, userID, groupID uuid.UUID) (member *models.Member, err error) {
	return member, r.db.WithContext(ctx).Preload("User").Where(&models.Member{UserID: userID, GroupID: groupID}).First(&member).Error
}

func (r *GroupsGormRepository) CreateMember(ctx context.Context, member *models.Member) (*models.Member, error) {
	if err := r.db.WithContext(ctx).Create(&member).Error; err != nil {
		return nil, err
	}
	return member, r.db.WithContext(ctx).Preload("User").First(&member, member.ID).Error
}

func (r *GroupsGormRepository) UpdateMember(ctx context.Context, member *models.Member) (*models.Member, error) {
	return member, r.db.WithContext(ctx).Preload("User").Save(&member).Error
}

func (r *GroupsGormRepository) DeleteMember(ctx context.Context, memberID uuid.UUID) (member *models.Member, err error) {
	if err := r.db.WithContext(ctx).Preload("User").First(&member, memberID).Error; err != nil {
		return nil, err
	}
	return member, r.db.WithContext(ctx).Delete(&models.Member{}, memberID).Error
}
