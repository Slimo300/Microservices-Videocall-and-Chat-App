package database

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/models"
	"github.com/google/uuid"
)

type UsersRepository interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)

	RegisterUser(ctx context.Context, user *models.User, code *models.AuthorizationCode) error
	CreateAuthorizationCode(ctx context.Context, code *models.AuthorizationCode) error

	// UpdateUserByID passes updateFn to database to execute. It is a way to separate business logic from database implementation.
	// If updateFn returns no error, the transaction should end with error and rollback
	// If updateFn returns no error and no user, the transaction should be finished successfully without updating user
	UpdateUserByID(ctx context.Context, userID uuid.UUID, updateFn func(u *models.User) (*models.User, error)) error
	// UpdateUserByCode passes updateFn to database to execute. It is a way to separate business logic from database implementation.
	// If updateFn returns no error, the transaction should end with error and rollback
	// If updateFn returns no error and no user, the transaction should be finished successfully without updating user
	UpdateUserByCode(ctx context.Context, code uuid.UUID, codeType models.CodeType, updateFn func(u *models.User) (*models.User, error)) error

	DeleteUser(ctx context.Context, userID uuid.UUID) error
}
