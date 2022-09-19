package orm

import (
	"time"

	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/models"
	"github.com/google/uuid"
)

func (db *Database) AddMessage(memberID uuid.UUID, text string, when time.Time) error {
	message := models.Message{
		ID:       uuid.New(),
		Text:     text,
		MemberID: memberID,
		Posted:   when,
	}

	if err := db.Create(&message).Error; err != nil {
		return err
	}

	return nil
}

func (db *Database) GetGroupMessages(groupID uuid.UUID, offset, num int) (messages []models.Message, err error) {
	return messages, db.Joins("Member", db.Where(&models.Member{GroupID: groupID})).Order("posted desc").Offset(offset).Limit(num).
		Find(&messages, "`Member`.`group_id` = ?", groupID).Error
}
