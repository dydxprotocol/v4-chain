// Code generated by mockery v2.42.0. DO NOT EDIT.

package mocks

import (
	big "math/big"

	mock "github.com/stretchr/testify/mock"
)

// ExtendVoteIndexPriceCache is an autogenerated mock type for the ExtendVoteIndexPriceCache type
type ExtendVoteIndexPriceCache struct {
	mock.Mock
}

// GetVEEncodedPrice provides a mock function with given fields: price
func (_m *ExtendVoteIndexPriceCache) GetVEEncodedPrice(price *big.Int) ([]byte, error) {
	ret := _m.Called(price)

	if len(ret) == 0 {
		panic("no return value specified for GetVEEncodedPrice")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(*big.Int) ([]byte, error)); ok {
		return rf(price)
	}
	if rf, ok := ret.Get(0).(func(*big.Int) []byte); ok {
		r0 = rf(price)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(*big.Int) error); ok {
		r1 = rf(price)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewExtendVoteIndexPriceCache creates a new instance of ExtendVoteIndexPriceCache. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewExtendVoteIndexPriceCache(t interface {
	mock.TestingT
	Cleanup(func())
}) *ExtendVoteIndexPriceCache {
	mock := &ExtendVoteIndexPriceCache{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
