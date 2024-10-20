package query

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/apperrors"
	"github.com/google/uuid"
)

type GetUserGroups struct {
	UserID uuid.UUID
}

type GetUserGroupsHandler struct {
	repo database.GroupsRepository
}

func NewGetUserGroupsHandler(repo database.GroupsRepository) GetUserGroupsHandler {
	if repo == nil {
		panic("repo is nil")
	}
	return GetUserGroupsHandler{repo: repo}
}

func (h GetUserGroupsHandler) Handle(ctx context.Context, cmd GetUserGroups) ([]models.Group, error) {
	groups, err := h.repo.GetUserGroups(ctx, cmd.UserID)
	if err != nil {
		return nil, apperrors.NewNotFound("user not found")
	}
	return groups, nil
}
