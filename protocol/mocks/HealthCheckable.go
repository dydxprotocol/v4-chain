// Code generated by mockery v2.50.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// HealthCheckable is an autogenerated mock type for the HealthCheckable type
type HealthCheckable struct {
	mock.Mock
}

// HealthCheck provides a mock function with no fields
func (_m *HealthCheckable) HealthCheck() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for HealthCheck")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ReportFailure provides a mock function with given fields: err
func (_m *HealthCheckable) ReportFailure(err error) {
	_m.Called(err)
}

// ReportSuccess provides a mock function with no fields
func (_m *HealthCheckable) ReportSuccess() {
	_m.Called()
}

// ServiceName provides a mock function with no fields
func (_m *HealthCheckable) ServiceName() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ServiceName")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// NewHealthCheckable creates a new instance of HealthCheckable. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewHealthCheckable(t interface {
	mock.TestingT
	Cleanup(func())
}) *HealthCheckable {
	mock := &HealthCheckable{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
