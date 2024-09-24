package orm

import (
	"context"
	"errors"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *GroupsGormRepository) GetUserInvites(ctx context.Context, userID uuid.UUID, num, offset int) (invites []*models.Invite, err error) {
	return invites, r.db.WithContext(ctx).Order("modified DESC").Limit(num).Offset(offset).
		Where(models.Invite{TargetID: userID}).
		Or(models.Invite{IssId: userID}).
		Preload("Iss").Preload("Group").Preload("Target").Find(&invites).Error
}

func (r *GroupsGormRepository) GetInviteByID(ctx context.Context, inviteID uuid.UUID) (invite *models.Invite, err error) {
	return invite, r.db.WithContext(ctx).Preload("Iss").Preload("Group").Preload("Target").First(&invite, inviteID).Error
}

func (r *GroupsGormRepository) IsUserInvited(ctx context.Context, userID, groupID uuid.UUID) (bool, error) {
	var invite models.Invite
	if err := r.db.WithContext(ctx).Where(&models.Invite{TargetID: userID, GroupID: groupID, Status: models.INVITE_AWAITING}).First(&invite).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *GroupsGormRepository) CreateInvite(ctx context.Context, invite *models.Invite) (*models.Invite, error) {
	if err := r.db.WithContext(ctx).Create(&invite).Error; err != nil {
		return nil, err
	}
	return invite, r.db.WithContext(ctx).Preload("Iss").Preload("Group").Preload("Target").First(&invite, invite.ID).Error
}

func (r *GroupsGormRepository) UpdateInvite(ctx context.Context, invite *models.Invite) (*models.Invite, error) {
	if err := r.db.WithContext(ctx).Model(&invite).Updates(*invite).Error; err != nil {
		return nil, err
	}
	return invite, r.db.WithContext(ctx).Preload("Iss").Preload("Group").Preload("Target").First(&invite, invite.ID).Error
}

func (r *GroupsGormRepository) DeleteInvite(ctx context.Context, inviteID uuid.UUID) (invite *models.Invite, err error) {
	if err := r.db.WithContext(ctx).Preload("Iss").Preload("Group").Preload("Target").First(&invite, inviteID).Error; err != nil {
		return nil, err
	}
	return invite, r.db.WithContext(ctx).Delete(&models.Invite{}, inviteID).Error
}
