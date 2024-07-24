// Code generated by mockery v2.42.0. DO NOT EDIT.

package mocks

import (
	clobtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	mock "github.com/stretchr/testify/mock"

	process "github.com/StreamFinance-Protocol/stream-chain/protocol/app/process"

	types "github.com/cosmos/cosmos-sdk/types"
)

// ProcessClobKeeper is an autogenerated mock type for the ProcessClobKeeper type
type ProcessClobKeeper struct {
	mock.Mock
}

// RecordMevMetrics provides a mock function with given fields: ctx, stakingKeeper, perpetualKeeper, msgProposedOperations
func (_m *ProcessClobKeeper) RecordMevMetrics(ctx types.Context, stakingKeeper process.ProcessStakingKeeper, perpetualKeeper process.ProcessPerpetualKeeper, msgProposedOperations *clobtypes.MsgProposedOperations) {
	_m.Called(ctx, stakingKeeper, perpetualKeeper, msgProposedOperations)
}

// RecordMevMetricsIsEnabled provides a mock function with given fields:
func (_m *ProcessClobKeeper) RecordMevMetricsIsEnabled() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for RecordMevMetricsIsEnabled")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// NewProcessClobKeeper creates a new instance of ProcessClobKeeper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewProcessClobKeeper(t interface {
	mock.TestingT
	Cleanup(func())
}) *ProcessClobKeeper {
	mock := &ProcessClobKeeper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
