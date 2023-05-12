// Code generated by mockery v2.20.0. DO NOT EDIT.

package mock

import (
	events "github.com/Slimo300/MicroservicesChatApp/backend/lib/events"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// MockWsDB is an autogenerated mock type for the DBLayer type
type MockWsDB struct {
	mock.Mock
}

// CheckAccessCode provides a mock function with given fields: accessCode
func (_m *MockWsDB) CheckAccessCode(accessCode string) (uuid.UUID, error) {
	ret := _m.Called(accessCode)

	var r0 uuid.UUID
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (uuid.UUID, error)); ok {
		return rf(accessCode)
	}
	if rf, ok := ret.Get(0).(func(string) uuid.UUID); ok {
		r0 = rf(accessCode)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(accessCode)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteGroup provides a mock function with given fields: event
func (_m *MockWsDB) DeleteGroup(event events.GroupDeletedEvent) error {
	ret := _m.Called(event)

	var r0 error
	if rf, ok := ret.Get(0).(func(events.GroupDeletedEvent) error); ok {
		r0 = rf(event)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteMember provides a mock function with given fields: event
func (_m *MockWsDB) DeleteMember(event events.MemberDeletedEvent) error {
	ret := _m.Called(event)

	var r0 error
	if rf, ok := ret.Get(0).(func(events.MemberDeletedEvent) error); ok {
		r0 = rf(event)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetUserGroups provides a mock function with given fields: userID
func (_m *MockWsDB) GetUserGroups(userID uuid.UUID) ([]uuid.UUID, error) {
	ret := _m.Called(userID)

	var r0 []uuid.UUID
	var r1 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) ([]uuid.UUID, error)); ok {
		return rf(userID)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID) []uuid.UUID); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]uuid.UUID)
		}
	}

	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewAccessCode provides a mock function with given fields: userID, accessCode
func (_m *MockWsDB) NewAccessCode(userID uuid.UUID, accessCode string) error {
	ret := _m.Called(userID, accessCode)

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID, string) error); ok {
		r0 = rf(userID, accessCode)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMember provides a mock function with given fields: event
func (_m *MockWsDB) NewMember(event events.MemberCreatedEvent) error {
	ret := _m.Called(event)

	var r0 error
	if rf, ok := ret.Get(0).(func(events.MemberCreatedEvent) error); ok {
		r0 = rf(event)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewMockWsDB interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockWsDB creates a new instance of MockWsDB. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockWsDB(t mockConstructorTestingTNewMockWsDB) *MockWsDB {
	mock := &MockWsDB{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}