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
		openInterestLowerCap   uint64
		openInterestUpperCap   uint64
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
		"Failure: lower cap is larger than upper cap": {
			initialMarginPpm:       150_000,       // 15%
			maintenanceFractionPpm: 800_000,       // 80% of IM
			ImpactNotional:         3_333_000_000, // 3_333 USDC
			openInterestLowerCap:   1_000_000,
			openInterestUpperCap:   500_000,
			expectedError:          types.ErrOpenInterestLowerCapLargerThanUpperCap,
		},
		"Failure: lower cap is larger than upper cap (upper cap is zero)": {
			initialMarginPpm:       150_000,       // 15%
			maintenanceFractionPpm: 800_000,       // 80% of IM
			ImpactNotional:         3_333_000_000, // 3_333 USDC
			openInterestLowerCap:   1_000_000,
			openInterestUpperCap:   0,
			expectedError:          types.ErrOpenInterestLowerCapLargerThanUpperCap,
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			liquidityTier := &types.LiquidityTier{
				InitialMarginPpm:       tc.initialMarginPpm,
				MaintenanceFractionPpm: tc.maintenanceFractionPpm,
				ImpactNotional:         tc.ImpactNotional,
				OpenInterestLowerCap:   tc.openInterestLowerCap,
				OpenInterestUpperCap:   tc.openInterestUpperCap,
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

func BenchmarkGetInitialMarginQuoteQuantums(b *testing.B) {
	openInterestLowerCap := uint64(1_000_000)
	openInterestUpperCap := uint64(2_000_000)
	openInterestMiddle := (openInterestLowerCap + openInterestUpperCap) / 2
	liquidityTier := &types.LiquidityTier{
		OpenInterestLowerCap: openInterestLowerCap,
		OpenInterestUpperCap: openInterestUpperCap,
		InitialMarginPpm:     200_000, // 20%
	}
	bigQuoteQuantums := big.NewInt(500_000)
	oiLower := new(big.Int).SetUint64(openInterestLowerCap)
	oiUpper := new(big.Int).SetUint64(openInterestUpperCap)
	oiMiddle := new(big.Int).SetUint64(openInterestMiddle)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = liquidityTier.GetInitialMarginQuoteQuantums(bigQuoteQuantums, oiLower, big.NewInt(0))
		_ = liquidityTier.GetInitialMarginQuoteQuantums(bigQuoteQuantums, oiUpper, big.NewInt(0))
		_ = liquidityTier.GetInitialMarginQuoteQuantums(bigQuoteQuantums, oiMiddle, big.NewInt(0))
	}
}

func TestGetInitialMarginQuoteQuantums(t *testing.T) {
	tests := map[string]struct {
		initialMarginPpm                   uint32
		openInterestLowerCap               uint64
		openInterestUpperCap               uint64
		openInterestNotional               *big.Int
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
		"base IMF = 20%, no OIMF since lower_cap = upper_cap = 0": {
			initialMarginPpm:     uint32(200_000), // 20%
			bigQuoteQuantums:     big.NewInt(500_000),
			openInterestNotional: big.NewInt(1_500_000_000_000),
			openInterestLowerCap: 0,
			openInterestUpperCap: 0,
			// initial margin * quote quantums
			// = 20% * 500_000
			// = 100_000
			expectedInitialMarginQuoteQuantums: big.NewInt(100_000),
		},
		"base IMF = 20%, scaling_factor = 0.5": {
			initialMarginPpm:     uint32(200_000), // 20%
			bigQuoteQuantums:     big.NewInt(500_000),
			openInterestNotional: big.NewInt(1_500_000_000_000),
			openInterestLowerCap: 1_000_000_000_000,
			openInterestUpperCap: 2_000_000_000_000,
			// OIMF = 20% + 0.5 * (1 - 20%) = 20% + 0.5 * 80% = 60%
			// initial margin * quote quantums
			// = 60% * 100% * 500_000
			// = 300_000
			expectedInitialMarginQuoteQuantums: big.NewInt(300_000),
		},
		"base IMF = 10%, scaling_factor = 1 since open_interest >> upper_cap": {
			initialMarginPpm:     uint32(200_000), // 20%
			bigQuoteQuantums:     big.NewInt(500_000),
			openInterestNotional: big.NewInt(80_000_000_000_000),
			openInterestLowerCap: 25_000_000_000_000,
			openInterestUpperCap: 50_000_000_000_000,
			// OIMF = 100%
			// initial margin * quote quantums
			// = 100% * 100% * 500_000
			// = 500_000
			expectedInitialMarginQuoteQuantums: big.NewInt(500_000),
		},
		"base IMF = 10%, open_interest = lower_cap so scaling_factor = 0": {
			initialMarginPpm:     uint32(200_000), // 20%
			bigQuoteQuantums:     big.NewInt(500_000),
			openInterestNotional: big.NewInt(25_000_000_000_000),
			openInterestLowerCap: 25_000_000_000_000,
			openInterestUpperCap: 50_000_000_000_000,
			// OIMF = 20%
			// initial margin * quote quantums
			// = 20% * 100% * 500_000
			// = 100_000
			expectedInitialMarginQuoteQuantums: big.NewInt(100_000),
		},
		"base IMF = 10%, lower_cap < open_interest < upper_cap, realistic numbers": {
			initialMarginPpm:     uint32(200_000), // 20%
			bigQuoteQuantums:     big.NewInt(500_000),
			openInterestNotional: big.NewInt(28_123_456_789_123),
			openInterestLowerCap: 25_000_000_000_000,
			openInterestUpperCap: 60_000_000_000_000,
			// scaling_factor = (28.123 - 25) / (60 - 25) ~= 0.08924
			// OIMF ~= 0.08924 * 80% + 20%
			//      ~= 71392% + 20%
			//      ~= 27.1392%
			// initial margin * quote quantums
			// = 27.1392% * 500_000
			// = 135_697
			expectedInitialMarginQuoteQuantums: big.NewInt(135_697),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			liquidityTier := &types.LiquidityTier{
				InitialMarginPpm:     tc.initialMarginPpm,
				OpenInterestLowerCap: tc.openInterestLowerCap,
				OpenInterestUpperCap: tc.openInterestUpperCap,
			}

			openInterestNotional := big.NewInt(0)
			if tc.openInterestNotional != nil {
				openInterestNotional.Set(tc.openInterestNotional)
			}
			adjustedIMQuoteQuantums := liquidityTier.GetInitialMarginQuoteQuantums(
				tc.bigQuoteQuantums,
				openInterestNotional,
				big.NewInt(0), // no leverage configured
			)

			require.Equal(t, tc.expectedInitialMarginQuoteQuantums, adjustedIMQuoteQuantums)
		})
	}
}

func TestGetAdjustedInitialMarginPpm(t *testing.T) {
	tests := map[string]struct {
		initialMarginPpm     uint32
		openInterestLowerCap uint64
		openInterestUpperCap uint64
		openInterestNotional *big.Int
		expectedPpm          *big.Int
	}{
		"Zero open interest": {
			initialMarginPpm:     uint32(200_000),
			openInterestLowerCap: 1_000_000,
			openInterestUpperCap: 2_000_000,
			openInterestNotional: big.NewInt(0),
			expectedPpm:          big.NewInt(200_000),
		},
		"Open interest within bounds": {
			initialMarginPpm:     uint32(200_000),
			openInterestLowerCap: 1_000_000,
			openInterestUpperCap: 2_000_000,
			openInterestNotional: big.NewInt(1_500_000),
			expectedPpm:          big.NewInt(600_000),
		},
		"Open interest within bounds, rounded": {
			initialMarginPpm:     uint32(200_000),
			openInterestLowerCap: 1_000_000,
			openInterestUpperCap: 2_000_000,
			openInterestNotional: big.NewInt(1_234_567),
			// Base_IMF + OI_IMF_Adjustment =
			// 0.2 + (0.234_567 * 0.8) =
			// 0.387_653_6 (rounded down to 0.387_653)
			expectedPpm: big.NewInt(387_653),
		},
		"Open interest at lower bound": {
			initialMarginPpm:     uint32(200_000),
			openInterestLowerCap: 1_000_000,
			openInterestUpperCap: 2_000_000,
			openInterestNotional: big.NewInt(1_000_000),
			expectedPpm:          big.NewInt(200_000),
		},
		"Open interest at upper bound": {
			initialMarginPpm:     uint32(200_000),
			openInterestLowerCap: 1_000_000,
			openInterestUpperCap: 2_000_000,
			openInterestNotional: big.NewInt(2_000_000),
			expectedPpm:          big.NewInt(1_000_000),
		},
		"Open interest above upper bound": {
			initialMarginPpm:     uint32(200_000),
			openInterestLowerCap: 1_000_000,
			openInterestUpperCap: 2_000_000,
			openInterestNotional: big.NewInt(2_500_000),
			expectedPpm:          big.NewInt(1_000_000),
		},
		"No upper bound, no increase": {
			initialMarginPpm:     uint32(200_000),
			openInterestLowerCap: 0,
			openInterestUpperCap: 0,
			openInterestNotional: big.NewInt(1_500_000),
			expectedPpm:          big.NewInt(200_000),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			liquidityTier := &types.LiquidityTier{
				InitialMarginPpm:     tc.initialMarginPpm,
				OpenInterestLowerCap: tc.openInterestLowerCap,
				OpenInterestUpperCap: tc.openInterestUpperCap,
			}
			adjustedIMQuoteQuantums := liquidityTier.GetAdjustedInitialMarginPpm(tc.openInterestNotional)
			require.Equal(t, tc.expectedPpm, adjustedIMQuoteQuantums, "Adjusted initial margin ppm mismatch")
		})
	}
}
