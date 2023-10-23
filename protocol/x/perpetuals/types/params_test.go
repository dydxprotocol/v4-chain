package types_test

import (
	"math"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	"github.com/stretchr/testify/require"
)

func TestParamsValidate(t *testing.T) {
	tests := map[string]struct {
		fundingRateClampFactorPpm uint32
		premiumVoteClampFactorPpm uint32
		minNumVotesPerSample      uint32
		expectedError             error
	}{
		"Validates successfully": {
			fundingRateClampFactorPpm: 6_000_000,
			premiumVoteClampFactorPpm: 60_000_000,
			minNumVotesPerSample:      15,
			expectedError:             nil,
		},
		"Validates successfully: max values": {
			fundingRateClampFactorPpm: math.MaxUint32,
			premiumVoteClampFactorPpm: math.MaxUint32,
			minNumVotesPerSample:      math.MaxUint32,
			expectedError:             nil,
		},
		"Failure: funding rate clamp factor ppm is zero": {
			fundingRateClampFactorPpm: 0,
			premiumVoteClampFactorPpm: 60_000_000,
			minNumVotesPerSample:      15,
			expectedError:             types.ErrFundingRateClampFactorPpmIsZero,
		},
		"Failure: premium vote clamp factor ppm is zero": {
			fundingRateClampFactorPpm: 6_000_000,
			premiumVoteClampFactorPpm: 0,
			minNumVotesPerSample:      15,
			expectedError:             types.ErrPremiumVoteClampFactorPpmIsZero,
		},
		"Failure: MinNumVotesPerSample is zero": {
			fundingRateClampFactorPpm: 6_000_000,
			premiumVoteClampFactorPpm: 60_000_000,
			minNumVotesPerSample:      0,
			expectedError:             types.ErrMinNumVotesPerSampleIsZero,
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			params := &types.Params{
				FundingRateClampFactorPpm: tc.fundingRateClampFactorPpm,
				PremiumVoteClampFactorPpm: tc.premiumVoteClampFactorPpm,
				MinNumVotesPerSample:      tc.minNumVotesPerSample,
			}

			err := params.Validate()
			if tc.expectedError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
