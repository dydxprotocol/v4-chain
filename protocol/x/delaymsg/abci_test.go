package delaymsg_test

import (
	cometbfttypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	sdktest "github.com/dydxprotocol/v4-chain/protocol/testutil/sdk"
	bridgetypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg"
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestEndBlocker(t *testing.T) {
	k := &mocks.DelayMsgKeeper{}
	ctx := sdktest.NewContextWithBlockHeightAndTime(0, time.Now())
	k.On("DispatchMessagesForBlock", ctx).Return().Once()
	delaymsg.EndBlocker(ctx, k)
	k.AssertExpectations(t)
}

var (
	BridgeAuthority      = authtypes.NewModuleAddress(bridgetypes.ModuleName).String()
	BridgeAccountAddress = sdk.MustAccAddressFromBech32(BridgeAuthority)

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
	bridgeKeeper.On("GetBridgeAuthority").Return(BridgeAuthority)
	bridgeKeeper.On("CompleteBridge", ctx, mock.Anything).Return(nil).Times(len(constants.AllMsgs))

	// Dispatch messages for block 0.

	delaymsg.DispatchMessagesForBlock(k, ctx)

	_, found = k.GetBlockMessageIds(ctx, 0)
	require.False(t, found)

	require.True(t, bridgeKeeper.AssertExpectations(t))
}

func TestDispatchMessageForBlock_Mixed(t *testing.T) {
	tests := map[string]struct {
		setupMocks func(ctx sdk.Context, k *mocks.DelayMsgKeeper)
	}{
		"No messages - dispatch terminates with no action":                {},
		"Unexpected message not found does not affect remaining messages": {},
		"Decode failure does not affect remaining messages":               {},
		"Execution error does not affect remaining messages":              {},
		"Deletion failure does not affect deletion of remaining messages": {},
	}
}

// generateBridgeEventMsgBytes wraps bridge event in a MsgCompleteBridge and byte-encodes it.
func generateBridgeEventMsgBytes(t *testing.T, event bridgetypes.BridgeEvent) []byte {
	_, k, _, _, _, _ := keepertest.DelayMsgKeepers(t)
	msgCompleteBridge := bridgetypes.MsgCompleteBridge{
		Authority: authtypes.NewModuleAddress(bridgetypes.ModuleName).String(),
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

	aliceAccountAddress := sdk.MustAccAddressFromBech32(constants.BridgeEvent_Id0_Height0.Address)

	// Sanity check: at block 1, balances are as expected before the message is sent.
	expectAccountBalance(t, ctx, &tApp, BridgeAccountAddress, BridgeGenesisAccountBalance)
	expectAccountBalance(t, ctx, &tApp, aliceAccountAddress, AliceInitialAccountBalance)

	// Advance to block 2 and invoke delayed message to complete bridge.
	ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

	// Assert: balances have been updated to reflect the executed CompleteBridge message.
	expectAccountBalance(t, ctx, &tApp, BridgeAccountAddress, BridgeExpectedAccountBalance)
	expectAccountBalance(t, ctx, &tApp, aliceAccountAddress, AliceExpectedAccountBalance)
}
