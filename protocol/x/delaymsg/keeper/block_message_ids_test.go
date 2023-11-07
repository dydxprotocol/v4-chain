package keeper_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/stretchr/testify/require"
)

const (
	blockDelay1 = 10
	blockDelay2 = 15
)

var (
	testBlockDelays = []uint32{
		blockDelay1,
		blockDelay1,
		blockDelay2,
		blockDelay1,
		blockDelay2,
		blockDelay2,
	}
	expectedBlock1MessageIds = []uint32{0, 1, 3}
	expectedBlock2MessageIds = []uint32{2, 4, 5}
)

func TestGetBlockMessageIds_ZeroBlockHeight(t *testing.T) {
	ctx, delaymsg, _, _, _, _ := keeper.DelayMsgKeepers(t)
	_, found := delaymsg.GetBlockMessageIds(ctx, 0)
	require.False(t, found)
}

func TestGetBlockMessageIds_DeleteAllMgs(t *testing.T) {
	ctx, delaymsg, _, _, _, _ := keeper.DelayMsgKeepers(t)

	for _, delay := range testBlockDelays {
		_, err := delaymsg.DelayMessageByBlocks(ctx, constants.TestMsg1, delay)
		require.NoError(t, err)
	}

	// Delete all messages from block 10.
	for _, i := range expectedBlock1MessageIds {
		err := delaymsg.DeleteMessage(ctx, i)
		require.NoError(t, err)
	}
	_, found := delaymsg.GetBlockMessageIds(ctx, blockDelay1)
	require.False(t, found)

	// Delete all messages from block 15.
	for _, i := range expectedBlock2MessageIds {
		err := delaymsg.DeleteMessage(ctx, i)
		require.NoError(t, err)
	}
	_, found = delaymsg.GetBlockMessageIds(ctx, blockDelay2)
	require.False(t, found)
}

func TestGetBlockMessageIds_DeleteWithMultipleIds(t *testing.T) {
	tests := map[string]struct {
		idToDelete           uint32
		expectedRemainingIds []uint32
	}{
		"delete first id": {
			idToDelete:           0,
			expectedRemainingIds: []uint32{1, 2},
		},
		"delete middle id": {
			idToDelete:           1,
			expectedRemainingIds: []uint32{0, 2},
		},
		"delete last id": {
			idToDelete:           2,
			expectedRemainingIds: []uint32{0, 1},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup - add all messages to the same block.
			ctx, delaymsg, _, _, _, _ := keeper.DelayMsgKeepers(t)

			for _, msg := range constants.AllMsgs {
				_, err := delaymsg.DelayMessageByBlocks(ctx, msg, 10)
				require.NoError(t, err)
			}

			// Act - delete the first message.
			err := delaymsg.DeleteMessage(ctx, tc.idToDelete)
			require.NoError(t, err)

			// Assert - message is gone, removed from block message ids, and next id is unchanged.
			_, found := delaymsg.GetMessage(ctx, tc.idToDelete)
			require.False(t, found)

			blockMessageIds, found := delaymsg.GetBlockMessageIds(ctx, 10)
			require.True(t, found)
			require.Equal(t, tc.expectedRemainingIds, blockMessageIds.Ids)

			require.Equal(t, uint32(len(constants.AllMsgs)), delaymsg.GetNextDelayedMessageId(ctx))
		})
	}
}
