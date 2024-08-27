// Code generated by mockery v2.42.0. DO NOT EDIT.

package mocks

import (
	clobtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	mock "github.com/stretchr/testify/mock"

	types "github.com/cosmos/cosmos-sdk/types"
)

// ExtendVoteClobKeeper is an autogenerated mock type for the ExtendVoteClobKeeper type
type ExtendVoteClobKeeper struct {
	mock.Mock
}

// GetClobPair provides a mock function with given fields: ctx, id
func (_m *ExtendVoteClobKeeper) GetClobPair(ctx types.Context, id clobtypes.ClobPairId) (clobtypes.ClobPair, bool) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetClobPair")
	}

	var r0 clobtypes.ClobPair
	var r1 bool
	if rf, ok := ret.Get(0).(func(types.Context, clobtypes.ClobPairId) (clobtypes.ClobPair, bool)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(types.Context, clobtypes.ClobPairId) clobtypes.ClobPair); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(clobtypes.ClobPair)
	}

	if rf, ok := ret.Get(1).(func(types.Context, clobtypes.ClobPairId) bool); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// GetSingleMarketClobMetadata provides a mock function with given fields: ctx, clobPair
func (_m *ExtendVoteClobKeeper) GetSingleMarketClobMetadata(ctx types.Context, clobPair clobtypes.ClobPair) clobtypes.ClobMetadata {
	ret := _m.Called(ctx, clobPair)

	if len(ret) == 0 {
		panic("no return value specified for GetSingleMarketClobMetadata")
	}

	var r0 clobtypes.ClobMetadata
	if rf, ok := ret.Get(0).(func(types.Context, clobtypes.ClobPair) clobtypes.ClobMetadata); ok {
		r0 = rf(ctx, clobPair)
	} else {
		r0 = ret.Get(0).(clobtypes.ClobMetadata)
	}

	return r0
}

// NewExtendVoteClobKeeper creates a new instance of ExtendVoteClobKeeper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewExtendVoteClobKeeper(t interface {
	mock.TestingT
	Cleanup(func())
}) *ExtendVoteClobKeeper {
	mock := &ExtendVoteClobKeeper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
