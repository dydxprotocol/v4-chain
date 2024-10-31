// Code generated by mockery v2.42.0. DO NOT EDIT.

package mocks

import (
	context "context"

	crypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	cosmos_sdktypes "github.com/cosmos/cosmos-sdk/types"

	math "cosmossdk.io/math"

	mock "github.com/stretchr/testify/mock"

	types "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// ValidatorStore is an autogenerated mock type for the ValidatorStore type
type ValidatorStore struct {
	mock.Mock
}

// GetAllValidators provides a mock function with given fields: ctx
func (_m *ValidatorStore) GetAllValidators(ctx context.Context) ([]types.Validator, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAllValidators")
	}

	var r0 []types.Validator
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]types.Validator, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []types.Validator); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]types.Validator)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPubKeyByConsAddr provides a mock function with given fields: _a0, _a1
func (_m *ValidatorStore) GetPubKeyByConsAddr(_a0 context.Context, _a1 cosmos_sdktypes.ConsAddress) (crypto.PublicKey, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetPubKeyByConsAddr")
	}

	var r0 crypto.PublicKey
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, cosmos_sdktypes.ConsAddress) (crypto.PublicKey, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, cosmos_sdktypes.ConsAddress) crypto.PublicKey); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(crypto.PublicKey)
	}

	if rf, ok := ret.Get(1).(func(context.Context, cosmos_sdktypes.ConsAddress) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetValidator provides a mock function with given fields: ctx, valAddr
func (_m *ValidatorStore) GetValidator(ctx context.Context, valAddr cosmos_sdktypes.ValAddress) (types.Validator, error) {
	ret := _m.Called(ctx, valAddr)

	if len(ret) == 0 {
		panic("no return value specified for GetValidator")
	}

	var r0 types.Validator
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, cosmos_sdktypes.ValAddress) (types.Validator, error)); ok {
		return rf(ctx, valAddr)
	}
	if rf, ok := ret.Get(0).(func(context.Context, cosmos_sdktypes.ValAddress) types.Validator); ok {
		r0 = rf(ctx, valAddr)
	} else {
		r0 = ret.Get(0).(types.Validator)
	}

	if rf, ok := ret.Get(1).(func(context.Context, cosmos_sdktypes.ValAddress) error); ok {
		r1 = rf(ctx, valAddr)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TotalBondedTokens provides a mock function with given fields: ctx
func (_m *ValidatorStore) TotalBondedTokens(ctx context.Context) (math.Int, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for TotalBondedTokens")
	}

	var r0 math.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (math.Int, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) math.Int); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(math.Int)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ValidatorByConsAddr provides a mock function with given fields: ctx, addr
func (_m *ValidatorStore) ValidatorByConsAddr(ctx context.Context, addr cosmos_sdktypes.ConsAddress) (types.ValidatorI, error) {
	ret := _m.Called(ctx, addr)

	if len(ret) == 0 {
		panic("no return value specified for ValidatorByConsAddr")
	}

	var r0 types.ValidatorI
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, cosmos_sdktypes.ConsAddress) (types.ValidatorI, error)); ok {
		return rf(ctx, addr)
	}
	if rf, ok := ret.Get(0).(func(context.Context, cosmos_sdktypes.ConsAddress) types.ValidatorI); ok {
		r0 = rf(ctx, addr)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.ValidatorI)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, cosmos_sdktypes.ConsAddress) error); ok {
		r1 = rf(ctx, addr)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewValidatorStore creates a new instance of ValidatorStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewValidatorStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *ValidatorStore {
	mock := &ValidatorStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
