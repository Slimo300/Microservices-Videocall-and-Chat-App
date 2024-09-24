package database

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/google/uuid"
)

type GroupsRepository interface {
	// Group model functionality

	GetGroupByID(ctx context.Context, groupID uuid.UUID) (*models.Group, error)
	GetUserGroups(ctx context.Context, userID uuid.UUID) ([]*models.Group, error)

	CreateGroup(ctx context.Context, group *models.Group) (*models.Group, error)
	UpdateGroup(ctx context.Context, group *models.Group) (*models.Group, error)
	DeleteGroup(ctx context.Context, groupID uuid.UUID) (*models.Group, error)

	// Member model functionality

	GetMemberByID(ctx context.Context, memberID uuid.UUID) (*models.Member, error)
	GetMemberByUserGroupID(ctx context.Context, userID, groupID uuid.UUID) (*models.Member, error)

	CreateMember(ctx context.Context, member *models.Member) (*models.Member, error)
	UpdateMember(ctx context.Context, member *models.Member) (*models.Member, error)
	DeleteMember(ctx context.Context, memberID uuid.UUID) (*models.Member, error)

	// User model functionality

	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) (*models.User, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) (*models.User, error)

	// Invite model functionality

	GetInviteByID(ctx context.Context, inviteID uuid.UUID) (*models.Invite, error)
	GetUserInvites(ctx context.Context, userID uuid.UUID, num, offset int) ([]*models.Invite, error)
	IsUserInvited(ctx context.Context, userID, groupID uuid.UUID) (bool, error)

	CreateInvite(ctx context.Context, invite *models.Invite) (*models.Invite, error)
	UpdateInvite(ctx context.Context, invite *models.Invite) (*models.Invite, error)
	DeleteInvite(ctx context.Context, inviteID uuid.UUID) (*models.Invite, error)

	BeginTransaction() (GroupsRepository, error)
	CommitTransaction() error
	RollbackTransaction() error
}
