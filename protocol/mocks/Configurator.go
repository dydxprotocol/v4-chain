// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	grpc "github.com/cosmos/gogoproto/grpc"
	google_golang_orggrpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"

	module "github.com/cosmos/cosmos-sdk/types/module"
)

// Configurator is an autogenerated mock type for the Configurator type
type Configurator struct {
	mock.Mock
}

// Error provides a mock function with given fields:
func (_m *Configurator) Error() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Error")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MsgServer provides a mock function with given fields:
func (_m *Configurator) MsgServer() grpc.Server {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for MsgServer")
	}

	var r0 grpc.Server
	if rf, ok := ret.Get(0).(func() grpc.Server); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(grpc.Server)
		}
	}

	return r0
}

// QueryServer provides a mock function with given fields:
func (_m *Configurator) QueryServer() grpc.Server {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for QueryServer")
	}

	var r0 grpc.Server
	if rf, ok := ret.Get(0).(func() grpc.Server); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(grpc.Server)
		}
	}

	return r0
}

// RegisterMigration provides a mock function with given fields: moduleName, fromVersion, handler
func (_m *Configurator) RegisterMigration(moduleName string, fromVersion uint64, handler module.MigrationHandler) error {
	ret := _m.Called(moduleName, fromVersion, handler)

	if len(ret) == 0 {
		panic("no return value specified for RegisterMigration")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, uint64, module.MigrationHandler) error); ok {
		r0 = rf(moduleName, fromVersion, handler)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RegisterService provides a mock function with given fields: sd, ss
func (_m *Configurator) RegisterService(sd *google_golang_orggrpc.ServiceDesc, ss interface{}) {
	_m.Called(sd, ss)
}

// NewConfigurator creates a new instance of Configurator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewConfigurator(t interface {
	mock.TestingT
	Cleanup(func())
}) *Configurator {
	mock := &Configurator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
