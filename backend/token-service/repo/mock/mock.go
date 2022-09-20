package mock

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type mockTokenRepository struct {
	mock.Mock
}

func NewMockTokenRepository() *mockTokenRepository {
	return new(mockTokenRepository)
}

func (mock mockTokenRepository) SaveToken(token string, expiration time.Duration) error {
	ret := mock.Called(token, expiration)
	return ret.Error(0)
}

func (mock mockTokenRepository) IsTokenValid(userID, tokenID string) (bool, error) {
	ret := mock.Called(userID, tokenID)
	return ret.Get(0).(bool), ret.Error(1)
}

func (mock mockTokenRepository) InvalidateTokens(userID, tokenID string) error {
	ret := mock.Called(userID, tokenID)
	return ret.Error(0)
}

func (mock mockTokenRepository) InvalidateToken(userID, tokenID string) error {
	ret := mock.Called(userID, tokenID)
	return ret.Error(0)
}
