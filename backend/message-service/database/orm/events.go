package orm

import (
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/events"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/models"
)

func (db *Database) NewMember(event events.MemberCreatedEvent) error {
	return db.Create(&models.Membership{
		MembershipID:     event.ID,
		GroupID:          event.GroupID,
		UserID:           event.UserID,
		Username:         event.User.UserName,
		Creator:          event.Creator,
		Admin:            event.Creator,
		DeletingMessages: event.Creator,
	}).Error
}

func (db *Database) ModifyMember(event events.MemberUpdatedEvent) error {
	return db.Model(&models.Membership{MembershipID: event.ID}).Updates(map[string]interface{}{"admin": event.Admin, "deleting_messages": event.DeletingMessages}).Error
}

func (db *Database) DeleteMember(event events.MemberDeletedEvent) error {
	return db.Delete(&models.Membership{}, event.ID).Error
}

func (db *Database) AddMessage(event events.MessageSentEvent) error {
	var files []models.MessageFile
	for _, f := range event.Files {
		files = append(files, models.MessageFile{Key: f.Key, Extention: f.Extension})
	}

	return db.Create(models.Message{
		ID:       event.ID,
		MemberID: event.MemberID,
		Text:     event.Text,
		Posted:   event.Posted,
		Files:    files,
	}).Error
}

func (db *Database) DeleteGroupMembers(event events.GroupDeletedEvent) error {
	return db.Where(models.Membership{GroupID: event.ID}).Delete(&models.Membership{}).Error
}
