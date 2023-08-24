package keeper_test

import (
	cometbfttypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	bridgetypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	BridgeAuthority = authtypes.NewModuleAddress(bridgetypes.ModuleName).String()

	BridgeGenesisAccountBalance = sdk.NewCoin("dv4tnt", sdk.NewInt(1000000000))
	AliceInitialAccountBalance  = sdk.NewCoin("dv4tnt", sdk.NewInt(99500000000))

	BridgeExpectdAccountBalance = sdk.NewCoin("dv4tnt", sdk.NewInt(1000000000-888))
	AliceExpectedAccountBalance = sdk.NewCoin("dv4tnt", sdk.NewInt(99500000000+888))
)

func TestDispatchMessagesForBlock(t *testing.T) {
	ctx, delaymsg, _, bridgeKeeper, _ := keepertest.DelayMsgKeeperWithMockBridgeKeeper(t)

	// Add messages to the keeper.
	for i, msg := range constants.AllMsgs {
		id, err := delaymsg.DelayMessageByBlocks(ctx, msg, 0)
		require.NoError(t, err)
		require.Equal(t, uint32(i), id)
	}

	// Sanity check: messages appear for block 0.
	blockMessageIds, found := delaymsg.GetBlockMessageIds(ctx, 0)
	require.True(t, found)
	require.Equal(t, []uint32{0, 1, 2}, blockMessageIds.Ids)

	// Mock the bridge keeper methods called by the bridge msg server.
	bridgeKeeper.On("GetBridgeAuthority").Return(BridgeAuthority).Once()
	bridgeKeeper.On("CompleteBridge", ctx, mock.Anything).Return(nil).Times(len(constants.AllMsgs))

	// Dispatch messages for block 0.
	delaymsg.DispatchMessagesForBlock(ctx)

	_, found = delaymsg.GetBlockMessageIds(ctx, 0)
	require.False(t, found)

	require.True(t, bridgeKeeper.AssertExpectations(t))
}

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

func TestSendDelayedCompleteBridgeMessage(t *testing.T) {
	// Create an encoded bridge event set to occur at block 10.
	// Expect that Alice's account will increase by 888 coins at block 1.
	// Bridge module account will also decrease by 888 coins at block 1.
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

	bridgeAccountAddress := sdk.MustAccAddressFromBech32(BridgeAuthority)
	aliceAccountAddress := sdk.MustAccAddressFromBech32(constants.BridgeEvent_Id0_Height0.Address)
	denom := constants.BridgeEvent_Id0_Height0.Coin.Denom

	bridgeAccountBalance := tApp.App.BankKeeper.GetBalance(ctx, bridgeAccountAddress, denom)
	aliceAccountBalance := tApp.App.BankKeeper.GetBalance(ctx, aliceAccountAddress, denom)
	t.Log("BridgeAccountAddress", bridgeAccountAddress)
	t.Log("AliceAccountAddress", constants.BridgeEvent_Id0_Height0.Address)
	t.Log("BridgeAccountBalance", bridgeAccountBalance)
	t.Log("AliceAccountBalance", aliceAccountBalance)

	// Sanity check: balances are as expected before the message is sent.
	require.Equal(t, BridgeGenesisAccountBalance.Amount, bridgeAccountBalance.Amount)
	require.Equal(t, BridgeGenesisAccountBalance.Denom, bridgeAccountBalance.Denom)
	require.Equal(t, AliceInitialAccountBalance.Amount, aliceAccountBalance.Amount)
	require.Equal(t, AliceInitialAccountBalance.Denom, aliceAccountBalance.Denom)

	tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{})

	bridgeAccountBalance = tApp.App.BankKeeper.GetBalance(ctx, bridgeAccountAddress, denom)
	aliceAccountBalance = tApp.App.BankKeeper.GetBalance(ctx, aliceAccountAddress, denom)

	t.Log("BridgeAccountBalance", bridgeAccountBalance)
	t.Log("AliceAccountBalance", aliceAccountBalance)

	// Sanity check: balances are as expected before the message is sent.
	//require.Equal(t, BridgeExpectdAccountBalance.Amount, bridgeAccountBalance.Amount)
	require.Equal(t, BridgeExpectdAccountBalance.Denom, bridgeAccountBalance.Denom)
	require.Equal(t, AliceExpectedAccountBalance.Amount, aliceAccountBalance.Amount)
	require.Equal(t, AliceExpectedAccountBalance.Denom, aliceAccountBalance.Denom)
}
