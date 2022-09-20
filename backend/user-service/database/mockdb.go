package database

import (
	"github.com/Slimo300/MicroservicesChatApp/backend/user-service/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type mockUserRepository struct {
	mock.Mock
}

func NewMockUserRepository() *mockUserRepository {
	return new(mockUserRepository)
}

func (mock mockUserRepository) IsEmailInDatabase(email string) bool {
	return mock.Called(email).Bool(0)
}
func (mock mockUserRepository) IsUsernameInDatabase(username string) bool {
	return mock.Called(username).Bool(0)
}
func (mock mockUserRepository) GetUserById(uid uuid.UUID) (models.User, error) {
	ret := mock.Called(uid)
	return ret.Get(0).(models.User), ret.Error(1)
}
func (mock mockUserRepository) GetUserByEmail(email string) (models.User, error) {
	ret := mock.Called(email)
	return ret.Get(0).(models.User), ret.Error(1)
}
func (mock mockUserRepository) GetUserByUsername(username string) (models.User, error) {
	ret := mock.Called(username)
	return ret.Get(0).(models.User), ret.Error(1)
}
func (mock mockUserRepository) RegisterUser(user models.User) (models.User, error) {
	ret := mock.Called(user)
	return ret.Get(0).(models.User), ret.Error(1)
}
func (mock mockUserRepository) SignInUser(id uuid.UUID) error {
	return mock.Called(id).Error(0)
}
func (mock mockUserRepository) SignOutUser(id uuid.UUID) error {
	return mock.Called(id).Error(0)
}
func (mock mockUserRepository) SetPassword(userID uuid.UUID, password string) error {
	return mock.Called(userID, password).Error(0)
}
func (mock mockUserRepository) GetProfilePictureURL(userID uuid.UUID) (string, error) {
	ret := mock.Called(userID)
	return ret.String(0), ret.Error(1)
}
func (mock mockUserRepository) SetProfilePicture(userID uuid.UUID, newURI string) error {
	return mock.Called(userID, newURI).Error(0)
}
func (mock mockUserRepository) DeleteProfilePicture(userID uuid.UUID) error {
	return mock.Called(userID).Error(0)
}
func (mock mockUserRepository) NewVerificationCode(userID uuid.UUID, code string) (models.VerificationCode, error) {
	ret := mock.Called(userID, code)
	return ret.Get(0).(models.VerificationCode), ret.Error(1)
}
func (mock mockUserRepository) VerifyCode(userID uuid.UUID, code string) error {
	return mock.Called(userID).Error(0)
}
