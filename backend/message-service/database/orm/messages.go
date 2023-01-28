package orm

import (
	"errors"

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
	if err := db.Order("posted desc").Offset(offset).Limit(num).Preload("Files").Preload("Deleters").Where(models.Message{GroupID: groupID}).Find(&messages).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.Message{}, nil
		}
		return []models.Message{}, apperrors.NewInternal()
	}
	for i, msg := range messages {
		for _, del := range msg.Deleters {
			if del.UserID == userID {
				messages[i].Text = ""
			}
		}
	}
	return messages, nil
}

func (db *Database) DeleteMessageForYourself(userID, messageID, groupID uuid.UUID) (models.Message, error) {
	// We first check if user is a member of group so we don't send him information that he has no right to see
	var membership models.Membership
	if err := db.Where(models.Membership{UserID: userID, GroupID: groupID}).First(&membership).Error; err != nil {
		return models.Message{}, apperrors.NewForbidden("User not in group")
	}

	// Here we find our message and if it belongs to a group user passed. Making user pass extra argument
	// groupID allows us to check if request sender is aware of message affiliations
	var message models.Message
	if err := db.Preload("Deleters").First(&message, messageID).Error; err != nil || message.GroupID != groupID {
		return models.Message{}, apperrors.NewNotFound("message", messageID.String())
	}

	// checking if user haven't already deleted this message
	for _, member := range message.Deleters {
		if member.UserID == userID {
			return models.Message{}, apperrors.NewConflict("deleted", userID.String())
		}
	}

	if err := db.Model(&message).Association("Deleters").Append(&membership); err != nil { // 500
		return models.Message{}, apperrors.NewInternal()
	}

	return message, nil
}

func (db *Database) DeleteMessageForEveryone(userID, messageID, groupID uuid.UUID) (models.Message, error) {
	var membership models.Membership
	if err := db.Where(models.Membership{UserID: userID, GroupID: groupID}).First(&membership).Error; err != nil {
		return models.Message{}, apperrors.NewForbidden("User not in group")
	}
	var message models.Message
	if err := db.First(&message, messageID).Error; err != nil {
		return models.Message{}, apperrors.NewNotFound("message", messageID.String())
	}
	if !membership.Creator && !membership.DeletingMessages && message.UserID != userID {
		return models.Message{}, apperrors.NewForbidden("User has no right to delete message")
	}
	if err := db.Model(&message).Update("text", "").Error; err != nil {
		return models.Message{}, apperrors.NewInternal()
	}
	return message, nil
}
