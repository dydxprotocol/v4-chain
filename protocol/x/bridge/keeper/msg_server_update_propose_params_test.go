package keeper_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/stretchr/testify/require"
)

func TestMsgServerUpdateProposeParams(t *testing.T) {
	k, ms, ctx := setupMsgServer(t)

	tests := map[string]struct {
		testMsg      types.MsgUpdateProposeParams
		expectedResp *types.MsgUpdateProposeParamsResponse
		expectedErr  string
	}{
		"Success": {
			testMsg: types.MsgUpdateProposeParams{
				Authority: k.GetGovAuthority(),
				Params: types.ProposeParams{
					MaxBridgesPerBlock:           3,
					ProposeDelayDuration:         time.Second,
					SkipRatePpm:                  600_000,
					SkipIfBlockDelayedByDuration: time.Second,
				},
			},
			expectedResp: &types.MsgUpdateProposeParamsResponse{},
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
				"invalid authority: expected %s, got %s",
				k.GetGovAuthority(),
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
