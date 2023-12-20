package database

import (
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/events"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/models"
	"github.com/google/uuid"
)

type DBLayer interface {
	NewMember(event events.MemberCreatedEvent) error
	ModifyMember(event events.MemberUpdatedEvent) error
	DeleteMember(event events.MemberDeletedEvent) error
	DeleteGroupMembers(event events.GroupDeletedEvent) error

	AddMessage(event events.MessageSentEvent) error

	GetGroupMembership(userID, groupID uuid.UUID) (models.Membership, error)
	GetGroupMessages(userID, groupID uuid.UUID, offset, num int) ([]models.Message, error)
	DeleteMessageForYourself(userID, messageID uuid.UUID) (*models.Message, error)
	DeleteMessageForEveryone(userID, messageID uuid.UUID) (*models.Message, error)
}
