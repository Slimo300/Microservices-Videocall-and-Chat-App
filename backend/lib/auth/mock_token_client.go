package auth

import (
	"crypto/rsa"

	"github.com/Slimo300/MicroservicesChatApp/backend/token-service/pb"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockTokenClient struct {
	mock.Mock
}

func NewMockTokenClient() *MockTokenClient {
	return new(MockTokenClient)
}

func (m MockTokenClient) NewPairFromUserID(userID uuid.UUID) (*pb.TokenPair, error) {
	ret := m.Called(userID)
	return ret.Get(0).(*pb.TokenPair), ret.Error(1)
}
func (m MockTokenClient) NewPairFromRefresh(refresh string) (*pb.TokenPair, error) {
	ret := m.Called(refresh)
	return ret.Get(0).(*pb.TokenPair), ret.Error(1)
}

func (m MockTokenClient) DeleteUserToken(refresh string) error {
	ret := m.Called(refresh)
	return ret.Error(0)
}

func (m MockTokenClient) GetPublicKey() *rsa.PublicKey {
	ret := m.Called()
	return ret.Get(0).(*rsa.PublicKey)
}
