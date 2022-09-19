package orm

import (
	"time"

	"github.com/Slimo300/MicrosevicesChatApp/backend/group-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (db *Database) GetUserGroups(id uuid.UUID) (groups []models.Group, err error) {

	var userGroupsIDs []uuid.UUID
	if err := db.Table("`groups`").Select("`groups`.id").
		Joins("inner join `members` on `members`.group_id = `groups`.id").
		Joins("inner join `users` on `users`.id = `members`.user_id").
		Where("`members`.deleted = false").
		Where("`users`.id = ?", id).Scan(&userGroupsIDs).Error; err != nil {
		return groups, err
	}

	if err := db.Where("id in (?)", userGroupsIDs).Preload("Members", "deleted is false").Preload("Members.User").Find(&groups).Error; err != nil {
		return groups, err
	}
	for _, group := range groups {
		for _, member := range group.Members {
			member.User.Pass = ""
		}
	}
	return groups, nil
}

func (db *Database) CreateGroup(userID uuid.UUID, name, desc string) (models.Group, error) {
	group := models.Group{ID: uuid.New(), Name: name, Desc: desc, Created: time.Now(), Picture: ""}

	var creator models.User
	if err := db.First(&creator, userID).Error; err != nil {
		return models.Group{}, err
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		creation := tx.Create(&group)
		if creation.Error != nil {
			return creation.Error
		}
		member := models.Member{ID: uuid.New(), UserID: userID, GroupID: group.ID, Adding: true, Deleting: true, Setting: true, Creator: true, Nick: creator.UserName}
		m_create := tx.Create(&member)
		if m_create.Error != nil {
			return m_create.Error
		}

		return nil
	}); err != nil {
		return models.Group{}, err
	}

	if err := db.Where(models.Group{ID: group.ID}).Preload("Members", "deleted is false").First(&group).Error; err != nil {
		return models.Group{}, err
	}
	return group, nil
}

func (db *Database) DeleteGroup(groupID uuid.UUID) (group models.Group, err error) {

	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where(models.Member{GroupID: groupID}).Delete(&models.Member{}).Error; err != nil {
			return err
		}
		group = models.Group{ID: groupID}
		if err := tx.Delete(&group).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return models.Group{}, err
	}

	return group, nil
}

func (db *Database) SetGroupProfilePicture(groupID uuid.UUID, newURI string) error {
	return db.First(&models.Group{}, groupID).Update("picture_url", newURI).Error
}

func (db *Database) DeleteGroupProfilePicture(groupID uuid.UUID) error {
	return db.First(&models.Group{}, groupID).Update("picture_url", "").Error
}

func (db *Database) GetGroupProfilePicture(groupID uuid.UUID) (string, error) {
	var group models.Group
	if err := db.First(&group, groupID).Error; err != nil {
		return "", err
	}
	return group.Picture, nil
}
