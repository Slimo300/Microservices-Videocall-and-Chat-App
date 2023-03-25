package orm

import (
	"fmt"
	"time"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/apperrors"
	"github.com/Slimo300/chat-groupservice/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (db *Database) GetUserInvites(userID uuid.UUID, num, offset int) (invites []models.Invite, err error) {
	return invites, db.Order("modified DESC").Limit(num).Offset(offset).
		Where(models.Invite{TargetID: userID}).
		Or(models.Invite{IssId: userID}).
		Preload("Iss").Preload("Group").Preload("Target").Find(&invites).Error
}

func (db *Database) AddInvite(issID, targetID, groupID uuid.UUID) (*models.Invite, error) {

	var member models.Member
	if err := db.Where(models.Member{UserID: issID, GroupID: groupID}).First(&member).Error; err != nil {
		return nil, apperrors.NewForbidden(fmt.Sprintf("User %v has no rights to add new members to group %v", issID, groupID))
	}
	if !member.Adding && !member.Admin && !member.Creator {
		return nil, apperrors.NewForbidden(fmt.Sprintf("User %v has no rights to add new members to group %v", issID, groupID))
	}
	if err := db.First(&models.User{}, targetID).Error; err != nil {
		return nil, apperrors.NewNotFound("user", targetID.String())
	}
	if err := db.Where(models.Member{UserID: targetID, GroupID: groupID}).First(&models.Member{}).Error; err != gorm.ErrRecordNotFound {
		return nil, apperrors.NewForbidden(fmt.Sprintf("User %v already is already a member of group %v", targetID, groupID))
	}
	if err := db.Where(models.Invite{GroupID: groupID, TargetID: targetID, Status: models.INVITE_AWAITING}).First(&models.Invite{}).Error; err != gorm.ErrRecordNotFound {
		return nil, apperrors.NewForbidden(fmt.Sprintf("User %v already invited to group %v", targetID, groupID))
	}
	invite := models.Invite{ID: uuid.New(), IssId: issID, TargetID: targetID, GroupID: groupID, Status: models.INVITE_AWAITING, Created: time.Now(), Modified: time.Now()}
	if err := db.Create(&invite).Error; err != nil {
		return nil, apperrors.NewInternal()
	}
	if err := db.Where(models.Invite{ID: invite.ID}).Preload("Iss").Preload("Group").Preload("Target").First(&invite).Error; err != nil {
		return nil, apperrors.NewInternal()
	}
	return &invite, nil
}

// AnswerInvite is a method for updating database after invite response providing that user has rights to answer the invite
// and answered invite is actually awaiting a response. Method returns updated Invite object, also Group and Member objects if user
// accepted invite. If user declined invite it will return nil.
func (db *Database) AnswerInvite(userID, inviteID uuid.UUID, answer bool) (*models.Invite, *models.Group, *models.Member, error) {
	// First we check if invite with provided ID exists in our database, then we check if user who answers it is
	// actually the one who was invited and if invite is waiting for response

	var invite models.Invite
	if err := db.Where(models.Invite{ID: inviteID}).Preload("Iss").Preload("Group").Preload("Target").First(&invite).Error; err != nil {
		return nil, nil, nil, apperrors.NewNotFound("invite", inviteID.String())
	}
	if invite.TargetID != userID {
		return nil, nil, nil, apperrors.NewNotFound("invite", inviteID.String())
	}
	if invite.Status != models.INVITE_AWAITING {
		return nil, nil, nil, apperrors.NewForbidden("invite already answered")
	}

	// if invite is declined we return just an invite with empty member as none was created
	if !answer {
		if err := db.Model(&invite).Updates(models.Invite{Status: models.INVITE_DECLINE, Modified: time.Now()}).Error; err != nil {
			return nil, nil, nil, apperrors.NewInternal()
		}
		return &invite, nil, nil, nil
	}

	memberID := uuid.New()
	// if invite is accepted we update invite status and create a new membership entry in our database
	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&models.Invite{}, inviteID).Updates(models.Invite{Status: models.INVITE_ACCEPT, Modified: time.Now()}).Error; err != nil {
			return err
		}
		if err := tx.Create(&models.Member{ID: memberID, UserID: userID, GroupID: invite.GroupID}).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, nil, nil, apperrors.NewInternal()
	}

	var member models.Member
	if err := db.Where(models.Member{ID: memberID}).Preload("User").First(&member, memberID).Error; err != nil {
		return nil, nil, nil, apperrors.NewInternal()
	}

	var group models.Group
	if err := db.Where(models.Group{ID: invite.GroupID}).Preload("Members").Preload("Members.User").First(&group, invite.GroupID).Error; err != nil {
		return nil, nil, nil, apperrors.NewInternal()
	}

	if err := db.Where(models.Invite{ID: inviteID}).Preload("Iss").Preload("Group").Preload("Target").First(&invite).Error; err != nil {
		return nil, nil, nil, apperrors.NewInternal()
	}

	return &invite, &group, &member, nil
}
