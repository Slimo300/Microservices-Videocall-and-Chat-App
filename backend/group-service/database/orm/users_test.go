package orm_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestHandleUser(t *testing.T) {
	user, err := db.CreateUser(context.Background(), &models.User{ID: uuid.New(), UserName: "1"})
	if err != nil {
		t.Fatalf("Error creating a user: %v", err)
	}

	user, err = db.GetUserByID(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("Error getting user by ID: %v", err)
	}

	user.HasPicture = true
	user, err = db.UpdateUser(context.Background(), user)
	if err != nil {
		t.Fatalf("Error updating user: %v", err)
	}

	user, err = db.GetUserByID(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("Error getting user by ID: %v", err)
	}
	if !user.HasPicture {
		t.Fatalf("User's picture parameter is not of expected value")
	}

	if _, err := db.DeleteUser(context.Background(), user.ID); err != nil {
		t.Fatalf("Error deleting user: %v", err)
	}

	if _, err = db.GetUserByID(context.Background(), user.ID); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("Getting user should return not found error, instead got: %v", err)
	}
}
