package database

import (
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/events"
)

type DBLayer interface {
	GetMember(userID string, groupID string) (string, error)

	NewMember(evt events.MemberCreatedEvent) error
	DeleteMember(evt events.MemberDeletedEvent) error
	DeleteGroup(evt events.GroupDeletedEvent) error

	NewAccessCode(groupID, userID, accessCode string) error
	CheckAccessCode(accessCode string) (string, string, error)
}
