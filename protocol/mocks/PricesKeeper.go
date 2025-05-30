// Code generated by mockery v2.50.0. DO NOT EDIT.

package mocks

import (
	log "cosmossdk.io/log"
	oracletypes "github.com/dydxprotocol/slinky/x/oracle/types"
	mock "github.com/stretchr/testify/mock"

	pkgtypes "github.com/dydxprotocol/slinky/pkg/types"

	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"

	types "github.com/cosmos/cosmos-sdk/types"
)

// PricesKeeper is an autogenerated mock type for the PricesKeeper type
type PricesKeeper struct {
	mock.Mock
}

// AddCurrencyPairIDToStore provides a mock function with given fields: ctx, id, cp
func (_m *PricesKeeper) AddCurrencyPairIDToStore(ctx types.Context, id uint32, cp pkgtypes.CurrencyPair) {
	_m.Called(ctx, id, cp)
}

// CreateMarket provides a mock function with given fields: ctx, param, price
func (_m *PricesKeeper) CreateMarket(ctx types.Context, param pricestypes.MarketParam, price pricestypes.MarketPrice) (pricestypes.MarketParam, error) {
	ret := _m.Called(ctx, param, price)

	if len(ret) == 0 {
		panic("no return value specified for CreateMarket")
	}

	var r0 pricestypes.MarketParam
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, pricestypes.MarketParam, pricestypes.MarketPrice) (pricestypes.MarketParam, error)); ok {
		return rf(ctx, param, price)
	}
	if rf, ok := ret.Get(0).(func(types.Context, pricestypes.MarketParam, pricestypes.MarketPrice) pricestypes.MarketParam); ok {
		r0 = rf(ctx, param, price)
	} else {
		r0 = ret.Get(0).(pricestypes.MarketParam)
	}

	if rf, ok := ret.Get(1).(func(types.Context, pricestypes.MarketParam, pricestypes.MarketPrice) error); ok {
		r1 = rf(ctx, param, price)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllMarketParamPrices provides a mock function with given fields: ctx
func (_m *PricesKeeper) GetAllMarketParamPrices(ctx types.Context) ([]pricestypes.MarketParamPrice, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAllMarketParamPrices")
	}

	var r0 []pricestypes.MarketParamPrice
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context) ([]pricestypes.MarketParamPrice, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(types.Context) []pricestypes.MarketParamPrice); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]pricestypes.MarketParamPrice)
		}
	}

	if rf, ok := ret.Get(1).(func(types.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllMarketParams provides a mock function with given fields: ctx
func (_m *PricesKeeper) GetAllMarketParams(ctx types.Context) []pricestypes.MarketParam {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAllMarketParams")
	}

	var r0 []pricestypes.MarketParam
	if rf, ok := ret.Get(0).(func(types.Context) []pricestypes.MarketParam); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]pricestypes.MarketParam)
		}
	}

	return r0
}

// GetAllMarketPrices provides a mock function with given fields: ctx
func (_m *PricesKeeper) GetAllMarketPrices(ctx types.Context) []pricestypes.MarketPrice {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAllMarketPrices")
	}

	var r0 []pricestypes.MarketPrice
	if rf, ok := ret.Get(0).(func(types.Context) []pricestypes.MarketPrice); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]pricestypes.MarketPrice)
		}
	}

	return r0
}

// GetCurrencyPairFromID provides a mock function with given fields: ctx, id
func (_m *PricesKeeper) GetCurrencyPairFromID(ctx types.Context, id uint64) (pkgtypes.CurrencyPair, bool) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetCurrencyPairFromID")
	}

	var r0 pkgtypes.CurrencyPair
	var r1 bool
	if rf, ok := ret.Get(0).(func(types.Context, uint64) (pkgtypes.CurrencyPair, bool)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(types.Context, uint64) pkgtypes.CurrencyPair); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(pkgtypes.CurrencyPair)
	}

	if rf, ok := ret.Get(1).(func(types.Context, uint64) bool); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// GetExponent provides a mock function with given fields: ctx, ticker
func (_m *PricesKeeper) GetExponent(ctx types.Context, ticker string) (int32, error) {
	ret := _m.Called(ctx, ticker)

	if len(ret) == 0 {
		panic("no return value specified for GetExponent")
	}

	var r0 int32
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, string) (int32, error)); ok {
		return rf(ctx, ticker)
	}
	if rf, ok := ret.Get(0).(func(types.Context, string) int32); ok {
		r0 = rf(ctx, ticker)
	} else {
		r0 = ret.Get(0).(int32)
	}

	if rf, ok := ret.Get(1).(func(types.Context, string) error); ok {
		r1 = rf(ctx, ticker)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetIDForCurrencyPair provides a mock function with given fields: ctx, cp
func (_m *PricesKeeper) GetIDForCurrencyPair(ctx types.Context, cp pkgtypes.CurrencyPair) (uint64, bool) {
	ret := _m.Called(ctx, cp)

	if len(ret) == 0 {
		panic("no return value specified for GetIDForCurrencyPair")
	}

	var r0 uint64
	var r1 bool
	if rf, ok := ret.Get(0).(func(types.Context, pkgtypes.CurrencyPair) (uint64, bool)); ok {
		return rf(ctx, cp)
	}
	if rf, ok := ret.Get(0).(func(types.Context, pkgtypes.CurrencyPair) uint64); ok {
		r0 = rf(ctx, cp)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func(types.Context, pkgtypes.CurrencyPair) bool); ok {
		r1 = rf(ctx, cp)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// GetMarketIdToValidIndexPrice provides a mock function with given fields: ctx
func (_m *PricesKeeper) GetMarketIdToValidIndexPrice(ctx types.Context) map[uint32]pricestypes.MarketPrice {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetMarketIdToValidIndexPrice")
	}

	var r0 map[uint32]pricestypes.MarketPrice
	if rf, ok := ret.Get(0).(func(types.Context) map[uint32]pricestypes.MarketPrice); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[uint32]pricestypes.MarketPrice)
		}
	}

	return r0
}

// GetMarketParam provides a mock function with given fields: ctx, id
func (_m *PricesKeeper) GetMarketParam(ctx types.Context, id uint32) (pricestypes.MarketParam, bool) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetMarketParam")
	}

	var r0 pricestypes.MarketParam
	var r1 bool
	if rf, ok := ret.Get(0).(func(types.Context, uint32) (pricestypes.MarketParam, bool)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(types.Context, uint32) pricestypes.MarketParam); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(pricestypes.MarketParam)
	}

	if rf, ok := ret.Get(1).(func(types.Context, uint32) bool); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// GetMarketPrice provides a mock function with given fields: ctx, id
func (_m *PricesKeeper) GetMarketPrice(ctx types.Context, id uint32) (pricestypes.MarketPrice, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetMarketPrice")
	}

	var r0 pricestypes.MarketPrice
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, uint32) (pricestypes.MarketPrice, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(types.Context, uint32) pricestypes.MarketPrice); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(pricestypes.MarketPrice)
	}

	if rf, ok := ret.Get(1).(func(types.Context, uint32) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPriceForCurrencyPair provides a mock function with given fields: ctx, cp
func (_m *PricesKeeper) GetPriceForCurrencyPair(ctx types.Context, cp pkgtypes.CurrencyPair) (oracletypes.QuotePrice, error) {
	ret := _m.Called(ctx, cp)

	if len(ret) == 0 {
		panic("no return value specified for GetPriceForCurrencyPair")
	}

	var r0 oracletypes.QuotePrice
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, pkgtypes.CurrencyPair) (oracletypes.QuotePrice, error)); ok {
		return rf(ctx, cp)
	}
	if rf, ok := ret.Get(0).(func(types.Context, pkgtypes.CurrencyPair) oracletypes.QuotePrice); ok {
		r0 = rf(ctx, cp)
	} else {
		r0 = ret.Get(0).(oracletypes.QuotePrice)
	}

	if rf, ok := ret.Get(1).(func(types.Context, pkgtypes.CurrencyPair) error); ok {
		r1 = rf(ctx, cp)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetValidMarketPriceUpdates provides a mock function with given fields: ctx
func (_m *PricesKeeper) GetValidMarketPriceUpdates(ctx types.Context) *pricestypes.MsgUpdateMarketPrices {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetValidMarketPriceUpdates")
	}

	var r0 *pricestypes.MsgUpdateMarketPrices
	if rf, ok := ret.Get(0).(func(types.Context) *pricestypes.MsgUpdateMarketPrices); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pricestypes.MsgUpdateMarketPrices)
		}
	}

	return r0
}

// HasAuthority provides a mock function with given fields: authority
func (_m *PricesKeeper) HasAuthority(authority string) bool {
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

// Logger provides a mock function with given fields: ctx
func (_m *PricesKeeper) Logger(ctx types.Context) log.Logger {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Logger")
	}

	var r0 log.Logger
	if rf, ok := ret.Get(0).(func(types.Context) log.Logger); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(log.Logger)
		}
	}

	return r0
}

// ModifyMarketParam provides a mock function with given fields: ctx, param
func (_m *PricesKeeper) ModifyMarketParam(ctx types.Context, param pricestypes.MarketParam) (pricestypes.MarketParam, error) {
	ret := _m.Called(ctx, param)

	if len(ret) == 0 {
		panic("no return value specified for ModifyMarketParam")
	}

	var r0 pricestypes.MarketParam
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, pricestypes.MarketParam) (pricestypes.MarketParam, error)); ok {
		return rf(ctx, param)
	}
	if rf, ok := ret.Get(0).(func(types.Context, pricestypes.MarketParam) pricestypes.MarketParam); ok {
		r0 = rf(ctx, param)
	} else {
		r0 = ret.Get(0).(pricestypes.MarketParam)
	}

	if rf, ok := ret.Get(1).(func(types.Context, pricestypes.MarketParam) error); ok {
		r1 = rf(ctx, param)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PerformStatefulPriceUpdateValidation provides a mock function with given fields: ctx, marketPriceUpdates, performNonDeterministicValidation
func (_m *PricesKeeper) PerformStatefulPriceUpdateValidation(ctx types.Context, marketPriceUpdates *pricestypes.MsgUpdateMarketPrices, performNonDeterministicValidation bool) error {
	ret := _m.Called(ctx, marketPriceUpdates, performNonDeterministicValidation)

	if len(ret) == 0 {
		panic("no return value specified for PerformStatefulPriceUpdateValidation")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, *pricestypes.MsgUpdateMarketPrices, bool) error); ok {
		r0 = rf(ctx, marketPriceUpdates, performNonDeterministicValidation)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetNextMarketID provides a mock function with given fields: ctx, nextID
func (_m *PricesKeeper) SetNextMarketID(ctx types.Context, nextID uint32) {
	_m.Called(ctx, nextID)
}

// UpdateMarketPrices provides a mock function with given fields: ctx, updates
func (_m *PricesKeeper) UpdateMarketPrices(ctx types.Context, updates []*pricestypes.MsgUpdateMarketPrices_MarketPrice) error {
	ret := _m.Called(ctx, updates)

	if len(ret) == 0 {
		panic("no return value specified for UpdateMarketPrices")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, []*pricestypes.MsgUpdateMarketPrices_MarketPrice) error); ok {
		r0 = rf(ctx, updates)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewPricesKeeper creates a new instance of PricesKeeper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPricesKeeper(t interface {
	mock.TestingT
	Cleanup(func())
}) *PricesKeeper {
	mock := &PricesKeeper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
