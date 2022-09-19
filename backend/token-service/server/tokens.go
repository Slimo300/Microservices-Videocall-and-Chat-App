package server

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type RefreshTokenData struct {
	Token string
	ID    uuid.UUID
}

func (srv TokenService) GenerateRefreshToken(userID string) (*RefreshTokenData, error) {
	currentTime := time.Now()
	tokenExp := currentTime.Add(srv.refreshTokenDuration)

	tokenID, err := uuid.NewRandom()
	if err != nil {
		log.Println("Failed to generate refresh token ID")
		return nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		IssuedAt:  currentTime.Unix(),
		ExpiresAt: tokenExp.Unix(),
		Id:        tokenID.String(),
		Subject:   userID,
	})

	tokenString, err := token.SignedString([]byte(srv.refreshTokenSecret))
	if err != nil {
		log.Println("Failed to sign refresh token string: ", err)
		return nil, err
	}

	return &RefreshTokenData{
		Token: tokenString,
		ID:    tokenID,
	}, nil
}

func (srv TokenService) GenerateAccessToken(userID string) (string, error) {
	currentTime := time.Now()
	tokenExp := currentTime.Add(srv.accessTokenDuration)

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, &jwt.StandardClaims{
		IssuedAt:  currentTime.Unix(),
		ExpiresAt: tokenExp.Unix(),
		Subject:   userID,
	})

	tokenString, err := token.SignedString(&srv.accessTokenPrivateKey)
	if err != nil {
		log.Println("Failed to sign access token string: ", err)
		return "", err
	}

	return tokenString, err
}
