package types

import (
	"fmt"
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

// GetInitialMarginQuoteQuantums returns the initial margin requirement (IMR) in quote quantums.
func (liquidityTier LiquidityTier) GetInitialMarginQuoteQuantums(
	quoteQuantums *big.Int,
	oiQuoteQuantums *big.Int,
) *big.Int {
	totalImfPpm := liquidityTier.GetAdjustedInitialMarginPpm(oiQuoteQuantums)
	return lib.BigMulPpm(
		quoteQuantums,
		totalImfPpm,
		true, // Round up initial margin.
	)
}

// GetAdjustedInitialMarginPpm returns the adjusted initial margin (in ppm) based on the current open interest.
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
func (liquidityTier LiquidityTier) GetAdjustedInitialMarginPpm(
	oiQuoteQuantums *big.Int,
) *big.Int {
	// If `open_interest_upper_cap` is 0, OIMF is disabled.
	if liquidityTier.OpenInterestUpperCap == 0 {
		return lib.BigU(liquidityTier.InitialMarginPpm)
	}

	// Or if `current_interest` <= `open_interest_lower_cap`, IMF is not scaled.
	oiCapLower := lib.BigU(liquidityTier.OpenInterestLowerCap)
	if oiQuoteQuantums.Cmp(oiCapLower) <= 0 {
		return lib.BigU(liquidityTier.InitialMarginPpm)
	}

	// If `open_interest` >= `open_interest_upper_cap` where `upper_cap` is non-zero, OIMF is 1.
	oiCapUpper := lib.BigU(liquidityTier.OpenInterestUpperCap)
	if oiQuoteQuantums.Cmp(oiCapUpper) >= 0 {
		return lib.BigU(lib.OneMillion)
	}

	// At this point, we have the strict ordering oiCapLower < oiQuoteQuantums < oiCapUpper.
	// Get the ratio of where the current OI is between the lower and upper caps.
	// The ratio should be between 0 and 1 (exclusive).
	capNum := new(big.Int).Sub(oiQuoteQuantums, oiCapLower)
	capDen := new(big.Int).Sub(oiCapUpper, oiCapLower)

	if capNum.Sign() <= 0 || capDen.Sign() <= 0 || capDen.Cmp(capNum) <= 0 {
		panic(fmt.Sprintf("invalid open interest values for liquidity tier %d", liquidityTier.Id))
	}
	if liquidityTier.InitialMarginPpm > lib.OneMillion {
		panic(errorsmod.Wrap(ErrInitialMarginPpmExceedsMax, lib.UintToString(liquidityTier.Id)))
	}

	// Total IMF.
	result := lib.BigU(lib.OneMillion - liquidityTier.InitialMarginPpm)
	result.Mul(result, capNum)
	result.Div(result, capDen)
	result.Add(result, lib.BigU(liquidityTier.InitialMarginPpm))
	return result
}
