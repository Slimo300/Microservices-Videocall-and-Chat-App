package handlers

import (
	"crypto/rand"
	"crypto/rsa"
	"time"

	"github.com/Slimo300/chat-tokenservice/internal/repo"
	"github.com/Slimo300/chat-tokenservice/pkg/client/pb"
	"github.com/go-redis/redis"
)

type TokenService struct {
	*pb.UnimplementedTokenServiceServer
	repo                  repo.TokenRepository
	refreshTokenSecret    string
	accessTokenPrivateKey rsa.PrivateKey
	accessTokenDuration   time.Duration
	refreshTokenDuration  time.Duration
}

// NewTokenService creates new token server
func NewTokenService(repo repo.TokenRepository, refreshSecret string,
	refreshDuration, accessDuration time.Duration) (*TokenService, error) {

	privKey, err := repo.GetPrivateKey()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if err == redis.Nil {
		privKey, err = rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return nil, err
		}

		if err = repo.SetPrivateKey(privKey); err != nil {
			return nil, err
		}
	}

	return &TokenService{
		repo:                  repo,
		refreshTokenSecret:    refreshSecret,
		accessTokenPrivateKey: *privKey,
		refreshTokenDuration:  refreshDuration,
		accessTokenDuration:   accessDuration,
	}, nil
}
