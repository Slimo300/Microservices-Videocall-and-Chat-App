package orm

import (
	"context"
	"fmt"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/apperrors"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *MessagesGormRepository) CreateMessage(ctx context.Context, message models.Message) error {
	m := unmarshalMessage(message)
	return r.DB.WithContext(ctx).Create(&m).Error
}

func (r *MessagesGormRepository) GetMessageByID(ctx context.Context, userID, messageID uuid.UUID) (models.Message, error) {
	var m Message
	if err := r.DB.WithContext(ctx).Preload("Deleters").Preload("Files").First(&m, messageID).Error; err != nil {
		return models.Message{}, err
	}
	var mem Member
	if err := r.DB.WithContext(ctx).Where(&Member{GroupID: m.GroupID, UserID: userID}).First(&mem).Error; err != nil {
		return models.Message{}, err
	}
	return m.marshalMessage(), nil
}

func (r *MessagesGormRepository) GetGroupMessages(ctx context.Context, userID, groupID uuid.UUID, offset, num int) ([]models.Message, error) {
	var mem Member
	if err := r.DB.WithContext(ctx).Where(Member{UserID: userID, GroupID: groupID}).First(&mem).Error; err != nil {
		return nil, err
	}
	var messages []Message
	if err := r.DB.WithContext(ctx).Model(&Message{}).
		Preload("Member").
		Joins("LEFT JOIN users_who_deleted uwd ON messages.id = uwd.message_id AND uwd.member_id = ?", mem.ID).
		Where("messages.group_id = ?", groupID).
		Where("uwd.message_id IS NULL").
		Order("messages.posted DESC").
		Offset(offset).
		Limit(num).
		Find(&messages).Error; err != nil {
		return nil, err
	}

	var res []models.Message
	for _, m := range messages {
		res = append(res, m.marshalMessage())
	}
	return res, nil
}

func (r *MessagesGormRepository) DeleteMessageForYourself(ctx context.Context, userID, messageID uuid.UUID) error {
	var msg Message
	if err := r.DB.WithContext(ctx).Preload("Member").Preload("Deleters").First(&msg, messageID).Error; err != nil {
		return apperrors.NewNotFound(fmt.Sprintf("message with id %s not found", messageID.String()))
	}
	var mem Member
	if err := r.DB.WithContext(ctx).Where(Member{UserID: userID, GroupID: msg.GroupID}).First(&mem).Error; err != nil {
		// Here we return not found not to give information about existance of message with given ID
		return apperrors.NewNotFound(fmt.Sprintf("message with id %s not found", messageID.String()))
	}
	if msg.marshalMessage().IsUserDeleter(userID) {
		return apperrors.NewConflict(fmt.Sprintf("message %v already deleted", messageID.String()))
	}
	return r.DB.WithContext(ctx).Model(&msg).Association("Deleters").Append(&mem)
}

func (r *MessagesGormRepository) DeleteMessageForEveryone(ctx context.Context, userID, messageID uuid.UUID) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var message Message
		if err := tx.Preload("Member").Preload("Files").First(&message, messageID).Error; err != nil {
			return apperrors.NewNotFound(fmt.Sprintf("message with id %s not found", messageID.String()))
		}
		msg := message.marshalMessage()
		var member Member
		if err := tx.Where(Member{UserID: userID, GroupID: message.Member.GroupID}).First(&member).Error; err != nil {
			// Here we return not found not to give information about existance of message with given ID
			return apperrors.NewNotFound(fmt.Sprintf("message with id %s not found", messageID.String()))
		}
		m := models.UnmarshalMemberFromDatabase(member.ID, member.UserID, member.GroupID, member.Username, member.Creator, member.Admin, member.DeletingMessages)
		if !m.CanDeleteMessage(&msg) {
			return apperrors.NewForbidden("user has no right to delete message")
		}
		if err := tx.Where(MessageFile{MessageID: message.ID}).Delete(&MessageFile{}).Error; err != nil {
			return err
		}
		return tx.Delete(&Message{ID: messageID}).Error
	})
}
