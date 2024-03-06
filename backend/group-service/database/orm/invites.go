package orm

import (
	"errors"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (db *Database) GetUserInvites(userID uuid.UUID, num, offset int) (invites []*models.Invite, err error) {
	return invites, db.Order("modified DESC").Limit(num).Offset(offset).
		Where(models.Invite{TargetID: userID}).
		Or(models.Invite{IssId: userID}).
		Preload("Iss").Preload("Group").Preload("Target").Find(&invites).Error
}

func (db *Database) GetInviteByID(inviteID uuid.UUID) (invite *models.Invite, err error) {
	return invite, db.Preload("Iss").Preload("Group").Preload("Target").First(&invite, inviteID).Error
}

func (db *Database) IsUserInvited(userID, groupID uuid.UUID) (bool, error) {

	var invite models.Invite
	if err := db.Where(&models.Invite{TargetID: userID, GroupID: groupID, Status: models.INVITE_AWAITING}).First(&invite).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (db *Database) CreateInvite(invite *models.Invite) (*models.Invite, error) {
	if err := db.Create(&invite).Error; err != nil {
		return nil, err
	}

	return invite, db.Preload("Iss").Preload("Group").Preload("Target").First(&invite, invite.ID).Error
}

func (db *Database) UpdateInvite(invite *models.Invite) (*models.Invite, error) {
	if err := db.Model(&invite).Updates(*invite).Error; err != nil {
		return nil, err
	}
	return invite, db.Preload("Iss").Preload("Group").Preload("Target").First(&invite, invite.ID).Error
}

func (db *Database) DeleteInvite(inviteID uuid.UUID) (invite *models.Invite, err error) {
	if err := db.Preload("Iss").Preload("Group").Preload("Target").First(&invite, inviteID).Error; err != nil {
		return nil, err
	}
	return invite, db.Delete(&models.Invite{}, inviteID).Error
}
