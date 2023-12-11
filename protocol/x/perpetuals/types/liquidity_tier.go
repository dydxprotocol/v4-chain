package types

import (
	"math/big"

	errorsmod "cosmossdk.io/errors"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
)

// - Initial margin is less than or equal to 1.
// - Maintenance fraction is less than or equal to 1.
func (liquidityTier LiquidityTier) Validate() error {
	if liquidityTier.InitialMarginPpm > MaxInitialMarginPpm {
		return errorsmod.Wrap(ErrInitialMarginPpmExceedsMax, lib.UintToString(liquidityTier.InitialMarginPpm))
	}

	if liquidityTier.MaintenanceFractionPpm > MaxMaintenanceFractionPpm {
		return errorsmod.Wrap(ErrMaintenanceFractionPpmExceedsMax,
			lib.UintToString(liquidityTier.MaintenanceFractionPpm))
	}

	if liquidityTier.ImpactNotional == 0 {
		return ErrImpactNotionalIsZero
	}

	return nil
}

// `GetMaintenanceMarginPpm` calculates maintenance margin ppm based on initial margin ppm
// and maintenance fraction ppm.
func (liquidityTier LiquidityTier) GetMaintenanceMarginPpm() uint32 {
	if liquidityTier.MaintenanceFractionPpm > MaxMaintenanceFractionPpm {
		// Invariant broken: `MaintenanceFractionPpm` should always be less than `MaxMaintenanceFractionPpm`,
		// which is checked in `SetLiquidityTier`.
		panic(errorsmod.Wrapf(ErrMaintenanceFractionPpmExceedsMax, "maintenance fraction ppm: %d",
			liquidityTier.MaintenanceFractionPpm))
	}
	// maintenance margin = initial margin * maintenance fraction
	bigMaintenanceMarginPpm := lib.BigIntMulPpm(
		new(big.Int).SetUint64(uint64(liquidityTier.InitialMarginPpm)),
		liquidityTier.MaintenanceFractionPpm,
	)
	// convert result to uint32 (which is fine because margin ppm never exceeds 1 million).
	return uint32(bigMaintenanceMarginPpm.Uint64())
}

// `GetMaxAbsFundingClampPpm` returns the maximum absolute value according to the funding clamp function:
// `|S| ≤ Clamp Factor * (Initial Margin - Maintenance Margin)`, which can be applied to both
// funding rate clamping and premium vote clamping, each having their own clamp factor.
func (liquidityTier LiquidityTier) GetMaxAbsFundingClampPpm(clampFactorPpm uint32) *big.Int {
	maintenanceMarginPpm := liquidityTier.GetMaintenanceMarginPpm()
	if maintenanceMarginPpm > liquidityTier.InitialMarginPpm {
		// Invariant broken: maintenance margin fraction should never be larger than initial margin fraction.
		panic(errorsmod.Wrapf(ErrMaintenanceMarginLargerThanInitialMargin, "maintenance fraction ppm: %d",
			liquidityTier.MaintenanceFractionPpm))
	}

	// Need to divide by 1 million (done by `BigIntMulPpm`) as both clamp factor and margin are in units of ppm.
	return lib.BigIntMulPpm(
		new(big.Int).SetUint64(uint64(clampFactorPpm)),
		liquidityTier.InitialMarginPpm-maintenanceMarginPpm,
	)
}

// GetInitialMarginQuoteQuantums returns initial margin in quote quantums.
//
// marginQuoteQuantums = initialMarginPpm * quoteQuantums / 1_000_000
//
// note: divisions are delayed for precision purposes.
func (liquidityTier LiquidityTier) GetInitialMarginQuoteQuantums(bigQuoteQuantums *big.Int) *big.Int {
	result := new(big.Rat).SetInt(bigQuoteQuantums)
	// Multiply above result with `initialMarginPpm` and divide by 1 million.
	result = lib.BigRatMulPpm(result, liquidityTier.InitialMarginPpm)
	return lib.BigRatRound(result, true) // Round up initial margin.
}
