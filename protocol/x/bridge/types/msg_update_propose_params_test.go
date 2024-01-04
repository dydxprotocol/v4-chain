package types_test

import (
	"testing"
	"time"

	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateProposeParams_ValidateBasic(t *testing.T) {
	tests := map[string]struct {
		msg         types.MsgUpdateProposeParams
		expectedErr string
	}{
		"Success": {
			msg: types.MsgUpdateProposeParams{
				Authority: validAuthority,
				Params: types.ProposeParams{
					MaxBridgesPerBlock:           5,
					ProposeDelayDuration:         time.Second,
					SkipRatePpm:                  800_000,
					SkipIfBlockDelayedByDuration: time.Minute,
				},
			},
		},
		"Failure: negative propose delay duration": {
			msg: types.MsgUpdateProposeParams{
				Authority: validAuthority,
				Params: types.ProposeParams{
					MaxBridgesPerBlock:           5,
					ProposeDelayDuration:         -time.Second,
					SkipRatePpm:                  800_000,
					SkipIfBlockDelayedByDuration: time.Minute,
				},
			},
			expectedErr: types.ErrNegativeDuration.Error(),
		},
		"Failure: negative skip if blocked delayed by duration": {
			msg: types.MsgUpdateProposeParams{
				Authority: validAuthority,
				Params: types.ProposeParams{
					MaxBridgesPerBlock:           5,
					ProposeDelayDuration:         time.Second,
					SkipRatePpm:                  800_000,
					SkipIfBlockDelayedByDuration: -time.Minute,
				},
			},
			expectedErr: types.ErrNegativeDuration.Error(),
		},
		"Failure: out-of-bound skip rate ppm": {
			msg: types.MsgUpdateProposeParams{
				Authority: validAuthority,
				Params: types.ProposeParams{
					MaxBridgesPerBlock:           5,
					ProposeDelayDuration:         time.Second,
					SkipRatePpm:                  1_000_001,
					SkipIfBlockDelayedByDuration: time.Minute,
				},
			},
			expectedErr: types.ErrRateOutOfBounds.Error(),
		},
		"Failure: empty authority": {
			msg: types.MsgUpdateProposeParams{
				Authority: "",
			},
			expectedErr: types.ErrInvalidAuthority.Error(),
		},
		"Failure: invalid authority": {
			msg: types.MsgUpdateProposeParams{
				Authority: "dydx1abc",
			},
			expectedErr: types.ErrInvalidAuthority.Error(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedErr)
			}
		})
	}
}
