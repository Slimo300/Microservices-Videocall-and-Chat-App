package query

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/apperrors"
	"github.com/google/uuid"
)

type GetGroup struct {
	UserID  uuid.UUID
	GroupID uuid.UUID
}

type GetGroupHandler struct {
	repo database.GroupsRepository
}

func NewGetGroupHandler(repo database.GroupsRepository) GetGroupHandler {
	if repo == nil {
		panic("repo is nil")
	}
	return GetGroupHandler{repo: repo}
}

func (h GetGroupHandler) Handle(ctx context.Context, cmd GetGroup) (models.Group, error) {
	groups, err := h.repo.GetGroupByID(ctx, cmd.UserID, cmd.GroupID)
	if err != nil {
		return models.Group{}, apperrors.NewNotFound("user not found")
	}
	return groups, nil
}
