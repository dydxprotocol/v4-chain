// Code generated by mockery v2.23.1. DO NOT EDIT.

package mocks

import (
	big "math/big"

	abcitypes "github.com/cometbft/cometbft/abci/types"

	mock "github.com/stretchr/testify/mock"

	types "github.com/cosmos/cosmos-sdk/types"
)

// ExtendVotePriceApplier is an autogenerated mock type for the ExtendVotePriceApplier type
type ExtendVotePriceApplier struct {
	mock.Mock
}

// ApplyPricesFromVoteExtensions provides a mock function with given fields: ctx, req
func (_m *ExtendVotePriceApplier) ApplyPricesFromVoteExtensions(ctx types.Context, req *abcitypes.RequestFinalizeBlock) (map[string]*big.Int, error) {
	ret := _m.Called(ctx, req)

	var r0 map[string]*big.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, *abcitypes.RequestFinalizeBlock) (map[string]*big.Int, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(types.Context, *abcitypes.RequestFinalizeBlock) map[string]*big.Int); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func(types.Context, *abcitypes.RequestFinalizeBlock) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewExtendVotePriceApplier interface {
	mock.TestingT
	Cleanup(func())
}

// NewExtendVotePriceApplier creates a new instance of ExtendVotePriceApplier. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewExtendVotePriceApplier(t mockConstructorTestingTNewExtendVotePriceApplier) *ExtendVotePriceApplier {
	mock := &ExtendVotePriceApplier{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
