package models

import "github.com/google/uuid"

type User struct {
	ID         uuid.UUID `json:"ID"`
	Username   string    `json:"username"`
	HasPicture bool      `json:"hasPicture"`
}
