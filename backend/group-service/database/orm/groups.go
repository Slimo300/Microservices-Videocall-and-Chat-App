package orm

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	merrors "github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *GroupsGormRepository) GetGroupByID(ctx context.Context, userID, groupID uuid.UUID) (models.Group, error) {
	var g Group
	if err := r.db.WithContext(ctx).Preload("Members.User").First(&g, groupID).Error; err != nil {
		return models.Group{}, err
	}
	group := unmarshalGroup(g)
	_, ok := group.GetMemberByUserID(userID)
	if !ok {
		return models.Group{}, database.ErrUserNotInGroup
	}
	return group, nil
}

func (r *GroupsGormRepository) GetUserGroups(ctx context.Context, userID uuid.UUID) ([]models.Group, error) {
	var userGroupsIDs []uuid.UUID
	if err := r.db.WithContext(ctx).Table("`groups`").Select("`groups`.id").
		Joins("inner join `members` on `members`.group_id = `groups`.id").
		Joins("inner join `users` on `users`.id = `members`.user_id").
		Where("`users`.id = ?", userID).Scan(&userGroupsIDs).Error; err != nil {
		return nil, err
	}
	var gs []Group
	if err := r.db.WithContext(ctx).Where("id in (?)", userGroupsIDs).Preload("Members").Preload("Members.User").Find(&gs).Error; err != nil {
		return nil, err
	}
	return unmarshalGroups(gs), nil
}

func (r *GroupsGormRepository) CreateGroup(ctx context.Context, group models.Group) (returnedGroup models.Group, err error) {
	// We create group in transaction to ensure atomicity between creating group and member inside of it
	return returnedGroup, r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		g := marshalGroup(group)
		if err := tx.Create(&g).Error; err != nil {
			return err
		}
		if err := tx.Preload("Members.User").First(&g, group.ID()).Error; err != nil {
			return err
		}
		returnedGroup = unmarshalGroup(g)
		return nil
	})
}

func (r *GroupsGormRepository) UpdateGroup(ctx context.Context, groupID uuid.UUID, updateFn func(g *models.Group) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var g Group
		if err := tx.Preload("Members").First(&g, groupID).Error; err != nil {
			return err
		}
		group := unmarshalGroup(g)
		if err := updateFn(&group); err != nil {
			return err
		}
		g = marshalGroup(group)
		return tx.Save(&g).Error
	})
}

func (r *GroupsGormRepository) DeleteGroup(ctx context.Context, userID, groupID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var m Member
		if err := tx.Where(Member{UserID: userID, GroupID: groupID}).First(&m).Error; err != nil {
			return err
		}
		member := unmarshalMember(m)
		if !member.CanDeleteGroup() {
			return merrors.NewMemberUnauthorizedError(groupID.String(), merrors.UpdateGroupAction())
		}
		if err := tx.Where(Invite{GroupID: groupID}).Delete(&Invite{}).Error; err != nil {
			return err
		}
		if err := tx.Where(Member{GroupID: groupID}).Delete(&Member{}).Error; err != nil {
			return err
		}
		return tx.Delete(&Group{ID: groupID}).Error
	})
}
