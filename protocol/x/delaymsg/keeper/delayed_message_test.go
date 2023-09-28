package keeper_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/encoding"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	bridgetypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
	"github.com/stretchr/testify/require"
)

func TestDelayMessageByBlocks(t *testing.T) {
	tests := map[string]struct {
		testDelayedMsgs []struct {
			msg   sdk.Msg
			delay uint32
		}
		expectedBlockToMessageIds map[int64]types.BlockMessageIds
	}{
		"single message": {
			testDelayedMsgs: []struct {
				msg   sdk.Msg
				delay uint32
			}{
				{
					msg:   constants.TestMsg1,
					delay: blockDelay1,
				},
			},
			expectedBlockToMessageIds: map[int64]types.BlockMessageIds{
				int64(blockDelay1): {
					Ids: []uint32{0},
				},
			},
		},
		"multiple messages": {
			testDelayedMsgs: []struct {
				msg   sdk.Msg
				delay uint32
			}{
				{
					msg:   constants.TestMsg1,
					delay: blockDelay1,
				},
				{
					msg:   constants.TestMsg2,
					delay: blockDelay2,
				},
			},
			expectedBlockToMessageIds: map[int64]types.BlockMessageIds{
				int64(blockDelay1): {
					Ids: []uint32{0},
				},
				int64(blockDelay2): {
					Ids: []uint32{1},
				},
			},
		},
		"multiple messages per block": {
			testDelayedMsgs: []struct {
				msg   sdk.Msg
				delay uint32
			}{
				{
					msg:   constants.TestMsg1,
					delay: blockDelay1,
				},
				{
					msg:   constants.TestMsg2,
					delay: blockDelay2,
				},
				{
					msg:   constants.TestMsg3,
					delay: blockDelay1,
				},
			},
			expectedBlockToMessageIds: map[int64]types.BlockMessageIds{
				int64(blockDelay1): {
					Ids: []uint32{0, 2},
				},
				int64(blockDelay2): {
					Ids: []uint32{1},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, delaymsg, _, _, _, _ := keepertest.DelayMsgKeepers(t)

			// Act - add messages.
			for i, testDelayedMsg := range tc.testDelayedMsgs {
				id, err := delaymsg.DelayMessageByBlocks(ctx, testDelayedMsg.msg, testDelayedMsg.delay)
				require.Equal(t, uint32(i), id)
				require.NoError(t, err)
			}

			idToDelayedMsg := make(map[uint32]types.DelayedMessage)
			for i, testDelayedMsg := range tc.testDelayedMsgs {
				idToDelayedMsg[uint32(i)] = types.DelayedMessage{
					Msg:         encoding.EncodeMessageToAny(t, testDelayedMsg.msg),
					BlockHeight: int64(testDelayedMsg.delay),
				}
			}

			// Assert - messages are added.
			expectDelayedMessagesAndBlockIds(
				t,
				ctx,
				delaymsg,
				idToDelayedMsg,
				tc.expectedBlockToMessageIds,
				uint32(len(tc.testDelayedMsgs)),
			)
		})
	}
}

func TestDelayMessageByBlocks_NoHandlerFound(t *testing.T) {
	ctx, delaymsg, _, _, _, _ := keepertest.DelayMsgKeepers(t)
	_, err := delaymsg.DelayMessageByBlocks(ctx, constants.InvalidMsg, blockDelay1)
	require.ErrorContains(t, err, "/testpb.TestMsg: Message not recognized by router")
}

func TestDelayMsgByBlocks_InvalidSigners(t *testing.T) {
	invalidSignerMsg := &bridgetypes.MsgCompleteBridge{
		Authority: authtypes.NewModuleAddress(bridgetypes.ModuleName).String(),
		Event:     constants.BridgeEvent_Id0_Height0,
	}
	ctx, delaymsg, _, _, _, _ := keepertest.DelayMsgKeepers(t)
	_, err := delaymsg.DelayMessageByBlocks(ctx, invalidSignerMsg, blockDelay1)
	require.ErrorContains(t, err, "message signer must be delaymsg module address: Invalid signer")
}

func TestDeleteMessage_NotFound(t *testing.T) {
	ctx, delaymsg, _, _, _, _ := keepertest.DelayMsgKeepers(t)

	err := delaymsg.DeleteMessage(ctx, 0)
	require.EqualError(t, err, "failed to delete message: message with id 0 not found: Invalid input")
}

func TestDeleteMessage(t *testing.T) {
	// Setup - add a message.
	ctx, delaymsg, _, _, _, _ := keepertest.DelayMsgKeepers(t)

	id, err := delaymsg.DelayMessageByBlocks(ctx, constants.TestMsg1, 10)
	require.Equal(t, uint32(0), id)
	require.NoError(t, err)

	// Act - delete the message.
	err = delaymsg.DeleteMessage(ctx, 0)
	require.NoError(t, err)

	// Assert - message is gone, removed from block message ids, and num messages is 0.
	_, found := delaymsg.GetMessage(ctx, 0)
	require.False(t, found)

	// Since this was the only message, the block message ids should be empty.
	_, found = delaymsg.GetBlockMessageIds(ctx, 10)
	require.False(t, found)

	// Message count unaffected.
	require.Equal(t, uint32(1), delaymsg.GetNumMessages(ctx))
}

func TestGetNumMessages(t *testing.T) {
	ctx, delaymsg, _, _, _, _ := keepertest.DelayMsgKeepers(t)

	// No messages.
	require.Equal(t, uint32(0), delaymsg.GetNumMessages(ctx))
}

func expectDelayedMessagesAndBlockIds(
	t *testing.T,
	ctx sdk.Context,
	delayMsg *keeper.Keeper,
	delayedMsgs map[uint32]types.DelayedMessage,
	blockMessageIds map[int64]types.BlockMessageIds,
	expectedNumMessages uint32,
) {
	for i, testDelayedMsg := range delayedMsgs {
		delayedMsg, found := delayMsg.GetMessage(ctx, uint32(i))
		require.True(t, found)
		require.Equal(t, testDelayedMsg.Msg, delayedMsg.Msg)
		require.Equal(t, testDelayedMsg.BlockHeight, delayedMsg.BlockHeight)
	}

	for block, expectedMessageIds := range blockMessageIds {
		actualMessageIds, found := delayMsg.GetBlockMessageIds(ctx, block)
		require.True(t, found)
		require.Equal(t, expectedMessageIds, actualMessageIds)
	}

	require.Equal(t, expectedNumMessages, delayMsg.GetNumMessages(ctx))
}

func TestGetNumMessages_AddAndDeleteMessages(t *testing.T) {
	ctx, delaymsg, _, _, _, _ := keepertest.DelayMsgKeepers(t)

	// No messages.
	require.Equal(t, uint32(0), delaymsg.GetNumMessages(ctx))

	// Add a message, then delete it.
	_, err := delaymsg.DelayMessageByBlocks(ctx, constants.TestMsg1, 10)
	require.NoError(t, err)

	expectDelayedMessagesAndBlockIds(
		t,
		ctx,
		delaymsg,
		map[uint32]types.DelayedMessage{
			0: {
				Msg:         encoding.EncodeMessageToAny(t, constants.TestMsg1),
				BlockHeight: 10,
			},
		},
		map[int64]types.BlockMessageIds{
			10: {
				Ids: []uint32{0},
			},
		},
		1,
	)

	err = delaymsg.DeleteMessage(ctx, 0)
	require.NoError(t, err)

	// No messages.
	expectDelayedMessagesAndBlockIds(
		t,
		ctx,
		delaymsg,
		map[uint32]types.DelayedMessage{},
		map[int64]types.BlockMessageIds{},
		1, // Message count unaffected.
	)

	// Add another message - expect an incremented id.
	_, err = delaymsg.DelayMessageByBlocks(ctx, constants.TestMsg1, 10)
	require.NoError(t, err)

	// Expect a single delayed message with id 1.
	expectDelayedMessagesAndBlockIds(
		t,
		ctx,
		delaymsg,
		map[uint32]types.DelayedMessage{
			1: { // Id incremented.
				Msg:         encoding.EncodeMessageToAny(t, constants.TestMsg1),
				BlockHeight: 10,
			},
		},
		map[int64]types.BlockMessageIds{
			10: {
				Ids: []uint32{1}, // Id incremented.
			},
		},
		2, // Message count incremented.
	)
}

func TestGetMessage_NotFound(t *testing.T) {
	ctx, delaymsg, _, _, _, _ := keepertest.DelayMsgKeepers(t)

	delayedMsg, found := delaymsg.GetMessage(ctx, 0)
	require.False(t, found)
	require.Zero(t, delayedMsg)
}

func TestSetDelayedMessage_Errors(t *testing.T) {
	tests := map[string]struct {
		msg    types.DelayedMessage
		expErr error
	}{
		"invalid block height": {
			msg: types.DelayedMessage{
				Id:          0,
				Msg:         encoding.EncodeMessageToAny(t, constants.TestMsg1),
				BlockHeight: -1,
			},
			expErr: fmt.Errorf("failed to delay message: block height -1 is in the past: Invalid input"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, delaymsg, _, _, _, _ := keepertest.DelayMsgKeepers(t)
			err := delaymsg.SetDelayedMessage(ctx, &tc.msg)
			require.EqualError(t, tc.expErr, err.Error())
		})
	}
}
