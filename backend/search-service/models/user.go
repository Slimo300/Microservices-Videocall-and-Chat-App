package models

import "github.com/google/uuid"

type User struct {
	ID         uuid.UUID `json:"ID"`
	Username   string    `json:"username"`
	PictureURL string    `json:"picture"`
}
