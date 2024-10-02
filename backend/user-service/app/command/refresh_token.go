package command

import (
	"context"
	"errors"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/auth"
)

var ErrTokenBlacklisted = errors.New("token blacklisted")

type RefreshToken struct {
	RefreshToken string
}

type RefreshTokenHandler struct {
	tokenService auth.TokenServiceClient
}

func NewRefreshTokenHandler(tokenService auth.TokenServiceClient) RefreshTokenHandler {
	if tokenService == nil {
		panic("tokenService is nil")
	}
	return RefreshTokenHandler{tokenService: tokenService}
}

func (h RefreshTokenHandler) Handle(ctx context.Context, cmd RefreshToken) (string, string, error) {
	tokens, err := h.tokenService.NewPairFromRefresh(context.TODO(), &auth.RefreshToken{Token: cmd.RefreshToken})
	if err != nil {
		return "", "", err
	}
	if tokens.Error == "Token Blacklisted" {
		return "", "", ErrTokenBlacklisted
	}
	if tokens.Error != "" {
		return "", "", errors.New(tokens.Error)
	}
	return tokens.AccessToken, tokens.RefreshToken, nil
}
