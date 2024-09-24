package orm_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestHandleMember(t *testing.T) {
	user, _ := db.CreateUser(context.Background(), &models.User{ID: uuid.New(), UserName: "1"})
	group, _ := db.CreateGroup(context.Background(), &models.Group{ID: uuid.New(), Name: "1", Created: time.Now()})

	t.Cleanup(func() {
		_, _ = db.DeleteUser(context.Background(), user.ID)
		_, _ = db.DeleteGroup(context.Background(), group.ID)
	})

	member, err := db.CreateMember(context.Background(), &models.Member{ID: uuid.New(), GroupID: group.ID, UserID: user.ID, Adding: true})
	if err != nil {
		t.Fatalf("Error creating a member: %v", err)
	}

	member, err = db.GetMemberByID(context.Background(), member.ID)
	if err != nil {
		t.Fatalf("Error getting member by ID: %v", err)
	}

	member.ApplyRights(models.MemberRights{Adding: false})

	member, err = db.UpdateMember(context.Background(), member)
	if err != nil {
		t.Fatalf("Error updating member: %v", err)
	}
	if member.User.UserName != "1" {
		t.Fatalf("Username after update should be \"1\", it is %s", member.User.UserName)
	}

	member, err = db.GetMemberByUserGroupID(context.Background(), user.ID, group.ID)
	if err != nil {
		t.Fatalf("Error getting member by its group and user IDs: %v", err)
	}

	if member.Adding {
		t.Fatalf("Member parameter adding should be false after update")
	}

	if _, err := db.DeleteMember(context.Background(), member.ID); err != nil {
		t.Fatalf("Error deleting member: %v", err)
	}

	if _, err := db.GetMemberByID(context.Background(), member.ID); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("Getting member should return not found error, instead got: %v", err)
	}
}
