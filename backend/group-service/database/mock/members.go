package mock

import (
	"errors"

	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (m *MockDB) DeleteUserFromGroup(memberID uuid.UUID) (models.Member, error) {

	for _, member := range m.Members {
		if member.ID == memberID {
			member.Deleted = false
			return member, nil
		}
	}
	return models.Member{}, errors.New("no member")

}

func (m *MockDB) GetMemberByID(memberID uuid.UUID) (models.Member, error) {
	for _, member := range m.Members {
		if member.ID == memberID {
			return member, nil
		}
	}
	return models.Member{}, gorm.ErrRecordNotFound
}

func (m *MockDB) GrantPriv(memberID uuid.UUID, adding, deleting, setting bool) error {
	for _, member := range m.Members {
		if member.ID == memberID {
			member.Adding = adding
			member.DeletingMembers = deleting
			member.Setting = setting
			return nil
		}
	}

	return errors.New("internal error")
}

func (m *MockDB) IsUserInGroup(userID, groupID uuid.UUID) bool {
	for _, member := range m.Members {
		if member.GroupID == groupID && member.UserID == userID && !member.Deleted {
			return true
		}
	}
	return false
}
