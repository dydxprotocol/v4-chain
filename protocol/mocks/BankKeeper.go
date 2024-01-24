// Code generated by mockery v2.23.1. DO NOT EDIT.

package mocks

import (
	context "context"

	cosmos_sdktypes "github.com/cosmos/cosmos-sdk/types"
	keeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	mock "github.com/stretchr/testify/mock"

	query "github.com/cosmos/cosmos-sdk/types/query"

	types "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// BankKeeper is an autogenerated mock type for the Keeper type
type BankKeeper struct {
	mock.Mock
}

// AllBalances provides a mock function with given fields: _a0, _a1
func (_m *BankKeeper) AllBalances(_a0 context.Context, _a1 *types.QueryAllBalancesRequest) (*types.QueryAllBalancesResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *types.QueryAllBalancesResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryAllBalancesRequest) (*types.QueryAllBalancesResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryAllBalancesRequest) *types.QueryAllBalancesResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.QueryAllBalancesResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.QueryAllBalancesRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AppendSendRestriction provides a mock function with given fields: restriction
func (_m *BankKeeper) AppendSendRestriction(restriction types.SendRestrictionFn) {
	_m.Called(restriction)
}

// Balance provides a mock function with given fields: _a0, _a1
func (_m *BankKeeper) Balance(_a0 context.Context, _a1 *types.QueryBalanceRequest) (*types.QueryBalanceResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *types.QueryBalanceResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryBalanceRequest) (*types.QueryBalanceResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryBalanceRequest) *types.QueryBalanceResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.QueryBalanceResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.QueryBalanceRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BlockedAddr provides a mock function with given fields: addr
func (_m *BankKeeper) BlockedAddr(addr cosmos_sdktypes.AccAddress) bool {
	ret := _m.Called(addr)

	var r0 bool
	if rf, ok := ret.Get(0).(func(cosmos_sdktypes.AccAddress) bool); ok {
		r0 = rf(addr)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// BurnCoins provides a mock function with given fields: ctx, moduleName, amt
func (_m *BankKeeper) BurnCoins(ctx context.Context, moduleName string, amt cosmos_sdktypes.Coins) error {
	ret := _m.Called(ctx, moduleName, amt)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, cosmos_sdktypes.Coins) error); ok {
		r0 = rf(ctx, moduleName, amt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ClearSendRestriction provides a mock function with given fields:
func (_m *BankKeeper) ClearSendRestriction() {
	_m.Called()
}

// DelegateCoins provides a mock function with given fields: ctx, delegatorAddr, moduleAccAddr, amt
func (_m *BankKeeper) DelegateCoins(ctx context.Context, delegatorAddr cosmos_sdktypes.AccAddress, moduleAccAddr cosmos_sdktypes.AccAddress, amt cosmos_sdktypes.Coins) error {
	ret := _m.Called(ctx, delegatorAddr, moduleAccAddr, amt)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, cosmos_sdktypes.AccAddress, cosmos_sdktypes.AccAddress, cosmos_sdktypes.Coins) error); ok {
		r0 = rf(ctx, delegatorAddr, moduleAccAddr, amt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DelegateCoinsFromAccountToModule provides a mock function with given fields: ctx, senderAddr, recipientModule, amt
func (_m *BankKeeper) DelegateCoinsFromAccountToModule(ctx context.Context, senderAddr cosmos_sdktypes.AccAddress, recipientModule string, amt cosmos_sdktypes.Coins) error {
	ret := _m.Called(ctx, senderAddr, recipientModule, amt)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, cosmos_sdktypes.AccAddress, string, cosmos_sdktypes.Coins) error); ok {
		r0 = rf(ctx, senderAddr, recipientModule, amt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteSendEnabled provides a mock function with given fields: ctx, denoms
func (_m *BankKeeper) DeleteSendEnabled(ctx context.Context, denoms ...string) {
	_va := make([]interface{}, len(denoms))
	for _i := range denoms {
		_va[_i] = denoms[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// DenomMetadata provides a mock function with given fields: _a0, _a1
func (_m *BankKeeper) DenomMetadata(_a0 context.Context, _a1 *types.QueryDenomMetadataRequest) (*types.QueryDenomMetadataResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *types.QueryDenomMetadataResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryDenomMetadataRequest) (*types.QueryDenomMetadataResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryDenomMetadataRequest) *types.QueryDenomMetadataResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.QueryDenomMetadataResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.QueryDenomMetadataRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DenomMetadataByQueryString provides a mock function with given fields: _a0, _a1
func (_m *BankKeeper) DenomMetadataByQueryString(_a0 context.Context, _a1 *types.QueryDenomMetadataByQueryStringRequest) (*types.QueryDenomMetadataByQueryStringResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *types.QueryDenomMetadataByQueryStringResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryDenomMetadataByQueryStringRequest) (*types.QueryDenomMetadataByQueryStringResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryDenomMetadataByQueryStringRequest) *types.QueryDenomMetadataByQueryStringResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.QueryDenomMetadataByQueryStringResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.QueryDenomMetadataByQueryStringRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DenomOwners provides a mock function with given fields: _a0, _a1
func (_m *BankKeeper) DenomOwners(_a0 context.Context, _a1 *types.QueryDenomOwnersRequest) (*types.QueryDenomOwnersResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *types.QueryDenomOwnersResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryDenomOwnersRequest) (*types.QueryDenomOwnersResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryDenomOwnersRequest) *types.QueryDenomOwnersResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.QueryDenomOwnersResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.QueryDenomOwnersRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DenomOwnersByQuery provides a mock function with given fields: _a0, _a1
func (_m *BankKeeper) DenomOwnersByQuery(_a0 context.Context, _a1 *types.QueryDenomOwnersByQueryRequest) (*types.QueryDenomOwnersByQueryResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *types.QueryDenomOwnersByQueryResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryDenomOwnersByQueryRequest) (*types.QueryDenomOwnersByQueryResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryDenomOwnersByQueryRequest) *types.QueryDenomOwnersByQueryResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.QueryDenomOwnersByQueryResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.QueryDenomOwnersByQueryRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DenomsMetadata provides a mock function with given fields: _a0, _a1
func (_m *BankKeeper) DenomsMetadata(_a0 context.Context, _a1 *types.QueryDenomsMetadataRequest) (*types.QueryDenomsMetadataResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *types.QueryDenomsMetadataResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryDenomsMetadataRequest) (*types.QueryDenomsMetadataResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryDenomsMetadataRequest) *types.QueryDenomsMetadataResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.QueryDenomsMetadataResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.QueryDenomsMetadataRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ExportGenesis provides a mock function with given fields: _a0
func (_m *BankKeeper) ExportGenesis(_a0 context.Context) *types.GenesisState {
	ret := _m.Called(_a0)

	var r0 *types.GenesisState
	if rf, ok := ret.Get(0).(func(context.Context) *types.GenesisState); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.GenesisState)
		}
	}

	return r0
}

// GetAccountsBalances provides a mock function with given fields: ctx
func (_m *BankKeeper) GetAccountsBalances(ctx context.Context) []types.Balance {
	ret := _m.Called(ctx)

	var r0 []types.Balance
	if rf, ok := ret.Get(0).(func(context.Context) []types.Balance); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]types.Balance)
		}
	}

	return r0
}

// GetAllBalances provides a mock function with given fields: ctx, addr
func (_m *BankKeeper) GetAllBalances(ctx context.Context, addr cosmos_sdktypes.AccAddress) cosmos_sdktypes.Coins {
	ret := _m.Called(ctx, addr)

	var r0 cosmos_sdktypes.Coins
	if rf, ok := ret.Get(0).(func(context.Context, cosmos_sdktypes.AccAddress) cosmos_sdktypes.Coins); ok {
		r0 = rf(ctx, addr)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(cosmos_sdktypes.Coins)
		}
	}

	return r0
}

// GetAllDenomMetaData provides a mock function with given fields: ctx
func (_m *BankKeeper) GetAllDenomMetaData(ctx context.Context) []types.Metadata {
	ret := _m.Called(ctx)

	var r0 []types.Metadata
	if rf, ok := ret.Get(0).(func(context.Context) []types.Metadata); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]types.Metadata)
		}
	}

	return r0
}

// GetAllSendEnabledEntries provides a mock function with given fields: ctx
func (_m *BankKeeper) GetAllSendEnabledEntries(ctx context.Context) []types.SendEnabled {
	ret := _m.Called(ctx)

	var r0 []types.SendEnabled
	if rf, ok := ret.Get(0).(func(context.Context) []types.SendEnabled); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]types.SendEnabled)
		}
	}

	return r0
}

// GetAuthority provides a mock function with given fields:
func (_m *BankKeeper) GetAuthority() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetBalance provides a mock function with given fields: ctx, addr, denom
func (_m *BankKeeper) GetBalance(ctx context.Context, addr cosmos_sdktypes.AccAddress, denom string) cosmos_sdktypes.Coin {
	ret := _m.Called(ctx, addr, denom)

	var r0 cosmos_sdktypes.Coin
	if rf, ok := ret.Get(0).(func(context.Context, cosmos_sdktypes.AccAddress, string) cosmos_sdktypes.Coin); ok {
		r0 = rf(ctx, addr, denom)
	} else {
		r0 = ret.Get(0).(cosmos_sdktypes.Coin)
	}

	return r0
}

// GetBlockedAddresses provides a mock function with given fields:
func (_m *BankKeeper) GetBlockedAddresses() map[string]bool {
	ret := _m.Called()

	var r0 map[string]bool
	if rf, ok := ret.Get(0).(func() map[string]bool); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]bool)
		}
	}

	return r0
}

// GetDenomMetaData provides a mock function with given fields: ctx, denom
func (_m *BankKeeper) GetDenomMetaData(ctx context.Context, denom string) (types.Metadata, bool) {
	ret := _m.Called(ctx, denom)

	var r0 types.Metadata
	var r1 bool
	if rf, ok := ret.Get(0).(func(context.Context, string) (types.Metadata, bool)); ok {
		return rf(ctx, denom)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) types.Metadata); ok {
		r0 = rf(ctx, denom)
	} else {
		r0 = ret.Get(0).(types.Metadata)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) bool); ok {
		r1 = rf(ctx, denom)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// GetPaginatedTotalSupply provides a mock function with given fields: ctx, pagination
func (_m *BankKeeper) GetPaginatedTotalSupply(ctx context.Context, pagination *query.PageRequest) (cosmos_sdktypes.Coins, *query.PageResponse, error) {
	ret := _m.Called(ctx, pagination)

	var r0 cosmos_sdktypes.Coins
	var r1 *query.PageResponse
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, *query.PageRequest) (cosmos_sdktypes.Coins, *query.PageResponse, error)); ok {
		return rf(ctx, pagination)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *query.PageRequest) cosmos_sdktypes.Coins); ok {
		r0 = rf(ctx, pagination)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(cosmos_sdktypes.Coins)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *query.PageRequest) *query.PageResponse); ok {
		r1 = rf(ctx, pagination)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*query.PageResponse)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context, *query.PageRequest) error); ok {
		r2 = rf(ctx, pagination)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetParams provides a mock function with given fields: ctx
func (_m *BankKeeper) GetParams(ctx context.Context) types.Params {
	ret := _m.Called(ctx)

	var r0 types.Params
	if rf, ok := ret.Get(0).(func(context.Context) types.Params); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(types.Params)
	}

	return r0
}

// GetSendEnabledEntry provides a mock function with given fields: ctx, denom
func (_m *BankKeeper) GetSendEnabledEntry(ctx context.Context, denom string) (types.SendEnabled, bool) {
	ret := _m.Called(ctx, denom)

	var r0 types.SendEnabled
	var r1 bool
	if rf, ok := ret.Get(0).(func(context.Context, string) (types.SendEnabled, bool)); ok {
		return rf(ctx, denom)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) types.SendEnabled); ok {
		r0 = rf(ctx, denom)
	} else {
		r0 = ret.Get(0).(types.SendEnabled)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) bool); ok {
		r1 = rf(ctx, denom)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// GetSupply provides a mock function with given fields: ctx, denom
func (_m *BankKeeper) GetSupply(ctx context.Context, denom string) cosmos_sdktypes.Coin {
	ret := _m.Called(ctx, denom)

	var r0 cosmos_sdktypes.Coin
	if rf, ok := ret.Get(0).(func(context.Context, string) cosmos_sdktypes.Coin); ok {
		r0 = rf(ctx, denom)
	} else {
		r0 = ret.Get(0).(cosmos_sdktypes.Coin)
	}

	return r0
}

// HasBalance provides a mock function with given fields: ctx, addr, amt
func (_m *BankKeeper) HasBalance(ctx context.Context, addr cosmos_sdktypes.AccAddress, amt cosmos_sdktypes.Coin) bool {
	ret := _m.Called(ctx, addr, amt)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, cosmos_sdktypes.AccAddress, cosmos_sdktypes.Coin) bool); ok {
		r0 = rf(ctx, addr, amt)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// HasDenomMetaData provides a mock function with given fields: ctx, denom
func (_m *BankKeeper) HasDenomMetaData(ctx context.Context, denom string) bool {
	ret := _m.Called(ctx, denom)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(ctx, denom)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// HasSupply provides a mock function with given fields: ctx, denom
func (_m *BankKeeper) HasSupply(ctx context.Context, denom string) bool {
	ret := _m.Called(ctx, denom)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(ctx, denom)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// InitGenesis provides a mock function with given fields: _a0, _a1
func (_m *BankKeeper) InitGenesis(_a0 context.Context, _a1 *types.GenesisState) {
	_m.Called(_a0, _a1)
}

// InputOutputCoins provides a mock function with given fields: ctx, input, outputs
func (_m *BankKeeper) InputOutputCoins(ctx context.Context, input types.Input, outputs []types.Output) error {
	ret := _m.Called(ctx, input, outputs)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, types.Input, []types.Output) error); ok {
		r0 = rf(ctx, input, outputs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IsSendEnabledCoin provides a mock function with given fields: ctx, coin
func (_m *BankKeeper) IsSendEnabledCoin(ctx context.Context, coin cosmos_sdktypes.Coin) bool {
	ret := _m.Called(ctx, coin)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, cosmos_sdktypes.Coin) bool); ok {
		r0 = rf(ctx, coin)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// IsSendEnabledCoins provides a mock function with given fields: ctx, coins
func (_m *BankKeeper) IsSendEnabledCoins(ctx context.Context, coins ...cosmos_sdktypes.Coin) error {
	_va := make([]interface{}, len(coins))
	for _i := range coins {
		_va[_i] = coins[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ...cosmos_sdktypes.Coin) error); ok {
		r0 = rf(ctx, coins...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IsSendEnabledDenom provides a mock function with given fields: ctx, denom
func (_m *BankKeeper) IsSendEnabledDenom(ctx context.Context, denom string) bool {
	ret := _m.Called(ctx, denom)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(ctx, denom)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// IterateAccountBalances provides a mock function with given fields: ctx, addr, cb
func (_m *BankKeeper) IterateAccountBalances(ctx context.Context, addr cosmos_sdktypes.AccAddress, cb func(cosmos_sdktypes.Coin) bool) {
	_m.Called(ctx, addr, cb)
}

// IterateAllBalances provides a mock function with given fields: ctx, cb
func (_m *BankKeeper) IterateAllBalances(ctx context.Context, cb func(cosmos_sdktypes.AccAddress, cosmos_sdktypes.Coin) bool) {
	_m.Called(ctx, cb)
}

// IterateAllDenomMetaData provides a mock function with given fields: ctx, cb
func (_m *BankKeeper) IterateAllDenomMetaData(ctx context.Context, cb func(types.Metadata) bool) {
	_m.Called(ctx, cb)
}

// IterateSendEnabledEntries provides a mock function with given fields: ctx, cb
func (_m *BankKeeper) IterateSendEnabledEntries(ctx context.Context, cb func(string, bool) bool) {
	_m.Called(ctx, cb)
}

// IterateTotalSupply provides a mock function with given fields: ctx, cb
func (_m *BankKeeper) IterateTotalSupply(ctx context.Context, cb func(cosmos_sdktypes.Coin) bool) {
	_m.Called(ctx, cb)
}

// LockedCoins provides a mock function with given fields: ctx, addr
func (_m *BankKeeper) LockedCoins(ctx context.Context, addr cosmos_sdktypes.AccAddress) cosmos_sdktypes.Coins {
	ret := _m.Called(ctx, addr)

	var r0 cosmos_sdktypes.Coins
	if rf, ok := ret.Get(0).(func(context.Context, cosmos_sdktypes.AccAddress) cosmos_sdktypes.Coins); ok {
		r0 = rf(ctx, addr)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(cosmos_sdktypes.Coins)
		}
	}

	return r0
}

// MintCoins provides a mock function with given fields: ctx, moduleName, amt
func (_m *BankKeeper) MintCoins(ctx context.Context, moduleName string, amt cosmos_sdktypes.Coins) error {
	ret := _m.Called(ctx, moduleName, amt)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, cosmos_sdktypes.Coins) error); ok {
		r0 = rf(ctx, moduleName, amt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Params provides a mock function with given fields: _a0, _a1
func (_m *BankKeeper) Params(_a0 context.Context, _a1 *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *types.QueryParamsResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryParamsRequest) (*types.QueryParamsResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryParamsRequest) *types.QueryParamsResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.QueryParamsResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.QueryParamsRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PrependSendRestriction provides a mock function with given fields: restriction
func (_m *BankKeeper) PrependSendRestriction(restriction types.SendRestrictionFn) {
	_m.Called(restriction)
}

// SendCoins provides a mock function with given fields: ctx, fromAddr, toAddr, amt
func (_m *BankKeeper) SendCoins(ctx context.Context, fromAddr cosmos_sdktypes.AccAddress, toAddr cosmos_sdktypes.AccAddress, amt cosmos_sdktypes.Coins) error {
	ret := _m.Called(ctx, fromAddr, toAddr, amt)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, cosmos_sdktypes.AccAddress, cosmos_sdktypes.AccAddress, cosmos_sdktypes.Coins) error); ok {
		r0 = rf(ctx, fromAddr, toAddr, amt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendCoinsFromAccountToModule provides a mock function with given fields: ctx, senderAddr, recipientModule, amt
func (_m *BankKeeper) SendCoinsFromAccountToModule(ctx context.Context, senderAddr cosmos_sdktypes.AccAddress, recipientModule string, amt cosmos_sdktypes.Coins) error {
	ret := _m.Called(ctx, senderAddr, recipientModule, amt)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, cosmos_sdktypes.AccAddress, string, cosmos_sdktypes.Coins) error); ok {
		r0 = rf(ctx, senderAddr, recipientModule, amt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendCoinsFromModuleToAccount provides a mock function with given fields: ctx, senderModule, recipientAddr, amt
func (_m *BankKeeper) SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr cosmos_sdktypes.AccAddress, amt cosmos_sdktypes.Coins) error {
	ret := _m.Called(ctx, senderModule, recipientAddr, amt)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, cosmos_sdktypes.AccAddress, cosmos_sdktypes.Coins) error); ok {
		r0 = rf(ctx, senderModule, recipientAddr, amt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendCoinsFromModuleToModule provides a mock function with given fields: ctx, senderModule, recipientModule, amt
func (_m *BankKeeper) SendCoinsFromModuleToModule(ctx context.Context, senderModule string, recipientModule string, amt cosmos_sdktypes.Coins) error {
	ret := _m.Called(ctx, senderModule, recipientModule, amt)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, cosmos_sdktypes.Coins) error); ok {
		r0 = rf(ctx, senderModule, recipientModule, amt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendEnabled provides a mock function with given fields: _a0, _a1
func (_m *BankKeeper) SendEnabled(_a0 context.Context, _a1 *types.QuerySendEnabledRequest) (*types.QuerySendEnabledResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *types.QuerySendEnabledResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.QuerySendEnabledRequest) (*types.QuerySendEnabledResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.QuerySendEnabledRequest) *types.QuerySendEnabledResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.QuerySendEnabledResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.QuerySendEnabledRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetAllSendEnabled provides a mock function with given fields: ctx, sendEnableds
func (_m *BankKeeper) SetAllSendEnabled(ctx context.Context, sendEnableds []*types.SendEnabled) {
	_m.Called(ctx, sendEnableds)
}

// SetDenomMetaData provides a mock function with given fields: ctx, denomMetaData
func (_m *BankKeeper) SetDenomMetaData(ctx context.Context, denomMetaData types.Metadata) {
	_m.Called(ctx, denomMetaData)
}

// SetParams provides a mock function with given fields: ctx, params
func (_m *BankKeeper) SetParams(ctx context.Context, params types.Params) error {
	ret := _m.Called(ctx, params)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, types.Params) error); ok {
		r0 = rf(ctx, params)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetSendEnabled provides a mock function with given fields: ctx, denom, value
func (_m *BankKeeper) SetSendEnabled(ctx context.Context, denom string, value bool) {
	_m.Called(ctx, denom, value)
}

// SpendableBalanceByDenom provides a mock function with given fields: _a0, _a1
func (_m *BankKeeper) SpendableBalanceByDenom(_a0 context.Context, _a1 *types.QuerySpendableBalanceByDenomRequest) (*types.QuerySpendableBalanceByDenomResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *types.QuerySpendableBalanceByDenomResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.QuerySpendableBalanceByDenomRequest) (*types.QuerySpendableBalanceByDenomResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.QuerySpendableBalanceByDenomRequest) *types.QuerySpendableBalanceByDenomResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.QuerySpendableBalanceByDenomResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.QuerySpendableBalanceByDenomRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SpendableBalances provides a mock function with given fields: _a0, _a1
func (_m *BankKeeper) SpendableBalances(_a0 context.Context, _a1 *types.QuerySpendableBalancesRequest) (*types.QuerySpendableBalancesResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *types.QuerySpendableBalancesResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.QuerySpendableBalancesRequest) (*types.QuerySpendableBalancesResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.QuerySpendableBalancesRequest) *types.QuerySpendableBalancesResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.QuerySpendableBalancesResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.QuerySpendableBalancesRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SpendableCoin provides a mock function with given fields: ctx, addr, denom
func (_m *BankKeeper) SpendableCoin(ctx context.Context, addr cosmos_sdktypes.AccAddress, denom string) cosmos_sdktypes.Coin {
	ret := _m.Called(ctx, addr, denom)

	var r0 cosmos_sdktypes.Coin
	if rf, ok := ret.Get(0).(func(context.Context, cosmos_sdktypes.AccAddress, string) cosmos_sdktypes.Coin); ok {
		r0 = rf(ctx, addr, denom)
	} else {
		r0 = ret.Get(0).(cosmos_sdktypes.Coin)
	}

	return r0
}

// SpendableCoins provides a mock function with given fields: ctx, addr
func (_m *BankKeeper) SpendableCoins(ctx context.Context, addr cosmos_sdktypes.AccAddress) cosmos_sdktypes.Coins {
	ret := _m.Called(ctx, addr)

	var r0 cosmos_sdktypes.Coins
	if rf, ok := ret.Get(0).(func(context.Context, cosmos_sdktypes.AccAddress) cosmos_sdktypes.Coins); ok {
		r0 = rf(ctx, addr)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(cosmos_sdktypes.Coins)
		}
	}

	return r0
}

// SupplyOf provides a mock function with given fields: _a0, _a1
func (_m *BankKeeper) SupplyOf(_a0 context.Context, _a1 *types.QuerySupplyOfRequest) (*types.QuerySupplyOfResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *types.QuerySupplyOfResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.QuerySupplyOfRequest) (*types.QuerySupplyOfResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.QuerySupplyOfRequest) *types.QuerySupplyOfResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.QuerySupplyOfResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.QuerySupplyOfRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TotalSupply provides a mock function with given fields: _a0, _a1
func (_m *BankKeeper) TotalSupply(_a0 context.Context, _a1 *types.QueryTotalSupplyRequest) (*types.QueryTotalSupplyResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *types.QueryTotalSupplyResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryTotalSupplyRequest) (*types.QueryTotalSupplyResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryTotalSupplyRequest) *types.QueryTotalSupplyResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.QueryTotalSupplyResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.QueryTotalSupplyRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UndelegateCoins provides a mock function with given fields: ctx, moduleAccAddr, delegatorAddr, amt
func (_m *BankKeeper) UndelegateCoins(ctx context.Context, moduleAccAddr cosmos_sdktypes.AccAddress, delegatorAddr cosmos_sdktypes.AccAddress, amt cosmos_sdktypes.Coins) error {
	ret := _m.Called(ctx, moduleAccAddr, delegatorAddr, amt)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, cosmos_sdktypes.AccAddress, cosmos_sdktypes.AccAddress, cosmos_sdktypes.Coins) error); ok {
		r0 = rf(ctx, moduleAccAddr, delegatorAddr, amt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UndelegateCoinsFromModuleToAccount provides a mock function with given fields: ctx, senderModule, recipientAddr, amt
func (_m *BankKeeper) UndelegateCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr cosmos_sdktypes.AccAddress, amt cosmos_sdktypes.Coins) error {
	ret := _m.Called(ctx, senderModule, recipientAddr, amt)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, cosmos_sdktypes.AccAddress, cosmos_sdktypes.Coins) error); ok {
		r0 = rf(ctx, senderModule, recipientAddr, amt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ValidateBalance provides a mock function with given fields: ctx, addr
func (_m *BankKeeper) ValidateBalance(ctx context.Context, addr cosmos_sdktypes.AccAddress) error {
	ret := _m.Called(ctx, addr)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, cosmos_sdktypes.AccAddress) error); ok {
		r0 = rf(ctx, addr)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WithMintCoinsRestriction provides a mock function with given fields: _a0
func (_m *BankKeeper) WithMintCoinsRestriction(_a0 types.MintingRestrictionFn) keeper.BaseKeeper {
	ret := _m.Called(_a0)

	var r0 keeper.BaseKeeper
	if rf, ok := ret.Get(0).(func(types.MintingRestrictionFn) keeper.BaseKeeper); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(keeper.BaseKeeper)
	}

	return r0
}

type mockConstructorTestingTNewBankKeeper interface {
	mock.TestingT
	Cleanup(func())
}

// NewBankKeeper creates a new instance of BankKeeper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewBankKeeper(t mockConstructorTestingTNewBankKeeper) *BankKeeper {
	mock := &BankKeeper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
