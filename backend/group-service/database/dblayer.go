package database

import (
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/google/uuid"
)

type DBLayer interface {
	// Group model functionality

	GetGroupByID(groupID uuid.UUID) (*models.Group, error)
	GetUserGroups(userID uuid.UUID) ([]*models.Group, error)

	CreateGroup(group *models.Group) (*models.Group, error)
	UpdateGroup(group *models.Group) (*models.Group, error)
	DeleteGroup(groupID uuid.UUID) (*models.Group, error)

	// Member model functionality

	GetMemberByID(memberID uuid.UUID) (*models.Member, error)
	GetMemberByUserGroupID(userID, groupID uuid.UUID) (*models.Member, error)

	CreateMember(member *models.Member) (*models.Member, error)
	UpdateMember(member *models.Member) (*models.Member, error)
	DeleteMember(memberID uuid.UUID) (*models.Member, error)

	// User model functionality

	GetUserByID(userID uuid.UUID) (*models.User, error)
	CreateUser(user *models.User) (*models.User, error)
	UpdateUser(user *models.User) (*models.User, error)
	DeleteUser(userID uuid.UUID) (*models.User, error)

	// Invite model functionality

	GetInviteByID(inviteID uuid.UUID) (*models.Invite, error)
	GetUserInvites(userID uuid.UUID, num, offset int) ([]*models.Invite, error)
	IsUserInvited(userID, groupID uuid.UUID) (bool, error)

	CreateInvite(invite *models.Invite) (*models.Invite, error)
	UpdateInvite(invite *models.Invite) (*models.Invite, error)
	DeleteInvite(inviteID uuid.UUID) (*models.Invite, error)
}
