package database

import (
	"context"
	"errors"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/google/uuid"
)

var ErrUserNotInGroup = errors.New("user doesn't belong to group")

type GroupsRepository interface {
	GetGroupByID(ctx context.Context, userID, groupID uuid.UUID) (models.Group, error)
	GetMemberByID(ctx context.Context, userID, memberID uuid.UUID) (models.Member, error)
	GetInviteByID(ctx context.Context, userID, inviteID uuid.UUID) (models.Invite, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (models.User, error)

	GetUserGroups(ctx context.Context, userID uuid.UUID) ([]models.Group, error)
	GetUserInvites(ctx context.Context, userID uuid.UUID, num, offset int) ([]models.Invite, error)

	CreateGroup(ctx context.Context, group models.Group) (models.Group, error)
	UpdateGroup(ctx context.Context, userID, groupID uuid.UUID, updateFn func(g *models.Group) error) error
	DeleteGroup(ctx context.Context, userID, groupID uuid.UUID) error

	UpdateMember(ctx context.Context, userID, memberID uuid.UUID, updateFn func(m *models.Member) error) error
	DeleteMember(ctx context.Context, userID, memberID uuid.UUID) error

	CreateInvite(ctx context.Context, invite models.Invite) (models.Invite, error)
	UpdateInvite(ctx context.Context, inviteID uuid.UUID, updateFn func(i *models.Invite) (*models.Member, error)) error

	CreateUser(ctx context.Context, user models.User) error
	UpdateUser(ctx context.Context, userID uuid.UUID, updateFn func(*models.User) error) error
}
