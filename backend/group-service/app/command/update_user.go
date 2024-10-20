package command

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/google/uuid"
)

type UpdateUserCommand struct {
	UserID     uuid.UUID
	HasPicture bool
}

type UpdateUserHandler struct {
	repo database.GroupsRepository
}

func NewUpdateUserHandler(repo database.GroupsRepository) UpdateUserHandler {
	if repo == nil {
		panic("repo is nil")
	}
	return UpdateUserHandler{repo: repo}
}

func (h UpdateUserHandler) Handle(ctx context.Context, cmd UpdateUserCommand) error {
	return h.repo.UpdateUser(context.Background(), cmd.UserID, func(u *models.User) error {
		u.UpdatePictureState(cmd.HasPicture)
		return nil
	})
}
