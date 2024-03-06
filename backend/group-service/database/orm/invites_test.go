package orm_test

import (
	"errors"
	"testing"
	"time"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestHandleInvite(t *testing.T) {

	user1, _ := db.CreateUser(&models.User{ID: uuid.New(), UserName: "1"})
	user2, _ := db.CreateUser(&models.User{ID: uuid.New(), UserName: "2"})
	group, _ := db.CreateGroup(&models.Group{ID: uuid.New(), Created: time.Now()})

	t.Cleanup(func() {
		_, _ = db.DeleteGroup(group.ID)
		_, _ = db.DeleteUser(user1.ID)
		_, _ = db.DeleteUser(user2.ID)
	})

	invite, err := db.CreateInvite(&models.Invite{
		ID:       uuid.New(),
		IssId:    user1.ID,
		TargetID: user2.ID,
		GroupID:  group.ID,
		Status:   models.INVITE_AWAITING,
		Created:  time.Now(),
		Modified: time.Now(),
	})
	if err != nil {
		t.Fatalf("error creating invite: %v", err)
	}

	// checking whether invited user is invited
	ok, err := db.IsUserInvited(invite.TargetID, invite.GroupID)
	if err != nil {
		t.Fatalf("Error checking if user is invited: %v", err)
	}
	if !ok {
		t.Fatalf("Invite not found")
	}

	// checking whether issuing user is invited - should return false
	ok, err = db.IsUserInvited(invite.IssId, invite.GroupID)
	if err != nil {
		t.Fatalf("Error checking if user is invited: %v", err)
	}
	if ok {
		t.Fatalf("User shouldn't be invited")
	}

	invite, err = db.UpdateInvite(&models.Invite{
		ID:       invite.ID,
		Status:   models.INVITE_ACCEPT,
		Modified: time.Now(),
	})
	if err != nil {
		t.Fatalf("Error updating invite: %v", err)
	}
	if invite.Iss.UserName == "" {
		t.Fatalf("users not preloaded")
	}

	invite, err = db.GetInviteByID(invite.ID)
	if err != nil {
		t.Fatalf("Error getting invite: %v", err)
	}

	if invite.Status != models.INVITE_ACCEPT {
		t.Fatalf("Invite status not correct")
	}

	if _, err = db.DeleteInvite(invite.ID); err != nil {
		t.Fatalf("Error when deleting invite: %v", err)
	}

	if _, err := db.GetInviteByID(invite.ID); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("Getting invite should return not found error, instead got %v", err)
	}
}

func TestGetUserInvites(t *testing.T) {
	user1, _ := db.CreateUser(&models.User{ID: uuid.New(), UserName: "1"})
	user2, _ := db.CreateUser(&models.User{ID: uuid.New(), UserName: "2"})
	user3, _ := db.CreateUser(&models.User{ID: uuid.New(), UserName: "3"})
	group1, _ := db.CreateGroup((&models.Group{ID: uuid.New(), Name: "Group", Created: time.Now()}))
	group2, _ := db.CreateGroup((&models.Group{ID: uuid.New(), Name: "Group2", Created: time.Now()}))
	invite1, _ := db.CreateInvite(&models.Invite{ID: uuid.New(), TargetID: user1.ID, IssId: user2.ID, GroupID: group1.ID, Status: models.INVITE_DECLINE, Created: time.Now(), Modified: time.Now()})
	invite2, _ := db.CreateInvite(&models.Invite{ID: uuid.New(), TargetID: user1.ID, IssId: user2.ID, GroupID: group1.ID, Status: models.INVITE_ACCEPT, Created: time.Now(), Modified: time.Now()})
	invite3, _ := db.CreateInvite(&models.Invite{ID: uuid.New(), TargetID: user3.ID, IssId: user1.ID, GroupID: group1.ID, Status: models.INVITE_AWAITING, Created: time.Now(), Modified: time.Now()})
	invite4, _ := db.CreateInvite(&models.Invite{ID: uuid.New(), TargetID: user2.ID, IssId: user3.ID, GroupID: group2.ID, Status: models.INVITE_AWAITING, Created: time.Now(), Modified: time.Now()})

	t.Cleanup(func() {
		_, _ = db.DeleteInvite(invite1.ID)
		_, _ = db.DeleteInvite(invite2.ID)
		_, _ = db.DeleteInvite(invite3.ID)
		_, _ = db.DeleteInvite(invite4.ID)
		_, _ = db.DeleteGroup(group1.ID)
		_, _ = db.DeleteGroup(group2.ID)
		_, _ = db.DeleteUser(user1.ID)
		_, _ = db.DeleteUser(user2.ID)
		_, _ = db.DeleteUser(user3.ID)
	})

	invites, err := db.GetUserInvites(user1.ID, 4, 0)
	if err != nil {
		t.Fatalf("Getting user invites failed with errors: %v", err)
	}
	if len(invites) != 3 {
		t.Fatalf("Wrong number of invites, should be 3, there are: %d", len(invites))
	}

	expectedInvites := map[uuid.UUID]bool{
		invite1.ID: true,
		invite2.ID: true,
		invite3.ID: true,
	}

	for _, invite := range invites {
		if expectedInvites[invite.ID] {
			delete(expectedInvites, invite.ID)
		} else {
			t.Fatalf("Unexpected invite received: %v", invite.ID)
		}
	}

	if len(expectedInvites) != 0 {
		t.Fatalf("Not all invites found: %v", expectedInvites)
	}

}
