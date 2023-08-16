package keeper_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/stretchr/testify/require"
)

func TestMsgServerCompleteBridge(t *testing.T) {
	k, ms, ctx := setupMsgServer(t)

	tests := map[string]struct {
		testMsg      types.MsgCompleteBridge
		expectedResp *types.MsgCompleteBridgeResponse
		expectedErr  string
	}{
		"Success": {
			testMsg: types.MsgCompleteBridge{
				Authority: k.GetSelfAuthority(),
				Event:     constants.BridgeEvent_Id0_Height0,
			},
			expectedResp: &types.MsgCompleteBridgeResponse{},
		},
		"Failure: invalid address to mint to": {
			testMsg: types.MsgCompleteBridge{
				Authority: k.GetSelfAuthority(),
				Event: types.BridgeEvent{
					Id:             0,
					Coin:           sdk.NewCoin("dv4tnt", sdk.NewInt(1)),
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
				"expected %s, got %s: Authority is invalid",
				k.GetSelfAuthority(),
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
