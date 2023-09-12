package keeper_test

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/delaymsg"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNumMessages(t *testing.T) {
	tests := map[string]struct {
		delayedMessages []sdk.Msg
	}{
		"No messages": {},
		"Two messages": {
			delayedMessages: []sdk.Msg{
				constants.TestMsg1,
				constants.TestMsg2,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, delaymsg, _, _, _, _ := keepertest.DelayMsgKeepers(t)
			for _, msg := range tc.delayedMessages {
				_, err := delaymsg.DelayMessageByBlocks(ctx, msg, 1)
				require.NoError(t, err)
			}

			wctx := sdk.WrapSDKContext(ctx)
			res, err := delaymsg.NumMessages(wctx, &types.QueryNumMessagesRequest{})
			require.NoError(t, err)
			require.Equal(t, uint32(len(tc.delayedMessages)), res.NumMessages)
		})
	}
}

func TestMessage(t *testing.T) {
	tests := map[string]struct {
		delayedMessage sdk.Msg
		expectedMsg    *codectypes.Any
	}{
		"Not found": {},
		"Found": {
			delayedMessage: constants.TestMsg1,
			expectedMsg:    delaymsg.EncodeMessageToAny(t, constants.TestMsg1),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, delaymsg, _, _, _, _ := keepertest.DelayMsgKeepers(t)
			if tc.delayedMessage != nil {
				_, err := delaymsg.DelayMessageByBlocks(ctx, tc.delayedMessage, 1)
				require.NoError(t, err)
			}
			wctx := sdk.WrapSDKContext(ctx)
			resp, err := delaymsg.Message(wctx, &types.QueryMessageRequest{Id: 0})
			if tc.delayedMessage == nil {
				require.Error(t, err)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.Equal(t, &types.DelayedMessage{
					Id:          0,
					Msg:         tc.expectedMsg,
					BlockHeight: 1,
				}, resp.Message)
			}
		})
	}
}

func TestBlockMessageIds(t *testing.T) {
	tests := map[string]struct {
		delayedMessages []sdk.Msg
	}{
		"Not found": {},
		"Two messages": {
			delayedMessages: []sdk.Msg{
				constants.TestMsg1,
				constants.TestMsg2,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, delaymsg, _, _, _, _ := keepertest.DelayMsgKeepers(t)
			for _, msg := range tc.delayedMessages {
				_, err := delaymsg.DelayMessageByBlocks(ctx, msg, 1)
				require.NoError(t, err)
			}

			wctx := sdk.WrapSDKContext(ctx)
			res, err := delaymsg.BlockMessageIds(wctx, &types.QueryBlockMessageIdsRequest{BlockHeight: 1})

			// Not found.
			if len(tc.delayedMessages) == 0 {
				require.Nil(t, res)
				require.Error(t, err)
			} else {
				// Found: check ids.
				// Construct expected id list.
				expectedIds := make([]uint32, len(tc.delayedMessages))
				for i := range tc.delayedMessages {
					expectedIds[i] = uint32(i)
				}
				require.Equal(
					t,
					expectedIds,
					res.MessageIds,
				)
				require.NoError(t, err)
			}
		})
	}
}

func TestBlockMessageIds_InvalidHeight(t *testing.T) {
	ctx, delaymsg, _, _, _, _ := keepertest.DelayMsgKeepers(t)
	_, err := delaymsg.BlockMessageIds(sdk.WrapSDKContext(ctx), &types.QueryBlockMessageIdsRequest{BlockHeight: -1})
	require.ErrorContains(t, err, "invalid block height")
}
