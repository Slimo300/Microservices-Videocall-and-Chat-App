package orm

import (
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/events"
	"github.com/Slimo300/MicroservicesChatApp/backend/message-service/models"
	"github.com/google/uuid"
)

func (db *Database) NewMember(event events.MemberCreatedEvent) error {
	return db.Create(&models.Membership{
		MembershipID:     event.ID,
		GroupID:          event.GroupID,
		UserID:           event.UserID,
		Creator:          event.Creator,
		DeletingMessages: false,
	}).Error
}

func (db *Database) ModifyMember(event events.MemberUpdatedEvent) error {
	return db.Where(models.Membership{MembershipID: event.ID}).Update("deleting_messages", event.DeletingMessages).Error
}

func (db *Database) DeleteMember(event events.MemberDeletedEvent) error {
	return db.Delete(&models.Membership{}, event.ID).Error
}

func (db *Database) AddMessage(event events.MessageSentEvent) error {
	return db.Create(models.Message{
		ID:      uuid.New(),
		GroupID: event.GroupID,
		UserID:  event.UserID,
		Text:    event.Text,
		Nick:    event.Nick,
		Posted:  event.Posted,
	}).Error
}

func (db *Database) DeleteGroupMembers(event events.GroupDeletedEvent) error {
	return db.Where(models.Membership{GroupID: event.ID}).Delete(&models.Membership{}).Error
}
