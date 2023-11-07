package keeper_test

import (
	"fmt"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"testing"
	"time"

	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/stretchr/testify/require"
)

func TestMsgServerUpdateProposeParams(t *testing.T) {
	_, ms, ctx := setupMsgServer(t)

	tests := map[string]struct {
		testMsg      types.MsgUpdateProposeParams
		expectedResp *types.MsgUpdateProposeParamsResponse
		expectedErr  string
	}{
		"Success": {
			testMsg: types.MsgUpdateProposeParams{
				Authority: lib.GovModuleAddress.String(),
				Params: types.ProposeParams{
					MaxBridgesPerBlock:           3,
					ProposeDelayDuration:         time.Second,
					SkipRatePpm:                  600_000,
					SkipIfBlockDelayedByDuration: time.Second,
				},
			},
			expectedResp: &types.MsgUpdateProposeParamsResponse{},
		},
		"Failure: invalid params": {
			testMsg: types.MsgUpdateProposeParams{
				Authority: lib.GovModuleAddress.String(),
				Params: types.ProposeParams{
					MaxBridgesPerBlock:           3,
					ProposeDelayDuration:         -time.Second, // invalid
					SkipRatePpm:                  600_000,
					SkipIfBlockDelayedByDuration: time.Second,
				},
			},
			expectedErr: "Duration is negative",
		},
		"Failure: invalid authority": {
			testMsg: types.MsgUpdateProposeParams{
				Authority: "12345",
				Params: types.ProposeParams{
					MaxBridgesPerBlock:           3,
					ProposeDelayDuration:         time.Second,
					SkipRatePpm:                  600_000,
					SkipIfBlockDelayedByDuration: time.Second,
				},
			},
			expectedErr: fmt.Sprintf(
				"message authority %s is not valid for sending update propose params messages",
				"12345",
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			resp, err := ms.UpdateProposeParams(ctx, &tc.testMsg)

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
