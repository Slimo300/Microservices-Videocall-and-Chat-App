package command

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/email"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/models"
)

type RegisterUser struct {
	Email    string
	Username string
	Password string
}

type RegisterUserHandler struct {
	repo         database.UsersRepository
	emailService email.EmailServiceClient
}

func NewRegisterUserHandler(repo database.UsersRepository, emailService email.EmailServiceClient) RegisterUserHandler {
	if repo == nil {
		panic("repo is nil")
	}
	if emailService == nil {
		panic("emailService is nil")
	}
	return RegisterUserHandler{repo: repo}
}

func (h RegisterUserHandler) Handle(ctx context.Context, cmd RegisterUser) error {
	user, err := models.NewUser(cmd.Email, cmd.Username, cmd.Password)
	if err != nil {
		return err
	}
	verificationCode := models.NewAuthorizationCode(user.ID(), models.EmailVerificationCode)
	if err := h.repo.RegisterUser(ctx, user, verificationCode); err != nil {
		return err
	}
	if _, err := h.emailService.SendVerificationEmail(context.TODO(), &email.EmailData{
		Email: user.Email(),
		Name:  user.Username(),
		Code:  verificationCode.Code().String(),
	}); err != nil {
		return err
	}
	return nil
}
