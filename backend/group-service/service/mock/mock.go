// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mock

import (
	context "context"
	multipart "mime/multipart"

	models "github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// GroupsMockService is an autogenerated mock type for the Service type
type GroupsMockService struct {
	mock.Mock
}

// AddInvite provides a mock function with given fields: ctx, issID, targetID, groupID
func (_m *GroupsMockService) AddInvite(ctx context.Context, issID uuid.UUID, targetID uuid.UUID, groupID uuid.UUID) (*models.Invite, error) {
	ret := _m.Called(ctx, issID, targetID, groupID)

	var r0 *models.Invite
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID, uuid.UUID) *models.Invite); ok {
		r0 = rf(ctx, issID, targetID, groupID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Invite)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, uuid.UUID, uuid.UUID) error); ok {
		r1 = rf(ctx, issID, targetID, groupID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateGroup provides a mock function with given fields: ctx, userID, name
func (_m *GroupsMockService) CreateGroup(ctx context.Context, userID uuid.UUID, name string) (*models.Group, error) {
	ret := _m.Called(ctx, userID, name)

	var r0 *models.Group
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, string) *models.Group); ok {
		r0 = rf(ctx, userID, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Group)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, string) error); ok {
		r1 = rf(ctx, userID, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteGroup provides a mock function with given fields: ctx, userID, groupID
func (_m *GroupsMockService) DeleteGroup(ctx context.Context, userID uuid.UUID, groupID uuid.UUID) (*models.Group, error) {
	ret := _m.Called(ctx, userID, groupID)

	var r0 *models.Group
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) *models.Group); ok {
		r0 = rf(ctx, userID, groupID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Group)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, uuid.UUID) error); ok {
		r1 = rf(ctx, userID, groupID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteGroupPicture provides a mock function with given fields: ctx, userID, groupID
func (_m *GroupsMockService) DeleteGroupPicture(ctx context.Context, userID uuid.UUID, groupID uuid.UUID) error {
	ret := _m.Called(ctx, userID, groupID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) error); ok {
		r0 = rf(ctx, userID, groupID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteMember provides a mock function with given fields: ctx, userID, memberID
func (_m *GroupsMockService) DeleteMember(ctx context.Context, userID uuid.UUID, memberID uuid.UUID) (*models.Member, error) {
	ret := _m.Called(ctx, userID, memberID)

	var r0 *models.Member
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) *models.Member); ok {
		r0 = rf(ctx, userID, memberID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Member)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, uuid.UUID) error); ok {
		r1 = rf(ctx, userID, memberID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserGroups provides a mock function with given fields: ctx, id
func (_m *GroupsMockService) GetUserGroups(ctx context.Context, id uuid.UUID) ([]*models.Group, error) {
	ret := _m.Called(ctx, id)

	var r0 []*models.Group
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) []*models.Group); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Group)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserInvites provides a mock function with given fields: ctx, userID, num, offset
func (_m *GroupsMockService) GetUserInvites(ctx context.Context, userID uuid.UUID, num int, offset int) ([]*models.Invite, error) {
	ret := _m.Called(ctx, userID, num, offset)

	var r0 []*models.Invite
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, int, int) []*models.Invite); ok {
		r0 = rf(ctx, userID, num, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Invite)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, int, int) error); ok {
		r1 = rf(ctx, userID, num, offset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GrantRights provides a mock function with given fields: ctx, userID, memberID, rights
func (_m *GroupsMockService) GrantRights(ctx context.Context, userID uuid.UUID, memberID uuid.UUID, rights models.MemberRights) (*models.Member, error) {
	ret := _m.Called(ctx, userID, memberID, rights)

	var r0 *models.Member
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID, models.MemberRights) *models.Member); ok {
		r0 = rf(ctx, userID, memberID, rights)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Member)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, uuid.UUID, models.MemberRights) error); ok {
		r1 = rf(ctx, userID, memberID, rights)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RespondInvite provides a mock function with given fields: ctx, userID, inviteID, answer
func (_m *GroupsMockService) RespondInvite(ctx context.Context, userID uuid.UUID, inviteID uuid.UUID, answer bool) (*models.Invite, *models.Group, error) {
	ret := _m.Called(ctx, userID, inviteID, answer)

	var r0 *models.Invite
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID, bool) *models.Invite); ok {
		r0 = rf(ctx, userID, inviteID, answer)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Invite)
		}
	}

	var r1 *models.Group
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, uuid.UUID, bool) *models.Group); ok {
		r1 = rf(ctx, userID, inviteID, answer)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*models.Group)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, uuid.UUID, uuid.UUID, bool) error); ok {
		r2 = rf(ctx, userID, inviteID, answer)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// SetGroupPicture provides a mock function with given fields: ctx, userID, groupID, file
func (_m *GroupsMockService) SetGroupPicture(ctx context.Context, userID uuid.UUID, groupID uuid.UUID, file multipart.File) error {
	ret := _m.Called(ctx, userID, groupID, file)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID, multipart.File) error); ok {
		r0 = rf(ctx, userID, groupID, file)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewGroupsMockService interface {
	mock.TestingT
	Cleanup(func())
}

// NewGroupsMockService creates a new instance of GroupsMockService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewGroupsMockService(t mockConstructorTestingTNewGroupsMockService) *GroupsMockService {
	mock := &GroupsMockService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
