// Code generated by mockery v2.23.1. DO NOT EDIT.

package mocks

import (
	api "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/bridge/api"
	clobtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"

	context "context"

	grpc "google.golang.org/grpc"

	liquidationapi "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/liquidation/api"

	mock "github.com/stretchr/testify/mock"

	perpetualstypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"

	pricefeedapi "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/pricefeed/api"

	pricestypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"

	subaccountstypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"

	types "github.com/StreamFinance-Protocol/stream-chain/protocol/x/blocktime/types"
)

// QueryClient is an autogenerated mock type for the QueryClient type
type QueryClient struct {
	mock.Mock
}

// AddBridgeEvents provides a mock function with given fields: ctx, in, opts
func (_m *QueryClient) AddBridgeEvents(ctx context.Context, in *api.AddBridgeEventsRequest, opts ...grpc.CallOption) (*api.AddBridgeEventsResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *api.AddBridgeEventsResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *api.AddBridgeEventsRequest, ...grpc.CallOption) (*api.AddBridgeEventsResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *api.AddBridgeEventsRequest, ...grpc.CallOption) *api.AddBridgeEventsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.AddBridgeEventsResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *api.AddBridgeEventsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AllDowntimeInfo provides a mock function with given fields: ctx, in, opts
func (_m *QueryClient) AllDowntimeInfo(ctx context.Context, in *types.QueryAllDowntimeInfoRequest, opts ...grpc.CallOption) (*types.QueryAllDowntimeInfoResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *types.QueryAllDowntimeInfoResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryAllDowntimeInfoRequest, ...grpc.CallOption) (*types.QueryAllDowntimeInfoResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryAllDowntimeInfoRequest, ...grpc.CallOption) *types.QueryAllDowntimeInfoResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.QueryAllDowntimeInfoResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.QueryAllDowntimeInfoRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AllLiquidityTiers provides a mock function with given fields: ctx, in, opts
func (_m *QueryClient) AllLiquidityTiers(ctx context.Context, in *perpetualstypes.QueryAllLiquidityTiersRequest, opts ...grpc.CallOption) (*perpetualstypes.QueryAllLiquidityTiersResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *perpetualstypes.QueryAllLiquidityTiersResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *perpetualstypes.QueryAllLiquidityTiersRequest, ...grpc.CallOption) (*perpetualstypes.QueryAllLiquidityTiersResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *perpetualstypes.QueryAllLiquidityTiersRequest, ...grpc.CallOption) *perpetualstypes.QueryAllLiquidityTiersResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*perpetualstypes.QueryAllLiquidityTiersResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *perpetualstypes.QueryAllLiquidityTiersRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AllMarketParams provides a mock function with given fields: ctx, in, opts
func (_m *QueryClient) AllMarketParams(ctx context.Context, in *pricestypes.QueryAllMarketParamsRequest, opts ...grpc.CallOption) (*pricestypes.QueryAllMarketParamsResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *pricestypes.QueryAllMarketParamsResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *pricestypes.QueryAllMarketParamsRequest, ...grpc.CallOption) (*pricestypes.QueryAllMarketParamsResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *pricestypes.QueryAllMarketParamsRequest, ...grpc.CallOption) *pricestypes.QueryAllMarketParamsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pricestypes.QueryAllMarketParamsResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *pricestypes.QueryAllMarketParamsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AllMarketPrices provides a mock function with given fields: ctx, in, opts
func (_m *QueryClient) AllMarketPrices(ctx context.Context, in *pricestypes.QueryAllMarketPricesRequest, opts ...grpc.CallOption) (*pricestypes.QueryAllMarketPricesResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *pricestypes.QueryAllMarketPricesResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *pricestypes.QueryAllMarketPricesRequest, ...grpc.CallOption) (*pricestypes.QueryAllMarketPricesResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *pricestypes.QueryAllMarketPricesRequest, ...grpc.CallOption) *pricestypes.QueryAllMarketPricesResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pricestypes.QueryAllMarketPricesResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *pricestypes.QueryAllMarketPricesRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AllPerpetuals provides a mock function with given fields: ctx, in, opts
func (_m *QueryClient) AllPerpetuals(ctx context.Context, in *perpetualstypes.QueryAllPerpetualsRequest, opts ...grpc.CallOption) (*perpetualstypes.QueryAllPerpetualsResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *perpetualstypes.QueryAllPerpetualsResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *perpetualstypes.QueryAllPerpetualsRequest, ...grpc.CallOption) (*perpetualstypes.QueryAllPerpetualsResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *perpetualstypes.QueryAllPerpetualsRequest, ...grpc.CallOption) *perpetualstypes.QueryAllPerpetualsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*perpetualstypes.QueryAllPerpetualsResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *perpetualstypes.QueryAllPerpetualsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BlockRateLimitConfiguration provides a mock function with given fields: ctx, in, opts
func (_m *QueryClient) BlockRateLimitConfiguration(ctx context.Context, in *clobtypes.QueryBlockRateLimitConfigurationRequest, opts ...grpc.CallOption) (*clobtypes.QueryBlockRateLimitConfigurationResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *clobtypes.QueryBlockRateLimitConfigurationResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *clobtypes.QueryBlockRateLimitConfigurationRequest, ...grpc.CallOption) (*clobtypes.QueryBlockRateLimitConfigurationResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *clobtypes.QueryBlockRateLimitConfigurationRequest, ...grpc.CallOption) *clobtypes.QueryBlockRateLimitConfigurationResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*clobtypes.QueryBlockRateLimitConfigurationResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *clobtypes.QueryBlockRateLimitConfigurationRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ClobPair provides a mock function with given fields: ctx, in, opts
func (_m *QueryClient) ClobPair(ctx context.Context, in *clobtypes.QueryGetClobPairRequest, opts ...grpc.CallOption) (*clobtypes.QueryClobPairResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *clobtypes.QueryClobPairResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *clobtypes.QueryGetClobPairRequest, ...grpc.CallOption) (*clobtypes.QueryClobPairResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *clobtypes.QueryGetClobPairRequest, ...grpc.CallOption) *clobtypes.QueryClobPairResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*clobtypes.QueryClobPairResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *clobtypes.QueryGetClobPairRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ClobPairAll provides a mock function with given fields: ctx, in, opts
func (_m *QueryClient) ClobPairAll(ctx context.Context, in *clobtypes.QueryAllClobPairRequest, opts ...grpc.CallOption) (*clobtypes.QueryClobPairAllResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *clobtypes.QueryClobPairAllResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *clobtypes.QueryAllClobPairRequest, ...grpc.CallOption) (*clobtypes.QueryClobPairAllResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *clobtypes.QueryAllClobPairRequest, ...grpc.CallOption) *clobtypes.QueryClobPairAllResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*clobtypes.QueryClobPairAllResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *clobtypes.QueryAllClobPairRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DowntimeParams provides a mock function with given fields: ctx, in, opts
func (_m *QueryClient) DowntimeParams(ctx context.Context, in *types.QueryDowntimeParamsRequest, opts ...grpc.CallOption) (*types.QueryDowntimeParamsResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *types.QueryDowntimeParamsResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryDowntimeParamsRequest, ...grpc.CallOption) (*types.QueryDowntimeParamsResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryDowntimeParamsRequest, ...grpc.CallOption) *types.QueryDowntimeParamsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.QueryDowntimeParamsResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.QueryDowntimeParamsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EquityTierLimitConfiguration provides a mock function with given fields: ctx, in, opts
func (_m *QueryClient) EquityTierLimitConfiguration(ctx context.Context, in *clobtypes.QueryEquityTierLimitConfigurationRequest, opts ...grpc.CallOption) (*clobtypes.QueryEquityTierLimitConfigurationResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *clobtypes.QueryEquityTierLimitConfigurationResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *clobtypes.QueryEquityTierLimitConfigurationRequest, ...grpc.CallOption) (*clobtypes.QueryEquityTierLimitConfigurationResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *clobtypes.QueryEquityTierLimitConfigurationRequest, ...grpc.CallOption) *clobtypes.QueryEquityTierLimitConfigurationResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*clobtypes.QueryEquityTierLimitConfigurationResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *clobtypes.QueryEquityTierLimitConfigurationRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetWithdrawalAndTransfersBlockedInfo provides a mock function with given fields: ctx, in, opts
func (_m *QueryClient) GetWithdrawalAndTransfersBlockedInfo(ctx context.Context, in *subaccountstypes.QueryGetWithdrawalAndTransfersBlockedInfoRequest, opts ...grpc.CallOption) (*subaccountstypes.QueryGetWithdrawalAndTransfersBlockedInfoResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *subaccountstypes.QueryGetWithdrawalAndTransfersBlockedInfoResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *subaccountstypes.QueryGetWithdrawalAndTransfersBlockedInfoRequest, ...grpc.CallOption) (*subaccountstypes.QueryGetWithdrawalAndTransfersBlockedInfoResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *subaccountstypes.QueryGetWithdrawalAndTransfersBlockedInfoRequest, ...grpc.CallOption) *subaccountstypes.QueryGetWithdrawalAndTransfersBlockedInfoResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*subaccountstypes.QueryGetWithdrawalAndTransfersBlockedInfoResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *subaccountstypes.QueryGetWithdrawalAndTransfersBlockedInfoRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LiquidateSubaccounts provides a mock function with given fields: ctx, in, opts
func (_m *QueryClient) LiquidateSubaccounts(ctx context.Context, in *liquidationapi.LiquidateSubaccountsRequest, opts ...grpc.CallOption) (*liquidationapi.LiquidateSubaccountsResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *liquidationapi.LiquidateSubaccountsResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *liquidationapi.LiquidateSubaccountsRequest, ...grpc.CallOption) (*liquidationapi.LiquidateSubaccountsResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *liquidationapi.LiquidateSubaccountsRequest, ...grpc.CallOption) *liquidationapi.LiquidateSubaccountsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*liquidationapi.LiquidateSubaccountsResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *liquidationapi.LiquidateSubaccountsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LiquidationsConfiguration provides a mock function with given fields: ctx, in, opts
func (_m *QueryClient) LiquidationsConfiguration(ctx context.Context, in *clobtypes.QueryLiquidationsConfigurationRequest, opts ...grpc.CallOption) (*clobtypes.QueryLiquidationsConfigurationResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *clobtypes.QueryLiquidationsConfigurationResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *clobtypes.QueryLiquidationsConfigurationRequest, ...grpc.CallOption) (*clobtypes.QueryLiquidationsConfigurationResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *clobtypes.QueryLiquidationsConfigurationRequest, ...grpc.CallOption) *clobtypes.QueryLiquidationsConfigurationResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*clobtypes.QueryLiquidationsConfigurationResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *clobtypes.QueryLiquidationsConfigurationRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MarketParam provides a mock function with given fields: ctx, in, opts
func (_m *QueryClient) MarketParam(ctx context.Context, in *pricestypes.QueryMarketParamRequest, opts ...grpc.CallOption) (*pricestypes.QueryMarketParamResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *pricestypes.QueryMarketParamResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *pricestypes.QueryMarketParamRequest, ...grpc.CallOption) (*pricestypes.QueryMarketParamResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *pricestypes.QueryMarketParamRequest, ...grpc.CallOption) *pricestypes.QueryMarketParamResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pricestypes.QueryMarketParamResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *pricestypes.QueryMarketParamRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MarketPrice provides a mock function with given fields: ctx, in, opts
func (_m *QueryClient) MarketPrice(ctx context.Context, in *pricestypes.QueryMarketPriceRequest, opts ...grpc.CallOption) (*pricestypes.QueryMarketPriceResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *pricestypes.QueryMarketPriceResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *pricestypes.QueryMarketPriceRequest, ...grpc.CallOption) (*pricestypes.QueryMarketPriceResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *pricestypes.QueryMarketPriceRequest, ...grpc.CallOption) *pricestypes.QueryMarketPriceResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pricestypes.QueryMarketPriceResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *pricestypes.QueryMarketPriceRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MevNodeToNodeCalculation provides a mock function with given fields: ctx, in, opts
func (_m *QueryClient) MevNodeToNodeCalculation(ctx context.Context, in *clobtypes.MevNodeToNodeCalculationRequest, opts ...grpc.CallOption) (*clobtypes.MevNodeToNodeCalculationResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *clobtypes.MevNodeToNodeCalculationResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *clobtypes.MevNodeToNodeCalculationRequest, ...grpc.CallOption) (*clobtypes.MevNodeToNodeCalculationResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *clobtypes.MevNodeToNodeCalculationRequest, ...grpc.CallOption) *clobtypes.MevNodeToNodeCalculationResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*clobtypes.MevNodeToNodeCalculationResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *clobtypes.MevNodeToNodeCalculationRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Params provides a mock function with given fields: ctx, in, opts
func (_m *QueryClient) Params(ctx context.Context, in *perpetualstypes.QueryParamsRequest, opts ...grpc.CallOption) (*perpetualstypes.QueryParamsResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *perpetualstypes.QueryParamsResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *perpetualstypes.QueryParamsRequest, ...grpc.CallOption) (*perpetualstypes.QueryParamsResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *perpetualstypes.QueryParamsRequest, ...grpc.CallOption) *perpetualstypes.QueryParamsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*perpetualstypes.QueryParamsResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *perpetualstypes.QueryParamsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Perpetual provides a mock function with given fields: ctx, in, opts
func (_m *QueryClient) Perpetual(ctx context.Context, in *perpetualstypes.QueryPerpetualRequest, opts ...grpc.CallOption) (*perpetualstypes.QueryPerpetualResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *perpetualstypes.QueryPerpetualResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *perpetualstypes.QueryPerpetualRequest, ...grpc.CallOption) (*perpetualstypes.QueryPerpetualResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *perpetualstypes.QueryPerpetualRequest, ...grpc.CallOption) *perpetualstypes.QueryPerpetualResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*perpetualstypes.QueryPerpetualResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *perpetualstypes.QueryPerpetualRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PremiumSamples provides a mock function with given fields: ctx, in, opts
func (_m *QueryClient) PremiumSamples(ctx context.Context, in *perpetualstypes.QueryPremiumSamplesRequest, opts ...grpc.CallOption) (*perpetualstypes.QueryPremiumSamplesResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *perpetualstypes.QueryPremiumSamplesResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *perpetualstypes.QueryPremiumSamplesRequest, ...grpc.CallOption) (*perpetualstypes.QueryPremiumSamplesResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *perpetualstypes.QueryPremiumSamplesRequest, ...grpc.CallOption) *perpetualstypes.QueryPremiumSamplesResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*perpetualstypes.QueryPremiumSamplesResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *perpetualstypes.QueryPremiumSamplesRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PremiumVotes provides a mock function with given fields: ctx, in, opts
func (_m *QueryClient) PremiumVotes(ctx context.Context, in *perpetualstypes.QueryPremiumVotesRequest, opts ...grpc.CallOption) (*perpetualstypes.QueryPremiumVotesResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *perpetualstypes.QueryPremiumVotesResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *perpetualstypes.QueryPremiumVotesRequest, ...grpc.CallOption) (*perpetualstypes.QueryPremiumVotesResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *perpetualstypes.QueryPremiumVotesRequest, ...grpc.CallOption) *perpetualstypes.QueryPremiumVotesResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*perpetualstypes.QueryPremiumVotesResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *perpetualstypes.QueryPremiumVotesRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PreviousBlockInfo provides a mock function with given fields: ctx, in, opts
func (_m *QueryClient) PreviousBlockInfo(ctx context.Context, in *types.QueryPreviousBlockInfoRequest, opts ...grpc.CallOption) (*types.QueryPreviousBlockInfoResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *types.QueryPreviousBlockInfoResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryPreviousBlockInfoRequest, ...grpc.CallOption) (*types.QueryPreviousBlockInfoResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryPreviousBlockInfoRequest, ...grpc.CallOption) *types.QueryPreviousBlockInfoResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.QueryPreviousBlockInfoResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.QueryPreviousBlockInfoRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StreamOrderbookUpdates provides a mock function with given fields: ctx, in, opts
func (_m *QueryClient) StreamOrderbookUpdates(ctx context.Context, in *clobtypes.StreamOrderbookUpdatesRequest, opts ...grpc.CallOption) (clobtypes.Query_StreamOrderbookUpdatesClient, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 clobtypes.Query_StreamOrderbookUpdatesClient
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *clobtypes.StreamOrderbookUpdatesRequest, ...grpc.CallOption) (clobtypes.Query_StreamOrderbookUpdatesClient, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *clobtypes.StreamOrderbookUpdatesRequest, ...grpc.CallOption) clobtypes.Query_StreamOrderbookUpdatesClient); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(clobtypes.Query_StreamOrderbookUpdatesClient)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *clobtypes.StreamOrderbookUpdatesRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Subaccount provides a mock function with given fields: ctx, in, opts
func (_m *QueryClient) Subaccount(ctx context.Context, in *subaccountstypes.QueryGetSubaccountRequest, opts ...grpc.CallOption) (*subaccountstypes.QuerySubaccountResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *subaccountstypes.QuerySubaccountResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *subaccountstypes.QueryGetSubaccountRequest, ...grpc.CallOption) (*subaccountstypes.QuerySubaccountResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *subaccountstypes.QueryGetSubaccountRequest, ...grpc.CallOption) *subaccountstypes.QuerySubaccountResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*subaccountstypes.QuerySubaccountResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *subaccountstypes.QueryGetSubaccountRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SubaccountAll provides a mock function with given fields: ctx, in, opts
func (_m *QueryClient) SubaccountAll(ctx context.Context, in *subaccountstypes.QueryAllSubaccountRequest, opts ...grpc.CallOption) (*subaccountstypes.QuerySubaccountAllResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *subaccountstypes.QuerySubaccountAllResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *subaccountstypes.QueryAllSubaccountRequest, ...grpc.CallOption) (*subaccountstypes.QuerySubaccountAllResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *subaccountstypes.QueryAllSubaccountRequest, ...grpc.CallOption) *subaccountstypes.QuerySubaccountAllResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*subaccountstypes.QuerySubaccountAllResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *subaccountstypes.QueryAllSubaccountRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateMarketPrices provides a mock function with given fields: ctx, in, opts
func (_m *QueryClient) UpdateMarketPrices(ctx context.Context, in *pricefeedapi.UpdateMarketPricesRequest, opts ...grpc.CallOption) (*pricefeedapi.UpdateMarketPricesResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *pricefeedapi.UpdateMarketPricesResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *pricefeedapi.UpdateMarketPricesRequest, ...grpc.CallOption) (*pricefeedapi.UpdateMarketPricesResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *pricefeedapi.UpdateMarketPricesRequest, ...grpc.CallOption) *pricefeedapi.UpdateMarketPricesResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pricefeedapi.UpdateMarketPricesResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *pricefeedapi.UpdateMarketPricesRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewQueryClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewQueryClient creates a new instance of QueryClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewQueryClient(t mockConstructorTestingTNewQueryClient) *QueryClient {
	mock := &QueryClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
