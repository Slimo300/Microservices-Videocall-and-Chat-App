package orm

import (
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/google/uuid"
)

func (db *Database) GetUserGroups(userID uuid.UUID) (groups []*models.Group, err error) {

	var userGroupsIDs []uuid.UUID
	if err := db.Table("`groups`").Select("`groups`.id").
		Joins("inner join `members` on `members`.group_id = `groups`.id").
		Joins("inner join `users` on `users`.id = `members`.user_id").
		Where("`users`.id = ?", userID).Scan(&userGroupsIDs).Error; err != nil {
		return nil, err
	}

	if err := db.Where("id in (?)", userGroupsIDs).Preload("Members").Preload("Members.User").Find(&groups).Error; err != nil {
		return groups, err
	}
	return groups, nil
}

func (db *Database) GetGroupByID(groupID uuid.UUID) (group *models.Group, err error) {
	return group, db.Preload("Members").Preload("Members.User").First(&group, groupID).Error
}

func (db *Database) CreateGroup(group *models.Group) (*models.Group, error) {
	return group, db.Preload("Members").Preload("Members.User").Create(&group).Error
}

func (db *Database) UpdateGroup(group *models.Group) (*models.Group, error) {
	return group, db.Preload("Members").Preload("Members.User").Model(&group).Updates(*group).Error
}

func (db *Database) DeleteGroup(groupID uuid.UUID) (group *models.Group, err error) {
	if err := db.Where(&models.Invite{GroupID: groupID}).Delete(&models.Invite{}).Error; err != nil {
		return nil, err
	}
	if err := db.Where(&models.Member{GroupID: groupID}).Delete(&models.Member{}).Error; err != nil {
		return nil, err
	}
	if err := db.Preload("Members").Preload("Members.User").First(&group, groupID).Error; err != nil {
		return nil, err
	}

	return group, db.Delete(&models.Group{}, groupID).Error
}
