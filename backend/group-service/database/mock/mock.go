// Code generated by mockery v2.42.0. DO NOT EDIT.

package mock

import (
	models "github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// MockGroupsDB is an autogenerated mock type for the DBLayer type
type MockGroupsDB struct {
	mock.Mock
}

// CreateGroup provides a mock function with given fields: group
func (_m *MockGroupsDB) CreateGroup(group *models.Group) (*models.Group, error) {
	ret := _m.Called(group)

	if len(ret) == 0 {
		panic("no return value specified for CreateGroup")
	}

	var r0 *models.Group
	var r1 error
	if rf, ok := ret.Get(0).(func(*models.Group) (*models.Group, error)); ok {
		return rf(group)
	}
	if rf, ok := ret.Get(0).(func(*models.Group) *models.Group); ok {
		r0 = rf(group)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Group)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.Group) error); ok {
		r1 = rf(group)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateInvite provides a mock function with given fields: invite
func (_m *MockGroupsDB) CreateInvite(invite *models.Invite) (*models.Invite, error) {
	ret := _m.Called(invite)

	if len(ret) == 0 {
		panic("no return value specified for CreateInvite")
	}

	var r0 *models.Invite
	var r1 error
	if rf, ok := ret.Get(0).(func(*models.Invite) (*models.Invite, error)); ok {
		return rf(invite)
	}
	if rf, ok := ret.Get(0).(func(*models.Invite) *models.Invite); ok {
		r0 = rf(invite)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Invite)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.Invite) error); ok {
		r1 = rf(invite)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateMember provides a mock function with given fields: member
func (_m *MockGroupsDB) CreateMember(member *models.Member) (*models.Member, error) {
	ret := _m.Called(member)

	if len(ret) == 0 {
		panic("no return value specified for CreateMember")
	}

	var r0 *models.Member
	var r1 error
	if rf, ok := ret.Get(0).(func(*models.Member) (*models.Member, error)); ok {
		return rf(member)
	}
	if rf, ok := ret.Get(0).(func(*models.Member) *models.Member); ok {
		r0 = rf(member)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Member)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.Member) error); ok {
		r1 = rf(member)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateUser provides a mock function with given fields: user
func (_m *MockGroupsDB) CreateUser(user *models.User) (*models.User, error) {
	ret := _m.Called(user)

	if len(ret) == 0 {
		panic("no return value specified for CreateUser")
	}

	var r0 *models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(*models.User) (*models.User, error)); ok {
		return rf(user)
	}
	if rf, ok := ret.Get(0).(func(*models.User) *models.User); ok {
		r0 = rf(user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.User) error); ok {
		r1 = rf(user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteGroup provides a mock function with given fields: groupID
func (_m *MockGroupsDB) DeleteGroup(groupID uuid.UUID) (*models.Group, error) {
	ret := _m.Called(groupID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteGroup")
	}

	var r0 *models.Group
	var r1 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) (*models.Group, error)); ok {
		return rf(groupID)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID) *models.Group); ok {
		r0 = rf(groupID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Group)
		}
	}

	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(groupID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteInvite provides a mock function with given fields: inviteID
func (_m *MockGroupsDB) DeleteInvite(inviteID uuid.UUID) (*models.Invite, error) {
	ret := _m.Called(inviteID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteInvite")
	}

	var r0 *models.Invite
	var r1 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) (*models.Invite, error)); ok {
		return rf(inviteID)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID) *models.Invite); ok {
		r0 = rf(inviteID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Invite)
		}
	}

	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(inviteID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteMember provides a mock function with given fields: memberID
func (_m *MockGroupsDB) DeleteMember(memberID uuid.UUID) (*models.Member, error) {
	ret := _m.Called(memberID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteMember")
	}

	var r0 *models.Member
	var r1 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) (*models.Member, error)); ok {
		return rf(memberID)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID) *models.Member); ok {
		r0 = rf(memberID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Member)
		}
	}

	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(memberID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteUser provides a mock function with given fields: userID
func (_m *MockGroupsDB) DeleteUser(userID uuid.UUID) (*models.User, error) {
	ret := _m.Called(userID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteUser")
	}

	var r0 *models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) (*models.User, error)); ok {
		return rf(userID)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID) *models.User); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetGroupByID provides a mock function with given fields: groupID
func (_m *MockGroupsDB) GetGroupByID(groupID uuid.UUID) (*models.Group, error) {
	ret := _m.Called(groupID)

	if len(ret) == 0 {
		panic("no return value specified for GetGroupByID")
	}

	var r0 *models.Group
	var r1 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) (*models.Group, error)); ok {
		return rf(groupID)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID) *models.Group); ok {
		r0 = rf(groupID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Group)
		}
	}

	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(groupID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetInviteByID provides a mock function with given fields: inviteID
func (_m *MockGroupsDB) GetInviteByID(inviteID uuid.UUID) (*models.Invite, error) {
	ret := _m.Called(inviteID)

	if len(ret) == 0 {
		panic("no return value specified for GetInviteByID")
	}

	var r0 *models.Invite
	var r1 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) (*models.Invite, error)); ok {
		return rf(inviteID)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID) *models.Invite); ok {
		r0 = rf(inviteID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Invite)
		}
	}

	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(inviteID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMemberByID provides a mock function with given fields: memberID
func (_m *MockGroupsDB) GetMemberByID(memberID uuid.UUID) (*models.Member, error) {
	ret := _m.Called(memberID)

	if len(ret) == 0 {
		panic("no return value specified for GetMemberByID")
	}

	var r0 *models.Member
	var r1 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) (*models.Member, error)); ok {
		return rf(memberID)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID) *models.Member); ok {
		r0 = rf(memberID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Member)
		}
	}

	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(memberID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMemberByUserGroupID provides a mock function with given fields: userID, groupID
func (_m *MockGroupsDB) GetMemberByUserGroupID(userID uuid.UUID, groupID uuid.UUID) (*models.Member, error) {
	ret := _m.Called(userID, groupID)

	if len(ret) == 0 {
		panic("no return value specified for GetMemberByUserGroupID")
	}

	var r0 *models.Member
	var r1 error
	if rf, ok := ret.Get(0).(func(uuid.UUID, uuid.UUID) (*models.Member, error)); ok {
		return rf(userID, groupID)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID, uuid.UUID) *models.Member); ok {
		r0 = rf(userID, groupID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Member)
		}
	}

	if rf, ok := ret.Get(1).(func(uuid.UUID, uuid.UUID) error); ok {
		r1 = rf(userID, groupID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByID provides a mock function with given fields: userID
func (_m *MockGroupsDB) GetUserByID(userID uuid.UUID) (*models.User, error) {
	ret := _m.Called(userID)

	if len(ret) == 0 {
		panic("no return value specified for GetUserByID")
	}

	var r0 *models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) (*models.User, error)); ok {
		return rf(userID)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID) *models.User); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserGroups provides a mock function with given fields: userID
func (_m *MockGroupsDB) GetUserGroups(userID uuid.UUID) ([]*models.Group, error) {
	ret := _m.Called(userID)

	if len(ret) == 0 {
		panic("no return value specified for GetUserGroups")
	}

	var r0 []*models.Group
	var r1 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) ([]*models.Group, error)); ok {
		return rf(userID)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID) []*models.Group); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Group)
		}
	}

	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserInvites provides a mock function with given fields: userID, num, offset
func (_m *MockGroupsDB) GetUserInvites(userID uuid.UUID, num int, offset int) ([]*models.Invite, error) {
	ret := _m.Called(userID, num, offset)

	if len(ret) == 0 {
		panic("no return value specified for GetUserInvites")
	}

	var r0 []*models.Invite
	var r1 error
	if rf, ok := ret.Get(0).(func(uuid.UUID, int, int) ([]*models.Invite, error)); ok {
		return rf(userID, num, offset)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID, int, int) []*models.Invite); ok {
		r0 = rf(userID, num, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Invite)
		}
	}

	if rf, ok := ret.Get(1).(func(uuid.UUID, int, int) error); ok {
		r1 = rf(userID, num, offset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsUserInvited provides a mock function with given fields: userID, groupID
func (_m *MockGroupsDB) IsUserInvited(userID uuid.UUID, groupID uuid.UUID) (bool, error) {
	ret := _m.Called(userID, groupID)

	if len(ret) == 0 {
		panic("no return value specified for IsUserInvited")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(uuid.UUID, uuid.UUID) (bool, error)); ok {
		return rf(userID, groupID)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID, uuid.UUID) bool); ok {
		r0 = rf(userID, groupID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(uuid.UUID, uuid.UUID) error); ok {
		r1 = rf(userID, groupID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateGroup provides a mock function with given fields: group
func (_m *MockGroupsDB) UpdateGroup(group *models.Group) (*models.Group, error) {
	ret := _m.Called(group)

	if len(ret) == 0 {
		panic("no return value specified for UpdateGroup")
	}

	var r0 *models.Group
	var r1 error
	if rf, ok := ret.Get(0).(func(*models.Group) (*models.Group, error)); ok {
		return rf(group)
	}
	if rf, ok := ret.Get(0).(func(*models.Group) *models.Group); ok {
		r0 = rf(group)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Group)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.Group) error); ok {
		r1 = rf(group)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateInvite provides a mock function with given fields: invite
func (_m *MockGroupsDB) UpdateInvite(invite *models.Invite) (*models.Invite, error) {
	ret := _m.Called(invite)

	if len(ret) == 0 {
		panic("no return value specified for UpdateInvite")
	}

	var r0 *models.Invite
	var r1 error
	if rf, ok := ret.Get(0).(func(*models.Invite) (*models.Invite, error)); ok {
		return rf(invite)
	}
	if rf, ok := ret.Get(0).(func(*models.Invite) *models.Invite); ok {
		r0 = rf(invite)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Invite)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.Invite) error); ok {
		r1 = rf(invite)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateMember provides a mock function with given fields: member
func (_m *MockGroupsDB) UpdateMember(member *models.Member) (*models.Member, error) {
	ret := _m.Called(member)

	if len(ret) == 0 {
		panic("no return value specified for UpdateMember")
	}

	var r0 *models.Member
	var r1 error
	if rf, ok := ret.Get(0).(func(*models.Member) (*models.Member, error)); ok {
		return rf(member)
	}
	if rf, ok := ret.Get(0).(func(*models.Member) *models.Member); ok {
		r0 = rf(member)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Member)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.Member) error); ok {
		r1 = rf(member)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateUser provides a mock function with given fields: user
func (_m *MockGroupsDB) UpdateUser(user *models.User) (*models.User, error) {
	ret := _m.Called(user)

	if len(ret) == 0 {
		panic("no return value specified for UpdateUser")
	}

	var r0 *models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(*models.User) (*models.User, error)); ok {
		return rf(user)
	}
	if rf, ok := ret.Get(0).(func(*models.User) *models.User); ok {
		r0 = rf(user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.User) error); ok {
		r1 = rf(user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockGroupsDB creates a new instance of MockGroupsDB. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockGroupsDB(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockGroupsDB {
	mock := &MockGroupsDB{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
