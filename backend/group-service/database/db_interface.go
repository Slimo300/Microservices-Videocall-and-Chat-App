package database

import (
	"errors"

	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/models"
	"github.com/google/uuid"
)

type DBlayer interface {
	GetUserByUsername(username string) (models.User, error)
	GetUserGroups(id uuid.UUID) ([]models.Group, error)

	GetMemberByID(memberID uuid.UUID) (models.Member, error)
	GetUserGroupMember(userID, groupID uuid.UUID) (models.Member, error)

	CreateGroup(userID uuid.UUID, name, desc string) (models.Group, error)
	DeleteUserFromGroup(memberID uuid.UUID) (models.Member, error)
	GrantPriv(memberID uuid.UUID, adding, deleting, setting bool) error

	DeleteGroup(groupID uuid.UUID) (models.Group, error)

	GetGroupProfilePicture(groupID uuid.UUID) (string, error)
	SetGroupProfilePicture(groupID uuid.UUID, newURI string) error
	DeleteGroupProfilePicture(groupID uuid.UUID) error

	GetUserInvites(userID uuid.UUID) ([]models.Invite, error)
	AddInvite(issID, targetID, groupID uuid.UUID) (models.Invite, error)

	DeclineInvite(inviteID uuid.UUID) error
	AcceptInvite(invite models.Invite) (models.Group, error)

	IsUserInGroup(userID, groupID uuid.UUID) bool
	IsUserInvited(userID, groupID uuid.UUID) bool

	GetInviteByID(inviteID uuid.UUID) (models.Invite, error)
}

const INVITE_AWAITING = 1
const INVITE_ACCEPT = 2
const INVITE_DECLINE = 3

var ErrINVALIDPASSWORD = errors.New("invalid password")
var ErrNoPrivilages = errors.New("insufficient privilages")
var ErrInternal = errors.New("internal transaction error")
