package database

import (
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/events"
	"github.com/google/uuid"
)

type DBLayer interface {
	GetUserGroups(userID uuid.UUID) ([]uuid.UUID, error)

	NewMember(event events.MemberCreatedEvent) error
	DeleteMember(event events.MemberDeletedEvent) error
	DeleteGroupMembers(event events.GroupDeletedEvent) error
}
