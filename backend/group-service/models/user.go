package models

import "github.com/google/uuid"

type User struct {
	id         uuid.UUID
	username   string
	hasPicture bool
}

func (u User) ID() uuid.UUID    { return u.id }
func (u User) Username() string { return u.username }
func (u User) HasPicture() bool { return u.hasPicture }

func (u *User) UpdatePictureState(hasPicture bool) {
	u.hasPicture = hasPicture
}

func NewUser(userID uuid.UUID, username string) User {
	return User{
		id:         userID,
		username:   username,
		hasPicture: false,
	}
}

func UnmarshalUserFromDatabase(id uuid.UUID, username string, hasPicture bool) User {
	return User{
		id:         id,
		username:   username,
		hasPicture: hasPicture,
	}
}
