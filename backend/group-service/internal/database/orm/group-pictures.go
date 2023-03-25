package orm

import (
	"fmt"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/apperrors"
	"github.com/Slimo300/chat-groupservice/internal/models"
	"github.com/google/uuid"
)

func (db *Database) GetGroupProfilePictureURL(userID, groupID uuid.UUID) (string, error) {
	var member models.Member
	if err := db.Where(models.Member{UserID: userID, GroupID: groupID}).First(&member).Error; err != nil {
		return "", apperrors.NewForbidden(fmt.Sprintf("User %v has no rights to set in group %v", userID, groupID))
	}

	if !member.Admin {
		return "", apperrors.NewForbidden(fmt.Sprintf("User %v has no rights to set in group %v", userID, groupID))
	}

	var group models.Group
	if err := db.First(&group, groupID).Error; err != nil {
		// TODO: Error here is only possible if there would exist membership to unexisting group. This should be internal error
		return "", apperrors.NewForbidden(fmt.Sprintf("User %v has no rights to set in group %v", userID, groupID))
	}

	if group.Picture == "" {
		newPictureURL := uuid.NewString()
		if err := db.Model(&group).Update("picture_url", newPictureURL).Error; err != nil {
			return "", apperrors.NewInternal()
		}
		return newPictureURL, nil
	}

	return group.Picture, nil
}

func (db *Database) DeleteGroupProfilePicture(userID, groupID uuid.UUID) (string, error) {

	var member models.Member
	if err := db.Where(models.Member{UserID: userID, GroupID: groupID}).First(&member).Error; err != nil {
		return "", apperrors.NewForbidden(fmt.Sprintf("User %v has no rights to set in group %v", userID, groupID))
	}

	if !member.Admin {
		return "", apperrors.NewForbidden(fmt.Sprintf("User %v has no rights to set in group %v", userID, groupID))
	}

	var group models.Group
	// TODO: Error here is only possible if there would exist membership to unexisting group. This should be internal error
	if err := db.First(&group, groupID).Error; err != nil {
		return "", apperrors.NewNotFound("group", groupID.String())
	}

	if group.Picture == "" {
		return "", apperrors.NewForbidden(fmt.Sprintf("group %v has no profile picture", groupID))
	}

	if err := db.Model(&group).Update("picture_url", "").Error; err != nil {
		return "", apperrors.NewInternal()
	}
	return group.Picture, nil

}
