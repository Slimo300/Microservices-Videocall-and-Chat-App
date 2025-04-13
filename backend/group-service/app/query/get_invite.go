package query

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/apperrors"
	"github.com/google/uuid"
)

type GetInvite struct {
	UserID   uuid.UUID
	InviteID uuid.UUID
}

type GetInviteHandler struct {
	repo database.GroupsRepository
}

func NewGetInviteHandler(repo database.GroupsRepository) GetInviteHandler {
	if repo == nil {
		panic("repo is nil")
	}
	return GetInviteHandler{repo: repo}
}

func (h GetInviteHandler) Handle(ctx context.Context, cmd GetInvite) (models.Invite, error) {
	invite, err := h.repo.GetInviteByID(ctx, cmd.UserID, cmd.InviteID)
	if err != nil {
		return models.Invite{}, apperrors.NewNotFound("invite not found")
	}
	return invite, nil
}
