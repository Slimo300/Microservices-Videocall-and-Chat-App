package handlers

import (
	"context"
	"errors"
	"fmt"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/auth"
	"github.com/Slimo300/MicroservicesChatApp/backend/token-service/repo"
	"github.com/golang-jwt/jwt"
)

func (srv *TokenService) NewPairFromUserID(ctx context.Context, userID *auth.UserID) (*auth.TokenPair, error) {

	refreshData, err := srv.GenerateRefreshToken(userID.ID)
	if err != nil {
		return &auth.TokenPair{}, err
	}

	if err := srv.repo.SaveToken(fmt.Sprintf("%s:%s", userID.ID, refreshData.ID.String()), srv.refreshTokenDuration); err != nil {
		return &auth.TokenPair{}, err
	}

	access, err := srv.GenerateAccessToken(userID.ID)
	if err != nil {
		return &auth.TokenPair{}, err
	}
	return &auth.TokenPair{
		AccessToken:  access,
		RefreshToken: refreshData.Token,
	}, nil
}

func (srv *TokenService) NewPairFromRefresh(ctx context.Context, refresh *auth.RefreshToken) (*auth.TokenPair, error) {

	token, err := jwt.ParseWithClaims(refresh.GetToken(), &jwt.StandardClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(srv.refreshTokenSecret), nil
		})
	if err != nil {
		return &auth.TokenPair{
			Error: err.Error(),
		}, nil
	}
	userID := token.Claims.(*jwt.StandardClaims).Subject
	tokenID := token.Claims.(*jwt.StandardClaims).Id

	ok, err := srv.repo.IsTokenValid(userID, tokenID)
	if err != nil {
		if errors.Is(err, repo.TokenBlacklistedError) {
			if err := srv.repo.InvalidateTokens(userID, tokenID); err != nil {
				return &auth.TokenPair{
					Error: err.Error(),
				}, err
			}
		}
		return &auth.TokenPair{
			Error: err.Error(),
		}, nil
	}
	if !ok {
		return &auth.TokenPair{
			Error: "Invalid Token",
		}, nil
	}

	if err := srv.repo.InvalidateToken(userID, tokenID); err != nil {
		if errors.Is(err, repo.TokenNotFoundError) {
			return &auth.TokenPair{
				Error: err.Error(),
			}, nil
		}
		return &auth.TokenPair{}, err
	}

	refreshData, err := srv.GenerateRefreshToken(userID)
	if err != nil {
		return &auth.TokenPair{}, err
	}

	if err := srv.repo.SaveToken(fmt.Sprintf("%s:%s:%s", userID, tokenID, refreshData.ID.String()), srv.refreshTokenDuration); err != nil {
		return &auth.TokenPair{}, err
	}

	access, err := srv.GenerateAccessToken(userID)
	if err != nil {
		return &auth.TokenPair{}, err
	}

	return &auth.TokenPair{
		AccessToken:  access,
		RefreshToken: refreshData.Token,
	}, nil
}

func (srv *TokenService) DeleteUserToken(ctx context.Context, refresh *auth.RefreshToken) (*auth.Msg, error) {

	token, err := jwt.ParseWithClaims(refresh.GetToken(), &jwt.StandardClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(srv.refreshTokenSecret), nil
		})
	if err != nil {
		return &auth.Msg{}, err
	}

	userID := token.Claims.(*jwt.StandardClaims).Subject
	tokenID := token.Claims.(*jwt.StandardClaims).Id

	if err := srv.repo.InvalidateToken(userID, tokenID); err != nil {
		if errors.Is(err, repo.TokenNotFoundError) {
			return &auth.Msg{
				Error: err.Error(),
			}, nil
		}
		return &auth.Msg{}, err
	}
	return &auth.Msg{}, nil
}
