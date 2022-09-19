package auth

import (
	"crypto/rsa"

	"github.com/Slimo300/MicrosevicesChatApp/backend/token-service/pb"
	"github.com/google/uuid"
)

type TokenClient interface {
	NewPairFromUserID(userID uuid.UUID) (*pb.TokenPair, error)
	NewPairFromRefresh(refresh string) (*pb.TokenPair, error)
	DeleteUserToken(refresh string) error
	GetPublicKey() *rsa.PublicKey
}
