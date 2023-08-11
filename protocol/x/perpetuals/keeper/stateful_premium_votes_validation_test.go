package keeper_test

import (
	"testing"

	keepertest "github.com/dydxprotocol/v4/testutil/keeper"
	"github.com/dydxprotocol/v4/x/perpetuals/types"
	"github.com/stretchr/testify/require"
)

func TestPerformStatefulPremiumVotesValidation(t *testing.T) {
	// In below test cases, perptual 0 is associated with liquidity tier 0,
	// perpetual 1 is associated with liquidity tier 1, etc.
	// `maxAbsPremiumVotePpm` for each liquidity tier is:
	// liquidity tier 0: 60_000_000 * (100% - 100%) = 0
	// liquidity tier 1: 60_000_000 * (100% - 75%) = 15_000_000
	// liquidity tier 2: 60_000_000 * (100% - 0%) = 60_000_000
	// liquidity tier 3: 60_000_000 * (20% - 10%) = 6_000_000
	// liquidity tier 4: 60_000_000 * (50% - 40%) = 6_000_000
	tests := map[string]struct {
		// Setup.
		votes         []types.FundingPremium
		numPerpetuals int
		expectedErr   error
	}{
		"Valid: empty votes": {
			votes:         []types.FundingPremium{},
			numPerpetuals: 1,
		},
		"Valid: votes on some perpetuals": {
			votes: []types.FundingPremium{
				{
					PerpetualId: 3,
					PremiumPpm:  1000,
				},
			},
			numPerpetuals: 5,
		},
		"Valid: votes on all perpetuals": {
			votes: []types.FundingPremium{
				{
					PerpetualId: 0,
					PremiumPpm:  0,
				},
				{
					PerpetualId: 2,
					PremiumPpm:  60_000_000,
				},
				{
					PerpetualId: 2,
					PremiumPpm:  -60_000_000,
				},
				{
					PerpetualId: 1,
					PremiumPpm:  -15_000_000,
				},
				{
					PerpetualId: 3,
					PremiumPpm:  -6_000_000,
				},
				{
					PerpetualId: 4,
					PremiumPpm:  6_000_000,
				},
				{
					PerpetualId: 4,
					PremiumPpm:  -20_000,
				},
				{
					PerpetualId: 3,
					PremiumPpm:  6_000_000,
				},
			},
			numPerpetuals: 5,
		},
		"Error: perpetual Id does not exist": {
			votes: []types.FundingPremium{
				{
					PerpetualId: 1,
					PremiumPpm:  1000,
				},
				{
					PerpetualId: 5, // invalid
					PremiumPpm:  -1000,
				},
			},
			numPerpetuals: 2,
			expectedErr:   types.ErrPerpetualDoesNotExist,
		},
		"Error: proposed premium vote is not upward clamped - perpetual 0": {
			votes: []types.FundingPremium{
				{
					PerpetualId: 0,
					PremiumPpm:  1,
				},
			},
			numPerpetuals: 1,
			expectedErr:   types.ErrPremiumVoteNotClamped,
		},
		"Error: proposed premium vote is not downward clamped - perpetual 1": {
			votes: []types.FundingPremium{
				{
					PerpetualId: 1,
					PremiumPpm:  -15_000_000 - 1,
				},
			},
			numPerpetuals: 2,
			expectedErr:   types.ErrPremiumVoteNotClamped,
		},
		"Error: proposed premium vote is not upward clamped - perpetual 3": {
			votes: []types.FundingPremium{
				{
					PerpetualId: 3,
					PremiumPpm:  6_000_000 + 1,
				},
			},
			numPerpetuals: 4,
			expectedErr:   types.ErrPremiumVoteNotClamped,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup.
			ctx, k, pricesKeeper, _, _ := keepertest.PerpetualsKeepers(t)

			_, err := createLiquidityTiersAndNPerpetuals(t, ctx, k, pricesKeeper, tc.numPerpetuals)
			require.NoError(t, err)

			// Run.
			msg := &types.MsgAddPremiumVotes{
				Votes: tc.votes,
			}

			err = k.PerformStatefulPremiumVotesValidation(ctx, msg)
			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
				return
			}

			require.NoError(t, err)
		})
	}
}
