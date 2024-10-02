package command

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/events"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/models"
	"github.com/google/uuid"
)

type VerifyEmail struct {
	Code uuid.UUID
}

type VerifyEmailHandler struct {
	repo   database.UsersRepository
	emiter msgqueue.EventEmiter
}

func NewVerifyEmailHandler(repo database.UsersRepository, emiter msgqueue.EventEmiter) VerifyEmailHandler {
	if repo == nil {
		panic("tokenService is nil")
	}
	if emiter == nil {
		panic("emiter is nil")
	}
	return VerifyEmailHandler{repo: repo, emiter: emiter}
}

func (h VerifyEmailHandler) Handle(ctx context.Context, cmd VerifyEmail) error {
	var user *models.User
	if err := h.repo.UpdateUserByCode(ctx, cmd.Code, models.EmailVerificationCode, func(u *models.User) (*models.User, error) {
		u.Verify()
		user = u
		return u, nil
	}); err != nil {
		return err
	}
	if err := h.emiter.Emit(events.UserRegisteredEvent{
		ID:       user.ID(),
		Username: user.Username(),
	}); err != nil {
		return err
	}
	return nil
}
