package keeper_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/mocks"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/encoding"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	bridgetypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
	"github.com/stretchr/testify/require"
)

// FakeRoutableMsg is a mock sdk.Msg that fools the router into thinking it is a registered message type.
type FakeRoutableMsg struct {
	mocks.Msg
}

// setting XXX_MessageName on the FakeRoutableMsg causes the router to incorrectly return the handler for the
// registered CompleteBridge message type. This is done so that we can bypass the handler check and trigger
// the ValidateBasic error.
func (msg *FakeRoutableMsg) XXX_MessageName() string {
	return "dydxprotocol.bridge.MsgCompleteBridge"
}

// implementing XXX_Size along with XXX_Marshal proto interface methods allows us to simulate an encoding failure.
func (msg *FakeRoutableMsg) XXX_Size() int {
	return 0
}

// implementing XXX_Marshal along with XXX_Size proto interface methods allows us to simulate an encoding failure.
func (msg *FakeRoutableMsg) XXX_Marshal([]byte, bool) ([]byte, error) {
	return nil, fmt.Errorf("Invalid input")
}

// routableInvalidSdkMsg returns a mock sdk.Msg that fools the router into thinking it is a registered message type,
// then fails ValidateBasic.
func routableInvalidSdkMsg() sdk.Msg {
	msg := &FakeRoutableMsg{}
	msg.On("ValidateBasic").Return(fmt.Errorf("Invalid msg"))
	return msg
}

// unencodableSdkMsg returns a mock sdk.Msg that fools the router into thinking it is a registered message type,
// passes ValidateBasic, passes validateSigners, then fails to encode.
func unencodableSdkMsg() sdk.Msg {
	msg := &FakeRoutableMsg{}
	msg.On("ValidateBasic").Return(nil)
	msg.On("GetSigners").Return([]sdk.AccAddress{types.ModuleAddress})
	return msg
}

func TestDelayMessageByBlocks(t *testing.T) {
	tests := map[string]struct {
		testDelayedMsgs []struct {
			msg   sdk.Msg
			delay uint32
		}
		expectedBlockToMessageIds map[uint32]types.BlockMessageIds
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
			expectedBlockToMessageIds: map[uint32]types.BlockMessageIds{
				uint32(blockDelay1): {
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
			expectedBlockToMessageIds: map[uint32]types.BlockMessageIds{
				uint32(blockDelay1): {
					Ids: []uint32{0},
				},
				uint32(blockDelay2): {
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
			expectedBlockToMessageIds: map[uint32]types.BlockMessageIds{
				uint32(blockDelay1): {
					Ids: []uint32{0, 2},
				},
				uint32(blockDelay2): {
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
					BlockHeight: testDelayedMsg.delay,
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

func TestDelayMessageByBlocks_Failures(t *testing.T) {
	tests := map[string]struct {
		msg           sdk.Msg
		expectedError string
		overflow      bool
	}{
		"No handler found": {
			msg:           constants.NoHandlerMsg,
			expectedError: "/testpb.TestMsg: Message not recognized by router",
		},
		"Message fails validation": {
			msg: &bridgetypes.MsgCompleteBridge{
				Authority: bridgetypes.ModuleAddress.String(),
				Event:     constants.BridgeEvent_Id0_Height0,
			},
			expectedError: "message signer must be delaymsg module address: Invalid signer",
		},
		"Message fails to encode": {
			msg:           unencodableSdkMsg(),
			expectedError: "failed to convert message to Any: Invalid input",
		},
		"Block number overflows": {
			msg:           constants.TestMsg1,
			overflow:      true,
			expectedError: "failed to add block delay",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, delaymsg, _, _, _, _ := keepertest.DelayMsgKeepers(t)

			delay := uint32(blockDelay1)
			if tc.overflow {
				ctx = ctx.WithBlockHeight(math.MaxInt64 - blockDelay1 - 1)
				delay = math.MaxUint32
			}

			_, err := delaymsg.DelayMessageByBlocks(ctx, tc.msg, delay)
			require.ErrorContains(t, err, tc.expectedError)
		})
	}
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

	// Assert - message is gone, removed from block message ids.
	_, found := delaymsg.GetMessage(ctx, 0)
	require.False(t, found)

	// Since this was the only message, the block message ids should be empty.
	_, found = delaymsg.GetBlockMessageIds(ctx, 10)
	require.False(t, found)

	// Next delayed message id unaffected.
	require.Equal(t, uint32(1), delaymsg.GetNextDelayedMessageId(ctx))
}

func TestGetNextDelayedMessageId(t *testing.T) {
	ctx, delaymsg, _, _, _, _ := keepertest.DelayMsgKeepers(t)

	// Next delayed message id should be 0.
	require.Equal(t, uint32(0), delaymsg.GetNextDelayedMessageId(ctx))
}

func expectDelayedMessagesAndBlockIds(
	t *testing.T,
	ctx sdk.Context,
	delayMsg *keeper.Keeper,
	delayedMsgs map[uint32]types.DelayedMessage,
	blockMessageIds map[uint32]types.BlockMessageIds,
	expectedNextDelayedMessageId uint32,
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

	require.Equal(t, expectedNextDelayedMessageId, delayMsg.GetNextDelayedMessageId(ctx))
}

func TestGetNextDelayedMessageId_AddAndDeleteMessages(t *testing.T) {
	ctx, delaymsg, _, _, _, _ := keepertest.DelayMsgKeepers(t)

	// No messages.
	require.Equal(t, uint32(0), delaymsg.GetNextDelayedMessageId(ctx))

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
		map[uint32]types.BlockMessageIds{
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
		map[uint32]types.BlockMessageIds{},
		1, // Next delayed message id unaffected.
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
		map[uint32]types.BlockMessageIds{
			10: {
				Ids: []uint32{1}, // Id incremented.
			},
		},
		2, // Next delayed message id incremented.
	)
}

func TestGetMessage_NotFound(t *testing.T) {
	ctx, delaymsg, _, _, _, _ := keepertest.DelayMsgKeepers(t)

	delayedMsg, found := delaymsg.GetMessage(ctx, 0)
	require.False(t, found)
	require.Zero(t, delayedMsg)
}

func TestValidateMsg(t *testing.T) {
	tests := map[string]struct {
		msg           sdk.Msg
		signer        []byte
		expectedError string
	}{
		"No handler found": {
			msg:           constants.NoHandlerMsg,
			signer:        types.ModuleAddress,
			expectedError: "/testpb.TestMsg: Message not recognized by router",
		},
		"Message fails ValidateBasic": {
			msg:           routableInvalidSdkMsg(),
			signer:        types.ModuleAddress,
			expectedError: "message failed basic validation: Invalid msg: Invalid input",
		},
		"Message fails validateSigners": {
			msg: &bridgetypes.MsgCompleteBridge{
				Authority: bridgetypes.ModuleAddress.String(),
				Event:     constants.BridgeEvent_Id0_Height0,
			},
			signer:        []byte("other signer"),
			expectedError: "message signer must be delaymsg module address: Invalid signer",
		},
		"Valid message": {
			msg:    constants.TestMsg1,
			signer: types.ModuleAddress,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, delaymsg, _, _, _, _ := keepertest.DelayMsgKeepers(t)
			err := delaymsg.ValidateMsg(tc.msg, [][]byte{tc.signer})
			if tc.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedError)
			}
		})
	}
}

func TestSetDelayedMessage(t *testing.T) {
	tests := map[string]struct {
		msg    types.DelayedMessage
		expErr string
	}{
		"Success": {
			msg: types.DelayedMessage{
				Id:          0,
				Msg:         encoding.EncodeMessageToAny(t, constants.TestMsg1),
				BlockHeight: 1,
			},
		},
		"nil msg": {
			msg:    types.DelayedMessage{},
			expErr: "failed to delay msg: failed to get message with error 'Delayed msg is nil': Invalid input",
		},
		"invalid msg": {
			msg: types.DelayedMessage{
				Id: 0,
				Msg: encoding.EncodeMessageToAny(
					t,
					&bridgetypes.MsgCompleteBridge{
						Authority: bridgetypes.ModuleAddress.String(),
						Event:     constants.BridgeEvent_Id0_Height0,
					},
				),
				BlockHeight: 1,
			},
			expErr: "failed to delay message: failed to validate with error 'message signer must be delaymsg",
		},
		"invalid block height": {
			msg: types.DelayedMessage{
				Id:          0,
				Msg:         encoding.EncodeMessageToAny(t, constants.TestMsg1),
				BlockHeight: 0,
			},
			expErr: "failed to delay message: block height 0 is in the past: Invalid input",
		},
		"duplicate id": {
			msg: types.DelayedMessage{
				Id:          1,
				Msg:         encoding.EncodeMessageToAny(t, constants.TestMsg1),
				BlockHeight: 1,
			},
			expErr: "failed to delay message: message with id 1 already exists: Invalid input",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, delaymsg, _, _, _, _ := keepertest.DelayMsgKeepers(t)

			// Add a message to the store to test for duplicate message id insertion.
			err := delaymsg.SetDelayedMessage(ctx, &types.DelayedMessage{
				Id:          1,
				Msg:         encoding.EncodeMessageToAny(t, constants.TestMsg1),
				BlockHeight: 0,
			})
			require.NoError(t, err)

			// Setup block height to test past message.
			ctx = ctx.WithBlockHeight(1)

			err = delaymsg.SetDelayedMessage(ctx, &tc.msg)
			if tc.expErr == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expErr)
			}
		})
	}
}
