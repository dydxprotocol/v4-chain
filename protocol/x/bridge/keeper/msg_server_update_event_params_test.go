package keeper_test

import (
	"fmt"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/stretchr/testify/require"
)

func TestMsgServerUpdateEventParams(t *testing.T) {
	_, ms, ctx := setupMsgServer(t)

	tests := map[string]struct {
		testMsg      types.MsgUpdateEventParams
		expectedResp *types.MsgUpdateEventParamsResponse
		expectedErr  string
	}{
		"Success": {
			testMsg: types.MsgUpdateEventParams{
				Authority: lib.GovModuleAddress.String(),
				Params: types.EventParams{
					Denom:      "denom",
					EthChainId: 1,
					EthAddress: "ethAddress",
				},
			},
			expectedResp: &types.MsgUpdateEventParamsResponse{},
		},
		"Failure: invalid params": {
			testMsg: types.MsgUpdateEventParams{
				Authority: lib.GovModuleAddress.String(),
				Params: types.EventParams{
					Denom:      "7coin",
					EthChainId: 1,
					EthAddress: "ethAddress",
				},
			},
			expectedErr: "invalid denom",
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
				"message authority %s is not valid for sending update event params messages",
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
