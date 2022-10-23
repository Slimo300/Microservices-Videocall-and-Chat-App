package server

import (
	"crypto/rsa"
	"time"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/auth/pb"
	"github.com/Slimo300/MicroservicesChatApp/backend/token-service/repo"
	"github.com/google/uuid"
)

type TokenService struct {
	*pb.UnimplementedTokenServiceServer
	iteration             uuid.UUID
	repo                  repo.TokenRepository
	refreshTokenSecret    string
	accessTokenPrivateKey rsa.PrivateKey
	accessTokenDuration   time.Duration
	refreshTokenDuration  time.Duration
}

func NewTokenService(repo repo.TokenRepository, refreshSecret string, accessPrivKey rsa.PrivateKey,
	refreshDuration, accessDuration time.Duration) *TokenService {

	return &TokenService{
		iteration:             uuid.New(),
		repo:                  repo,
		refreshTokenSecret:    refreshSecret,
		accessTokenPrivateKey: accessPrivKey,
		refreshTokenDuration:  refreshDuration,
		accessTokenDuration:   accessDuration,
	}
}
