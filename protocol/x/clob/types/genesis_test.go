package types_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/dtypes"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	tests := map[string]struct {
		genState      *types.GenesisState
		expectedError error
	}{
		"default is valid": {
			genState:      types.DefaultGenesis(),
			expectedError: nil,
		},
		"valid genesis state": {
			genState: &types.GenesisState{
				BlockRateLimitConfig: types.BlockRateLimitConfiguration{
					MaxShortTermOrdersAndCancelsPerNBlocks: []types.MaxPerNBlocksRateLimit{
						{
							NumBlocks: 1,
							Limit:     1,
						},
						{
							NumBlocks: types.MaxShortTermOrdersAndCancelsPerNBlocksNumBlocks,
							Limit:     types.MaxShortTermOrdersAndCancelsPerNBlocksLimit,
						},
					},
					MaxStatefulOrdersPerNBlocks: []types.MaxPerNBlocksRateLimit{
						{
							NumBlocks: 1,
							Limit:     1,
						},
						{
							NumBlocks: types.MaxStatefulOrdersPerNBlocksNumBlocks,
							Limit:     types.MaxStatefulOrdersPerNBlocksLimit,
						},
					},
				},
				ClobPairs: []types.ClobPair{
					{
						Id: uint32(0),
					},
					{
						Id: uint32(1),
					},
				},
				EquityTierLimitConfig: types.EquityTierLimitConfiguration{
					ShortTermOrderEquityTiers: []types.EquityTierLimit{
						{
							UsdTncRequired: dtypes.NewInt(1),
							Limit:          1,
						},
						{
							UsdTncRequired: dtypes.NewInt(2),
							Limit:          types.MaxShortTermOrdersForEquityTier,
						},
					},
					StatefulOrderEquityTiers: []types.EquityTierLimit{
						{
							UsdTncRequired: dtypes.NewInt(1),
							Limit:          1,
						},
						{
							UsdTncRequired: dtypes.NewInt(2),
							Limit:          types.MaxStatefulOrdersForEquityTier,
						},
					},
				},
				LiquidationsConfig: types.LiquidationsConfig{
					InsuranceFundFeePpm: 100_00,
					ValidatorFeePpm:     200_000,
					LiquidityFeePpm:     800_000,
					FillablePriceConfig: types.FillablePriceConfig{
						BankruptcyAdjustmentPpm:           lib.OneMillion + 1,
						SpreadToMaintenanceMarginRatioPpm: 1,
					},
					MaxCumulativeInsuranceFundDelta: uint64(1_000_000_000_000),
				},
			},
			expectedError: nil,
		},
		"duplicated clobPair": {
			genState: &types.GenesisState{
				ClobPairs: []types.ClobPair{
					{
						Id: uint32(0),
					},
					{
						Id: uint32(0),
					},
				},
			},
			expectedError: errors.New("duplicated id for clobPair"),
		},
		"gap in clobPair": {
			genState: &types.GenesisState{
				ClobPairs: []types.ClobPair{
					{
						Id: uint32(0),
					},
					{
						Id: uint32(2),
					},
				},
			},
			expectedError: errors.New("found gap in clobPair id"),
		},
		"spread to maintenance margin ratio of 0 is invalid": {
			genState: &types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					InsuranceFundFeePpm: 100_00,
					ValidatorFeePpm:     200_000,
					LiquidityFeePpm:     800_000,
					FillablePriceConfig: types.FillablePriceConfig{
						BankruptcyAdjustmentPpm:           lib.OneMillion,
						SpreadToMaintenanceMarginRatioPpm: 0,
					},
					MaxCumulativeInsuranceFundDelta: uint64(1_000_000_000_000),
				},
			},
			expectedError: errors.New(
				"0 is not a valid SpreadToMaintenanceMarginRatioPpm: Proposed LiquidationsConfig is invalid"),
		},
		"spread to maintenance margin ratio of greater than one million is valid": {
			genState: &types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					InsuranceFundFeePpm: 100_00,
					ValidatorFeePpm:     200_000,
					LiquidityFeePpm:     800_000,
					FillablePriceConfig: types.FillablePriceConfig{
						BankruptcyAdjustmentPpm:           lib.OneMillion,
						SpreadToMaintenanceMarginRatioPpm: lib.OneMillion + 1,
					},
					MaxCumulativeInsuranceFundDelta: uint64(1_000_000_000_000),
				},
			},
		},
		"bankruptcy adjustment ppm of 0 is invalid": {
			genState: &types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					InsuranceFundFeePpm: 100_00,
					ValidatorFeePpm:     200_000,
					LiquidityFeePpm:     800_000,
					FillablePriceConfig: types.FillablePriceConfig{
						BankruptcyAdjustmentPpm:           0,
						SpreadToMaintenanceMarginRatioPpm: lib.OneMillion,
					},
					MaxCumulativeInsuranceFundDelta: uint64(1_000_000_000_000),
				},
			},
			expectedError: errors.New(
				"0 is not a valid BankruptcyAdjustmentPpm: Proposed LiquidationsConfig is invalid"),
		},
		"bankruptcy adjustment ppm of less than one million is invalid": {
			genState: &types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					InsuranceFundFeePpm: 100_00,
					ValidatorFeePpm:     200_000,
					LiquidityFeePpm:     800_000,
					FillablePriceConfig: types.FillablePriceConfig{
						BankruptcyAdjustmentPpm:           lib.OneMillion - 1,
						SpreadToMaintenanceMarginRatioPpm: lib.OneMillion,
					},
					MaxCumulativeInsuranceFundDelta: uint64(1_000_000_000_000),
				},
			},
			expectedError: errors.New(
				"999999 is not a valid BankruptcyAdjustmentPpm: Proposed LiquidationsConfig is invalid"),
		},
		"max liquidation fee ppm of zero is invalid": {
			genState: &types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					InsuranceFundFeePpm:             0,
					ValidatorFeePpm:                 200_000,
					LiquidityFeePpm:                 800_000,
					FillablePriceConfig:             constants.FillablePriceConfig_Default,
					MaxCumulativeInsuranceFundDelta: uint64(1_000_000_000_000),
				},
			},
			expectedError: errors.New(
				"0 is not a valid InsuranceFundFeePpm: Proposed LiquidationsConfig is invalid"),
		},
		"max liquidation fee ppm of greater than one million is invalid": {
			genState: &types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					InsuranceFundFeePpm:             lib.OneMillion + 1,
					ValidatorFeePpm:                 200_000,
					LiquidityFeePpm:                 800_000,
					FillablePriceConfig:             constants.FillablePriceConfig_Default,
					MaxCumulativeInsuranceFundDelta: uint64(1_000_000_000_000),
				},
			},
			expectedError: errors.New(
				"1000001 is not a valid InsuranceFundFeePpm: Proposed LiquidationsConfig is invalid"),
		},
		"max num blocks for short term order rate limit is zero": {
			genState: &types.GenesisState{
				BlockRateLimitConfig: types.BlockRateLimitConfiguration{
					MaxShortTermOrdersAndCancelsPerNBlocks: []types.MaxPerNBlocksRateLimit{
						{
							NumBlocks: 0,
							Limit:     1,
						},
					},
				},
			},
			expectedError: errors.New("0 is not a valid NumBlocks for MaxShortTermOrdersAndCancelsPerNBlocks"),
		},
		"max limit for short term order rate limit is zero": {
			genState: &types.GenesisState{
				BlockRateLimitConfig: types.BlockRateLimitConfiguration{
					MaxShortTermOrdersAndCancelsPerNBlocks: []types.MaxPerNBlocksRateLimit{
						{
							NumBlocks: 1,
							Limit:     0,
						},
					},
				},
			},
			expectedError: errors.New("0 is not a valid Limit for MaxShortTermOrdersAndCancelsPerNBlocks"),
		},
		"max num blocks for stateful order rate limit is zero": {
			genState: &types.GenesisState{
				BlockRateLimitConfig: types.BlockRateLimitConfiguration{
					MaxStatefulOrdersPerNBlocks: []types.MaxPerNBlocksRateLimit{
						{
							NumBlocks: 0,
							Limit:     1,
						},
					},
				},
			},
			expectedError: errors.New("0 is not a valid NumBlocks for MaxStatefulOrdersPerNBlocks"),
		},
		"max limit for stateful order rate limit is zero": {
			genState: &types.GenesisState{
				BlockRateLimitConfig: types.BlockRateLimitConfiguration{
					MaxStatefulOrdersPerNBlocks: []types.MaxPerNBlocksRateLimit{
						{
							NumBlocks: 1,
							Limit:     0,
						},
					},
				},
			},
			expectedError: errors.New("0 is not a valid Limit for MaxStatefulOrdersPerNBlocks"),
		},
		"max num blocks for short term order rate limit is greater than max": {
			genState: &types.GenesisState{
				BlockRateLimitConfig: types.BlockRateLimitConfiguration{
					MaxShortTermOrdersAndCancelsPerNBlocks: []types.MaxPerNBlocksRateLimit{
						{
							NumBlocks: types.MaxShortTermOrdersAndCancelsPerNBlocksNumBlocks + 1,
							Limit:     1,
						},
					},
				},
			},
			expectedError: fmt.Errorf("%d is not a valid NumBlocks for MaxShortTermOrdersAndCancelsPerNBlocks",
				types.MaxShortTermOrdersAndCancelsPerNBlocksNumBlocks+1),
		},
		"max limit for short term order rate limit is greater than max": {
			genState: &types.GenesisState{
				BlockRateLimitConfig: types.BlockRateLimitConfiguration{
					MaxShortTermOrdersAndCancelsPerNBlocks: []types.MaxPerNBlocksRateLimit{
						{
							NumBlocks: 1,
							Limit:     types.MaxShortTermOrdersAndCancelsPerNBlocksLimit + 1,
						},
					},
				},
			},
			expectedError: fmt.Errorf("%d is not a valid Limit for MaxShortTermOrdersAndCancelsPerNBlocks",
				types.MaxShortTermOrdersAndCancelsPerNBlocksLimit+1),
		},
		"max num blocks for stateful order rate limit is greater than max": {
			genState: &types.GenesisState{
				BlockRateLimitConfig: types.BlockRateLimitConfiguration{
					MaxStatefulOrdersPerNBlocks: []types.MaxPerNBlocksRateLimit{
						{
							NumBlocks: types.MaxStatefulOrdersPerNBlocksNumBlocks + 1,
							Limit:     1,
						},
					},
				},
			},
			expectedError: fmt.Errorf("%d is not a valid NumBlocks for MaxStatefulOrdersPerNBlocks",
				types.MaxStatefulOrdersPerNBlocksNumBlocks+1),
		},
		"max limit for stateful order is greater than max": {
			genState: &types.GenesisState{
				BlockRateLimitConfig: types.BlockRateLimitConfiguration{
					MaxStatefulOrdersPerNBlocks: []types.MaxPerNBlocksRateLimit{
						{
							NumBlocks: 1,
							Limit:     types.MaxStatefulOrdersPerNBlocksLimit + 1,
						},
					},
				},
			},
			expectedError: fmt.Errorf("%d is not a valid Limit for MaxStatefulOrdersPerNBlocks",
				types.MaxStatefulOrdersPerNBlocksLimit+1),
		},
		"duplicate short term order rate limit NumBlocks not allowed": {
			genState: &types.GenesisState{
				BlockRateLimitConfig: types.BlockRateLimitConfiguration{
					MaxShortTermOrdersAndCancelsPerNBlocks: []types.MaxPerNBlocksRateLimit{
						{
							NumBlocks: 1,
							Limit:     1,
						},
						{
							NumBlocks: 1,
							Limit:     2,
						},
					},
				},
			},
			expectedError: fmt.Errorf("Multiple rate limits"),
		},
		"duplicate stateful order rate limit NumBlocks not allowed": {
			genState: &types.GenesisState{
				BlockRateLimitConfig: types.BlockRateLimitConfiguration{
					MaxStatefulOrdersPerNBlocks: []types.MaxPerNBlocksRateLimit{
						{
							NumBlocks: 1,
							Limit:     1,
						},
						{
							NumBlocks: 1,
							Limit:     2,
						},
					},
				},
			},
			expectedError: fmt.Errorf("Multiple rate limits"),
		},
		"out of order short term order equity tier limit UsdTncRequired not allowed": {
			genState: &types.GenesisState{
				EquityTierLimitConfig: types.EquityTierLimitConfiguration{
					ShortTermOrderEquityTiers: []types.EquityTierLimit{
						{
							UsdTncRequired: dtypes.NewInt(1),
							Limit:          5,
						},
						{
							UsdTncRequired: dtypes.NewInt(1),
							Limit:          3,
						},
					},
				},
			},
			expectedError: fmt.Errorf("Expected ShortTermOrderEquityTiers equity tier UsdTncRequired to be strictly ascending."),
		},
		"out of order stateful order equity tier limit UsdTncRequired not allowed": {
			genState: &types.GenesisState{
				EquityTierLimitConfig: types.EquityTierLimitConfiguration{
					StatefulOrderEquityTiers: []types.EquityTierLimit{
						{
							UsdTncRequired: dtypes.NewInt(1),
							Limit:          5,
						},
						{
							UsdTncRequired: dtypes.NewInt(1),
							Limit:          3,
						},
					},
				},
			},
			expectedError: fmt.Errorf("Expected StatefulOrderEquityTiers equity tier UsdTncRequired to be strictly ascending."),
		},
		"short term order equity tier limit UsdTncRequired is nil": {
			genState: &types.GenesisState{
				EquityTierLimitConfig: types.EquityTierLimitConfiguration{
					ShortTermOrderEquityTiers: []types.EquityTierLimit{
						{
							UsdTncRequired: dtypes.SerializableInt{},
							Limit:          5,
						},
					},
				},
			},
			expectedError: fmt.Errorf("not a valid UsdTncRequired"),
		},
		"short term order equity tier limit UsdTncRequired is negative": {
			genState: &types.GenesisState{
				EquityTierLimitConfig: types.EquityTierLimitConfiguration{
					ShortTermOrderEquityTiers: []types.EquityTierLimit{
						{
							UsdTncRequired: dtypes.NewInt(-1),
							Limit:          5,
						},
					},
				},
			},
			expectedError: fmt.Errorf("not a valid UsdTncRequired"),
		},
		"stateful order equity tier limit UsdTncRequired is nil": {
			genState: &types.GenesisState{
				EquityTierLimitConfig: types.EquityTierLimitConfiguration{
					StatefulOrderEquityTiers: []types.EquityTierLimit{
						{
							UsdTncRequired: dtypes.SerializableInt{},
							Limit:          5,
						},
					},
				},
			},
			expectedError: fmt.Errorf("not a valid UsdTncRequired"),
		},
		"stateful order equity tier limit UsdTncRequired is negative": {
			genState: &types.GenesisState{
				EquityTierLimitConfig: types.EquityTierLimitConfiguration{
					StatefulOrderEquityTiers: []types.EquityTierLimit{
						{
							UsdTncRequired: dtypes.NewInt(-1),
							Limit:          5,
						},
					},
				},
			},
			expectedError: fmt.Errorf("not a valid UsdTncRequired"),
		},
		"short term order equity tier limit Limit is greater than max": {
			genState: &types.GenesisState{
				EquityTierLimitConfig: types.EquityTierLimitConfiguration{
					ShortTermOrderEquityTiers: []types.EquityTierLimit{
						{
							UsdTncRequired: dtypes.NewInt(1),
							Limit:          types.MaxShortTermOrdersForEquityTier + 1,
						},
					},
				},
			},
			expectedError: fmt.Errorf("not a valid Limit"),
		},
		"stateful order equity tier limit Limit is greater than max": {
			genState: &types.GenesisState{
				EquityTierLimitConfig: types.EquityTierLimitConfiguration{
					StatefulOrderEquityTiers: []types.EquityTierLimit{
						{
							UsdTncRequired: dtypes.NewInt(1),
							Limit:          types.MaxStatefulOrdersForEquityTier + 1,
						},
					},
				},
			},
			expectedError: fmt.Errorf("not a valid Limit"),
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
