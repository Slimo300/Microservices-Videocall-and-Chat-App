package handlers

import (
	"crypto/rsa"
	"time"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/auth"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/token-service/database"
)

type TokenService struct {
	*auth.UnimplementedTokenServiceServer
	db                    database.TokenDB
	refreshTokenSecret    string
	accessTokenPrivateKey *rsa.PrivateKey
	accessTokenDuration   time.Duration
	refreshTokenDuration  time.Duration
}

// NewTokenService creates new token server
func NewTokenService(db database.TokenDB, privKey *rsa.PrivateKey, refreshSecret string,
	refreshDuration, accessDuration time.Duration) *TokenService {

	return &TokenService{
		db:                    db,
		refreshTokenSecret:    refreshSecret,
		accessTokenPrivateKey: privKey,
		refreshTokenDuration:  refreshDuration,
		accessTokenDuration:   accessDuration,
	}
}
