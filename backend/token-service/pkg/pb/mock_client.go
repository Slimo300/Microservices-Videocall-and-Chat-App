// Code generated by mockery v2.20.0. DO NOT EDIT.

package pb

import (
	context "context"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"
)

// MockTokenServiceClient is an autogenerated mock type for the TokenServiceClient type
type MockTokenServiceClient struct {
	mock.Mock
}

// DeleteUserToken provides a mock function with given fields: ctx, in, opts
func (_m *MockTokenServiceClient) DeleteUserToken(ctx context.Context, in *RefreshToken, opts ...grpc.CallOption) (*Msg, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *Msg
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *RefreshToken, ...grpc.CallOption) (*Msg, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *RefreshToken, ...grpc.CallOption) *Msg); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Msg)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *RefreshToken, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPublicKey provides a mock function with given fields: ctx, in, opts
func (_m *MockTokenServiceClient) GetPublicKey(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*PublicKey, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *PublicKey
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *Empty, ...grpc.CallOption) (*PublicKey, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *Empty, ...grpc.CallOption) *PublicKey); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*PublicKey)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *Empty, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewPairFromRefresh provides a mock function with given fields: ctx, in, opts
func (_m *MockTokenServiceClient) NewPairFromRefresh(ctx context.Context, in *RefreshToken, opts ...grpc.CallOption) (*TokenPair, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *TokenPair
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *RefreshToken, ...grpc.CallOption) (*TokenPair, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *RefreshToken, ...grpc.CallOption) *TokenPair); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*TokenPair)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *RefreshToken, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewPairFromUserID provides a mock function with given fields: ctx, in, opts
func (_m *MockTokenServiceClient) NewPairFromUserID(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*TokenPair, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *TokenPair
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *UserID, ...grpc.CallOption) (*TokenPair, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *UserID, ...grpc.CallOption) *TokenPair); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*TokenPair)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *UserID, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMockTokenServiceClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockTokenServiceClient creates a new instance of MockTokenServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockTokenServiceClient(t mockConstructorTestingTNewMockTokenServiceClient) *MockTokenServiceClient {
	mock := &MockTokenServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}