// Code generated by mockery v2.50.0. DO NOT EDIT.

package mocks

import (
	big "math/big"

	perpetualstypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	mock "github.com/stretchr/testify/mock"

	types "github.com/cosmos/cosmos-sdk/types"
)

// PerpetualsKeeper is an autogenerated mock type for the PerpetualsKeeper type
type PerpetualsKeeper struct {
	mock.Mock
}

// AddPremiumVotes provides a mock function with given fields: ctx, votes
func (_m *PerpetualsKeeper) AddPremiumVotes(ctx types.Context, votes []perpetualstypes.FundingPremium) error {
	ret := _m.Called(ctx, votes)

	if len(ret) == 0 {
		panic("no return value specified for AddPremiumVotes")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, []perpetualstypes.FundingPremium) error); ok {
		r0 = rf(ctx, votes)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreatePerpetual provides a mock function with given fields: ctx, id, ticker, marketId, atomicResolution, defaultFundingPpm, liquidityTier, marketType
func (_m *PerpetualsKeeper) CreatePerpetual(ctx types.Context, id uint32, ticker string, marketId uint32, atomicResolution int32, defaultFundingPpm int32, liquidityTier uint32, marketType perpetualstypes.PerpetualMarketType) (perpetualstypes.Perpetual, error) {
	ret := _m.Called(ctx, id, ticker, marketId, atomicResolution, defaultFundingPpm, liquidityTier, marketType)

	if len(ret) == 0 {
		panic("no return value specified for CreatePerpetual")
	}

	var r0 perpetualstypes.Perpetual
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, uint32, string, uint32, int32, int32, uint32, perpetualstypes.PerpetualMarketType) (perpetualstypes.Perpetual, error)); ok {
		return rf(ctx, id, ticker, marketId, atomicResolution, defaultFundingPpm, liquidityTier, marketType)
	}
	if rf, ok := ret.Get(0).(func(types.Context, uint32, string, uint32, int32, int32, uint32, perpetualstypes.PerpetualMarketType) perpetualstypes.Perpetual); ok {
		r0 = rf(ctx, id, ticker, marketId, atomicResolution, defaultFundingPpm, liquidityTier, marketType)
	} else {
		r0 = ret.Get(0).(perpetualstypes.Perpetual)
	}

	if rf, ok := ret.Get(1).(func(types.Context, uint32, string, uint32, int32, int32, uint32, perpetualstypes.PerpetualMarketType) error); ok {
		r1 = rf(ctx, id, ticker, marketId, atomicResolution, defaultFundingPpm, liquidityTier, marketType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAddPremiumVotes provides a mock function with given fields: ctx
func (_m *PerpetualsKeeper) GetAddPremiumVotes(ctx types.Context) *perpetualstypes.MsgAddPremiumVotes {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAddPremiumVotes")
	}

	var r0 *perpetualstypes.MsgAddPremiumVotes
	if rf, ok := ret.Get(0).(func(types.Context) *perpetualstypes.MsgAddPremiumVotes); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*perpetualstypes.MsgAddPremiumVotes)
		}
	}

	return r0
}

// GetAllLiquidityTiers provides a mock function with given fields: ctx
func (_m *PerpetualsKeeper) GetAllLiquidityTiers(ctx types.Context) []perpetualstypes.LiquidityTier {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAllLiquidityTiers")
	}

	var r0 []perpetualstypes.LiquidityTier
	if rf, ok := ret.Get(0).(func(types.Context) []perpetualstypes.LiquidityTier); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]perpetualstypes.LiquidityTier)
		}
	}

	return r0
}

// GetAllPerpetuals provides a mock function with given fields: ctx
func (_m *PerpetualsKeeper) GetAllPerpetuals(ctx types.Context) []perpetualstypes.Perpetual {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAllPerpetuals")
	}

	var r0 []perpetualstypes.Perpetual
	if rf, ok := ret.Get(0).(func(types.Context) []perpetualstypes.Perpetual); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]perpetualstypes.Perpetual)
		}
	}

	return r0
}

// GetNetCollateral provides a mock function with given fields: ctx, id, bigQuantums
func (_m *PerpetualsKeeper) GetNetCollateral(ctx types.Context, id uint32, bigQuantums *big.Int) (*big.Int, error) {
	ret := _m.Called(ctx, id, bigQuantums)

	if len(ret) == 0 {
		panic("no return value specified for GetNetCollateral")
	}

	var r0 *big.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, uint32, *big.Int) (*big.Int, error)); ok {
		return rf(ctx, id, bigQuantums)
	}
	if rf, ok := ret.Get(0).(func(types.Context, uint32, *big.Int) *big.Int); ok {
		r0 = rf(ctx, id, bigQuantums)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func(types.Context, uint32, *big.Int) error); ok {
		r1 = rf(ctx, id, bigQuantums)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetNetNotional provides a mock function with given fields: ctx, id, bigQuantums
func (_m *PerpetualsKeeper) GetNetNotional(ctx types.Context, id uint32, bigQuantums *big.Int) (*big.Int, error) {
	ret := _m.Called(ctx, id, bigQuantums)

	if len(ret) == 0 {
		panic("no return value specified for GetNetNotional")
	}

	var r0 *big.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, uint32, *big.Int) (*big.Int, error)); ok {
		return rf(ctx, id, bigQuantums)
	}
	if rf, ok := ret.Get(0).(func(types.Context, uint32, *big.Int) *big.Int); ok {
		r0 = rf(ctx, id, bigQuantums)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func(types.Context, uint32, *big.Int) error); ok {
		r1 = rf(ctx, id, bigQuantums)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetNotionalInBaseQuantums provides a mock function with given fields: ctx, id, bigQuoteQuantums
func (_m *PerpetualsKeeper) GetNotionalInBaseQuantums(ctx types.Context, id uint32, bigQuoteQuantums *big.Int) (*big.Int, error) {
	ret := _m.Called(ctx, id, bigQuoteQuantums)

	if len(ret) == 0 {
		panic("no return value specified for GetNotionalInBaseQuantums")
	}

	var r0 *big.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, uint32, *big.Int) (*big.Int, error)); ok {
		return rf(ctx, id, bigQuoteQuantums)
	}
	if rf, ok := ret.Get(0).(func(types.Context, uint32, *big.Int) *big.Int); ok {
		r0 = rf(ctx, id, bigQuoteQuantums)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func(types.Context, uint32, *big.Int) error); ok {
		r1 = rf(ctx, id, bigQuoteQuantums)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPerpetual provides a mock function with given fields: ctx, id
func (_m *PerpetualsKeeper) GetPerpetual(ctx types.Context, id uint32) (perpetualstypes.Perpetual, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetPerpetual")
	}

	var r0 perpetualstypes.Perpetual
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, uint32) (perpetualstypes.Perpetual, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(types.Context, uint32) perpetualstypes.Perpetual); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(perpetualstypes.Perpetual)
	}

	if rf, ok := ret.Get(1).(func(types.Context, uint32) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// HasAuthority provides a mock function with given fields: authority
func (_m *PerpetualsKeeper) HasAuthority(authority string) bool {
	ret := _m.Called(authority)

	if len(ret) == 0 {
		panic("no return value specified for HasAuthority")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(authority)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MaybeProcessNewFundingSampleEpoch provides a mock function with given fields: ctx
func (_m *PerpetualsKeeper) MaybeProcessNewFundingSampleEpoch(ctx types.Context) {
	_m.Called(ctx)
}

// MaybeProcessNewFundingTickEpoch provides a mock function with given fields: ctx
func (_m *PerpetualsKeeper) MaybeProcessNewFundingTickEpoch(ctx types.Context) {
	_m.Called(ctx)
}

// ModifyOpenInterest provides a mock function with given fields: ctx, perpetualId, openInterestDeltaBaseQuantums
func (_m *PerpetualsKeeper) ModifyOpenInterest(ctx types.Context, perpetualId uint32, openInterestDeltaBaseQuantums *big.Int) error {
	ret := _m.Called(ctx, perpetualId, openInterestDeltaBaseQuantums)

	if len(ret) == 0 {
		panic("no return value specified for ModifyOpenInterest")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, uint32, *big.Int) error); ok {
		r0 = rf(ctx, perpetualId, openInterestDeltaBaseQuantums)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ModifyPerpetual provides a mock function with given fields: ctx, id, ticker, marketId, defaultFundingPpm, liquidityTier
func (_m *PerpetualsKeeper) ModifyPerpetual(ctx types.Context, id uint32, ticker string, marketId uint32, defaultFundingPpm int32, liquidityTier uint32) (perpetualstypes.Perpetual, error) {
	ret := _m.Called(ctx, id, ticker, marketId, defaultFundingPpm, liquidityTier)

	if len(ret) == 0 {
		panic("no return value specified for ModifyPerpetual")
	}

	var r0 perpetualstypes.Perpetual
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, uint32, string, uint32, int32, uint32) (perpetualstypes.Perpetual, error)); ok {
		return rf(ctx, id, ticker, marketId, defaultFundingPpm, liquidityTier)
	}
	if rf, ok := ret.Get(0).(func(types.Context, uint32, string, uint32, int32, uint32) perpetualstypes.Perpetual); ok {
		r0 = rf(ctx, id, ticker, marketId, defaultFundingPpm, liquidityTier)
	} else {
		r0 = ret.Get(0).(perpetualstypes.Perpetual)
	}

	if rf, ok := ret.Get(1).(func(types.Context, uint32, string, uint32, int32, uint32) error); ok {
		r1 = rf(ctx, id, ticker, marketId, defaultFundingPpm, liquidityTier)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PerformStatefulPremiumVotesValidation provides a mock function with given fields: ctx, msg
func (_m *PerpetualsKeeper) PerformStatefulPremiumVotesValidation(ctx types.Context, msg *perpetualstypes.MsgAddPremiumVotes) error {
	ret := _m.Called(ctx, msg)

	if len(ret) == 0 {
		panic("no return value specified for PerformStatefulPremiumVotesValidation")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, *perpetualstypes.MsgAddPremiumVotes) error); ok {
		r0 = rf(ctx, msg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetLiquidityTier provides a mock function with given fields: ctx, id, name, initialMarginPpm, maintenanceFractionPpm, impactNotional, openInterestLowerCap, openInterestUpperCap
func (_m *PerpetualsKeeper) SetLiquidityTier(ctx types.Context, id uint32, name string, initialMarginPpm uint32, maintenanceFractionPpm uint32, impactNotional uint64, openInterestLowerCap uint64, openInterestUpperCap uint64) (perpetualstypes.LiquidityTier, error) {
	ret := _m.Called(ctx, id, name, initialMarginPpm, maintenanceFractionPpm, impactNotional, openInterestLowerCap, openInterestUpperCap)

	if len(ret) == 0 {
		panic("no return value specified for SetLiquidityTier")
	}

	var r0 perpetualstypes.LiquidityTier
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, uint32, string, uint32, uint32, uint64, uint64, uint64) (perpetualstypes.LiquidityTier, error)); ok {
		return rf(ctx, id, name, initialMarginPpm, maintenanceFractionPpm, impactNotional, openInterestLowerCap, openInterestUpperCap)
	}
	if rf, ok := ret.Get(0).(func(types.Context, uint32, string, uint32, uint32, uint64, uint64, uint64) perpetualstypes.LiquidityTier); ok {
		r0 = rf(ctx, id, name, initialMarginPpm, maintenanceFractionPpm, impactNotional, openInterestLowerCap, openInterestUpperCap)
	} else {
		r0 = ret.Get(0).(perpetualstypes.LiquidityTier)
	}

	if rf, ok := ret.Get(1).(func(types.Context, uint32, string, uint32, uint32, uint64, uint64, uint64) error); ok {
		r1 = rf(ctx, id, name, initialMarginPpm, maintenanceFractionPpm, impactNotional, openInterestLowerCap, openInterestUpperCap)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetNextPerpetualID provides a mock function with given fields: ctx, nextID
func (_m *PerpetualsKeeper) SetNextPerpetualID(ctx types.Context, nextID uint32) {
	_m.Called(ctx, nextID)
}

// SetParams provides a mock function with given fields: ctx, params
func (_m *PerpetualsKeeper) SetParams(ctx types.Context, params perpetualstypes.Params) error {
	ret := _m.Called(ctx, params)

	if len(ret) == 0 {
		panic("no return value specified for SetParams")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, perpetualstypes.Params) error); ok {
		r0 = rf(ctx, params)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetPerpetualMarketType provides a mock function with given fields: ctx, id, marketType
func (_m *PerpetualsKeeper) SetPerpetualMarketType(ctx types.Context, id uint32, marketType perpetualstypes.PerpetualMarketType) (perpetualstypes.Perpetual, error) {
	ret := _m.Called(ctx, id, marketType)

	if len(ret) == 0 {
		panic("no return value specified for SetPerpetualMarketType")
	}

	var r0 perpetualstypes.Perpetual
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, uint32, perpetualstypes.PerpetualMarketType) (perpetualstypes.Perpetual, error)); ok {
		return rf(ctx, id, marketType)
	}
	if rf, ok := ret.Get(0).(func(types.Context, uint32, perpetualstypes.PerpetualMarketType) perpetualstypes.Perpetual); ok {
		r0 = rf(ctx, id, marketType)
	} else {
		r0 = ret.Get(0).(perpetualstypes.Perpetual)
	}

	if rf, ok := ret.Get(1).(func(types.Context, uint32, perpetualstypes.PerpetualMarketType) error); ok {
		r1 = rf(ctx, id, marketType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ValidateAndSetPerpetual provides a mock function with given fields: ctx, perpetual
func (_m *PerpetualsKeeper) ValidateAndSetPerpetual(ctx types.Context, perpetual perpetualstypes.Perpetual) error {
	ret := _m.Called(ctx, perpetual)

	if len(ret) == 0 {
		panic("no return value specified for ValidateAndSetPerpetual")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, perpetualstypes.Perpetual) error); ok {
		r0 = rf(ctx, perpetual)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewPerpetualsKeeper creates a new instance of PerpetualsKeeper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPerpetualsKeeper(t interface {
	mock.TestingT
	Cleanup(func())
}) *PerpetualsKeeper {
	mock := &PerpetualsKeeper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
