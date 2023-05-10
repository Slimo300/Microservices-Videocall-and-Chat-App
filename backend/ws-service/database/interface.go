package database

import (
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/events"
	"github.com/google/uuid"
)

type DBLayer interface {
	GetUserGroups(userID uuid.UUID) ([]uuid.UUID, error)

	NewMember(event events.MemberCreatedEvent) error
	DeleteMember(event events.MemberDeletedEvent) error
	DeleteGroup(event events.GroupDeletedEvent) error

	NewAccessCode(userID uuid.UUID, accessCode string) error
	CheckAccessCode(accessCode string) (uuid.UUID, error)
}
