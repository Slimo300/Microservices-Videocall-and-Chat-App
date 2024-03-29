package database

import (
	"errors"
	"time"
)

type TokenValue string

const TOKEN_VALID TokenValue = "1"
const TOKEN_BLACKLISTED TokenValue = "2"

type TokenDB interface {
	SaveToken(token string, expiration time.Duration) error
	IsTokenValid(userID, tokenID string) (bool, error)
	InvalidateToken(userID, tokenID string) error
	InvalidateTokens(userID, tokenID string) error
}

var ErrTokenBlacklisted = errors.New("Token Blacklisted")
var ErrTokenNotFound = errors.New("Token not found")
var ErrUnexpectedTokenValue = errors.New("Unexpected token value")
