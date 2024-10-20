package query

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/apperrors"
	"github.com/google/uuid"
)

type GetUserInvites struct {
	UserID uuid.UUID
	Num    int
	Offset int
}

type GetUserInvitesHandler struct {
	repo database.GroupsRepository
}

func NewGetUserInvitesHandler(repo database.GroupsRepository) GetUserInvitesHandler {
	if repo == nil {
		panic("repo is nil")
	}
	return GetUserInvitesHandler{repo: repo}
}

func (h GetUserInvitesHandler) Handle(ctx context.Context, cmd GetUserInvites) ([]models.Invite, error) {
	invites, err := h.repo.GetUserInvites(ctx, cmd.UserID, cmd.Num, cmd.Offset)
	if err != nil {
		return nil, apperrors.NewNotFound("user not found")
	}
	return invites, nil
}
