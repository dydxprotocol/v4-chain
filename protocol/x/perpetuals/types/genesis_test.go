package types_test

import (
	"errors"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	tests := map[string]struct {
		genState      *types.GenesisState
		expectedError error
	}{
		"valid: default": {
			genState:      types.DefaultGenesis(),
			expectedError: nil,
		},
		"valid": {
			genState: &types.GenesisState{
				Perpetuals: []types.Perpetual{
					{
						Params: types.PerpetualParams{
							Id:            0,
							Ticker:        "EXAM-USD",
							LiquidityTier: 0,
						},
						FundingIndex: dtypes.ZeroInt(),
					},
				},
				LiquidityTiers: []types.LiquidityTier{
					{
						Id:                     0,
						Name:                   "Large-Cap",
						InitialMarginPpm:       500_000,
						MaintenanceFractionPpm: 750_000,
						ImpactNotional:         1_000_000_000,
					},
				},
				Params: types.Params{
					FundingRateClampFactorPpm: 3_000_000,
					PremiumVoteClampFactorPpm: 30_000_000,
					MinNumVotesPerSample:      15,
				},
			},
			expectedError: nil,
		},
		"invalid: duplicate perpetual ids": {
			genState: &types.GenesisState{
				Perpetuals: []types.Perpetual{
					{
						Params: types.PerpetualParams{
							Id:            0,
							Ticker:        "EXAM-USD",
							LiquidityTier: 0,
						},
						FundingIndex: dtypes.ZeroInt(),
					},
					{
						Params: types.PerpetualParams{
							Id:            0, // duplicate
							Ticker:        "PERP-USD",
							LiquidityTier: 0,
						},
						FundingIndex: dtypes.ZeroInt(),
					},
				},
				LiquidityTiers: []types.LiquidityTier{
					{
						Id:                     0,
						Name:                   "Large-Cap",
						InitialMarginPpm:       500_000,
						MaintenanceFractionPpm: 750_000,
						ImpactNotional:         1_000_000_000,
					},
				},
				Params: types.Params{
					FundingRateClampFactorPpm: 6_000_000,
					PremiumVoteClampFactorPpm: 60_000_000,
					MinNumVotesPerSample:      15,
				},
			},
			expectedError: errors.New("duplicated perpetual id"),
		},
		"invalid: empty ticker": {
			genState: &types.GenesisState{
				Perpetuals: []types.Perpetual{
					{
						Params: types.PerpetualParams{
							Id:            0,
							Ticker:        "",
							LiquidityTier: 0,
						},
						FundingIndex: dtypes.ZeroInt(),
					},
				},
				LiquidityTiers: []types.LiquidityTier{
					{
						Id:                     0,
						Name:                   "Large-Cap",
						InitialMarginPpm:       500_000,
						MaintenanceFractionPpm: 750_000,
						ImpactNotional:         1_000_000_000,
					},
				},
				Params: types.Params{
					FundingRateClampFactorPpm: 6_000_000,
					PremiumVoteClampFactorPpm: 60_000_000,
					MinNumVotesPerSample:      15,
				},
			},
			expectedError: errors.New("Ticker must be non-empty string"),
		},
		"invalid: initial margin ppm > max": {
			genState: &types.GenesisState{
				Perpetuals: []types.Perpetual{
					{
						Params: types.PerpetualParams{
							Id:            0,
							Ticker:        "EXAM-USD",
							LiquidityTier: 0,
						},
						FundingIndex: dtypes.ZeroInt(),
					},
				},
				LiquidityTiers: []types.LiquidityTier{
					{
						Id:                     0,
						Name:                   "Large-Cap",
						InitialMarginPpm:       1_000_001,
						MaintenanceFractionPpm: 750_000,
						ImpactNotional:         1_000_000_000,
					},
				},
				Params: types.Params{
					FundingRateClampFactorPpm: 6_000_000,
					PremiumVoteClampFactorPpm: 60_000_000,
					MinNumVotesPerSample:      15,
				},
			},
			expectedError: errors.New("InitialMarginPpm exceeds maximum value of 1e6"),
		},
		"invalid: maintenance fraction ppm > max": {
			genState: &types.GenesisState{
				Perpetuals: []types.Perpetual{
					{
						Params: types.PerpetualParams{
							Id:            0,
							Ticker:        "EXAM-USD",
							LiquidityTier: 0,
						},
						FundingIndex: dtypes.ZeroInt(),
					},
				},
				LiquidityTiers: []types.LiquidityTier{
					{
						Id:                     0,
						Name:                   "Large-Cap",
						InitialMarginPpm:       1_000,
						MaintenanceFractionPpm: 1_000_001,
						ImpactNotional:         1_000_000_000,
					},
				},
				Params: types.Params{
					FundingRateClampFactorPpm: 6_000_000,
					PremiumVoteClampFactorPpm: 60_000_000,
					MinNumVotesPerSample:      15,
				},
			},
			expectedError: errors.New("MaintenanceFractionPpm exceeds maximum value of 1e6"),
		},
		"invalid: funding rate clamp factor ppm is zero": {
			genState: &types.GenesisState{
				Perpetuals: []types.Perpetual{
					{
						Params: types.PerpetualParams{
							Id:            0,
							Ticker:        "EXAM-USD",
							LiquidityTier: 0,
						},
						FundingIndex: dtypes.ZeroInt(),
					},
				},
				LiquidityTiers: []types.LiquidityTier{
					{
						Id:                     0,
						Name:                   "Large-Cap",
						InitialMarginPpm:       200_000,
						MaintenanceFractionPpm: 1_000_000,
						ImpactNotional:         2_500_000_000,
					},
				},
				Params: types.Params{
					FundingRateClampFactorPpm: 0,
					PremiumVoteClampFactorPpm: 60_000_000,
					MinNumVotesPerSample:      15,
				},
			},
			expectedError: errors.New("Funding rate clamp factor ppm is zero"),
		},
		"invalid: premium vote clamp factor ppm is zero": {
			genState: &types.GenesisState{
				Perpetuals: []types.Perpetual{
					{
						Params: types.PerpetualParams{
							Id:            0,
							Ticker:        "EXAM-USD",
							LiquidityTier: 0,
						},
						FundingIndex: dtypes.ZeroInt(),
					},
				},
				LiquidityTiers: []types.LiquidityTier{
					{
						Id:                     0,
						Name:                   "Large-Cap",
						InitialMarginPpm:       200_000,
						MaintenanceFractionPpm: 1_000_000,
						ImpactNotional:         2_500_000_000,
					},
				},
				Params: types.Params{
					FundingRateClampFactorPpm: 6_000_000,
					PremiumVoteClampFactorPpm: 0,
					MinNumVotesPerSample:      15,
				},
			},
			expectedError: errors.New("Premium vote clamp factor ppm is zero"),
		},
		"invalid: min num votes per sample": {
			genState: &types.GenesisState{
				Perpetuals: []types.Perpetual{
					{
						Params: types.PerpetualParams{
							Id:            0,
							Ticker:        "EXAM-USD",
							LiquidityTier: 0,
						},
						FundingIndex: dtypes.ZeroInt(),
					},
				},
				LiquidityTiers: []types.LiquidityTier{
					{
						Id:                     0,
						Name:                   "Large-Cap",
						InitialMarginPpm:       200_000,
						MaintenanceFractionPpm: 1_000_000,
						ImpactNotional:         2_500_000_000,
					},
				},
				Params: types.Params{
					FundingRateClampFactorPpm: 6_000_000,
					PremiumVoteClampFactorPpm: 60_000_000,
					MinNumVotesPerSample:      0,
				},
			},
			expectedError: errors.New("MinNumVotesPerSample is zero"),
		},
		"invalid: impact notional is zero": {
			genState: &types.GenesisState{
				Perpetuals: []types.Perpetual{
					{
						Params: types.PerpetualParams{
							Id:            0,
							Ticker:        "EXAM-USD",
							LiquidityTier: 0,
						},
						FundingIndex: dtypes.ZeroInt(),
					},
				},
				LiquidityTiers: []types.LiquidityTier{
					{
						Id:                     0,
						Name:                   "Large-Cap",
						InitialMarginPpm:       200_000,
						MaintenanceFractionPpm: 1_000_000,
						ImpactNotional:         0,
					},
				},
				Params: types.Params{
					FundingRateClampFactorPpm: 6_000_000,
					PremiumVoteClampFactorPpm: 60_000_000,
					MinNumVotesPerSample:      15,
				},
			},
			expectedError: errors.New("Impact notional is zero"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.expectedError == nil {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedError.Error())
			}
		})
	}
}
