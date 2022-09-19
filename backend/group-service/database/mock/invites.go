package mock

import (
	"errors"
	"time"

	"github.com/Slimo300/ChatApp/backend/src/models"
	"github.com/Slimo300/MicrosevicesChatApp/backend/group-service/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (m *MockDB) AddInvite(issID, targetID, groupID uuid.UUID) (models.Invite, error) {
	invite := models.Invite{ID: uuid.New(), IssId: issID, TargetID: targetID, GroupID: groupID,
		Status: database.INVITE_AWAITING, Created: time.Now(), Modified: time.Now()}
	m.Invites = append(m.Invites, invite)
	return invite, nil

}

func (mock *MockDB) GetUserInvites(userID uuid.UUID) ([]models.Invite, error) {

	userInvites := []models.Invite{}
	for _, invite := range mock.Invites {
		if invite.TargetID == userID && invite.Status == database.INVITE_AWAITING {
			userInvites = append(userInvites, invite)
		}
	}
	return userInvites, nil
}

func (mock *MockDB) GetInviteByID(inviteID uuid.UUID) (models.Invite, error) {
	for _, invite := range mock.Invites {
		if invite.ID == inviteID {
			return invite, nil
		}
	}
	return models.Invite{}, gorm.ErrRecordNotFound
}

func (mock *MockDB) AcceptInvite(invite models.Invite) (models.Group, error) {

	if err := mock.createOrRestoreMember(invite.TargetID, invite.GroupID); err != nil {
		return models.Group{}, gorm.ErrRecordNotFound
	}

	for _, inv := range mock.Invites {
		if inv.ID == invite.ID {
			inv.Status = database.INVITE_ACCEPT
			inv.Modified = time.Now()
		}
	}

	for _, group := range mock.Groups {
		if group.ID == invite.GroupID {
			return group, nil
		}
	}

	return models.Group{}, errors.New("Something went wrong")
}

func (mock *MockDB) createOrRestoreMember(userID, groupID uuid.UUID) error {
	for _, member := range mock.Members {
		if member.UserID == userID && member.GroupID == groupID && member.Deleted {
			member.Deleted = false
			return nil
		}
	}

	for _, user := range mock.Users {
		if user.ID == userID {
			mock.Members = append(mock.Members, models.Member{
				ID:       uuid.New(),
				UserID:   userID,
				GroupID:  groupID,
				Nick:     user.UserName,
				Adding:   false,
				Deleting: false,
				Setting:  false,
				Creator:  false,
				Deleted:  false,
			})
			return nil
		}
	}

	return errors.New("error")
}

func (m *MockDB) IsUserInvited(userID, groupID uuid.UUID) bool {
	for _, invite := range m.Invites {
		if invite.TargetID == userID && invite.GroupID == groupID && invite.Status == database.INVITE_AWAITING {
			return true
		}
	}
	return false
}

func (m *MockDB) DeclineInvite(inviteID uuid.UUID) error {
	for _, invite := range m.Invites {
		if invite.ID == inviteID {
			invite.Status = database.INVITE_DECLINE
			invite.Modified = time.Now()
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}
