package keeper_test

import (
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	bridgemoduletypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
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
	bridgeKeeper.On("GetBridgeAuthority").Return(authtypes.NewModuleAddress(bridgemoduletypes.ModuleName).String())
	bridgeKeeper.On("CompleteBridge", ctx, mock.Anything).Return(nil).Times(len(constants.AllMsgs))

	// Dispatch messages for block 0.
	delaymsg.DispatchMessagesForBlock(ctx)

	_, found = delaymsg.GetBlockMessageIds(ctx, 0)
	require.False(t, found)

	require.True(t, bridgeKeeper.AssertExpectations(t))
}
