package service

import (
	"context"
	"mime/multipart"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/storage"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
	"github.com/google/uuid"
)

type Service interface {
	GetUserGroups(ctx context.Context, id uuid.UUID) ([]*models.Group, error)

	CreateGroup(ctx context.Context, userID uuid.UUID, name string) (*models.Group, error)
	DeleteGroup(ctx context.Context, userID, groupID uuid.UUID) (*models.Group, error)

	DeleteMember(ctx context.Context, userID, memberID uuid.UUID) (*models.Member, error)
	GrantRights(ctx context.Context, userID, memberID uuid.UUID, rights models.MemberRights) (*models.Member, error)

	SetGroupPicture(ctx context.Context, userID, groupID uuid.UUID, file multipart.File) (string, error)
	DeleteGroupPicture(ctx context.Context, userID, groupID uuid.UUID) error

	GetUserInvites(ctx context.Context, userID uuid.UUID, num, offset int) ([]*models.Invite, error)
	AddInvite(ctx context.Context, issID, targetID, groupID uuid.UUID) (*models.Invite, error)
	RespondInvite(ctx context.Context, userID, inviteID uuid.UUID, answer bool) (*models.Invite, *models.Group, error)
}

type GroupService struct {
	DB      database.GroupsRepository
	Storage storage.StorageLayer
	Emitter msgqueue.EventEmiter
}

func NewService(DB database.GroupsRepository, storage storage.StorageLayer, emiter msgqueue.EventEmiter) Service {
	return &GroupService{
		DB:      DB,
		Storage: storage,
		Emitter: emiter,
	}
}
