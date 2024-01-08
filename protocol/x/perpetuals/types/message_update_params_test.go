package types_test

import (
	"testing"

	types "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateParams_ValidateBasic(t *testing.T) {
	tests := map[string]struct {
		msg         types.MsgUpdateParams
		expectedErr string
	}{
		"Success": {
			msg: types.MsgUpdateParams{
				Authority: validAuthority,
				Params: types.Params{
					FundingRateClampFactorPpm: 400_000,
					PremiumVoteClampFactorPpm: 400_000,
					MinNumVotesPerSample:      5,
				},
			},
		},
		"Failure: Invalid authority": {
			msg: types.MsgUpdateParams{
				Authority: "",
			},
			expectedErr: "Authority is invalid",
		},
		"Failure: 0 FundingRateClampFactorPpm": {
			msg: types.MsgUpdateParams{
				Authority: validAuthority,
				Params: types.Params{
					FundingRateClampFactorPpm: 0,
					PremiumVoteClampFactorPpm: 400_000,
					MinNumVotesPerSample:      5,
				},
			},
			expectedErr: "Funding rate clamp factor ppm is zero",
		},
		"Failure: 0 PremiumVoteClampFactorPpm": {
			msg: types.MsgUpdateParams{
				Authority: validAuthority,
				Params: types.Params{
					FundingRateClampFactorPpm: 400_000,
					PremiumVoteClampFactorPpm: 0,
					MinNumVotesPerSample:      5,
				},
			},
			expectedErr: "Premium vote clamp factor ppm is zero",
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
