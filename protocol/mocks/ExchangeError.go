// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ExchangeError is an autogenerated mock type for the ExchangeError type
type ExchangeError struct {
	mock.Mock
}

// Error provides a mock function with given fields:
func (_m *ExchangeError) Error() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Error")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetExchangeId provides a mock function with given fields:
func (_m *ExchangeError) GetExchangeId() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetExchangeId")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// NewExchangeError creates a new instance of ExchangeError. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewExchangeError(t interface {
	mock.TestingT
	Cleanup(func())
}) *ExchangeError {
	mock := &ExchangeError{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
