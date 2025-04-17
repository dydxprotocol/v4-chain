// Code generated by mockery v2.50.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// SidecarVersionChecker is an autogenerated mock type for the SidecarVersionChecker type
type SidecarVersionChecker struct {
	mock.Mock
}

// CheckSidecarVersion provides a mock function with given fields: _a0
func (_m *SidecarVersionChecker) CheckSidecarVersion(_a0 context.Context) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for CheckSidecarVersion")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Start provides a mock function with given fields: ctx
func (_m *SidecarVersionChecker) Start(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Start")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Stop provides a mock function with no fields
func (_m *SidecarVersionChecker) Stop() {
	_m.Called()
}

// NewSidecarVersionChecker creates a new instance of SidecarVersionChecker. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSidecarVersionChecker(t interface {
	mock.TestingT
	Cleanup(func())
}) *SidecarVersionChecker {
	mock := &SidecarVersionChecker{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
