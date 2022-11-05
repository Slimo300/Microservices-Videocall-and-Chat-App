package orm

import (
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue/events"
	"github.com/Slimo300/MicroservicesChatApp/backend/message-service/models"
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
	return db.Delete(&models.Membership{MembershipID: event.ID}).Error
}
