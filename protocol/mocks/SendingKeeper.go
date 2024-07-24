// Code generated by mockery v2.42.0. DO NOT EDIT.

package mocks

import (
	cosmos_sdktypes "github.com/cosmos/cosmos-sdk/types"
	mock "github.com/stretchr/testify/mock"

	types "github.com/StreamFinance-Protocol/stream-chain/protocol/x/sending/types"
)

// SendingKeeper is an autogenerated mock type for the SendingKeeper type
type SendingKeeper struct {
	mock.Mock
}

// HasAuthority provides a mock function with given fields: authority
func (_m *SendingKeeper) HasAuthority(authority string) bool {
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

// ProcessDepositToSubaccount provides a mock function with given fields: ctx, msgDepositToSubaccount
func (_m *SendingKeeper) ProcessDepositToSubaccount(ctx cosmos_sdktypes.Context, msgDepositToSubaccount *types.MsgDepositToSubaccount) error {
	ret := _m.Called(ctx, msgDepositToSubaccount)

	if len(ret) == 0 {
		panic("no return value specified for ProcessDepositToSubaccount")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(cosmos_sdktypes.Context, *types.MsgDepositToSubaccount) error); ok {
		r0 = rf(ctx, msgDepositToSubaccount)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ProcessTransfer provides a mock function with given fields: ctx, transfer
func (_m *SendingKeeper) ProcessTransfer(ctx cosmos_sdktypes.Context, transfer *types.Transfer) error {
	ret := _m.Called(ctx, transfer)

	if len(ret) == 0 {
		panic("no return value specified for ProcessTransfer")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(cosmos_sdktypes.Context, *types.Transfer) error); ok {
		r0 = rf(ctx, transfer)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ProcessWithdrawFromSubaccount provides a mock function with given fields: ctx, msgWithdrawFromSubaccount
func (_m *SendingKeeper) ProcessWithdrawFromSubaccount(ctx cosmos_sdktypes.Context, msgWithdrawFromSubaccount *types.MsgWithdrawFromSubaccount) error {
	ret := _m.Called(ctx, msgWithdrawFromSubaccount)

	if len(ret) == 0 {
		panic("no return value specified for ProcessWithdrawFromSubaccount")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(cosmos_sdktypes.Context, *types.MsgWithdrawFromSubaccount) error); ok {
		r0 = rf(ctx, msgWithdrawFromSubaccount)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendFromModuleToAccount provides a mock function with given fields: ctx, msg
func (_m *SendingKeeper) SendFromModuleToAccount(ctx cosmos_sdktypes.Context, msg *types.MsgSendFromModuleToAccount) error {
	ret := _m.Called(ctx, msg)

	if len(ret) == 0 {
		panic("no return value specified for SendFromModuleToAccount")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(cosmos_sdktypes.Context, *types.MsgSendFromModuleToAccount) error); ok {
		r0 = rf(ctx, msg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewSendingKeeper creates a new instance of SendingKeeper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSendingKeeper(t interface {
	mock.TestingT
	Cleanup(func())
}) *SendingKeeper {
	mock := &SendingKeeper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
