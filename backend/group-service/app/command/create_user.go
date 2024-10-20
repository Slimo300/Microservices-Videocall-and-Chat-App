package command

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/google/uuid"
)

type CreateUserCommand struct {
	UserID   uuid.UUID
	Username string
}

type CreateUserHandler struct {
	repo database.GroupsRepository
}

func NewCreateUserHandler(repo database.GroupsRepository) CreateUserHandler {
	if repo == nil {
		panic("repo is nil")
	}
	return CreateUserHandler{repo: repo}
}

func (h CreateUserHandler) Handle(ctx context.Context, cmd CreateUserCommand) error {
	user := models.NewUser(cmd.UserID, cmd.Username)
	return h.repo.CreateUser(context.Background(), user)
}
