package models

import (
	"errors"
	"regexp"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

var ErrInvalidEmail = errors.New("invalid email")
var ErrInvalidUsername = errors.New("invalid username")
var ErrInvalidPassword = errors.New("invalid password")

type User struct {
	id         uuid.UUID
	userName   string
	email      string
	password   string
	hasPicture bool
	verified   bool
}

func (u User) ID() uuid.UUID        { return u.id }
func (u User) Username() string     { return u.userName }
func (u User) Email() string        { return u.email }
func (u User) PasswordHash() string { return u.password }
func (u User) HasPicture() bool     { return u.hasPicture }
func (u User) Verified() bool       { return u.verified }

func (u User) CheckPassword(incomingPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.password), []byte(incomingPassword)) == nil
}

func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.password = string(hashedPassword)
	return nil
}

// ChangePictureStateIfIncorrect checks whether passed state is different than User's hasPicture field and if it is, it updates
// User's field and returns information whether change was applied
func (u *User) ChangePictureStateIfIncorrect(state bool) bool {
	if u.hasPicture == state {
		return false
	}
	u.hasPicture = state
	return true
}

func (u *User) Verify() { u.verified = true }

func NewUser(email, username, password string) (*User, error) {
	if !validateEmail(email) {
		return nil, ErrInvalidEmail
	}
	if !validateUsername(username) {
		return nil, ErrInvalidUsername
	}
	if !validatePassword(password) {
		return nil, ErrInvalidPassword
	}
	user := &User{
		id:         uuid.New(),
		userName:   username,
		email:      email,
		hasPicture: false,
		verified:   false,
	}
	if err := user.SetPassword(password); err != nil {
		return nil, err
	}
	return user, nil
}

func MustNewUser(email, username, password string) *User {
	user, err := NewUser(email, username, password)
	if err != nil {
		panic(err)
	}
	return user
}

// UnmarshalUserFromDatabase should be only called from database layer, as it calls no logic to verify
// attributes inserted into User struct
func UnmarshalUserFromDatabase(userID uuid.UUID, username, email, password string, hasPicture, verified bool) *User {
	return &User{
		id:         userID,
		userName:   username,
		email:      email,
		password:   password,
		hasPicture: hasPicture,
		verified:   verified,
	}
}

func validatePassword(password string) bool {
	return len(password) >= 6
}

func validateUsername(username string) bool {
	return len(username) >= 3
}

func validateEmail(e string) bool {
	return emailRegex.MatchString(e)
}
