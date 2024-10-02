package command

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/models"
	"github.com/google/uuid"
)

type ResetForgottenPassword struct {
	Code        uuid.UUID
	NewPassword string
}

type ResetForgottenPasswordHandler struct {
	repo database.UsersRepository
}

func NewResetForgottenPasswordHandler(repo database.UsersRepository) ResetForgottenPasswordHandler {
	if repo == nil {
		panic("repo is nil")
	}
	return ResetForgottenPasswordHandler{repo: repo}
}

func (h ResetForgottenPasswordHandler) Handle(ctx context.Context, cmd ResetForgottenPassword) error {
	if err := h.repo.UpdateUserByCode(ctx, cmd.Code, models.ResetPasswordCode, func(u *models.User) (*models.User, error) {
		if err := u.SetPassword(cmd.NewPassword); err != nil {
			return nil, err
		}
		return u, nil
	}); err != nil {
		return err
	}
	return nil
}
