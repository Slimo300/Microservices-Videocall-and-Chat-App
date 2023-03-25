package database

import (
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/events"
	"github.com/Slimo300/chat-groupservice/internal/models"
	"github.com/google/uuid"
)

type DBLayer interface {
	GetUserGroups(id uuid.UUID) ([]models.Group, error)

	CreateGroup(userID uuid.UUID, name string) (models.Group, error)
	DeleteMember(userID, groupID, memberID uuid.UUID) (*models.Member, error)
	GrantRights(userID, groupID, memberID uuid.UUID, rights models.MemberRights) (*models.Member, error)
	DeleteGroup(userID, groupID uuid.UUID) (models.Group, error)

	GetGroupProfilePictureURL(userID, groupID uuid.UUID) (string, error)
	DeleteGroupProfilePicture(userID, groupID uuid.UUID) (string, error)

	GetUserInvites(userID uuid.UUID, num, offset int) ([]models.Invite, error)
	AddInvite(issID, targetID, groupID uuid.UUID) (*models.Invite, error)
	AnswerInvite(userID, inviteID uuid.UUID, answer bool) (*models.Invite, *models.Group, *models.Member, error)

	NewUser(event events.UserRegisteredEvent) error
	UpdateUserProfilePictureURL(event events.UserPictureModifiedEvent) error
}
