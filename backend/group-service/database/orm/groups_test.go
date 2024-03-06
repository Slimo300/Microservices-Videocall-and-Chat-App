package orm_test

import (
	"errors"
	"testing"
	"time"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestHandleGroup(t *testing.T) {
	group, err := db.CreateGroup(&models.Group{
		ID:      uuid.New(),
		Name:    "New Group",
		Created: time.Now(),
	})
	if err != nil {
		t.Fatalf("Creating group failed with error: %v", err)
	}

	group, err = db.GetGroupByID(group.ID)
	if err != nil {
		t.Fatalf("Error getting newly created group: %v", err)
	}
	if group.Name != "New Group" {
		t.Fatal("group name not correct")
	}

	user, _ := db.CreateUser(&models.User{ID: uuid.New(), UserName: "SLimo"})
	_, _ = db.CreateMember(&models.Member{ID: uuid.New(), GroupID: group.ID, UserID: user.ID})

	group, err = db.UpdateGroup(&models.Group{
		ID:      group.ID,
		Picture: "Picture",
	})
	if err != nil {
		t.Fatalf("Error updating group: %v", err)
	}
	if len(group.Members) == 0 {
		t.Fatal("group has 0 members after update")
	}

	group, err = db.GetGroupByID(group.ID)
	if err != nil {
		t.Fatalf("Error getting newly created group: %v", err)
	}
	if group.Picture != "Picture" {
		t.Fatal("group picture not correct")
	}

	_, err = db.DeleteGroup(group.ID)
	if err != nil {
		t.Fatalf("Error deleting group: %v", err)
	}
	_, err = db.GetGroupByID(group.ID)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("Error received from getting deleted group: %v", err)
	}
}

func TestGetUserGroups(t *testing.T) {

	user, _ := db.CreateUser(&models.User{ID: uuid.New(), UserName: "Slimo300"})
	group1, _ := db.CreateGroup(&models.Group{ID: uuid.New(), Name: "New Group", Created: time.Now()})
	group2, _ := db.CreateGroup(&models.Group{ID: uuid.New(), Name: "New Group2", Created: time.Now()})
	member1, _ := db.CreateMember(&models.Member{ID: uuid.New(), UserID: user.ID, GroupID: group1.ID})
	member2, _ := db.CreateMember(&models.Member{ID: uuid.New(), UserID: user.ID, GroupID: group2.ID})

	t.Cleanup(func() {
		_, _ = db.DeleteMember(member1.ID)
		_, _ = db.DeleteMember(member2.ID)
		_, _ = db.DeleteGroup(group1.ID)
		_, _ = db.DeleteGroup(group2.ID)
		_, _ = db.DeleteUser(user.ID)
	})

	expectedGroups := map[uuid.UUID]bool{
		group1.ID: true,
		group2.ID: true,
	}

	groups, err := db.GetUserGroups(user.ID)
	if err != nil {
		t.Fatalf("Error getting user groups: %v", err)
	}
	if len(groups) != 2 {
		t.Fatalf("Wrong number of returned groups: %d, should be 2", len(groups))
	}

	for _, group := range groups {
		if !expectedGroups[group.ID] {
			t.Fatalf("Unexpected group returned: %v", group)
		}
		// we delete the group ID from map to mark that it has already been found
		delete(expectedGroups, group.ID)
	}
	if len(expectedGroups) != 0 {
		t.Fatalf("Not all groups found: %v", expectedGroups)
	}
}
