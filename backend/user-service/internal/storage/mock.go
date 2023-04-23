// Code generated by mockery v2.20.0. DO NOT EDIT.

package storage

import (
	multipart "mime/multipart"

	mock "github.com/stretchr/testify/mock"
)

// MockStorage is an autogenerated mock type for the StorageLayer type
type MockStorage struct {
	mock.Mock
}

// DeleteFile provides a mock function with given fields: key
func (_m *MockStorage) DeleteFile(key string) error {
	ret := _m.Called(key)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UploadFile provides a mock function with given fields: img, key
func (_m *MockStorage) UploadFile(img multipart.File, key string) error {
	ret := _m.Called(img, key)

	var r0 error
	if rf, ok := ret.Get(0).(func(multipart.File, string) error); ok {
		r0 = rf(img, key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewMockStorage interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockStorage creates a new instance of MockStorage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockStorage(t mockConstructorTestingTNewMockStorage) *MockStorage {
	mock := &MockStorage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}