package command

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/email"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/models"
)

type ForgotPassword struct {
	Email string
}

type ForgotPasswordHandler struct {
	repo        database.UsersRepository
	emailClient email.EmailServiceClient
}

func NewForgotPasswordHandler(repo database.UsersRepository, emailClient email.EmailServiceClient) ForgotPasswordHandler {
	if repo == nil {
		panic("repo is nil")
	}
	if emailClient == nil {
		panic("emailClient is nil")
	}
	return ForgotPasswordHandler{repo: repo, emailClient: emailClient}
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
	if _, err := h.emailClient.SendResetPasswordEmail(ctx, &email.EmailData{
		Email: user.Email(),
		Name:  user.Username(),
		Code:  verificationCode.Code().String(),
	}); err != nil {
		return err
	}
	return nil
}
