package email

import "github.com/stretchr/testify/mock"

type MockEmailService struct {
	mock.Mock
}

func NewMockEmailService() *MockEmailService {
	return new(MockEmailService)
}

func (m MockEmailService) SendVerificationEmail(data EmailData) error {
	ret := m.Called(data)
	return ret.Error(0)
}

func (m MockEmailService) SendResetPasswordEmail(data EmailData) error {
	ret := m.Called(data)
	return ret.Error(0)
}
