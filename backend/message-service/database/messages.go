package database

import (
	"time"

	"github.com/Slimo300/MicroservicesChatApp/backend/message-service/models"
	"github.com/google/uuid"
)

func (db *Database) AddMessage(memberID uuid.UUID, nick string, text string, when time.Time) error {
	message := models.Message{
		ID:       uuid.New(),
		Text:     text,
		MemberID: memberID,
		Posted:   when,
		Nick:     nick,
	}

	if err := db.Create(&message).Error; err != nil {
		return err
	}

	return nil
}

func (db *Database) GetGroupMessages(groupID uuid.UUID, offset, num int) (messages []models.Message, err error) {
	return messages, db.Order("posted desc").Offset(offset).Limit(num).Where(models.Message{GroupID: groupID}).Find(&messages).Error
}

func (db *Database) IsUserInGroup(userID, groupID uuid.UUID) bool {
	var membership models.Membership
	err := db.Where(models.Membership{UserID: userID, GroupID: groupID}).First(&membership).Error
	if err != nil || membership.Deleted {
		return false
	}
	return true
}
