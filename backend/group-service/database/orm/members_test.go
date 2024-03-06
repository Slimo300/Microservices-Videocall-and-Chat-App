package orm_test

import (
	"errors"
	"testing"
	"time"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestHandleMember(t *testing.T) {

	user, _ := db.CreateUser(&models.User{ID: uuid.New(), UserName: "1"})
	group, _ := db.CreateGroup(&models.Group{ID: uuid.New(), Name: "1", Created: time.Now()})

	t.Cleanup(func() {
		_, _ = db.DeleteUser(user.ID)
		_, _ = db.DeleteGroup(group.ID)
	})

	member, err := db.CreateMember(&models.Member{ID: uuid.New(), GroupID: group.ID, UserID: user.ID, Adding: true})
	if err != nil {
		t.Fatalf("Error creating a member: %v", err)
	}

	member, err = db.GetMemberByID(member.ID)
	if err != nil {
		t.Fatalf("Error getting member by ID: %v", err)
	}

	if err := member.ApplyRights(models.MemberRights{
		Adding: models.REVOKE,
	}); err != nil {
		t.Fatalf("Error applying rights to member: %v", err)
	}
	_, err = db.UpdateMember(member)
	if err != nil {
		t.Fatalf("Error updating member: %v", err)
	}

	member, err = db.GetMemberByUserGroupID(user.ID, group.ID)
	if err != nil {
		t.Fatalf("Error getting member by its group and user IDs: %v", err)
	}

	if member.Adding {
		t.Fatalf("Member parameter adding should be false after update")
	}

	if _, err := db.DeleteMember(member.ID); err != nil {
		t.Fatalf("Error deleting member: %v", err)
	}

	if _, err := db.GetMemberByID(member.ID); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("Getting member should return not found error, instead got: %v", err)
	}
}
