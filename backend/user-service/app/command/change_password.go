package command

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/apperrors"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/models"
	"github.com/google/uuid"
)

type ChangePassword struct {
	UserID      uuid.UUID
	OldPassword string
	NewPassword string
}

type ChangePasswordHandler struct {
	repo database.UsersRepository
}

func NewChangePasswordHandler(repo database.UsersRepository) ChangePasswordHandler {
	if repo == nil {
		panic("repo is nil")
	}
	return ChangePasswordHandler{repo: repo}
}

func (h ChangePasswordHandler) Handle(ctx context.Context, cmd ChangePassword) error {
	if err := h.repo.UpdateUserByID(ctx, cmd.UserID, func(u *models.User) (*models.User, error) {
		if !u.CheckPassword(cmd.OldPassword) {
			return nil, apperrors.NewForbidden("incorrect password")
		}
		if err := u.SetPassword(cmd.NewPassword); err != nil {
			return nil, err
		}
		return u, nil
	}); err != nil {
		return err
	}
	return nil
}
