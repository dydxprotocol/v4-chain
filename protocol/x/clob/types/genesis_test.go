package types_test

import (
	"errors"
	"fmt"
	"github.com/dydxprotocol/v4/dtypes"
	"testing"

	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/testutil/constants"
	"github.com/dydxprotocol/v4/x/clob/types"
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
					MaxShortTermOrdersPerMarketPerNBlocks: []types.MaxPerNBlocksRateLimit{
						{
							NumBlocks: 1,
							Limit:     1,
						},
						{
							NumBlocks: types.MaxShortTermOrdersPerMarketPerNBlocksNumBlocks,
							Limit:     types.MaxShortTermOrdersPerMarketPerNBlocksLimit,
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
					MaxShortTermOrderCancellationsPerMarketPerNBlocks: []types.MaxPerNBlocksRateLimit{
						{
							NumBlocks: 1,
							Limit:     1,
						},
						{
							NumBlocks: types.MaxShortTermOrderCancellationsPerMarketPerNBlocksNumBlocks,
							Limit:     types.MaxShortTermOrderCancellationsPerMarketPerNBlocksLimit,
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
					MaxLiquidationFeePpm: 100_00,
					FillablePriceConfig: types.FillablePriceConfig{
						BankruptcyAdjustmentPpm:           lib.OneMillion + 1,
						SpreadToMaintenanceMarginRatioPpm: 1,
					},
					PositionBlockLimits:   constants.PositionBlockLimits_Default,
					SubaccountBlockLimits: constants.SubaccountBlockLimits_Default,
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
					MaxLiquidationFeePpm: 100_00,
					FillablePriceConfig: types.FillablePriceConfig{
						BankruptcyAdjustmentPpm:           lib.OneMillion,
						SpreadToMaintenanceMarginRatioPpm: 0,
					},
					PositionBlockLimits:   constants.PositionBlockLimits_Default,
					SubaccountBlockLimits: constants.SubaccountBlockLimits_Default,
				},
			},
			expectedError: errors.New(
				"0 is not a valid SpreadToMaintenanceMarginRatioPpm: Proposed LiquidationsConfig is invalid"),
		},
		"spread to maintenance margin ratio of greater than one million is invalid": {
			genState: &types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					MaxLiquidationFeePpm: 100_00,
					FillablePriceConfig: types.FillablePriceConfig{
						BankruptcyAdjustmentPpm:           lib.OneMillion,
						SpreadToMaintenanceMarginRatioPpm: lib.OneMillion + 1,
					},
					PositionBlockLimits:   constants.PositionBlockLimits_Default,
					SubaccountBlockLimits: constants.SubaccountBlockLimits_Default,
				},
			},
			expectedError: errors.New(
				"1000001 is not a valid SpreadToMaintenanceMarginRatioPpm: Proposed LiquidationsConfig is invalid"),
		},
		"bankruptcy adjustment ppm of 0 is invalid": {
			genState: &types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					MaxLiquidationFeePpm: 100_00,
					FillablePriceConfig: types.FillablePriceConfig{
						BankruptcyAdjustmentPpm:           0,
						SpreadToMaintenanceMarginRatioPpm: lib.OneMillion,
					},
					PositionBlockLimits:   constants.PositionBlockLimits_Default,
					SubaccountBlockLimits: constants.SubaccountBlockLimits_Default,
				},
			},
			expectedError: errors.New(
				"0 is not a valid BankruptcyAdjustmentPpm: Proposed LiquidationsConfig is invalid"),
		},
		"bankruptcy adjustment ppm of less than one million is invalid": {
			genState: &types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					MaxLiquidationFeePpm: 100_00,
					FillablePriceConfig: types.FillablePriceConfig{
						BankruptcyAdjustmentPpm:           lib.OneMillion - 1,
						SpreadToMaintenanceMarginRatioPpm: lib.OneMillion,
					},
					PositionBlockLimits:   constants.PositionBlockLimits_Default,
					SubaccountBlockLimits: constants.SubaccountBlockLimits_Default,
				},
			},
			expectedError: errors.New(
				"999999 is not a valid BankruptcyAdjustmentPpm: Proposed LiquidationsConfig is invalid"),
		},
		"max liquidation fee ppm of zero is invalid": {
			genState: &types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					MaxLiquidationFeePpm:  0,
					FillablePriceConfig:   constants.FillablePriceConfig_Default,
					PositionBlockLimits:   constants.PositionBlockLimits_Default,
					SubaccountBlockLimits: constants.SubaccountBlockLimits_Default,
				},
			},
			expectedError: errors.New(
				"0 is not a valid MaxLiquidationFeePpm: Proposed LiquidationsConfig is invalid"),
		},
		"max liquidation fee ppm of greater than one million is invalid": {
			genState: &types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					MaxLiquidationFeePpm:  lib.OneMillion + 1,
					FillablePriceConfig:   constants.FillablePriceConfig_Default,
					PositionBlockLimits:   constants.PositionBlockLimits_Default,
					SubaccountBlockLimits: constants.SubaccountBlockLimits_Default,
				},
			},
			expectedError: errors.New(
				"1000001 is not a valid MaxLiquidationFeePpm: Proposed LiquidationsConfig is invalid"),
		},
		"max position portion liquidated ppm of 0 is invalid": {
			genState: &types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					MaxLiquidationFeePpm: lib.OneMillion,
					FillablePriceConfig:  constants.FillablePriceConfig_Default,
					PositionBlockLimits: types.PositionBlockLimits{
						MinPositionNotionalLiquidated:   1_000,
						MaxPositionPortionLiquidatedPpm: 0,
					},
					SubaccountBlockLimits: constants.SubaccountBlockLimits_Default,
				},
			},
			expectedError: errors.New(
				"0 is not a valid MaxPositionPortionLiquidatedPpm: Proposed LiquidationsConfig is invalid"),
		},
		"max position portion liquidated ppm of greater than one million is invalid": {
			genState: &types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					MaxLiquidationFeePpm: lib.OneMillion,
					FillablePriceConfig:  constants.FillablePriceConfig_Default,
					PositionBlockLimits: types.PositionBlockLimits{
						MinPositionNotionalLiquidated:   1_000,
						MaxPositionPortionLiquidatedPpm: lib.OneMillion + 1,
					},
					SubaccountBlockLimits: constants.SubaccountBlockLimits_Default,
				},
			},
			expectedError: errors.New(
				"1000001 is not a valid MaxPositionPortionLiquidatedPpm: Proposed LiquidationsConfig is invalid"),
		},
		"max notional liquidated of 0 is invalid": {
			genState: &types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					MaxLiquidationFeePpm: lib.OneMillion,
					FillablePriceConfig:  constants.FillablePriceConfig_Default,
					PositionBlockLimits:  constants.PositionBlockLimits_Default,
					SubaccountBlockLimits: types.SubaccountBlockLimits{
						MaxNotionalLiquidated:    0,
						MaxQuantumsInsuranceLost: 100_000_000_000_000,
					},
				},
			},
			expectedError: errors.New(
				"0 is not a valid MaxNotionalLiquidated: Proposed LiquidationsConfig is invalid"),
		},
		"max quantums insurance lost of 0 is invalid": {
			genState: &types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					MaxLiquidationFeePpm: lib.OneMillion,
					FillablePriceConfig:  constants.FillablePriceConfig_Default,
					PositionBlockLimits:  constants.PositionBlockLimits_Default,
					SubaccountBlockLimits: types.SubaccountBlockLimits{
						MaxNotionalLiquidated:    100_000_000_000_000,
						MaxQuantumsInsuranceLost: 0,
					},
				},
			},
			expectedError: errors.New(
				"0 is not a valid MaxQuantumsInsuranceLost: Proposed LiquidationsConfig is invalid"),
		},
		"max num blocks for short term order rate limit is zero": {
			genState: &types.GenesisState{
				BlockRateLimitConfig: types.BlockRateLimitConfiguration{
					MaxShortTermOrdersPerMarketPerNBlocks: []types.MaxPerNBlocksRateLimit{
						{
							NumBlocks: 0,
							Limit:     1,
						},
					},
				},
			},
			expectedError: errors.New("0 is not a valid NumBlocks for MaxShortTermOrdersPerMarketPerNBlocks"),
		},
		"max limit for short term order rate limit is zero": {
			genState: &types.GenesisState{
				BlockRateLimitConfig: types.BlockRateLimitConfiguration{
					MaxShortTermOrdersPerMarketPerNBlocks: []types.MaxPerNBlocksRateLimit{
						{
							NumBlocks: 1,
							Limit:     0,
						},
					},
				},
			},
			expectedError: errors.New("0 is not a valid Limit for MaxShortTermOrdersPerMarketPerNBlocks"),
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
		"max num blocks for short term order cancellation rate limit is zero": {
			genState: &types.GenesisState{
				BlockRateLimitConfig: types.BlockRateLimitConfiguration{
					MaxShortTermOrderCancellationsPerMarketPerNBlocks: []types.MaxPerNBlocksRateLimit{
						{
							NumBlocks: 0,
							Limit:     1,
						},
					},
				},
			},
			expectedError: errors.New("0 is not a valid NumBlocks for MaxShortTermOrderCancellationsPerMarketPerNBlocks"),
		},
		"max limit for short term order cancellation rate limit is zero": {
			genState: &types.GenesisState{
				BlockRateLimitConfig: types.BlockRateLimitConfiguration{
					MaxShortTermOrderCancellationsPerMarketPerNBlocks: []types.MaxPerNBlocksRateLimit{
						{
							NumBlocks: 1,
							Limit:     0,
						},
					},
				},
			},
			expectedError: errors.New("0 is not a valid Limit for MaxShortTermOrderCancellationsPerMarketPerNBlocks"),
		},
		"max num blocks for short term order rate limit is greater than max": {
			genState: &types.GenesisState{
				BlockRateLimitConfig: types.BlockRateLimitConfiguration{
					MaxShortTermOrdersPerMarketPerNBlocks: []types.MaxPerNBlocksRateLimit{
						{
							NumBlocks: types.MaxShortTermOrdersPerMarketPerNBlocksNumBlocks + 1,
							Limit:     1,
						},
					},
				},
			},
			expectedError: fmt.Errorf("%d is not a valid NumBlocks for MaxShortTermOrdersPerMarketPerNBlocks",
				types.MaxShortTermOrdersPerMarketPerNBlocksNumBlocks+1),
		},
		"max limit for short term order rate limit is greater than max": {
			genState: &types.GenesisState{
				BlockRateLimitConfig: types.BlockRateLimitConfiguration{
					MaxShortTermOrdersPerMarketPerNBlocks: []types.MaxPerNBlocksRateLimit{
						{
							NumBlocks: 1,
							Limit:     types.MaxShortTermOrdersPerMarketPerNBlocksLimit + 1,
						},
					},
				},
			},
			expectedError: fmt.Errorf("%d is not a valid Limit for MaxShortTermOrdersPerMarketPerNBlocks",
				types.MaxShortTermOrdersPerMarketPerNBlocksLimit+1),
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
		"max num blocks for short term order cancellation rate limit is greater than max": {
			genState: &types.GenesisState{
				BlockRateLimitConfig: types.BlockRateLimitConfiguration{
					MaxShortTermOrderCancellationsPerMarketPerNBlocks: []types.MaxPerNBlocksRateLimit{
						{
							NumBlocks: types.MaxShortTermOrderCancellationsPerMarketPerNBlocksNumBlocks + 1,
							Limit:     1,
						},
					},
				},
			},
			expectedError: fmt.Errorf("%d is not a valid NumBlocks for MaxShortTermOrderCancellationsPerMarketPerNBlocks",
				types.MaxShortTermOrdersPerMarketPerNBlocksNumBlocks+1),
		},
		"max limit for short term order cancellation rate limit is greater than max": {
			genState: &types.GenesisState{
				BlockRateLimitConfig: types.BlockRateLimitConfiguration{
					MaxShortTermOrderCancellationsPerMarketPerNBlocks: []types.MaxPerNBlocksRateLimit{
						{
							NumBlocks: 1,
							Limit:     types.MaxShortTermOrderCancellationsPerMarketPerNBlocksLimit + 1,
						},
					},
				},
			},
			expectedError: fmt.Errorf("%d is not a valid Limit for MaxShortTermOrderCancellationsPerMarketPerNBlocks",
				types.MaxShortTermOrdersPerMarketPerNBlocksLimit+1),
		},
		"duplicate short term order rate limit NumBlocks not allowed": {
			genState: &types.GenesisState{
				BlockRateLimitConfig: types.BlockRateLimitConfiguration{
					MaxShortTermOrdersPerMarketPerNBlocks: []types.MaxPerNBlocksRateLimit{
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
		"duplicate short term order cancellation rate limit NumBlocks not allowed": {
			genState: &types.GenesisState{
				BlockRateLimitConfig: types.BlockRateLimitConfiguration{
					MaxShortTermOrderCancellationsPerMarketPerNBlocks: []types.MaxPerNBlocksRateLimit{
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
		"duplicate short term order equity tier limit UsdTncRequired not allowed": {
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
			expectedError: fmt.Errorf("Multiple equity tier limits"),
		},
		"duplicate stateful order equity tier limit UsdTncRequired not allowed": {
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
			expectedError: fmt.Errorf("Multiple equity tier limits"),
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
