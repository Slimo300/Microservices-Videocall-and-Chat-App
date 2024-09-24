package orm

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/google/uuid"
)

func (r *GroupsGormRepository) GetUserGroups(ctx context.Context, userID uuid.UUID) (groups []*models.Group, err error) {

	var userGroupsIDs []uuid.UUID
	if err := r.db.WithContext(ctx).Table("`groups`").Select("`groups`.id").
		Joins("inner join `members` on `members`.group_id = `groups`.id").
		Joins("inner join `users` on `users`.id = `members`.user_id").
		Where("`users`.id = ?", userID).Scan(&userGroupsIDs).Error; err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Where("id in (?)", userGroupsIDs).Preload("Members").Preload("Members.User").Find(&groups).Error; err != nil {
		return groups, err
	}
	return groups, nil
}

func (r *GroupsGormRepository) GetGroupByID(ctx context.Context, groupID uuid.UUID) (group *models.Group, err error) {
	return group, r.db.WithContext(ctx).Preload("Members").Preload("Members.User").First(&group, groupID).Error
}

func (r *GroupsGormRepository) CreateGroup(ctx context.Context, group *models.Group) (*models.Group, error) {
	if err := r.db.WithContext(ctx).Create(&group).Error; err != nil {
		return nil, err
	}
	return group, r.db.Preload("Members").Preload("Members.User").First(&group, group.ID).Error
}

func (r *GroupsGormRepository) UpdateGroup(ctx context.Context, group *models.Group) (*models.Group, error) {
	if err := r.db.WithContext(ctx).Model(&group).Updates(*group).Error; err != nil {
		return nil, err
	}
	return group, r.db.Preload("Members").Preload("Members.User").First(&group, group.ID).Error
}

func (r *GroupsGormRepository) DeleteGroup(ctx context.Context, groupID uuid.UUID) (group *models.Group, err error) {
	if err := r.db.WithContext(ctx).Where(&models.Invite{GroupID: groupID}).Delete(&models.Invite{}).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Where(&models.Member{GroupID: groupID}).Delete(&models.Member{}).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Preload("Members").Preload("Members.User").First(&group, groupID).Error; err != nil {
		return nil, err
	}

	return group, r.db.Delete(&models.Group{}, groupID).Error
}
