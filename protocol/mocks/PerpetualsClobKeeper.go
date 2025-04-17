// Code generated by mockery v2.50.0. DO NOT EDIT.

package mocks

import (
	types "github.com/cosmos/cosmos-sdk/types"
	perpetualstypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	mock "github.com/stretchr/testify/mock"
)

// PerpetualsClobKeeper is an autogenerated mock type for the PerpetualsClobKeeper type
type PerpetualsClobKeeper struct {
	mock.Mock
}

// GetPricePremiumForPerpetual provides a mock function with given fields: ctx, perpetualId, params
func (_m *PerpetualsClobKeeper) GetPricePremiumForPerpetual(ctx types.Context, perpetualId uint32, params perpetualstypes.GetPricePremiumParams) (int32, error) {
	ret := _m.Called(ctx, perpetualId, params)

	if len(ret) == 0 {
		panic("no return value specified for GetPricePremiumForPerpetual")
	}

	var r0 int32
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, uint32, perpetualstypes.GetPricePremiumParams) (int32, error)); ok {
		return rf(ctx, perpetualId, params)
	}
	if rf, ok := ret.Get(0).(func(types.Context, uint32, perpetualstypes.GetPricePremiumParams) int32); ok {
		r0 = rf(ctx, perpetualId, params)
	} else {
		r0 = ret.Get(0).(int32)
	}

	if rf, ok := ret.Get(1).(func(types.Context, uint32, perpetualstypes.GetPricePremiumParams) error); ok {
		r1 = rf(ctx, perpetualId, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsPerpetualClobPairActive provides a mock function with given fields: ctx, perpetualId
func (_m *PerpetualsClobKeeper) IsPerpetualClobPairActive(ctx types.Context, perpetualId uint32) (bool, error) {
	ret := _m.Called(ctx, perpetualId)

	if len(ret) == 0 {
		panic("no return value specified for IsPerpetualClobPairActive")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, uint32) (bool, error)); ok {
		return rf(ctx, perpetualId)
	}
	if rf, ok := ret.Get(0).(func(types.Context, uint32) bool); ok {
		r0 = rf(ctx, perpetualId)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(types.Context, uint32) error); ok {
		r1 = rf(ctx, perpetualId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewPerpetualsClobKeeper creates a new instance of PerpetualsClobKeeper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPerpetualsClobKeeper(t interface {
	mock.TestingT
	Cleanup(func())
}) *PerpetualsClobKeeper {
	mock := &PerpetualsClobKeeper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
