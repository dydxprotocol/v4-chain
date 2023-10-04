package keeper_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	"github.com/stretchr/testify/require"
)

type IsPerpetualClobPairActiveResp struct {
	isPerpetualClobPairActive    bool
	isPerpetualClobPairActiveErr error
}

func TestPerformStatefulPremiumVotesValidation(t *testing.T) {
	// In below test cases, perpetual 0 is associated with liquidity tier 0,
	// perpetual 1 is associated with liquidity tier 1, etc.
	// `maxAbsPremiumVotePpm` for each liquidity tier is:
	// liquidity tier 0: 60_000_000 * (100% - 100%) = 0
	// liquidity tier 1: 60_000_000 * (100% - 75%) = 15_000_000
	// liquidity tier 2: 60_000_000 * (100% - 0%) = 60_000_000
	// liquidity tier 3: 60_000_000 * (20% - 10%) = 6_000_000
	// liquidity tier 4: 60_000_000 * (50% - 40%) = 6_000_000
	tests := map[string]struct {
		// Setup.
		votes                         []types.FundingPremium
		isPerpetualClobPairActiveResp *IsPerpetualClobPairActiveResp
		numPerpetuals                 int
		expectedErr                   error
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
		"Error: fails to determine clob pair status": {
			votes: []types.FundingPremium{
				{
					PerpetualId: 0,
					PremiumPpm:  0,
				},
			},
			isPerpetualClobPairActiveResp: &IsPerpetualClobPairActiveResp{
				isPerpetualClobPairActiveErr: clobtypes.ErrInvalidClob,
			},
			numPerpetuals: 1,
			expectedErr:   clobtypes.ErrInvalidClob,
		},
		"Error: rejects the premium vote if the clob pair is initializing": {
			votes: []types.FundingPremium{
				{
					PerpetualId: 0,
					PremiumPpm:  1,
				},
			},
			isPerpetualClobPairActiveResp: &IsPerpetualClobPairActiveResp{
				isPerpetualClobPairActive:    false,
				isPerpetualClobPairActiveErr: nil,
			},
			numPerpetuals: 1,
			expectedErr:   types.ErrPremiumVoteForNonActiveMarket,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup.
			mockPerpetualsClobKeeper := &mocks.PerpetualsClobKeeper{}
			pc := keepertest.PerpetualsKeepersWithClobHelpers(
				t,
				mockPerpetualsClobKeeper,
			)

			// set mock expectations
			for _, vote := range tc.votes {
				isActive := true
				var err error
				if tc.isPerpetualClobPairActiveResp != nil {
					isActive = tc.isPerpetualClobPairActiveResp.isPerpetualClobPairActive
					err = tc.isPerpetualClobPairActiveResp.isPerpetualClobPairActiveErr
				}
				mockPerpetualsClobKeeper.On("IsPerpetualClobPairActive", pc.Ctx, vote.PerpetualId).Once().Return(
					isActive,
					err,
				)
			}

			_ = keepertest.CreateLiquidityTiersAndNPerpetuals(t, pc.Ctx, pc.PerpetualsKeeper, pc.PricesKeeper, tc.numPerpetuals)

			// Run.
			msg := &types.MsgAddPremiumVotes{
				Votes: tc.votes,
			}

			err := pc.PerpetualsKeeper.PerformStatefulPremiumVotesValidation(pc.Ctx, msg)
			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
				return
			}

			require.NoError(t, err)

			mockPerpetualsClobKeeper.AssertExpectations(t)
		})
	}
}
