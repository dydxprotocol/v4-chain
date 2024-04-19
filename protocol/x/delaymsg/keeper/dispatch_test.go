package keeper_test

import (
	"fmt"
	"testing"

	"cosmossdk.io/log"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	testapp "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/app"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/encoding"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/delaymsg/keeper"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/delaymsg/types"
	feetierstypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/feetiers/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cometbfttypes "github.com/cometbft/cometbft/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	DelayMsgAuthority = types.ModuleAddress
)

func setupMockKeeperNoMessages(t *testing.T, ctx sdk.Context, k *mocks.DelayMsgKeeper) {
	k.On("GetBlockMessageIds", ctx, uint32(0)).Return(types.BlockMessageIds{}, false).Once()
}

func HandlerSuccess(_ sdk.Context, _ sdk.Msg) (*sdk.Result, error) {
	return &sdk.Result{}, nil
}

func HandlerFailure(_ sdk.Context, _ sdk.Msg) (*sdk.Result, error) {
	return &sdk.Result{}, fmt.Errorf("failed to handle message")
}

// mockSuccessRouter returns a handler that succeeds on all calls.
func mockSuccessRouter(_ sdk.Context) *mocks.MsgRouter {
	router := &mocks.MsgRouter{}
	router.On("Handler", mock.Anything).Return(HandlerSuccess).Times(3)
	return router
}

// mockFailingRouter returns a handler that fails on the first call.
func mockFailingRouter(ctx sdk.Context) *mocks.MsgRouter {
	router := mocks.MsgRouter{}
	router.On("Handler", mock.Anything).Return(HandlerFailure).Once()
	return &router
}

// mockPanickingRouter returns a handler that panics on the first call.
func mockPanickingRouter(ctx sdk.Context) *mocks.MsgRouter {
	router := mocks.MsgRouter{}
	router.On("Handler", mock.Anything).Panic("panic").Once()
	return &router
}

func setupMockKeeperMessageNotFound(t *testing.T, ctx sdk.Context, k *mocks.DelayMsgKeeper) {
	k.On("GetBlockMessageIds", ctx, uint32(0)).Return(types.BlockMessageIds{
		Ids: []uint32{0, 1, 2},
	}, true).Once()

	// Second message is not found.
	k.On("GetMessage", ctx, uint32(0)).Return(types.DelayedMessage{
		Id:          0,
		Msg:         encoding.EncodeMessageToAny(t, constants.TestMsg1),
		BlockHeight: 0,
	}, true).Once()
	k.On("GetMessage", ctx, uint32(1)).Return(types.DelayedMessage{}, false).Once()
	k.On("GetMessage", ctx, uint32(2)).Return(types.DelayedMessage{
		Id:          2,
		Msg:         encoding.EncodeMessageToAny(t, constants.TestMsg3),
		BlockHeight: 0,
	}, true).Once()

	// 2 messages are routed.
	msgRouter := mockSuccessRouter(ctx)
	k.On("Router").Return(msgRouter).Times(2)

	// 2 message executions are persisted.
	cms := ctx.MultiStore().CacheMultiStore().(*mocks.CacheMultiStore)
	cms.On("Write").Return(nil).Times(2)

	// All deletes are called.
	k.On("DeleteMessage", ctx, uint32(0)).Return(nil).Once()
	k.On("DeleteMessage", ctx, uint32(1)).Return(nil).Once()
	k.On("DeleteMessage", ctx, uint32(2)).Return(nil).Once()
}

func setupMockKeeperExecutionFailure(t *testing.T, ctx sdk.Context, k *mocks.DelayMsgKeeper) {
	k.On("GetBlockMessageIds", ctx, uint32(0)).Return(types.BlockMessageIds{
		Ids: []uint32{0, 1, 2},
	}, true).Once()

	// All messages found.
	k.On("GetMessage", ctx, uint32(0)).Return(types.DelayedMessage{
		Id:          0,
		Msg:         encoding.EncodeMessageToAny(t, constants.TestMsg1),
		BlockHeight: 0,
	}, true).Once()
	k.On("GetMessage", ctx, uint32(1)).Return(types.DelayedMessage{
		Id:          1,
		Msg:         encoding.EncodeMessageToAny(t, constants.TestMsg2),
		BlockHeight: 0,
	}, true).Once()
	k.On("GetMessage", ctx, uint32(2)).Return(types.DelayedMessage{
		Id:          2,
		Msg:         encoding.EncodeMessageToAny(t, constants.TestMsg3),
		BlockHeight: 0,
	}, true).Once()

	// 1st message fails to execute. Following 2 succeed.
	successRouter := mockSuccessRouter(ctx)
	failureRouter := mockFailingRouter(ctx)
	k.On("Router").Return(failureRouter).Times(1)
	k.On("Router").Return(successRouter).Times(2)

	// 2 message executions are persisted.
	cms := ctx.MultiStore().CacheMultiStore().(*mocks.CacheMultiStore)
	cms.On("Write").Return(nil).Times(2)

	// All deletes are called.
	k.On("DeleteMessage", ctx, uint32(0)).Return(nil).Once()
	k.On("DeleteMessage", ctx, uint32(1)).Return(nil).Once()
	k.On("DeleteMessage", ctx, uint32(2)).Return(nil).Once()
}

func setupMockKeeperMessageHandlerPanic(t *testing.T, ctx sdk.Context, k *mocks.DelayMsgKeeper) {
	k.On("GetBlockMessageIds", ctx, uint32(0)).Return(types.BlockMessageIds{
		Ids: []uint32{0, 1, 2},
	}, true).Once()

	// All messages found.
	k.On("GetMessage", ctx, uint32(0)).Return(types.DelayedMessage{
		Id:          0,
		Msg:         encoding.EncodeMessageToAny(t, constants.TestMsg1),
		BlockHeight: 0,
	}, true).Once()
	k.On("GetMessage", ctx, uint32(1)).Return(types.DelayedMessage{
		Id:          1,
		Msg:         encoding.EncodeMessageToAny(t, constants.TestMsg2),
		BlockHeight: 0,
	}, true).Once()
	k.On("GetMessage", ctx, uint32(2)).Return(types.DelayedMessage{
		Id:          2,
		Msg:         encoding.EncodeMessageToAny(t, constants.TestMsg3),
		BlockHeight: 0,
	}, true).Once()

	// 1st message fails to execute. Following 2 succeed.
	successRouter := mockSuccessRouter(ctx)
	panicRouter := mockPanickingRouter(ctx)
	k.On("Router").Return(panicRouter).Times(1)
	k.On("Router").Return(successRouter).Times(2)

	// 2 message executions are persisted.
	cms := ctx.MultiStore().CacheMultiStore().(*mocks.CacheMultiStore)
	cms.On("Write").Return(nil).Times(2)

	// All deletes are called.
	k.On("DeleteMessage", ctx, uint32(0)).Return(nil).Once()
	k.On("DeleteMessage", ctx, uint32(1)).Return(nil).Once()
	k.On("DeleteMessage", ctx, uint32(2)).Return(nil).Once()
}

func setupMockKeeperDecodeFailure(t *testing.T, ctx sdk.Context, k *mocks.DelayMsgKeeper) {
	k.On("GetBlockMessageIds", ctx, uint32(0)).Return(types.BlockMessageIds{
		Ids: []uint32{0, 1, 2},
	}, true).Once()

	nonMsgAnyProto, err := codectypes.NewAnyWithValue(&types.BlockMessageIds{})
	require.NoError(t, err)

	// All messages found.
	k.On("GetMessage", ctx, uint32(0)).Return(types.DelayedMessage{
		Id:          0,
		Msg:         encoding.EncodeMessageToAny(t, constants.TestMsg1),
		BlockHeight: 0,
	}, true).Once()
	k.On("GetMessage", ctx, uint32(1)).Return(types.DelayedMessage{
		Id:          1,
		Msg:         nonMsgAnyProto,
		BlockHeight: 0,
	}, true).Once()
	k.On("GetMessage", ctx, uint32(2)).Return(types.DelayedMessage{
		Id:          2,
		Msg:         encoding.EncodeMessageToAny(t, constants.TestMsg3),
		BlockHeight: 0,
	}, true).Once()

	// 2 messages are routed.
	k.On("Router").Return(mockSuccessRouter(ctx)).Times(2)

	// 2 message executions are persisted.
	cms := ctx.MultiStore().CacheMultiStore().(*mocks.CacheMultiStore)
	cms.On("Write").Return(nil).Times(2)

	// All deletes are called. 2nd delete fails.
	k.On("DeleteMessage", ctx, uint32(0)).Return(nil).Once()
	k.On("DeleteMessage", ctx, uint32(1)).Return(nil).Once()
	k.On("DeleteMessage", ctx, uint32(2)).Return(nil).Once()
}

func setupMockKeeperDeletionFailure(t *testing.T, ctx sdk.Context, k *mocks.DelayMsgKeeper) {
	k.On("GetBlockMessageIds", ctx, uint32(0)).Return(types.BlockMessageIds{
		Ids: []uint32{0, 1, 2},
	}, true).Once()

	// All messages found.
	k.On("GetMessage", ctx, uint32(0)).Return(types.DelayedMessage{
		Id:          0,
		Msg:         encoding.EncodeMessageToAny(t, constants.TestMsg1),
		BlockHeight: 0,
	}, true).Once()
	k.On("GetMessage", ctx, uint32(1)).Return(types.DelayedMessage{
		Id:          1,
		Msg:         encoding.EncodeMessageToAny(t, constants.TestMsg2),
		BlockHeight: 0,
	}, true).Once()
	k.On("GetMessage", ctx, uint32(2)).Return(types.DelayedMessage{
		Id:          2,
		Msg:         encoding.EncodeMessageToAny(t, constants.TestMsg3),
		BlockHeight: 0,
	}, true).Once()

	// All messages are routed.
	k.On("Router").Return(mockSuccessRouter(ctx)).Times(3)

	// All message executions are persisted.
	cms := ctx.MultiStore().CacheMultiStore().(*mocks.CacheMultiStore)
	cms.On("Write").Return(nil).Times(3)

	// All deletes are called. 2nd delete fails.
	k.On("DeleteMessage", ctx, uint32(0)).Return(nil).Once()
	k.On("DeleteMessage", ctx, uint32(1)).Return(fmt.Errorf("Deletion failure")).Once()
	k.On("DeleteMessage", ctx, uint32(2)).Return(nil).Once()
}

func TestDispatchMessagesForBlock_Mixed(t *testing.T) {
	tests := map[string]struct {
		setupMocks func(t *testing.T, ctx sdk.Context, k *mocks.DelayMsgKeeper)
	}{
		"No messages - dispatch terminates with no action": {
			setupMocks: setupMockKeeperNoMessages,
		},
		"Unexpected message not found does not affect remaining messages": {
			setupMocks: setupMockKeeperMessageNotFound,
		},
		"Execution error does not affect remaining messages": {
			setupMocks: setupMockKeeperExecutionFailure,
		},
		"Execution panic does not affect remaining messages": {
			setupMocks: setupMockKeeperMessageHandlerPanic,
		},
		"Decode failure does not affect remaining messages": {
			setupMocks: setupMockKeeperDecodeFailure,
		},
		"Deletion failure does not affect deletion of remaining messages": {
			setupMocks: setupMockKeeperDeletionFailure,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			k := &mocks.DelayMsgKeeper{}
			ms := &mocks.MultiStore{}
			cms := &mocks.CacheMultiStore{}
			// Expect that the cached store is accessed 0 or more times.
			ms.On("CacheMultiStore").Return(cms).Maybe().Times(0)
			ctx := sdk.NewContext(ms, tmproto.Header{}, false, log.NewNopLogger())

			tc.setupMocks(t, ctx, k)

			keeper.DispatchMessagesForBlock(k, ctx)

			mock.AssertExpectationsForObjects(t, k, ms, cms)
		})
	}
}

// TestSendDelayedPerpetualFeeParamsUpdate tests that the delayed message testApp genesis state, which contains a
// message to update the x/feetiers perpetual fee params after ~120 days of blocks, is executed correctly. In this
// test, we modify the genesis state to apply the parameter update on block 2 to validate that the update is applied
// correctly.
func TestSendDelayedPerpetualFeeParamsUpdate(t *testing.T) {
	// TODO(CORE-858): Re-enable determinism checks once non-determinism issue is found and resolved.
	tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis cometbfttypes.GenesisDoc) {
		genesis = testapp.DefaultGenesis()
		// Update the genesis state to execute the perpetual fee params update at block 2.
		testapp.UpdateGenesisDocWithAppStateForModule(
			&genesis,
			func(genesisState *types.GenesisState) {
				// Update the default state to apply the first delayed message on block 2.
				// This is the PerpetualFeeParamsUpdate message.
				genesisState.DelayedMessages[0].BlockHeight = 2
			},
		)
		return genesis
	}).WithNonDeterminismChecksEnabled(false).Build()
	ctx := tApp.InitChain()

	resp, err := tApp.App.FeeTiersKeeper.PerpetualFeeParams(ctx, &feetierstypes.QueryPerpetualFeeParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, feetierstypes.PromotionalParams(), resp.Params)

	// Advance to block 2 and invoke delayed message to complete bridge. In this context, bridge seems to be
	// referring to an implicit bridging operation not reliant on the bridge module.
	ctx = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{})

	resp, err = tApp.App.FeeTiersKeeper.PerpetualFeeParams(ctx, &feetierstypes.QueryPerpetualFeeParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, feetierstypes.StandardParams(), resp.Params)
}
