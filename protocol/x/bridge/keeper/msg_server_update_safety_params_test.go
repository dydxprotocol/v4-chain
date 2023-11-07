package keeper_test

import (
	"fmt"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/stretchr/testify/require"
)

func TestMsgServerUpdateSafetyParams(t *testing.T) {
	_, ms, ctx := setupMsgServer(t)

	tests := map[string]struct {
		testMsg      types.MsgUpdateSafetyParams
		expectedResp *types.MsgUpdateSafetyParamsResponse
		expectedErr  string
	}{
		"Success": {
			testMsg: types.MsgUpdateSafetyParams{
				Authority: lib.GovModuleAddress.String(),
				Params: types.SafetyParams{
					IsDisabled:  false,
					DelayBlocks: 100,
				},
			},
			expectedResp: &types.MsgUpdateSafetyParamsResponse{},
		},
		"Failure: invalid authority": {
			testMsg: types.MsgUpdateSafetyParams{
				Authority: "12345",
				Params: types.SafetyParams{
					IsDisabled:  false,
					DelayBlocks: 100,
				},
			},
			expectedErr: fmt.Sprintf(
				"message authority %s is not valid for sending update safety params messages",
				"12345",
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			resp, err := ms.UpdateSafetyParams(ctx, &tc.testMsg)

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
