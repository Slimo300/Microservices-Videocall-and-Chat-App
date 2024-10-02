package models_test

import (
	"testing"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc          string
		email         string
		username      string
		password      string
		returnError   bool
		returnedError error
	}{
		{
			desc:          "invalid email",
			email:         "invalid",
			username:      "username",
			password:      "password",
			returnError:   true,
			returnedError: models.ErrInvalidEmail,
		},
		{
			desc:          "invalid username",
			email:         "host@net.pl",
			username:      "us",
			password:      "password",
			returnError:   true,
			returnedError: models.ErrInvalidUsername,
		},
		{
			desc:          "invalid password",
			email:         "host@net.pl",
			username:      "username",
			password:      "pass",
			returnError:   true,
			returnedError: models.ErrInvalidPassword,
		},
		{
			desc:          "success",
			email:         "host@net.pl",
			username:      "username",
			password:      "password",
			returnError:   false,
			returnedError: nil,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			user, err := models.NewUser(tC.email, tC.username, tC.password)
			if tC.returnError {
				assert.Equal(t, err, tC.returnedError)
			} else {
				assert.Equal(t, tC.email, user.Email())
				assert.Equal(t, tC.username, user.Username())
			}
		})
	}
}

func TestCheckPassword(t *testing.T) {
	t.Parallel()

	user, err := models.NewUser("host@net.pl", "username", "password")
	if err != nil {
		t.Errorf("error when creating new user: %v", err)
	}

	assert.False(t, user.CheckPassword("passworld"))
	assert.True(t, user.CheckPassword("password"))
}

func TestChangePictureIfStateIncorrect(t *testing.T) {
	user, err := models.NewUser("host@net.pl", "username", "password")
	if err != nil {
		t.Errorf("error when creating new user: %v", err)
	}

	require.False(t, user.HasPicture())

	assert.True(t, user.ChangePictureStateIfIncorrect(true))
	assert.True(t, user.HasPicture())

	assert.False(t, user.ChangePictureStateIfIncorrect(true))
	assert.True(t, user.HasPicture())

	assert.True(t, user.ChangePictureStateIfIncorrect(false))
	assert.False(t, user.HasPicture())

	assert.False(t, user.ChangePictureStateIfIncorrect(false))
	assert.False(t, user.HasPicture())
}
