package server

import (
	"context"
	"crypto/x509"
	"errors"
	"fmt"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/auth/pb"
	"github.com/Slimo300/MicroservicesChatApp/backend/token-service/repo"
	"github.com/golang-jwt/jwt"
)

func (srv *TokenService) NewPairFromUserID(ctx context.Context, userID *pb.UserID) (*pb.TokenPair, error) {

	refreshData, err := srv.GenerateRefreshToken(userID.ID)
	if err != nil {
		return &pb.TokenPair{}, err
	}

	if err := srv.repo.SaveToken(fmt.Sprintf("%s:%s", userID.ID, refreshData.ID.String()), srv.refreshTokenDuration); err != nil {
		return &pb.TokenPair{}, err
	}

	access, err := srv.GenerateAccessToken(userID.ID)
	if err != nil {
		return &pb.TokenPair{}, err
	}
	return &pb.TokenPair{
		AccessToken:  access,
		RefreshToken: refreshData.Token,
	}, nil
}

func (srv *TokenService) NewPairFromRefresh(ctx context.Context, refresh *pb.RefreshToken) (*pb.TokenPair, error) {

	token, err := jwt.ParseWithClaims(refresh.GetToken(), &jwt.StandardClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(srv.refreshTokenSecret), nil
		})
	if err != nil {
		return &pb.TokenPair{
			Error: err.Error(),
		}, nil
	}
	userID := token.Claims.(*jwt.StandardClaims).Subject
	tokenID := token.Claims.(*jwt.StandardClaims).Id

	ok, err := srv.repo.IsTokenValid(userID, tokenID)
	if err != nil {
		if errors.Is(err, repo.TokenBlacklistedError) {
			if err := srv.repo.InvalidateTokens(userID, tokenID); err != nil {
				return &pb.TokenPair{
					Error: repo.TokenBlacklistedError.Error(),
				}, err
			}
			return &pb.TokenPair{
				Error: repo.TokenBlacklistedError.Error(),
			}, nil
		}
		return &pb.TokenPair{
			Error: err.Error(),
		}, nil
	}
	if !ok {
		return &pb.TokenPair{
			Error: "Invalid Token",
		}, nil
	}

	if err := srv.repo.InvalidateToken(userID, tokenID); err != nil {
		if errors.Is(err, repo.TokenNotFoundError) || errors.Is(err, repo.TooManyTokensFoundError) {
			return &pb.TokenPair{
				Error: err.Error(),
			}, nil
		}
		return &pb.TokenPair{}, err
	}

	refreshData, err := srv.GenerateRefreshToken(userID)
	if err != nil {
		return &pb.TokenPair{}, err
	}

	if err := srv.repo.SaveToken(fmt.Sprintf("%s:%s:%s", userID, tokenID, refreshData.ID.String()), srv.refreshTokenDuration); err != nil {
		return &pb.TokenPair{}, err
	}

	access, err := srv.GenerateAccessToken(userID)
	if err != nil {
		return &pb.TokenPair{}, err
	}
	return &pb.TokenPair{
		AccessToken:  access,
		RefreshToken: refreshData.Token,
	}, nil
}

func (srv *TokenService) DeleteUserToken(ctx context.Context, refresh *pb.RefreshToken) (*pb.Msg, error) {

	token, err := jwt.ParseWithClaims(refresh.GetToken(), &jwt.StandardClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(srv.refreshTokenSecret), nil
		})
	if err != nil {
		return &pb.Msg{}, err
	}
	userID := token.Claims.(*jwt.StandardClaims).Subject
	tokenID := token.Claims.(*jwt.StandardClaims).Id
	if err := srv.repo.InvalidateToken(userID, tokenID); err != nil {
		if errors.Is(err, repo.TokenNotFoundError) || errors.Is(err, repo.TooManyTokensFoundError) {
			return &pb.Msg{
				Error: err.Error(),
			}, nil
		}
		return &pb.Msg{}, err
	}
	return &pb.Msg{}, nil
}

func (srv *TokenService) GetPublicKey(ctx context.Context, empty *pb.Empty) (*pb.PublicKey, error) {

	pubKey, err := x509.MarshalPKIXPublicKey(&srv.accessTokenPrivateKey.PublicKey)
	if err != nil {
		return &pb.PublicKey{Error: err.Error()}, err
	}

	return &pb.PublicKey{
		PublicKey: pubKey,
	}, nil
}
