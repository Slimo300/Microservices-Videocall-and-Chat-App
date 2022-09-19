package email

import "github.com/stretchr/testify/mock"

type MockEmailService struct {
	mock.Mock
}

func NewMockEmailService() *MockEmailService {
	return new(MockEmailService)
}

func (m MockEmailService) SendVerificationEmail(data VerificationEmailData) error {
	ret := m.Called(data)
	return ret.Error(0)
}
