package types

import (
	errorsmod "cosmossdk.io/errors"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/int256"
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

	if liquidityTier.OpenInterestLowerCap > liquidityTier.OpenInterestUpperCap {
		return errorsmod.Wrapf(
			ErrOpenInterestLowerCapLargerThanUpperCap,
			"open_interest_lower_cap: %d, open_interest_upper_cap: %d",
			liquidityTier.OpenInterestLowerCap,
			liquidityTier.OpenInterestUpperCap,
		)
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
	result := int256.NewUnsignedInt(uint64(liquidityTier.InitialMarginPpm))
	result.MulPpm(result, liquidityTier.MaintenanceFractionPpm)
	// convert result to uint32 (which is fine because margin ppm never exceeds 1 million).
	return uint32(result.Uint64())
}

// `GetMaxAbsFundingClampPpm` returns the maximum absolute value according to the funding clamp function:
// `|S| ≤ Clamp Factor * (Initial Margin - Maintenance Margin)`, which can be applied to both
// funding rate clamping and premium vote clamping, each having their own clamp factor.
func (liquidityTier LiquidityTier) GetMaxAbsFundingClampPpm(clampFactorPpm uint32) *int256.Int {
	maintenanceMarginPpm := liquidityTier.GetMaintenanceMarginPpm()
	if maintenanceMarginPpm > liquidityTier.InitialMarginPpm {
		// Invariant broken: maintenance margin fraction should never be larger than initial margin fraction.
		panic(errorsmod.Wrapf(ErrMaintenanceMarginLargerThanInitialMargin, "maintenance fraction ppm: %d",
			liquidityTier.MaintenanceFractionPpm))
	}

	// Need to divide by 1 million (done by `MulPpm`) as both clamp factor and margin are in units of ppm.
	ret := new(int256.Int).SetUint64(uint64(clampFactorPpm))
	return ret.MulPpm(ret, liquidityTier.InitialMarginPpm-maintenanceMarginPpm)
}

// GetInitialMarginQuoteQuantums returns initial margin requirement (IMR) in quote quantums.
//
// Now that OIMF is introduced, the calculation of IMR is as follows:
//
// - Each market has a `Lower Cap` and `Upper Cap` denominated in USDC.
// - Each market already has a `Base IMF`.
// - At any point in time, for each market:
//   - Define
//   - `Open Notional = Open Interest * Oracle Price`
//   - `Scaling Factor = (Open Notional - Lower Cap) / (Upper Cap - Lower Cap)`
//   - `IMF Increase = Scaling Factor * (1 - Base IMF)`
//   - Then a market’s `Effective IMF = Min(Base IMF + Max(IMF Increase, 0), 1.0)`
//
// - I.e. the effective IMF is the base IMF while the OI < lower cap, and increases linearly until OI = Upper Cap,
// at which point the IMF stays at 1.0 (requiring 1:1 collateral for trading).
// - initialMarginQuoteQuantums = scaledInitialMarginPpm * quoteQuantums / 1_000_000
//
// note:
// - divisions are delayed for precision purposes.
func (liquidityTier LiquidityTier) GetInitialMarginQuoteQuantums(
	quoteQuantums *int256.Int,
	openInterestQuoteQuantums *int256.Int,
) *int256.Int {
	openInterestUpperCap := int256.NewUnsignedInt(liquidityTier.OpenInterestUpperCap)

	// If `open_interest` >= `open_interest_upper_cap` where `upper_cap` is non-zero,
	// OIMF = 1.0 so return input quote quantums as the IMR.
	if openInterestQuoteQuantums.Cmp(
		openInterestUpperCap,
	) >= 0 && liquidityTier.OpenInterestUpperCap != 0 {
		return quoteQuantums
	}

	// If `open_interest_upper_cap` is 0, OIMF is disabled。
	// Or if `current_interest` <= `open_interest_lower_cap`, IMF is not scaled.
	// In both cases, use base IMF as OIMF.
	openInterestLowerCap := int256.NewUnsignedInt(liquidityTier.OpenInterestLowerCap)
	baseImr := new(int256.Int).MulPpmRoundUp(quoteQuantums, liquidityTier.InitialMarginPpm)
	if liquidityTier.OpenInterestUpperCap == 0 || openInterestQuoteQuantums.Cmp(
		openInterestLowerCap,
	) <= 0 {
		// Calculate base IMR: multiply `quoteQuantums` with `initialMarginPpm` and divide by 1 million.
		return baseImr
	}

	// If `open_interest_lower_cap` < `open_interest` <= `open_interest_upper_cap`, calculate the scaled OIMF.
	// `Scaling Factor = (Open Notional - Lower Cap) / (Upper Cap - Lower Cap)`
	additionalImr := new(int256.Int)
	additionalImr.Mul(quoteQuantums, additionalImr.Sub(openInterestQuoteQuantums, openInterestLowerCap))
	additionalImr.DivRoundUp(additionalImr, openInterestLowerCap.Sub(openInterestUpperCap, openInterestLowerCap))
	additionalImr.MulPpmRoundUp(additionalImr, lib.OneMillion-liquidityTier.InitialMarginPpm)
	additionalImr = additionalImr.Add(baseImr, additionalImr)

	// Return min(Effective IMR, Quote Quantums)
	if additionalImr.Cmp(quoteQuantums) >= 0 {
		return quoteQuantums
	}
	return additionalImr
}
