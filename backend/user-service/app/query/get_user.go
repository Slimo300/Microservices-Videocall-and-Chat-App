package query

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/apperrors"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/models"
	"github.com/google/uuid"
)

type GetUser struct {
	UserID uuid.UUID
}

type GetUserHandler struct {
	repo database.UsersRepository
}

func NewGetUserHandler(repo database.UsersRepository) GetUserHandler {
	if repo == nil {
		panic("repo is nil")
	}
	return GetUserHandler{repo: repo}
}

func (h GetUserHandler) Handle(ctx context.Context, cmd GetUser) (models.User, error) {
	user, err := h.repo.GetUserByID(ctx, cmd.UserID)
	if err != nil {
		return models.User{}, apperrors.NewNotFound("user not found")
	}
	return *user, nil
}
