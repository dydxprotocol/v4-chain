package keeper_test

import (
	"fmt"
	"reflect"
	"testing"

	sdkmath "cosmossdk.io/math"

	"cosmossdk.io/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cometbfttypes "github.com/cometbft/cometbft/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/encoding"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	bridgetypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
	feetierstypes "github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	BridgeAuthority      = bridgetypes.ModuleAddress.String()
	BridgeAccountAddress = sdk.MustAccAddressFromBech32(BridgeAuthority)

	DelayMsgAuthority = types.ModuleAddress

	testDenom = "adv4tnt"

	BridgeGenesisAccountBalance = sdk.NewCoin(testDenom, sdkmath.NewInt(1000000000))

	delta                        = constants.BridgeEvent_Id0_Height0.Coin.Amount.Int64()
	BridgeExpectedAccountBalance = sdk.NewCoin(testDenom,
		BridgeGenesisAccountBalance.Amount.Sub(
			constants.BridgeEvent_Id0_Height0.Coin.Amount,
		),
	)
)

func TestDispatchMessagesForBlock(t *testing.T) {
	ctx, k, _, bridgeKeeper, _ := keepertest.DelayMsgKeeperWithMockBridgeKeeper(t)

	// Add messages to the keeper.
	for i, msg := range constants.AllMsgs {
		id, err := k.DelayMessageByBlocks(ctx, msg, 0)
		require.NoError(t, err)
		require.Equal(t, uint32(i), id)
	}

	// Sanity check: messages appear for block 0.
	blockMessageIds, found := k.GetBlockMessageIds(ctx, 0)
	require.True(t, found)
	require.Equal(t, []uint32{0, 1, 2}, blockMessageIds.Ids)

	// Mock the bridge keeper methods called by the bridge msg server.
	bridgeKeeper.On("CompleteBridge", mock.AnythingOfType("types.Context"), mock.Anything).
		Return(nil).Times(len(constants.AllMsgs))
	bridgeKeeper.On("HasAuthority", DelayMsgAuthority.String()).Return(true).Times(len(constants.AllMsgs))

	// Dispatch messages for block 0.

	keeper.DispatchMessagesForBlock(k, ctx)

	_, found = k.GetBlockMessageIds(ctx, 0)
	require.False(t, found)

	require.True(t, bridgeKeeper.AssertExpectations(t))
}

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

// generateBridgeEventMsgAny wraps bridge event in a MsgCompleteBridge and encodes it into an Any.
func generateBridgeEventMsgAny(t *testing.T, event bridgetypes.BridgeEvent) *codectypes.Any {
	msgCompleteBridge := bridgetypes.MsgCompleteBridge{
		Authority: DelayMsgAuthority.String(),
		Event:     event,
	}
	any, err := codectypes.NewAnyWithValue(&msgCompleteBridge)
	require.NoError(t, err)
	return any
}

// expectAccountBalance checks that the specified account has the expected balance.
func expectAccountBalance(
	t *testing.T,
	ctx sdk.Context,
	tApp *testapp.TestApp,
	address sdk.AccAddress,
	expectedBalance sdk.Coin,
) {
	balance := tApp.App.BankKeeper.GetBalance(ctx, address, expectedBalance.Denom)
	require.Equal(t, expectedBalance.Amount, balance.Amount)
	require.Equal(t, expectedBalance.Denom, balance.Denom)
}

func TestSendDelayedCompleteBridgeMessage(t *testing.T) {
	// Create an encoded bridge event set to occur at block 2.
	// Expect that Alice's account will increase by 888 coins at block 2.
	// Bridge module account will also decrease by 888 coins at block 2.
	delayedMessage := types.DelayedMessage{
		Id:          0,
		Msg:         generateBridgeEventMsgAny(t, constants.BridgeEvent_Id0_Height0),
		BlockHeight: 2,
	}

	tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis cometbfttypes.GenesisDoc) {
		genesis = testapp.DefaultGenesis()
		// Add the delayed message to the genesis state.
		testapp.UpdateGenesisDocWithAppStateForModule(
			&genesis,
			func(genesisState *types.GenesisState) {
				genesisState.DelayedMessages = []*types.DelayedMessage{&delayedMessage}
				genesisState.NextDelayedMessageId = 1
			},
		)
		return genesis
	}).Build()
	ctx := tApp.InitChain()

	// Sanity check: the delayed message is in the keeper scheduled for block 2.
	blockMessageIds, found := tApp.App.DelayMsgKeeper.GetBlockMessageIds(ctx, 2)
	require.True(t, found)
	require.Equal(t, []uint32{0}, blockMessageIds.Ids)

	aliceAccountAddress := sdk.MustAccAddressFromBech32(constants.BridgeEvent_Id0_Height0.Address)

	// Sanity check: at block 1, expect bridge balance is genesis value before the message is sent.
	expectAccountBalance(t, ctx, tApp, BridgeAccountAddress, BridgeGenesisAccountBalance)

	// Get initial Alice balance
	aliceInitialBalance := tApp.App.BankKeeper.GetBalance(ctx, aliceAccountAddress, testDenom)
	// Calculate Alice's expected balance after complete bridge event.
	aliceExpectedAccountBalance := sdk.NewCoin(
		testDenom,
		aliceInitialBalance.Amount.Add(sdkmath.NewInt(delta)),
	)

	// Advance to block 2 and invoke delayed message to complete bridge.
	ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

	// Assert: balances have been updated to reflect the executed CompleteBridge message.
	expectAccountBalance(t, ctx, tApp, BridgeAccountAddress, BridgeExpectedAccountBalance)
	expectAccountBalance(t, ctx, tApp, aliceAccountAddress, aliceExpectedAccountBalance)

	// Assert: the message has been deleted from the keeper.
	_, found = tApp.App.DelayMsgKeeper.GetMessage(ctx, 0)
	require.False(t, found)

	// The block message ids have also been deleted.
	_, found = tApp.App.DelayMsgKeeper.GetBlockMessageIds(ctx, 2)
	require.False(t, found)
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

	// Advance to block 2 and invoke delayed message to complete bridge.
	ctx = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{})

	resp, err = tApp.App.FeeTiersKeeper.PerpetualFeeParams(ctx, &feetierstypes.QueryPerpetualFeeParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, feetierstypes.StandardParams(), resp.Params)
}

func TestSendDelayedCompleteBridgeMessage_Failure(t *testing.T) {
	// Create an encoded bridge event set to occur at block 2.
	// The bridge event is invalid and will not execute.
	// Expect no account balance changes, and the message to be deleted from the keeper.
	invalidBridgeEvent := bridgetypes.BridgeEvent{
		Id:             constants.BridgeEvent_Id0_Height0.Id,
		Address:        "INVALID",
		Coin:           constants.BridgeEvent_Id0_Height0.Coin,
		EthBlockHeight: constants.BridgeEvent_Id0_Height0.EthBlockHeight,
	}
	_, err := sdk.AccAddressFromBech32("INVALID")
	require.Error(t, err)

	delayedMessage := types.DelayedMessage{
		Id:          0,
		Msg:         generateBridgeEventMsgAny(t, invalidBridgeEvent),
		BlockHeight: 2,
	}

	tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis cometbfttypes.GenesisDoc) {
		genesis = testapp.DefaultGenesis()
		// Add the delayed message to the genesis state.
		testapp.UpdateGenesisDocWithAppStateForModule(
			&genesis,
			func(genesisState *types.GenesisState) {
				genesisState.DelayedMessages = []*types.DelayedMessage{&delayedMessage}
				genesisState.NextDelayedMessageId = 1
			},
		)
		return genesis
	}).Build()
	ctx := tApp.InitChain()

	// Sanity check: at block 1, balances are as expected before the message is sent.
	expectAccountBalance(t, ctx, tApp, BridgeAccountAddress, BridgeGenesisAccountBalance)

	// Sanity check: a message with this id exists within the keeper.
	_, found := tApp.App.DelayMsgKeeper.GetMessage(ctx, 0)
	require.True(t, found)

	// Sanity check: this message id is scheduled to be executed at block 2.
	messageIds, found := tApp.App.DelayMsgKeeper.GetBlockMessageIds(ctx, 2)
	require.True(t, found)
	require.Equal(t, []uint32{0}, messageIds.Ids)

	// Advance to block 2 and invoke delayed message to complete bridge.
	ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

	// Assert: balances have been updated to reflect the executed CompleteBridge message.
	expectAccountBalance(t, ctx, tApp, BridgeAccountAddress, BridgeGenesisAccountBalance)

	// Assert: the message has been deleted from the keeper.
	_, found = tApp.App.DelayMsgKeeper.GetMessage(ctx, 0)
	require.False(t, found)

	// The block message ids have also been deleted.
	_, found = tApp.App.DelayMsgKeeper.GetBlockMessageIds(ctx, 2)
	require.False(t, found)
}

// This test case verifies that events emitted from message executions are correctly
// propagated to the base context.
func TestDispatchMessagesForBlock_EventsArePropagated(t *testing.T) {
	ctx, k, _, _, bankKeeper, _ := keepertest.DelayMsgKeepers(t)
	// Mint coins to the bridge module account so that it has enough balance for completing bridges.
	err := bankKeeper.MintCoins(ctx, bridgetypes.ModuleName, sdk.NewCoins(BridgeGenesisAccountBalance))
	require.NoError(t, err)

	// Delay a complete bridge message, which calls bank transfer that emits a transfer event.
	bridgeEvent := bridgetypes.BridgeEvent{
		Id:             1,
		Coin:           sdk.NewCoin("adv4tnt", sdkmath.NewInt(1_000)),
		Address:        constants.AliceAccAddress.String(),
		EthBlockHeight: 0,
	}
	_, err = k.DelayMessageByBlocks(
		ctx,
		&bridgetypes.MsgCompleteBridge{
			Authority: DelayMsgAuthority.String(),
			Event:     bridgeEvent,
		},
		0,
	)
	require.NoError(t, err)

	// Sanity check: messages appear for block 0.
	blockMessageIds, found := k.GetBlockMessageIds(ctx, 0)
	require.True(t, found)
	require.Equal(t, []uint32{0}, blockMessageIds.Ids)

	// Dispatch messages for block 0.
	keeper.DispatchMessagesForBlock(k, ctx)

	_, found = k.GetBlockMessageIds(ctx, 0)
	require.False(t, found)

	emittedEvents := ctx.EventManager().Events()
	expectedTransferEvent := sdk.NewEvent(
		"transfer",
		sdk.NewAttribute("recipient", bridgeEvent.Address),
		sdk.NewAttribute("sender", BridgeAccountAddress.String()),
		sdk.NewAttribute("amount", bridgeEvent.Coin.String()),
	)

	// Verify that emitted events contains the expected transfer event exactly once.
	foundExpectedTransferEvent := 0
	for _, emittedEvent := range emittedEvents {
		if reflect.DeepEqual(expectedTransferEvent, emittedEvent) {
			foundExpectedTransferEvent++
		}
	}
	require.Equal(t, 1, foundExpectedTransferEvent)
}
