// Code generated by mockery v2.14.1. DO NOT EDIT.

package mock

import (
	events "github.com/Slimo300/MicroservicesChatApp/backend/lib/events"
	mock "github.com/stretchr/testify/mock"

	models "github.com/Slimo300/MicroservicesChatApp/backend/group-service/models"

	uuid "github.com/google/uuid"
)

// MockGroupsDB is an autogenerated mock type for the DBLayer type
type MockGroupsDB struct {
	mock.Mock
}

// AddInvite provides a mock function with given fields: issID, targetID, groupID
func (_m *MockGroupsDB) AddInvite(issID uuid.UUID, targetID uuid.UUID, groupID uuid.UUID) (*models.Invite, error) {
	ret := _m.Called(issID, targetID, groupID)

	var r0 *models.Invite
	if rf, ok := ret.Get(0).(func(uuid.UUID, uuid.UUID, uuid.UUID) *models.Invite); ok {
		r0 = rf(issID, targetID, groupID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Invite)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID, uuid.UUID, uuid.UUID) error); ok {
		r1 = rf(issID, targetID, groupID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AnswerInvite provides a mock function with given fields: userID, inviteID, answer
func (_m *MockGroupsDB) AnswerInvite(userID uuid.UUID, inviteID uuid.UUID, answer bool) (*models.Invite, *models.Group, *models.Member, error) {
	ret := _m.Called(userID, inviteID, answer)

	var r0 *models.Invite
	if rf, ok := ret.Get(0).(func(uuid.UUID, uuid.UUID, bool) *models.Invite); ok {
		r0 = rf(userID, inviteID, answer)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Invite)
		}
	}

	var r1 *models.Group
	if rf, ok := ret.Get(1).(func(uuid.UUID, uuid.UUID, bool) *models.Group); ok {
		r1 = rf(userID, inviteID, answer)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*models.Group)
		}
	}

	var r2 *models.Member
	if rf, ok := ret.Get(2).(func(uuid.UUID, uuid.UUID, bool) *models.Member); ok {
		r2 = rf(userID, inviteID, answer)
	} else {
		if ret.Get(2) != nil {
			r2 = ret.Get(2).(*models.Member)
		}
	}

	var r3 error
	if rf, ok := ret.Get(3).(func(uuid.UUID, uuid.UUID, bool) error); ok {
		r3 = rf(userID, inviteID, answer)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

// CreateGroup provides a mock function with given fields: userID, name
func (_m *MockGroupsDB) CreateGroup(userID uuid.UUID, name string) (models.Group, error) {
	ret := _m.Called(userID, name)

	var r0 models.Group
	if rf, ok := ret.Get(0).(func(uuid.UUID, string) models.Group); ok {
		r0 = rf(userID, name)
	} else {
		r0 = ret.Get(0).(models.Group)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID, string) error); ok {
		r1 = rf(userID, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteGroup provides a mock function with given fields: userID, groupID
func (_m *MockGroupsDB) DeleteGroup(userID uuid.UUID, groupID uuid.UUID) (models.Group, error) {
	ret := _m.Called(userID, groupID)

	var r0 models.Group
	if rf, ok := ret.Get(0).(func(uuid.UUID, uuid.UUID) models.Group); ok {
		r0 = rf(userID, groupID)
	} else {
		r0 = ret.Get(0).(models.Group)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID, uuid.UUID) error); ok {
		r1 = rf(userID, groupID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteGroupProfilePicture provides a mock function with given fields: userID, groupID
func (_m *MockGroupsDB) DeleteGroupProfilePicture(userID uuid.UUID, groupID uuid.UUID) (string, error) {
	ret := _m.Called(userID, groupID)

	var r0 string
	if rf, ok := ret.Get(0).(func(uuid.UUID, uuid.UUID) string); ok {
		r0 = rf(userID, groupID)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID, uuid.UUID) error); ok {
		r1 = rf(userID, groupID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteMember provides a mock function with given fields: userID, groupID, memberID
func (_m *MockGroupsDB) DeleteMember(userID uuid.UUID, groupID uuid.UUID, memberID uuid.UUID) error {
	ret := _m.Called(userID, groupID, memberID)

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID, uuid.UUID, uuid.UUID) error); ok {
		r0 = rf(userID, groupID, memberID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetGroupProfilePictureURL provides a mock function with given fields: userID, groupID
func (_m *MockGroupsDB) GetGroupProfilePictureURL(userID uuid.UUID, groupID uuid.UUID) (string, error) {
	ret := _m.Called(userID, groupID)

	var r0 string
	if rf, ok := ret.Get(0).(func(uuid.UUID, uuid.UUID) string); ok {
		r0 = rf(userID, groupID)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID, uuid.UUID) error); ok {
		r1 = rf(userID, groupID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserGroups provides a mock function with given fields: id
func (_m *MockGroupsDB) GetUserGroups(id uuid.UUID) ([]models.Group, error) {
	ret := _m.Called(id)

	var r0 []models.Group
	if rf, ok := ret.Get(0).(func(uuid.UUID) []models.Group); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Group)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserInvites provides a mock function with given fields: userID, num, offset
func (_m *MockGroupsDB) GetUserInvites(userID uuid.UUID, num int, offset int) ([]models.Invite, error) {
	ret := _m.Called(userID, num, offset)

	var r0 []models.Invite
	if rf, ok := ret.Get(0).(func(uuid.UUID, int, int) []models.Invite); ok {
		r0 = rf(userID, num, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Invite)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID, int, int) error); ok {
		r1 = rf(userID, num, offset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GrantRights provides a mock function with given fields: userID, groupID, memberID, rights
func (_m *MockGroupsDB) GrantRights(userID uuid.UUID, groupID uuid.UUID, memberID uuid.UUID, rights models.MemberRights) (*models.Member, error) {
	ret := _m.Called(userID, groupID, memberID, rights)

	var r0 *models.Member
	if rf, ok := ret.Get(0).(func(uuid.UUID, uuid.UUID, uuid.UUID, models.MemberRights) *models.Member); ok {
		r0 = rf(userID, groupID, memberID, rights)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Member)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID, uuid.UUID, uuid.UUID, models.MemberRights) error); ok {
		r1 = rf(userID, groupID, memberID, rights)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUser provides a mock function with given fields: event
func (_m *MockGroupsDB) NewUser(event events.UserRegisteredEvent) error {
	ret := _m.Called(event)

	var r0 error
	if rf, ok := ret.Get(0).(func(events.UserRegisteredEvent) error); ok {
		r0 = rf(event)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewMockGroupsDB interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockGroupsDB creates a new instance of MockGroupsDB. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockGroupsDB(t mockConstructorTestingTNewMockGroupsDB) *MockGroupsDB {
	mock := &MockGroupsDB{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
