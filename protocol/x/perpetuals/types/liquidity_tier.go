package types

import (
	errorsmod "cosmossdk.io/errors"
	"math/big"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
)

// - Initial margin is less than or equal to 1.
// - Maintenance fraction is less than or equal to 1.
// - Base position notional is not 0.
func (liquidityTier LiquidityTier) Validate() error {
	if liquidityTier.InitialMarginPpm > MaxInitialMarginPpm {
		return errorsmod.Wrap(ErrInitialMarginPpmExceedsMax, lib.Uint32ToString(liquidityTier.InitialMarginPpm))
	}

	if liquidityTier.MaintenanceFractionPpm > MaxMaintenanceFractionPpm {
		return errorsmod.Wrap(ErrMaintenanceFractionPpmExceedsMax,
			lib.Uint32ToString(liquidityTier.MaintenanceFractionPpm))
	}

	if liquidityTier.BasePositionNotional == 0 {
		return errorsmod.Wrap(ErrBasePositionNotionalIsZero, lib.Uint32ToString(0))
	}

	if liquidityTier.ImpactNotional == 0 {
		return errorsmod.Wrap(ErrImpactNotionalIsZero, lib.Uint32ToString(0))
	}

	return nil
}

// `GetMaintenanceMarginPpm` calculates maintenance margin ppm based on initial margin ppm
// and maintenance fraction ppm.
func (liquidityTier LiquidityTier) GetMaintenanceMarginPpm() uint32 {
	if liquidityTier.MaintenanceFractionPpm > MaxMaintenanceFractionPpm {
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
	// Need to divide by 1 million (done by `BigIntMulPpm`) as both clamp factor and margin are in units of ppm.
	return lib.BigIntMulPpm(
		new(big.Int).SetUint64(uint64(clampFactorPpm)),
		liquidityTier.InitialMarginPpm-liquidityTier.GetMaintenanceMarginPpm(),
	)
}

// GetMarginAdjustmentPpm calculates margin adjustment (in ppm) given quote quantums
// and `liquidityTier`'s base position notional.
//
// The idea is to have margin requirement increase as amount of notional increases. Adjustment
// is `1` for any position smaller than `basePositionNotional` and sqrt of position size
// for larger positions. Formula for marginAdjustmentPpm is:
//
// marginAdjustmentPpm = max(
//
//	oneMillion,
//	sqrt(
//		quoteQuantums * (oneMillion * oneMillion) / basePositionNotional
//	)
//
// )
func (liquidityTier LiquidityTier) GetMarginAdjustmentPpm(bigQuoteQuantums *big.Int) *big.Int {
	bigBasePositionNotional := new(big.Int).SetUint64(liquidityTier.BasePositionNotional)
	if bigQuoteQuantums.Cmp(bigBasePositionNotional) <= 0 {
		return lib.BigIntOneMillion()
	}
	adjustmentFactor := new(big.Int).Mul(bigQuoteQuantums, lib.BigIntOneTrillion())
	adjustmentFactor.Quo(adjustmentFactor, bigBasePositionNotional)
	return adjustmentFactor.Sqrt(adjustmentFactor)
}

// GetAdjustedInitialMarginQuoteQuantums returns adjusted initial margin in quote quantums
// (capped at 100% of notional).
//
// marginQuoteQuantums = adjustedMarginPpm * quoteQuantums / 1_000_000
// = min(1_000_000, marginAdjustmentPpm * marginPpm / 1_000_000) * quoteQuantums / 1_000_000
// = min(quoteQuantums, marginPpm * quoteQuantums * marginAdjustmentPpm / 1_000_000 / 1_000_000)
//
// note: divisions are delayed for precision purposes.
func (liquidityTier LiquidityTier) GetAdjustedInitialMarginQuoteQuantums(bigQuoteQuantums *big.Int) *big.Int {
	marginAdjustmentPpm := liquidityTier.GetMarginAdjustmentPpm(bigQuoteQuantums)

	result := new(big.Int).SetUint64(uint64(liquidityTier.InitialMarginPpm))
	// Multiply `initialMarginPpm` with `quoteQuantums`.
	result = result.Mul(result, bigQuoteQuantums)
	// Multiply above result with `marginAdjustmentPpm`.
	result = result.Mul(result, marginAdjustmentPpm)
	// Divide above result by 1 trillion.
	result = result.Quo(result, lib.BigIntOneTrillion())
	// Cap adjusted initial margin at 100% of notional.
	return lib.BigMin(bigQuoteQuantums, result)
}
