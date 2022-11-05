package database

import (
	"time"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue/events"
	"github.com/Slimo300/MicroservicesChatApp/backend/message-service/models"
	"github.com/google/uuid"
)

type DBLayer interface {
	NewMember(event events.MemberCreatedEvent) error
	ModifyMember(event events.MemberUpdatedEvent) error
	DeleteMember(event events.MemberDeletedEvent) error

	GetGroupMessages(userID, groupID uuid.UUID, offset, num int) ([]models.Message, error)
	AddMessage(userID, groupID uuid.UUID, nick string, text string, when time.Time) error
	DeleteMessageForYourself(userID, messageID uuid.UUID) (models.Message, error)
	DeleteMessageForEveryone(userID, messageID uuid.UUID) (models.Message, error)
}
