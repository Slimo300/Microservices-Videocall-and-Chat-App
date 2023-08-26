package database

import (
	"github.com/Slimo300/MicroservicesChatApp/backend/webrtc-service/models"
)

type DBLayer interface {
	GetMember(memberID string) (*models.Member, error)

	NewMember(member models.Member) error
	DeleteMember(memberID string) error
	DeleteGroup(groupID string) error

	NewAccessCode(accessCode, memberID string) error
	CheckAccessCode(accessCode string) (string, error)
}
