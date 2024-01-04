// Code generated by mockery v2.23.1. DO NOT EDIT.

package mocks

import (
	bridgetypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	mock "github.com/stretchr/testify/mock"

	types "github.com/cosmos/cosmos-sdk/types"
)

// ProcessBridgeKeeper is an autogenerated mock type for the ProcessBridgeKeeper type
type ProcessBridgeKeeper struct {
	mock.Mock
}

// GetAcknowledgedEventInfo provides a mock function with given fields: ctx
func (_m *ProcessBridgeKeeper) GetAcknowledgedEventInfo(ctx types.Context) bridgetypes.BridgeEventInfo {
	ret := _m.Called(ctx)

	var r0 bridgetypes.BridgeEventInfo
	if rf, ok := ret.Get(0).(func(types.Context) bridgetypes.BridgeEventInfo); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(bridgetypes.BridgeEventInfo)
	}

	return r0
}

// GetBridgeEventFromServer provides a mock function with given fields: ctx, id
func (_m *ProcessBridgeKeeper) GetBridgeEventFromServer(ctx types.Context, id uint32) (bridgetypes.BridgeEvent, bool) {
	ret := _m.Called(ctx, id)

	var r0 bridgetypes.BridgeEvent
	var r1 bool
	if rf, ok := ret.Get(0).(func(types.Context, uint32) (bridgetypes.BridgeEvent, bool)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(types.Context, uint32) bridgetypes.BridgeEvent); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(bridgetypes.BridgeEvent)
	}

	if rf, ok := ret.Get(1).(func(types.Context, uint32) bool); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// GetRecognizedEventInfo provides a mock function with given fields: ctx
func (_m *ProcessBridgeKeeper) GetRecognizedEventInfo(ctx types.Context) bridgetypes.BridgeEventInfo {
	ret := _m.Called(ctx)

	var r0 bridgetypes.BridgeEventInfo
	if rf, ok := ret.Get(0).(func(types.Context) bridgetypes.BridgeEventInfo); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(bridgetypes.BridgeEventInfo)
	}

	return r0
}

// GetSafetyParams provides a mock function with given fields: ctx
func (_m *ProcessBridgeKeeper) GetSafetyParams(ctx types.Context) bridgetypes.SafetyParams {
	ret := _m.Called(ctx)

	var r0 bridgetypes.SafetyParams
	if rf, ok := ret.Get(0).(func(types.Context) bridgetypes.SafetyParams); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(bridgetypes.SafetyParams)
	}

	return r0
}

type mockConstructorTestingTNewProcessBridgeKeeper interface {
	mock.TestingT
	Cleanup(func())
}

// NewProcessBridgeKeeper creates a new instance of ProcessBridgeKeeper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewProcessBridgeKeeper(t mockConstructorTestingTNewProcessBridgeKeeper) *ProcessBridgeKeeper {
	mock := &ProcessBridgeKeeper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
