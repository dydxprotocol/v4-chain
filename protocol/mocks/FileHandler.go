// Code generated by mockery v2.23.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// FileHandler is an autogenerated mock type for the FileHandler type
type FileHandler struct {
	mock.Mock
}

// RemoveAll provides a mock function with given fields: path
func (_m *FileHandler) RemoveAll(path string) error {
	ret := _m.Called(path)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(path)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewFileHandler interface {
	mock.TestingT
	Cleanup(func())
}

// NewFileHandler creates a new instance of FileHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewFileHandler(t mockConstructorTestingTNewFileHandler) *FileHandler {
	mock := &FileHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
