package repo

import (
	"errors"
	"time"
)

type TokenValue string

func StringToTokenValue(s string) TokenValue {
	switch s {
	case "1":
		return TOKEN_VALID
	case "2":
		return TOKEN_BLACKLISTED
	default:
		panic("Inalid conversion to TokenValue from string")
	}
}

const TOKEN_VALID TokenValue = "1"
const TOKEN_BLACKLISTED TokenValue = "2"

type TokenRepository interface {
	SaveToken(token string, expiration time.Duration) error
	IsTokenValid(userID, tokenID string) (bool, error)
	InvalidateToken(userID, tokenID string) error
	InvalidateTokens(userID, tokenID string) error
}

var TokenBlacklistedError = errors.New("Token Blacklisted")
var TooManyTokensFoundError = errors.New("Too many tokens")
var TokenNotFoundError = errors.New("Token not found")
