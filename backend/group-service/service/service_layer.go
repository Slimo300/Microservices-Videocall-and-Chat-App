package service

import (
	"mime/multipart"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/storage"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
	"github.com/google/uuid"
)

type ServiceLayer interface {
	GetUserGroups(id uuid.UUID) ([]*models.Group, error)

	CreateGroup(userID uuid.UUID, name string) (*models.Group, error)
	DeleteGroup(userID, groupID uuid.UUID) (*models.Group, error)

	DeleteMember(userID, memberID uuid.UUID) (*models.Member, error)
	GrantRights(userID, memberID uuid.UUID, rights models.MemberRights) (*models.Member, error)

	SetGroupPicture(userID, groupID uuid.UUID, file multipart.File) (string, error)
	DeleteGroupPicture(userID, groupID uuid.UUID) error

	GetUserInvites(userID uuid.UUID, num, offset int) ([]*models.Invite, error)
	AddInvite(issID, targetID, groupID uuid.UUID) (*models.Invite, error)
	RespondInvite(userID, inviteID uuid.UUID, answer bool) (*models.Invite, *models.Group, error)
}

type GroupService struct {
	DB      database.DBLayer
	Storage storage.StorageLayer
	Emitter msgqueue.EventEmiter
}

func NewService(DB database.DBLayer, storage storage.StorageLayer, emiter msgqueue.EventEmiter) ServiceLayer {
	return &GroupService{
		DB:      DB,
		Storage: storage,
		Emitter: emiter,
	}
}
