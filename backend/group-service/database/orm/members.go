package orm

import (
	"fmt"

	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/models"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/apperrors"
	"github.com/google/uuid"
)

func (db *Database) DeleteMember(userID, groupID, memberID uuid.UUID) (*models.Member, error) {
	var issuer models.Member
	if err := db.Where(models.Member{UserID: userID, GroupID: groupID}).First(&issuer).Error; err != nil {
		return nil, apperrors.NewForbidden(fmt.Sprintf("User %v has no right to delete members in group %v", userID, groupID))
	}
	var target models.Member
	if err := db.Where(models.Member{ID: memberID, GroupID: groupID}).First(&target).Error; err != nil {
		return nil, apperrors.NewNotFound("member", memberID.String())
	}
	if !issuer.CanDelete(target) {
		return nil, apperrors.NewForbidden(fmt.Sprintf("User %v cannot delete member %v", userID, memberID))
	}
	if err := db.Where(models.Member{ID: target.ID}).Delete(&models.Member{}).Error; err != nil {
		return nil, apperrors.NewInternal()
	}

	return &target, nil
}

func (db *Database) GrantRights(userID, groupID, memberID uuid.UUID, rights models.MemberRights) (*models.Member, error) {

	var issuer models.Member
	if err := db.Where(models.Member{UserID: userID, GroupID: groupID}).First(&issuer).Error; err != nil {
		return nil, apperrors.NewForbidden(fmt.Sprintf("User %v has no right to alter members in group %v", userID, groupID))
	}
	var target models.Member
	if err := db.Where(models.Member{ID: memberID, GroupID: groupID}).Preload("User").First(&target).Error; err != nil {
		return nil, apperrors.NewNotFound("member", memberID.String())
	}

	if !issuer.CanAlter(target) {
		return nil, apperrors.NewForbidden(fmt.Sprintf("User %v cannot alter member %v", userID, memberID))
	}

	if err := target.ApplyRights(rights); err != nil {
		return nil, apperrors.NewBadRequest(err.Error())
	}

	if err := db.Save(&target).Error; err != nil {
		return nil, apperrors.NewInternal()
	}
	return &target, nil
}
