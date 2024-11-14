package orm

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *MessagesGormRepository) CreateMember(ctx context.Context, member models.Member) error {
	m := unmarshalMember(member)
	return r.DB.WithContext(ctx).Create(&m).Error
}

func (r *MessagesGormRepository) UpdateMember(ctx context.Context, memberID uuid.UUID, updateFn func(m *models.Member) bool) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var m Member
		if err := tx.First(&m, memberID).Error; err != nil {
			return err
		}
		member := models.UnmarshalMemberFromDatabase(m.ID, m.UserID, m.GroupID, m.Username, m.Creator, m.Admin, m.DeletingMessages)
		if !updateFn(&member) {
			return nil
		}
		m = unmarshalMember(member)
		return r.DB.Save(&m).Error
	})
}

func (r *MessagesGormRepository) DeleteMember(ctx context.Context, memberID uuid.UUID) error {
	return r.DB.WithContext(ctx).Delete(&Member{}, memberID).Error
}

func (r *MessagesGormRepository) DeleteGroup(ctx context.Context, groupID uuid.UUID) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var msgs []Message
		if err := tx.Where(&Message{GroupID: groupID}).Find(&msgs).Error; err != nil {
			return err
		}
		for _, message := range msgs {
			if err := tx.Model(&message).Association("Deleters").Clear(); err != nil {
				return err
			}
		}
		if err := tx.Where("message_id IN (?)", tx.Model(&Message{}).Select("id").Where("group_id = ?", groupID)).
			Delete(&MessageFile{}).Error; err != nil {
			return err
		}
		if err := tx.Where(&Message{GroupID: groupID}).Delete(&Message{}).Error; err != nil {
			return err
		}
		return tx.Where(&Member{GroupID: groupID}).Delete(&Member{}).Error
	})
}

func (r *MessagesGormRepository) GetUserGroupMember(ctx context.Context, userID, groupID uuid.UUID) (models.Member, error) {
	var m Member
	if err := r.DB.WithContext(ctx).Where(Member{UserID: userID, GroupID: groupID}).First(&m).Error; err != nil {
		return models.Member{}, err
	}
	return models.UnmarshalMemberFromDatabase(m.ID, m.UserID, m.GroupID, m.Username, m.Creator, m.Admin, m.DeletingMessages), nil
}
