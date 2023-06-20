package database

import (
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/events"
)

type DBLayer interface {
	IsUserMember(userID string, groupID string) bool

	NewMember(evt events.MemberCreatedEvent) error
	DeleteMember(evt events.MemberDeletedEvent) error
	DeleteGroup(evt events.GroupDeletedEvent) error

	// CheckGroupSession(groupID string) (bool, error)
	// AddGroupSession(groupID string) error
	// DeleteGroupSession(groupID string) error

	NewAccessCode(userID string, accessCode string) error
	CheckAccessCode(accessCode string) (string, error)
}
