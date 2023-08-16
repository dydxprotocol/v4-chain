package keeper_test

import (
	"fmt"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/stretchr/testify/require"
)

func TestMsgServerUpdateEventParams(t *testing.T) {
	k, ms, ctx := setupMsgServer(t)

	tests := map[string]struct {
		testMsg      types.MsgUpdateEventParams
		expectedResp *types.MsgUpdateEventParamsResponse
		expectedErr  string
	}{
		"Success": {
			testMsg: types.MsgUpdateEventParams{
				Authority: k.GetAuthority(),
				Params: types.EventParams{
					Denom:      "denom",
					EthChainId: 1,
					EthAddress: "ethAddress",
				},
			},
			expectedResp: &types.MsgUpdateEventParamsResponse{},
		},
		"Failure: invalid authority": {
			testMsg: types.MsgUpdateEventParams{
				Authority: "12345",
				Params: types.EventParams{
					Denom:      "denom",
					EthChainId: 1,
					EthAddress: "ethAddress",
				},
			},
			expectedErr: fmt.Sprintf(
				"invalid authority: expected %s, got %s",
				k.GetAuthority(),
				"12345",
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			resp, err := ms.UpdateEventParams(ctx, &tc.testMsg)

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
