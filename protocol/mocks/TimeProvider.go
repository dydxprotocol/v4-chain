// Code generated by mockery v2.44.1. DO NOT EDIT.

package mocks

import (
	time "time"

	mock "github.com/stretchr/testify/mock"
)

// TimeProvider is an autogenerated mock type for the TimeProvider type
type TimeProvider struct {
	mock.Mock
}

// Now provides a mock function with given fields:
func (_m *TimeProvider) Now() time.Time {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Now")
	}

	var r0 time.Time
	if rf, ok := ret.Get(0).(func() time.Time); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	return r0
}

// NewTimeProvider creates a new instance of TimeProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTimeProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *TimeProvider {
	mock := &TimeProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
