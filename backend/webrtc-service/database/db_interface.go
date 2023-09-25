package database

import (
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/webrtc-service/models"
)

type DBLayer interface {
	GetMemberByID(memberID string) (*models.Member, error)
	GetMemberByGroupAndUserID(groupID, userID string) (*models.Member, error)

	NewMember(member models.Member) error
	DeleteMember(memberID string) error
	DeleteGroup(groupID string) error

	NewAccessCode(accessCode, memberID string) error
	CheckAccessCode(accessCode string) (string, error)
}
