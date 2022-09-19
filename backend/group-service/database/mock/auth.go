package mock

import (
	"fmt"
	"time"

	"github.com/Slimo300/MicrosevicesChatApp/backend/group-service/models"
	"github.com/google/uuid"
)

func (m *MockDB) GetUserById(uid uuid.UUID) (models.User, error) {
	for _, user := range m.Users {
		if user.ID == uid {
			return user, nil
		}
	}
	return models.User{}, fmt.Errorf("No user with id: %v", uid)
}

func (m *MockDB) GetUserByEmail(email string) (models.User, error) {
	for _, user := range m.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return models.User{}, fmt.Errorf("No user with email: %v", email)
}

func (m *MockDB) GetUserByUsername(username string) (models.User, error) {
	for _, user := range m.Users {
		if user.UserName == username {
			return user, nil
		}
	}
	return models.User{}, fmt.Errorf("No user with email: %s", username)
}

func (m *MockDB) RegisterUser(user models.User) (models.User, error) {
	user.ID = uuid.New()
	user.Active = time.Now()
	user.SignUp = time.Now()
	user.LoggedIn = false
	m.Users = append(m.Users, user)
	return user, nil
}

func (m *MockDB) SignInUser(id uuid.UUID) error {
	for _, user := range m.Users {
		if user.ID == id {
			user.LoggedIn = true
			return nil
		}
	}
	return fmt.Errorf("No user with id: %d", id)
}

func (m *MockDB) SignOutUser(id uuid.UUID) error {
	for _, user := range m.Users {
		if user.ID == id {
			user.LoggedIn = false
			return nil
		}
	}
	return fmt.Errorf("No user with id: %v", id.String())
}

func (m *MockDB) IsEmailInDatabase(email string) bool {
	for _, user := range m.Users {
		if user.Email == email {
			return true
		}
	}
	return false
}

func (m *MockDB) IsUsernameInDatabase(username string) bool {
	for _, user := range m.Users {
		if user.UserName == username {
			return true
		}
	}
	return false
}
