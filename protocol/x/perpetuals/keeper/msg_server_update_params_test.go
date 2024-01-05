package keeper_test

import (
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	perpkeeper "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	"github.com/stretchr/testify/require"
)

func TestUpdateParams(t *testing.T) {
	initialParams := types.Params{
		FundingRateClampFactorPpm: 6_000_000,
		PremiumVoteClampFactorPpm: 60_000_000,
		MinNumVotesPerSample:      15,
	}

	tests := map[string]struct {
		msg         *types.MsgUpdateParams
		expectedErr string
	}{
		"Success: modify funding rate clamp factor": {
			msg: &types.MsgUpdateParams{
				Authority: lib.GovModuleAddress.String(),
				Params: types.Params{
					FundingRateClampFactorPpm: 1_234,
					PremiumVoteClampFactorPpm: initialParams.PremiumVoteClampFactorPpm,
					MinNumVotesPerSample:      initialParams.MinNumVotesPerSample,
				},
			},
		},
		"Success: modify premium vote clamp factor and min num votes": {
			msg: &types.MsgUpdateParams{
				Authority: lib.GovModuleAddress.String(),
				Params: types.Params{
					FundingRateClampFactorPpm: initialParams.FundingRateClampFactorPpm,
					PremiumVoteClampFactorPpm: 1_234,
					MinNumVotesPerSample:      7,
				},
			},
		},
		"Failure: parameters are not valid": {
			msg: &types.MsgUpdateParams{
				Authority: lib.GovModuleAddress.String(),
				Params: types.Params{
					FundingRateClampFactorPpm: initialParams.FundingRateClampFactorPpm,
					PremiumVoteClampFactorPpm: 0, // invalid
					MinNumVotesPerSample:      initialParams.MinNumVotesPerSample,
				},
			},
			expectedErr: "Premium vote clamp factor ppm is zero",
		},
		"Failure: empty authority": {
			msg: &types.MsgUpdateParams{
				Authority: "",
				Params: types.Params{
					FundingRateClampFactorPpm: initialParams.FundingRateClampFactorPpm,
					PremiumVoteClampFactorPpm: 1_234,
					MinNumVotesPerSample:      7,
				},
			},
			expectedErr: "invalid authority",
		},
		"Failure: authority is not gov module": {
			msg: &types.MsgUpdateParams{
				Authority: constants.BobAccAddress.String(),
				Params: types.Params{
					FundingRateClampFactorPpm: initialParams.FundingRateClampFactorPpm,
					PremiumVoteClampFactorPpm: 1_234,
					MinNumVotesPerSample:      7,
				},
			},
			expectedErr: "invalid authority",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			pc := keepertest.PerpetualsKeepers(t)
			err := pc.PerpetualsKeeper.SetParams(pc.Ctx, initialParams)
			require.NoError(t, err)

			msgServer := perpkeeper.NewMsgServerImpl(pc.PerpetualsKeeper)

			_, err = msgServer.UpdateParams(pc.Ctx, tc.msg)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
				// Verify that params in state are unchanged.
				got := pc.PerpetualsKeeper.GetParams(pc.Ctx)
				require.Equal(t, initialParams, got)
			} else {
				require.NoError(t, err)
				// Verify that params in state are updated.
				got := pc.PerpetualsKeeper.GetParams(pc.Ctx)
				require.Equal(t, tc.msg.Params, got)
			}
		})
	}
}
