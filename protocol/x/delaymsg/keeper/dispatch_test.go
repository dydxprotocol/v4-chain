package keeper_test

import (
	"fmt"
	"github.com/cometbft/cometbft/libs/log"
	cometbfttypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	sdktest "github.com/dydxprotocol/v4-chain/protocol/testutil/sdk"
	bridgetypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	BridgeAuthority      = authtypes.NewModuleAddress(bridgetypes.ModuleName).String()
	BridgeAccountAddress = sdk.MustAccAddressFromBech32(BridgeAuthority)

	DelayMsgAuthority = authtypes.NewModuleAddress(types.ModuleName).String()

	BridgeGenesisAccountBalance = sdk.NewCoin("dv4tnt", sdk.NewInt(1000000000))
	AliceInitialAccountBalance  = sdk.NewCoin("dv4tnt", sdk.NewInt(99500000000))

	delta                        = constants.BridgeEvent_Id0_Height0.Coin.Amount.Int64()
	BridgeExpectedAccountBalance = sdk.NewCoin("dv4tnt", sdk.NewInt(1000000000-delta))
	AliceExpectedAccountBalance  = sdk.NewCoin("dv4tnt", sdk.NewInt(99500000000+delta))
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
	bridgeKeeper.On("CompleteBridge", ctx, mock.Anything).Return(nil).Times(len(constants.AllMsgs))
	bridgeKeeper.On("HasAuthority", DelayMsgAuthority).Return(true).Times(len(constants.AllMsgs))

	// Dispatch messages for block 0.

	keeper.DispatchMessagesForBlock(k, ctx)

	_, found = k.GetBlockMessageIds(ctx, 0)
	require.False(t, found)

	require.True(t, bridgeKeeper.AssertExpectations(t))
}

func setupMockKeeperNoMessages(ctx sdk.Context, k *mocks.DelayMsgKeeper) {
	k.On("GetBlockMessageIds", ctx, int64(0)).Return(types.BlockMessageIds{}, false).Once()
}

func HandlerSuccess(ctx sdk.Context, req sdk.Msg) (*sdk.Result, error) {
	return &sdk.Result{}, nil
}

func HandlerFailure(ctx sdk.Context, req sdk.Msg) (*sdk.Result, error) {
	return &sdk.Result{}, fmt.Errorf("failed to handle message")
}

// mockSuccessRouter returns a handler that succeeds on all calls.
func mockSuccessRouter(ctx sdk.Context) *mocks.MsgRouter {
	router := &mocks.MsgRouter{}
	router.On("Handler", mock.Anything).Return(HandlerSuccess).Times(3)
	return router
}

// mockFailingRouter returns a handler that fails on the first call, then returns a handler
// that succeeds on the next two.
func mockFailingRouter(ctx sdk.Context) *mocks.MsgRouter {
	router := mocks.MsgRouter{}
	router.On("Handler", mock.Anything).Return(HandlerFailure).Once()
	router.On("Handler", mock.Anything).Return(HandlerSuccess).Times(2)
	return &router
}

func setupMockKeeperMessageNotFound(ctx sdk.Context, k *mocks.DelayMsgKeeper) {
	k.On("GetBlockMessageIds", ctx, int64(0)).Return(types.BlockMessageIds{
		Ids: []uint32{0, 1, 2},
	}, true).Once()

	// Second message is not found.
	k.On("GetMessage", ctx, uint32(0)).Return(types.DelayedMessage{}, true).Once()
	k.On("GetMessage", ctx, uint32(1)).Return(types.DelayedMessage{}, false).Once()
	k.On("GetMessage", ctx, uint32(2)).Return(types.DelayedMessage{}, true).Once()

	// 2 messages are decoded and routed.
	k.On("DecodeMessage", mock.Anything, mock.Anything).Return(nil).Times(2)

	msgRouter := mockSuccessRouter(ctx)
	k.On("Router").Return(msgRouter).Times(2)

	// For error logging.
	k.On("Logger", ctx).Return(log.NewNopLogger()).Times(1)

	// All deletes are called.
	k.On("DeleteMessage", ctx, uint32(0)).Return(nil).Once()
	k.On("DeleteMessage", ctx, uint32(1)).Return(nil).Once()
	k.On("DeleteMessage", ctx, uint32(2)).Return(nil).Once()
}

func setupMockKeeperDecodeFailure(ctx sdk.Context, k *mocks.DelayMsgKeeper) {
	k.On("GetBlockMessageIds", ctx, int64(0)).Return(types.BlockMessageIds{
		Ids: []uint32{0, 1, 2},
	}, true).Once()

	// All messages found.
	k.On("GetMessage", ctx, mock.Anything).Return(types.DelayedMessage{}, true).Times(3)

	// First decode fails.
	k.On("DecodeMessage", mock.Anything, mock.Anything).Return(fmt.Errorf("failed to decode message")).Once()
	k.On("DecodeMessage", mock.Anything, mock.Anything).Return(nil).Times(2)

	// 2 messages are routed.
	msgRouter := mockSuccessRouter(ctx)
	k.On("Router").Return(msgRouter).Times(2)

	// For error logging.
	k.On("Logger", ctx).Return(log.NewNopLogger()).Times(1)

	// All deletes are called.
	k.On("DeleteMessage", ctx, uint32(0)).Return(nil).Once()
	k.On("DeleteMessage", ctx, uint32(1)).Return(nil).Once()
	k.On("DeleteMessage", ctx, uint32(2)).Return(nil).Once()
}

func setupMockKeeperExecutionFailure(ctx sdk.Context, k *mocks.DelayMsgKeeper) {
	k.On("GetBlockMessageIds", ctx, int64(0)).Return(types.BlockMessageIds{
		Ids: []uint32{0, 1, 2},
	}, true).Once()

	// All messages found.
	k.On("GetMessage", ctx, mock.Anything).Return(types.DelayedMessage{}, true).Times(3)

	// All messages are decoded.
	k.On("DecodeMessage", mock.Anything, mock.Anything).Return(nil).Times(3)

	// 1st message fails to execute. Following 2 succeed.
	successRouter := mockSuccessRouter(ctx)
	failureRouter := mockFailingRouter(ctx)
	k.On("Router").Return(failureRouter).Times(1)
	k.On("Router").Return(successRouter).Times(2)

	// For error logging.
	k.On("Logger", ctx).Return(log.NewNopLogger()).Times(1)

	// All deletes are called.
	k.On("DeleteMessage", ctx, uint32(0)).Return(nil).Once()
	k.On("DeleteMessage", ctx, uint32(1)).Return(nil).Once()
	k.On("DeleteMessage", ctx, uint32(2)).Return(nil).Once()
}

func setupMockKeeperDeletionFailure(ctx sdk.Context, k *mocks.DelayMsgKeeper) {
	k.On("GetBlockMessageIds", ctx, int64(0)).Return(types.BlockMessageIds{
		Ids: []uint32{0, 1, 2},
	}, true).Once()

	// All messages found.
	k.On("GetMessage", ctx, mock.Anything).Return(types.DelayedMessage{}, true).Times(3)

	// All messages are decoded.
	k.On("DecodeMessage", mock.Anything, mock.Anything).Return(nil).Times(3)

	// All messages are routed.
	k.On("Router").Return(mockSuccessRouter(ctx)).Times(3)

	// For error logging.
	k.On("Logger", ctx).Return(log.NewNopLogger()).Times(1)

	// All deletes are called. 2nd delete fails.
	k.On("DeleteMessage", ctx, uint32(0)).Return(nil).Once()
	k.On("DeleteMessage", ctx, uint32(1)).Return(fmt.Errorf("Deletion failure")).Once()
	k.On("DeleteMessage", ctx, uint32(2)).Return(nil).Once()
}

func TestDispatchMessageForBlock_Mixed(t *testing.T) {
	tests := map[string]struct {
		setupMocks func(ctx sdk.Context, k *mocks.DelayMsgKeeper)
	}{
		"No messages - dispatch terminates with no action": {
			setupMocks: setupMockKeeperNoMessages,
		},
		"Unexpected message not found does not affect remaining messages": {
			setupMocks: setupMockKeeperMessageNotFound,
		},
		"Decode failure does not affect remaining messages": {
			setupMocks: setupMockKeeperDecodeFailure,
		},
		"Execution error does not affect remaining messages": {
			setupMocks: setupMockKeeperExecutionFailure,
		},
		"Deletion failure does not affect deletion of remaining messages": {
			setupMocks: setupMockKeeperDeletionFailure,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			k := &mocks.DelayMsgKeeper{}
			ctx := sdktest.NewContextWithBlockHeightAndTime(0, time.Now())
			tc.setupMocks(ctx, k)

			keeper.DispatchMessagesForBlock(k, ctx)

			k.AssertExpectations(t)
		})
	}
}

// generateBridgeEventMsgBytes wraps bridge event in a MsgCompleteBridge and byte-encodes it.
func generateBridgeEventMsgBytes(t *testing.T, event bridgetypes.BridgeEvent) []byte {
	_, k, _, _, _, _ := keepertest.DelayMsgKeepers(t)
	msgCompleteBridge := bridgetypes.MsgCompleteBridge{
		Authority: authtypes.NewModuleAddress(types.ModuleName).String(),
		Event:     event,
	}
	bytes, err := k.EncodeMessage(&msgCompleteBridge)
	require.NoError(t, err)
	return bytes
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
		Msg:         generateBridgeEventMsgBytes(t, constants.BridgeEvent_Id0_Height0),
		BlockHeight: 2,
	}

	tApp := testapp.NewTestAppBuilder().WithGenesisDocFn(func() (genesis cometbfttypes.GenesisDoc) {
		genesis = testapp.DefaultGenesis()
		// Add the delayed message to the genesis state.
		testapp.UpdateGenesisDocWithAppStateForModule(
			&genesis,
			func(genesisState *types.GenesisState) {
				genesisState.DelayedMessages = []*types.DelayedMessage{&delayedMessage}
				genesisState.NumMessages = 1
			},
		)
		return genesis
	}).WithTesting(t).Build()
	ctx := tApp.InitChain()

	// Sanity check: the delayed message is in the keeper scheduled for block 2.
	blockMessageIds, found := tApp.App.DelayMsgKeeper.GetBlockMessageIds(ctx, 2)
	require.True(t, found)
	require.Equal(t, []uint32{0}, blockMessageIds.Ids)

	aliceAccountAddress := sdk.MustAccAddressFromBech32(constants.BridgeEvent_Id0_Height0.Address)

	// Sanity check: at block 1, balances are as expected before the message is sent.
	expectAccountBalance(t, ctx, &tApp, BridgeAccountAddress, BridgeGenesisAccountBalance)
	expectAccountBalance(t, ctx, &tApp, aliceAccountAddress, AliceInitialAccountBalance)

	// Advance to block 2 and invoke delayed message to complete bridge.
	ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

	// Assert: balances have been updated to reflect the executed CompleteBridge message.
	expectAccountBalance(t, ctx, &tApp, BridgeAccountAddress, BridgeExpectedAccountBalance)
	expectAccountBalance(t, ctx, &tApp, aliceAccountAddress, AliceExpectedAccountBalance)

	// Assert: the message has been deleted from the keeper.
	_, found = tApp.App.DelayMsgKeeper.GetMessage(ctx, 0)
	require.False(t, found)

	// The block message ids have also been deleted.
	_, found = tApp.App.DelayMsgKeeper.GetBlockMessageIds(ctx, 2)
	require.False(t, found)
}
