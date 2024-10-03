package command

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/apperrors"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/events"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/models"
)

type RegisterUser struct {
	Email    string
	Username string
	Password string
}

type RegisterUserHandler struct {
	repo    database.UsersRepository
	emitter msgqueue.EventEmiter
}

func NewRegisterUserHandler(repo database.UsersRepository, emitter msgqueue.EventEmiter) RegisterUserHandler {
	if repo == nil {
		panic("repo is nil")
	}
	if emitter == nil {
		panic("emitter is nil")
	}
	return RegisterUserHandler{repo: repo, emitter: emitter}
}

func (h RegisterUserHandler) Handle(ctx context.Context, cmd RegisterUser) error {
	user, err := models.NewUser(cmd.Email, cmd.Username, cmd.Password)
	if err != nil {
		return apperrors.NewInternal(err)
	}
	verificationCode := models.NewAuthorizationCode(user.ID(), models.EmailVerificationCode)
	if err := h.repo.RegisterUser(ctx, user, verificationCode); err != nil {
		return err
	}
	if err := h.emitter.Emit(events.UserRegisteredEvent{
		Email:    user.Email(),
		Username: user.Username(),
		Code:     verificationCode.Code().String(),
	}); err != nil {
		return err
	}
	return nil
}
