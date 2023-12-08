package types_test

import (
	"math"
	"math/big"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	"github.com/stretchr/testify/require"
)

func TestLiquidityTierValidate(t *testing.T) {
	tests := map[string]struct {
		initialMarginPpm       uint32
		maintenanceFractionPpm uint32
		ImpactNotional         uint64
		expectedError          error
	}{
		"Validates successfully": {
			initialMarginPpm:       150_000,       // 15%
			maintenanceFractionPpm: 800_000,       // 80% of IM
			ImpactNotional:         3_333_000_000, // 3_333 USDC
			expectedError:          nil,
		},
		"Failure: initial margin ppm exceeds max": {
			initialMarginPpm:       1_000_001,     // above 100%
			maintenanceFractionPpm: 800_000,       // 80% of IM
			ImpactNotional:         1_000_000_000, // 1_000 USDC
			expectedError:          types.ErrInitialMarginPpmExceedsMax,
		},
		"Failure: maintenance fraction ppm exceeds max": {
			initialMarginPpm:       1_000_000,     // 100%
			maintenanceFractionPpm: 1_000_001,     // above 100%
			ImpactNotional:         1_000_000_000, // 1_000 USDC
			expectedError:          types.ErrMaintenanceFractionPpmExceedsMax,
		},
		"Failure: impact notional is zero": {
			initialMarginPpm:       1_000_000, // 100%
			maintenanceFractionPpm: 1_000_000, // 100%
			ImpactNotional:         0,         // 0
			expectedError:          types.ErrImpactNotionalIsZero,
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			liquidityTier := &types.LiquidityTier{
				InitialMarginPpm:       tc.initialMarginPpm,
				MaintenanceFractionPpm: tc.maintenanceFractionPpm,
				ImpactNotional:         tc.ImpactNotional,
			}

			err := liquidityTier.Validate()
			if tc.expectedError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestLiquidityTierGetMaintenanceMarginPpm(t *testing.T) {
	tests := map[string]struct {
		initialMarginPpm             uint32
		maintenanceFractionPpm       uint32
		expectedMaintenanceMarginPpm uint32
	}{
		"15% initial margin, 80% maintenance fraction": {
			initialMarginPpm:       uint32(150_000), // 15%
			maintenanceFractionPpm: uint32(800_000), // 80% of IM
			// 15% * 80% = 12%
			expectedMaintenanceMarginPpm: uint32(120_000),
		},
		"45% initial margin, 77% maintenance fraction": {
			initialMarginPpm:       uint32(450_000), // 45%
			maintenanceFractionPpm: uint32(770_000), // 77% of IM
			// 45% * 77% = 34.65%
			expectedMaintenanceMarginPpm: uint32(346_500),
		},
		"63.7239% initial margin, 28.9284% maintenance fraction": {
			initialMarginPpm:       uint32(637_239), // 62.7239%
			maintenanceFractionPpm: uint32(289_284), // 28.9284% of IM
			// 63.7239% * 28.9284% ~= 18.4343%
			expectedMaintenanceMarginPpm: uint32(184_343),
		},
		"100% initial margin, 100% maintenance fraction": {
			initialMarginPpm:       uint32(1_000_000), // 100%
			maintenanceFractionPpm: uint32(1_000_000), // 100% of IM
			// 100% * 100% = 100%
			expectedMaintenanceMarginPpm: uint32(1_000_000),
		},
		"0% initial margin, 100% maintenance fraction": {
			initialMarginPpm:       uint32(0),         // 0%
			maintenanceFractionPpm: uint32(1_000_000), // 100% of IM
			// 0% * 100% = 0%
			expectedMaintenanceMarginPpm: uint32(0),
		},
		"0% initial margin, 0% maintenance fraction": {
			initialMarginPpm:       uint32(0), // 0%
			maintenanceFractionPpm: uint32(0), // 0% of IM
			// 0% * 0% = 0%
			expectedMaintenanceMarginPpm: uint32(0),
		},
		"100% initial margin, 0% maintenance fraction": {
			initialMarginPpm:       uint32(1_000_000), // 100%
			maintenanceFractionPpm: uint32(0),         // 0% of IM
			// 100% * 0% = 0%
			expectedMaintenanceMarginPpm: uint32(0),
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			liquidityTier := &types.LiquidityTier{
				InitialMarginPpm:       tc.initialMarginPpm,
				MaintenanceFractionPpm: tc.maintenanceFractionPpm,
			}

			maintenanceMarginPpm := liquidityTier.GetMaintenanceMarginPpm()
			require.Equal(t, tc.expectedMaintenanceMarginPpm, maintenanceMarginPpm)
		})
	}
}

func TestLiquidityTierGetMaxAbsFundingClampPpm(t *testing.T) {
	tests := map[string]struct {
		clampFactorPpm         uint32
		initialMarginPpm       uint32
		maintenanceFractionPpm uint32
		expectedUpperBoundPpm  *big.Int
	}{
		"600% clamp factor, 15% initial margin, 80% maintenance fraction": {
			clampFactorPpm:         uint32(6_000_000), // 600%
			initialMarginPpm:       uint32(150_000),   // 15%
			maintenanceFractionPpm: uint32(800_000),   // 80% of IM
			// 6_000_000 * (15% - 15% * 80%) = 18%
			expectedUpperBoundPpm: big.NewInt(180_000), //18%
		},
		"600% clamp factor, 45% initial margin, 100% maintenance fraction": {
			clampFactorPpm:         uint32(6_000_000), // 600%
			initialMarginPpm:       uint32(450_000),   // 45%
			maintenanceFractionPpm: uint32(1_000_000), // 100% of IM
			// 6_000_000 * (45% - 45% * 100%) = 0%
			expectedUpperBoundPpm: big.NewInt(0), // 0%
		},
		"600% clamp factor, 0% initial margin, 100% maintenance fraction": {
			clampFactorPpm:         uint32(6_000_000), // 600%
			initialMarginPpm:       uint32(0),         // 0%
			maintenanceFractionPpm: uint32(1_000_000), // 100% of IM
			// 6_000_000 * (0% - 0% * 100%) = 0%
			expectedUpperBoundPpm: big.NewInt(0), // 0%
		},
		"300% clamp factor, 100% initial margin, 0% maintenance fraction": {
			clampFactorPpm:         uint32(3_000_000), // 300%
			initialMarginPpm:       uint32(1_000_000), // 100%
			maintenanceFractionPpm: uint32(0),         // 0% of IM
			// 3_000_000 * (100% - 100% * 0%) = 300%
			expectedUpperBoundPpm: big.NewInt(3_000_000), // 300%
		},
		"0% clamp factor, 100% initial margin, 0% maintenance fraction": {
			clampFactorPpm:         uint32(0),         // 0%
			initialMarginPpm:       uint32(1_000_000), // 100%
			maintenanceFractionPpm: uint32(1_000_000), // 100% of IM
			// 0 * (100% - 100% * 100%) = 0%
			expectedUpperBoundPpm: big.NewInt(0), // 0%
		},
		"max clamp factor, 100% initial margin, 0% maintenance fraction": {
			clampFactorPpm:         uint32(math.MaxUint32), // max uint32 %
			initialMarginPpm:       uint32(1_000_000),      // 100%
			maintenanceFractionPpm: uint32(0),              // 0% of IM
			// max * (max - 100% * 0%) = max
			expectedUpperBoundPpm: big.NewInt(math.MaxUint32), // max uint32 %
		},
		"42_397.8721% clamp factor, 74.5345% initial margin, 56.4947% maintenance fraction": {
			clampFactorPpm:         uint32(423_978_721), // 42_397.8721%
			initialMarginPpm:       uint32(745_345),     // 74.5345%
			maintenanceFractionPpm: uint32(564_947),     // 56.4947% of IM
			// maintenance margin ppm = initial margin ppm * maintenance fraction ~= 421_080
			// 423_978_721 * (745_345 - 421_080) ~= 137_481_459
			expectedUpperBoundPpm: big.NewInt(137_481_459), // 13_748.1459%
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			liquidityTier := &types.LiquidityTier{
				InitialMarginPpm:       tc.initialMarginPpm,
				MaintenanceFractionPpm: tc.maintenanceFractionPpm,
			}

			fundingRateUpperBoundPpm := liquidityTier.GetMaxAbsFundingClampPpm(tc.clampFactorPpm)
			require.Equal(t, tc.expectedUpperBoundPpm, fundingRateUpperBoundPpm)
		})
	}
}

func TestGetInitialMarginQuoteQuantums(t *testing.T) {
	tests := map[string]struct {
		initialMarginPpm                   uint32
		bigQuoteQuantums                   *big.Int
		expectedInitialMarginQuoteQuantums *big.Int
	}{
		"initial margin 20%": {
			initialMarginPpm: uint32(200_000), // 20%
			bigQuoteQuantums: big.NewInt(500_000),
			// initial margin * quote quantums
			// = 20% * 100% * 500_000
			// = 100_000
			expectedInitialMarginQuoteQuantums: big.NewInt(100_000),
		},
		"initial margin 50%": {
			initialMarginPpm: uint32(500_000), // 50%
			bigQuoteQuantums: big.NewInt(1_000_000),
			// initial margin * quote quantums
			// = 50% * 100% * 1_000_000
			// = 500_000
			expectedInitialMarginQuoteQuantums: big.NewInt(500_000),
		},
		"initial margin 10%, quote quantums = 1, should round up to 1": {
			initialMarginPpm: uint32(100_000), // 10%
			bigQuoteQuantums: big.NewInt(1),
			// initial margin * quote quantums
			// = 10% * 1
			// = 0.1 -> round up to 1
			expectedInitialMarginQuoteQuantums: big.NewInt(1),
		},
		"initial margin 56.7243%, quote quantums = 123_456, should round up to 70_030": {
			initialMarginPpm: uint32(567_243), // 56.7243%
			bigQuoteQuantums: big.NewInt(123_456),
			// initial margin * quote quantums
			// = 56.7243% * 123_456
			// ~= 70029.5518 -> round up to 70030
			expectedInitialMarginQuoteQuantums: big.NewInt(70_030),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			liquidityTier := &types.LiquidityTier{
				InitialMarginPpm: tc.initialMarginPpm,
			}
			adjustedIMQuoteQuantums := liquidityTier.GetInitialMarginQuoteQuantums(tc.bigQuoteQuantums)

			require.Equal(t, tc.expectedInitialMarginQuoteQuantums, adjustedIMQuoteQuantums)
		})
	}
}
