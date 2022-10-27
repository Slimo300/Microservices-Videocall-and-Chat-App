package database

import "github.com/google/uuid"

type DBLayer interface {
	GetUserGroups(userID uuid.UUID) ([]uuid.UUID, error)
}
