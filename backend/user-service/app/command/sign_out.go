package command

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/auth"
)

type SignOut struct {
	RefreshToken string
}

type SignOutHandler struct {
	tokenService auth.TokenServiceClient
}

func NewSignOutHandler(tokenService auth.TokenServiceClient) SignOutHandler {
	if tokenService == nil {
		panic("tokenService is nil")
	}
	return SignOutHandler{tokenService: tokenService}
}

func (h SignOutHandler) Handle(ctx context.Context, cmd SignOut) error {
	if _, err := h.tokenService.DeleteUserToken(context.TODO(), &auth.RefreshToken{Token: cmd.RefreshToken}); err != nil {
		return err
	}
	return nil
}
