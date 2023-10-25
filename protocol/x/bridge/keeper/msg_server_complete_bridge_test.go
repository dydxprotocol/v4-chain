package keeper_test

import (
	"fmt"
	delaymsgtypes "github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
	"testing"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/stretchr/testify/require"
)

func TestMsgServerCompleteBridge(t *testing.T) {
	_, ms, ctx := setupMsgServer(t)

	tests := map[string]struct {
		testMsg      types.MsgCompleteBridge
		expectedResp *types.MsgCompleteBridgeResponse
		expectedErr  string
	}{
		"Success": {
			testMsg: types.MsgCompleteBridge{
				Authority: delaymsgtypes.ModuleAddress.String(),
				Event:     constants.BridgeEvent_Id0_Height0,
			},
			expectedResp: &types.MsgCompleteBridgeResponse{},
		},
		"Failure: invalid address": {
			testMsg: types.MsgCompleteBridge{
				Authority: delaymsgtypes.ModuleAddress.String(),
				Event: types.BridgeEvent{
					Id:             0,
					Coin:           sdk.NewCoin("adv4tnt", sdkmath.NewInt(1)),
					Address:        "invalid",
					EthBlockHeight: 1,
				},
			},
			expectedErr: "decoding bech32 failed",
		},
		"Failure: invalid authority": {
			testMsg: types.MsgCompleteBridge{
				Authority: "12345",
				Event:     constants.BridgeEvent_Id0_Height0,
			},
			expectedErr: fmt.Sprintf(
				"message authority %s is not valid for sending complete bridge messages",
				"12345",
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			resp, err := ms.CompleteBridge(ctx, &tc.testMsg)

			// Assert msg server response.
			require.Equal(t, tc.expectedResp, resp)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
