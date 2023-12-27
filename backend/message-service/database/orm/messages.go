package orm

import (
	"errors"
	"fmt"
	"log"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/apperrors"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (db *Database) GetGroupMessages(userID, groupID uuid.UUID, offset, num int) (messages []models.Message, err error) {
	var membership models.Membership
	err = db.Where(models.Membership{UserID: userID, GroupID: groupID}).First(&membership).Error
	if err != nil {
		return []models.Message{}, apperrors.NewForbidden("User not in group")
	}
	if err := db.Order("posted desc").Offset(offset).Limit(num).Preload("Member").Preload("Files").Preload("Deleters").Where(models.Message{GroupID: groupID}).Find(&messages).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.Message{}, nil
		}
		return []models.Message{}, err
	}
	for i, msg := range messages {
		for _, del := range msg.Deleters {
			if del.UserID == userID {
				messages[i].Text = ""
				messages[i].Files = []models.MessageFile{}
			}
		}
	}
	return messages, nil
}

func (db *Database) DeleteMessageForEveryone(userID, messageID uuid.UUID) (*models.Message, error) {

	var message models.Message
	if err := db.Preload("Member").Preload("Files").First(&message, messageID).Error; err != nil {
		return nil, apperrors.NewNotFound(fmt.Sprintf("message with id %s not found", messageID.String()))
	}

	var membership models.Membership
	if err := db.Where(models.Membership{UserID: userID, GroupID: message.Member.GroupID}).First(&membership).Error; err != nil {
		log.Println(err)
		// Here we return not found not to give information about existance of message with given ID
		return nil, apperrors.NewNotFound(fmt.Sprintf("message with id %s not found", messageID.String()))
	}
	if !membership.CanDeleteMessage(&message) {
		return nil, apperrors.NewForbidden("User has no right to delete message")
	}

	if err := db.Model(&message).Update("text", "").Error; err != nil {
		return nil, err
	}
	if err := db.Where(models.MessageFile{MessageID: message.ID.String()}).Delete(&models.MessageFile{}).Error; err != nil {
		return nil, err
	}
	return &message, nil
}

func (db *Database) DeleteMessageForYourself(userID, messageID uuid.UUID) (*models.Message, error) {
	var message models.Message
	if err := db.Preload("Member").Preload("Deleters").First(&message, messageID).Error; err != nil {
		log.Println(err)
		return nil, apperrors.NewNotFound(fmt.Sprintf("message with id %s not found", messageID.String()))
	}

	var membership models.Membership
	if err := db.Where(models.Membership{UserID: userID, GroupID: message.Member.GroupID}).First(&membership).Error; err != nil {
		// Here we return not found not to give information about existance of message with given ID
		log.Println(err)
		return nil, apperrors.NewNotFound(fmt.Sprintf("message with id %s not found", messageID.String()))
	}
	// checking if user haven't already deleted this message
	for _, member := range message.Deleters {
		if member.UserID == userID {
			log.Println("message already deleted")
			return nil, apperrors.NewConflict(fmt.Sprintf("Message %v already deleted", messageID.String()))
		}
	}

	if err := db.Model(&message).Association("Deleters").Append(&membership); err != nil { // 500
		log.Println(err)
		return nil, err
	}
	log.Println(message)

	return &message, nil
}
