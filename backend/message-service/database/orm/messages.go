package orm

import (
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/apperrors"
	"github.com/Slimo300/MicroservicesChatApp/backend/message-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (db *Database) GetGroupMessages(userID, groupID uuid.UUID, offset, num int) (messages []models.Message, err error) {
	var membership models.Membership
	err = db.Where(models.Membership{UserID: userID, GroupID: groupID}).First(&membership).Error
	if err != nil {
		return []models.Message{}, apperrors.NewForbidden("User not in group")
	}
	if err := db.Order("posted desc").Offset(offset).Limit(num).Where(models.Message{GroupID: groupID}).Find(&messages).Error; err != nil {
		return []models.Message{}, apperrors.NewInternal()
	}
	for _, msg := range messages {
		for _, del := range msg.Deleters {
			if del.UserID == userID {
				msg.Text = ""
			}
		}
	}
	return messages, nil
}

func (db *Database) DeleteMessageForYourself(userID, messageID uuid.UUID) (models.Message, error) {
	var message models.Message
	if err := db.First(&message, messageID).Error; err != nil {
		return models.Message{}, apperrors.NewNotFound("message", messageID.String())
	}

	var membership models.Membership
	if err := db.Where(models.Membership{UserID: userID, GroupID: message.GroupID}).First(&membership).Error; err != nil {
		return models.Message{}, apperrors.NewForbidden("User not in group")
	}

	// checking if user haven't already deleted this message
	if err := db.Model(&message).Where(models.Membership{UserID: userID}).Association("Deleters").DB.First(&models.Membership{}).Error; err != gorm.ErrRecordNotFound {
		if err == nil {
			return models.Message{}, apperrors.NewForbidden("User already blacklisted this message")
		}
		return models.Message{}, apperrors.NewInternal()
	}
	if err := db.Model(&message).Association("Deleters").Append(membership); err != nil { // 500
		return models.Message{}, apperrors.NewInternal()
	}

	return message, nil
}

func (db *Database) DeleteMessageForEveryone(userID, messageID uuid.UUID) (models.Message, error) {
	var message models.Message
	if err := db.First(&message, messageID).Error; err != nil {
		return models.Message{}, apperrors.NewNotFound("message", messageID.String())
	}
	var membership models.Membership
	if err := db.Where(models.Membership{UserID: userID, GroupID: message.GroupID}).First(&membership).Error; err != nil {
		return models.Message{}, apperrors.NewForbidden(err.Error())
	}
	if !membership.Creator && !membership.DeletingMessages && message.UserID != userID {
		return models.Message{}, apperrors.NewForbidden("User has no right to delete message")
	}
	if err := db.Delete(&message).Error; err != nil {
		return models.Message{}, apperrors.NewInternal()
	}
	return message, nil
}
