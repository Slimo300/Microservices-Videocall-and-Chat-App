package command

import (
	"context"
	"errors"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/apperrors"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/auth"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/database"
)

type SignIn struct {
	Email    string
	Password string
}

type SignInHandler struct {
	repo         database.UsersRepository
	tokenService auth.TokenServiceClient
}

func NewSignInHandler(repo database.UsersRepository, tokenService auth.TokenServiceClient) SignInHandler {
	if repo == nil {
		panic("repo is nil")
	}
	if tokenService == nil {
		panic("tokenService is nil")
	}
	return SignInHandler{repo: repo, tokenService: tokenService}
}

func (h SignInHandler) Handle(ctx context.Context, cmd SignIn) (string, string, error) {
	user, err := h.repo.GetUserByEmail(ctx, cmd.Email)
	if err != nil {
		return "", "", apperrors.NewNotFound("user not found")
	}
	if !user.CheckPassword(cmd.Password) {
		return "", "", apperrors.NewNotFound("user not found")
	}
	tokens, err := h.tokenService.NewPairFromUserID(ctx, &auth.UserID{ID: user.ID().String()})
	if err != nil {
		return "", "", err
	}
	if tokens.Error != "" {
		return "", "", errors.New(tokens.Error)
	}
	return tokens.AccessToken, tokens.RefreshToken, nil
}
