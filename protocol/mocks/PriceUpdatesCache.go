// Code generated by mockery v2.42.0. DO NOT EDIT.

package mocks

import (
	pricecache "github.com/StreamFinance-Protocol/stream-chain/protocol/caches/pricecache"
	types "github.com/cosmos/cosmos-sdk/types"
	mock "github.com/stretchr/testify/mock"
)

// PriceUpdatesCache is an autogenerated mock type for the PriceUpdatesCache type
type PriceUpdatesCache struct {
	mock.Mock
}

// GetPriceUpdates provides a mock function with given fields:
func (_m *PriceUpdatesCache) GetPriceUpdates() pricecache.PriceUpdates {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetPriceUpdates")
	}

	var r0 pricecache.PriceUpdates
	if rf, ok := ret.Get(0).(func() pricecache.PriceUpdates); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(pricecache.PriceUpdates)
		}
	}

	return r0
}

// HasValidValues provides a mock function with given fields: currTxHash
func (_m *PriceUpdatesCache) HasValidValues(currTxHash []byte) bool {
	ret := _m.Called(currTxHash)

	if len(ret) == 0 {
		panic("no return value specified for HasValidValues")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func([]byte) bool); ok {
		r0 = rf(currTxHash)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// SetPriceUpdates provides a mock function with given fields: ctx, updates, txHash
func (_m *PriceUpdatesCache) SetPriceUpdates(ctx types.Context, updates pricecache.PriceUpdates, txHash []byte) {
	_m.Called(ctx, updates, txHash)
}

// NewPriceUpdatesCache creates a new instance of PriceUpdatesCache. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPriceUpdatesCache(t interface {
	mock.TestingT
	Cleanup(func())
}) *PriceUpdatesCache {
	mock := &PriceUpdatesCache{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
