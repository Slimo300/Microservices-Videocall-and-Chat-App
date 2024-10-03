package command

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/events"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/models"
)

type ForgotPassword struct {
	Email string
}

type ForgotPasswordHandler struct {
	repo    database.UsersRepository
	emitter msgqueue.EventEmiter
}

func NewForgotPasswordHandler(repo database.UsersRepository, emitter msgqueue.EventEmiter) ForgotPasswordHandler {
	if repo == nil {
		panic("repo is nil")
	}
	if emitter == nil {
		panic("emitter is nil")
	}
	return ForgotPasswordHandler{repo: repo, emitter: emitter}
}

func (h ForgotPasswordHandler) Handle(ctx context.Context, cmd ForgotPassword) error {
	user, err := h.repo.GetUserByEmail(ctx, cmd.Email)
	if err != nil {
		return err
	}
	verificationCode := models.NewAuthorizationCode(user.ID(), models.ResetPasswordCode)
	if err := h.repo.CreateAuthorizationCode(ctx, verificationCode); err != nil {
		return err
	}
	if err := h.emitter.Emit(events.UserForgotPasswordEvent{
		Email:    user.Email(),
		Username: user.Username(),
		Code:     verificationCode.Code().String(),
	}); err != nil {
		return err
	}
	return nil
}
