package orm

import (
	"errors"
	"time"

	"github.com/Slimo300/MicrosevicesChatApp/backend/group-service/database"
	"github.com/Slimo300/MicrosevicesChatApp/backend/group-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (db *Database) AddInvite(issID, targetID, groupID uuid.UUID) (invite models.Invite, err error) {
	invite = models.Invite{ID: uuid.New(), IssId: issID, TargetID: targetID, GroupID: groupID, Status: database.INVITE_AWAITING, Created: time.Now(), Modified: time.Now()}
	return invite, db.Create(&invite).Error
}

func (db *Database) AcceptInvite(invite models.Invite) (group models.Group, err error) {

	if err = db.Transaction(func(tx *gorm.DB) error {
		if err := db.createOrRestoreMember(invite.TargetID, invite.GroupID); err != nil {
			return err
		}
		if err := db.First(&invite, invite.ID).Updates(models.Invite{Status: database.INVITE_ACCEPT, Modified: time.Now()}).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return models.Group{}, err
	}

	if err := db.Preload("Members", "deleted is false").First(&group, invite.GroupID).Error; err != nil {
		return models.Group{}, err
	}

	return group, nil
}

func (db *Database) GetUserInvites(userID uuid.UUID) (invites []models.Invite, err error) {
	return invites, db.Where(models.Invite{TargetID: userID, Status: database.INVITE_AWAITING}).Preload("Iss").Preload("Group").Find(&invites).Error
}

func (db *Database) GetInviteByID(inviteID uuid.UUID) (invite models.Invite, err error) {
	return invite, db.First(&invite, inviteID).Error
}

// helper for creating membership with id, it first find user to get his
// username and use it as member's nick
func (db *Database) createOrRestoreMember(userID, groupID uuid.UUID) error {

	if err := db.Where(models.Member{UserID: userID, GroupID: groupID, Deleted: true}).First(&models.Member{}).Update("deleted", false).Error; err == nil {
		return nil
	}

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return err
	}

	member := models.Member{ID: uuid.New(), GroupID: groupID, UserID: userID, Nick: user.UserName, Adding: false, Deleting: false, Setting: false, Creator: false, Deleted: false}
	if err := db.Create(&member).Error; err != nil {
		return err
	}

	return nil
}

func (db *Database) IsUserInvited(userID, groupID uuid.UUID) bool {
	return errors.Is(db.Where(models.Invite{GroupID: groupID, TargetID: userID, Status: database.INVITE_AWAITING}).First(&models.Invite{}).Error, nil)
}

func (db *Database) DeclineInvite(inviteID uuid.UUID) error {
	return db.First(&models.Invite{}, inviteID).Updates(models.Invite{Status: database.INVITE_DECLINE, Modified: time.Now()}).Error
}
