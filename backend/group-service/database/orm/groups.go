package orm

import (
	"time"

	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/models"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/apperrors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (db *Database) GetUserGroups(id uuid.UUID) (groups []models.Group, err error) {

	var userGroupsIDs []uuid.UUID
	if err := db.Table("`groups`").Select("`groups`.id").
		Joins("inner join `members` on `members`.group_id = `groups`.id").
		Joins("inner join `users` on `users`.id = `members`.user_id").
		Where("`users`.id = ?", id).Scan(&userGroupsIDs).Error; err != nil {
		return groups, err
	}

	if err := db.Where("id in (?)", userGroupsIDs).Preload("Members").Preload("Members.User").Find(&groups).Error; err != nil {
		return groups, err
	}
	return groups, nil
}

func (db *Database) CreateGroup(userID uuid.UUID, name string) (models.Group, error) {
	group := models.Group{ID: uuid.New(), Name: name, Created: time.Now(), Picture: ""}

	var creator models.User
	if err := db.First(&creator, userID).Error; err != nil {
		return models.Group{}, err
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&group).Error; err != nil {
			return err
		}
		member := models.Member{ID: uuid.New(), UserID: userID, GroupID: group.ID, Adding: true, DeletingMembers: true, Admin: true, Creator: true}
		if err := tx.Create(&member).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return models.Group{}, err
	}

	if err := db.Where(models.Group{ID: group.ID}).Preload("Members").Preload("Members.User").First(&group).Error; err != nil {
		return models.Group{}, err
	}
	return group, nil
}

func (db *Database) DeleteGroup(userID, groupID uuid.UUID) (models.Group, error) {

	var member models.Member
	if err := db.Where(models.Member{UserID: userID, GroupID: groupID}).First(&member).Error; err != nil {
		return models.Group{}, apperrors.NewForbidden("User has no right to delete group")
	}
	if !member.Creator {
		return models.Group{}, apperrors.NewForbidden("User has no right to delete group")
	}

	var group models.Group
	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where(models.Member{GroupID: groupID}).Delete(&models.Member{}).Error; err != nil {
			return err
		}
		if err := tx.Where(models.Invite{GroupID: groupID}).Delete(&models.Invite{}).Error; err != nil {
			return err
		}
		group = models.Group{ID: groupID}
		if err := tx.Delete(&group).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return models.Group{}, apperrors.NewInternal()
	}

	return group, nil
}
