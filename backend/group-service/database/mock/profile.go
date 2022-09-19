package mock

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (m *MockDB) SetPassword(userID uuid.UUID, password string) error {
	for _, user := range m.Users {
		if user.ID == userID {
			user.Pass = password
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

func (m *MockDB) DeleteProfilePicture(userID uuid.UUID) error {
	for _, user := range m.Users {
		if user.ID == userID {
			user.Picture = ""
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

func (m *MockDB) SetProfilePicture(userID uuid.UUID, newURI string) error {
	for _, user := range m.Users {
		if user.ID == userID {
			user.Picture = newURI
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

func (m *MockDB) GetProfilePictureURL(userID uuid.UUID) (string, error) {
	for _, user := range m.Users {
		if user.ID == userID {
			return user.Picture, nil
		}
	}
	return "", gorm.ErrRecordNotFound
}
