// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	indexer_manager "github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	mock "github.com/stretchr/testify/mock"

	msgsender "github.com/dydxprotocol/v4-chain/protocol/indexer/msgsender"

	types "github.com/cosmos/cosmos-sdk/types"
)

// IndexerEventManager is an autogenerated mock type for the IndexerEventManager type
type IndexerEventManager struct {
	mock.Mock
}

// AddBlockEvent provides a mock function with given fields: ctx, subType, blockEvent, version, dataBytes
func (_m *IndexerEventManager) AddBlockEvent(ctx types.Context, subType string, blockEvent indexer_manager.IndexerTendermintEvent_BlockEvent, version uint32, dataBytes []byte) {
	_m.Called(ctx, subType, blockEvent, version, dataBytes)
}

// AddTxnEvent provides a mock function with given fields: ctx, subType, version, dataByes
func (_m *IndexerEventManager) AddTxnEvent(ctx types.Context, subType string, version uint32, dataByes []byte) {
	_m.Called(ctx, subType, version, dataByes)
}

// ClearEvents provides a mock function with given fields: ctx
func (_m *IndexerEventManager) ClearEvents(ctx types.Context) {
	_m.Called(ctx)
}

// Enabled provides a mock function with given fields:
func (_m *IndexerEventManager) Enabled() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Enabled")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// ProduceBlock provides a mock function with given fields: ctx
func (_m *IndexerEventManager) ProduceBlock(ctx types.Context) *indexer_manager.IndexerTendermintBlock {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for ProduceBlock")
	}

	var r0 *indexer_manager.IndexerTendermintBlock
	if rf, ok := ret.Get(0).(func(types.Context) *indexer_manager.IndexerTendermintBlock); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*indexer_manager.IndexerTendermintBlock)
		}
	}

	return r0
}

// SendOffchainData provides a mock function with given fields: message
func (_m *IndexerEventManager) SendOffchainData(message msgsender.Message) {
	_m.Called(message)
}

// SendOnchainData provides a mock function with given fields: block
func (_m *IndexerEventManager) SendOnchainData(block *indexer_manager.IndexerTendermintBlock) {
	_m.Called(block)
}

// NewIndexerEventManager creates a new instance of IndexerEventManager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIndexerEventManager(t interface {
	mock.TestingT
	Cleanup(func())
}) *IndexerEventManager {
	mock := &IndexerEventManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
