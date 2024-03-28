package types

import (
	"fmt"
	"math/big"

	errorsmod "cosmossdk.io/errors"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/holiman/uint256"
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
	bigQuoteQuantums *big.Int,
	openInterestQuoteQuantums *big.Int,
) *big.Int {
	bigOpenInterestUpperCap := new(big.Int).SetUint64(liquidityTier.OpenInterestUpperCap)

	// If `open_interest` >= `open_interest_upper_cap` where `upper_cap` is non-zero,
	// OIMF = 1.0 so return input quote quantums as the IMR.
	if openInterestQuoteQuantums.Cmp(
		bigOpenInterestUpperCap,
	) >= 0 && liquidityTier.OpenInterestUpperCap != 0 {
		return bigQuoteQuantums
	}

	ratQuoteQuantums := new(big.Rat).SetInt(bigQuoteQuantums)

	// If `open_interest_upper_cap` is 0, OIMF is disabled。
	// Or if `current_interest` <= `open_interest_lower_cap`, IMF is not scaled.
	// In both cases, use base IMF as OIMF.
	bigOpenInterestLowerCap := new(big.Int).SetUint64(liquidityTier.OpenInterestLowerCap)
	if liquidityTier.OpenInterestUpperCap == 0 || openInterestQuoteQuantums.Cmp(
		bigOpenInterestLowerCap,
	) <= 0 {
		// Calculate base IMR: multiply `bigQuoteQuantums` with `initialMarginPpm` and divide by 1 million.
		ratBaseIMR := lib.BigRatMulPpm(ratQuoteQuantums, liquidityTier.InitialMarginPpm)
		return lib.BigRatRound(ratBaseIMR, true) // Round up initial margin.
	}

	// If `open_interest_lower_cap` < `open_interest` <= `open_interest_upper_cap`, calculate the scaled OIMF.
	// `Scaling Factor = (Open Notional - Lower Cap) / (Upper Cap - Lower Cap)`
	ratScalingFactor := new(big.Rat).SetFrac(
		new(big.Int).Sub(
			openInterestQuoteQuantums, // reuse pointer for memory efficiency
			bigOpenInterestLowerCap,
		),
		bigOpenInterestUpperCap.Sub(
			bigOpenInterestUpperCap, // reuse pointer for memory efficiency
			bigOpenInterestLowerCap,
		),
	)

	// `IMF Increase = Scaling Factor * (1 - Base IMF)`
	ratIMFIncrease := lib.BigRatMulPpm(
		ratScalingFactor,
		lib.OneMillion-liquidityTier.InitialMarginPpm, // >= 0, since we check in `liquidityTier.Validate()`
	)

	// Calculate `Max(IMF Increase, 0)`.
	if ratIMFIncrease.Sign() < 0 {
		panic(
			fmt.Sprintf(
				"GetInitialMarginQuoteQuantums: IMF Increase is negative (%s), liquidityTier: %+v, openInterestQuoteQuantums: %s",
				ratIMFIncrease.String(),
				liquidityTier,
				openInterestQuoteQuantums.String(),
			),
		)
	}

	// First, calculate base IMF in big.Rat
	ratBaseIMF := new(big.Rat).SetFrac64(
		int64(liquidityTier.InitialMarginPpm), // safe, since `InitialMargin` is uint32
		int64(lib.OneMillion),
	)

	// `Effective IMF = Min(Base IMF + Max(IMF Increase, 0), 1.0)`
	ratEffectiveIMF := ratBaseIMF.Add(
		ratBaseIMF, // reuse pointer for memory efficiency
		ratIMFIncrease,
	)

	// `Effective IMR = Effective IMF * Quote Quantums`
	bigIMREffective := lib.BigRatRound(
		ratEffectiveIMF.Mul(
			ratEffectiveIMF, // reuse pointer for memory efficiency
			ratQuoteQuantums,
		),
		true, // Round up initial margin.
	)

	// Return min(Effective IMR, Quote Quantums)
	if bigIMREffective.Cmp(bigQuoteQuantums) >= 0 {
		return bigQuoteQuantums
	}
	return bigIMREffective
}

func (liquidityTier LiquidityTier) GetInitialMarginQuoteQuantumsUint256(
	quoteQuantums *uint256.Int,
	openInterestQuoteQuantums *uint256.Int,
) *uint256.Int {
	openInterestUpperCap := uint256.NewInt(liquidityTier.OpenInterestUpperCap)

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
	openInterestLowerCap := uint256.NewInt(liquidityTier.OpenInterestLowerCap)
	if liquidityTier.OpenInterestUpperCap == 0 || openInterestQuoteQuantums.Cmp(
		openInterestLowerCap,
	) <= 0 {
		// Calculate base IMR: multiply `bigQuoteQuantums` with `initialMarginPpm` and divide by 1 million.
		return lib.MulPpmUint256(quoteQuantums, liquidityTier.InitialMarginPpm)
	}

	// If `open_interest_lower_cap` < `open_interest` <= `open_interest_upper_cap`, calculate the scaled OIMF.
	// `Scaling Factor = (Open Notional - Lower Cap) / (Upper Cap - Lower Cap)`
	scalingFactor := new(uint256.Int).Div(
		new(uint256.Int).Sub(
			openInterestQuoteQuantums, // reuse pointer for memory efficiency
			openInterestLowerCap,
		),
		openInterestLowerCap.Sub(
			openInterestUpperCap, // reuse pointer for memory efficiency
			openInterestLowerCap,
		),
	)

	// `IMF Increase = Scaling Factor * (1 - Base IMF)`
	imfIncrease := lib.MulPpmUint256(
		scalingFactor,
		lib.OneMillion-liquidityTier.InitialMarginPpm, // >= 0, since we check in `liquidityTier.Validate()`
	)

	// Calculate `Max(IMF Increase, 0)`.
	if imfIncrease.Sign() < 0 {
		panic(
			fmt.Sprintf(
				"GetInitialMarginQuoteQuantums: IMF Increase is negative (%s), liquidityTier: %+v, openInterestQuoteQuantums: %s",
				imfIncrease.String(),
				liquidityTier,
				openInterestQuoteQuantums.String(),
			),
		)
	}

	// First, calculate base IMF in big.Rat
	baseImf := new(uint256.Int).Div(
		uint256.NewInt(uint64(liquidityTier.InitialMarginPpm)), // safe, since `InitialMargin` is uint32
		uint256.NewInt(uint64(lib.OneMillion)),
	)

	// `Effective IMF = Min(Base IMF + Max(IMF Increase, 0), 1.0)`
	effectiveImf := baseImf.Add(
		baseImf, // reuse pointer for memory efficiency
		imfIncrease,
	)

	// `Effective IMR = Effective IMF * Quote Quantums`
	effectiveImf.Mul(
		effectiveImf, // reuse pointer for memory efficiency
		quoteQuantums,
	)

	// Return min(Effective IMR, Quote Quantums)
	if effectiveImf.Cmp(quoteQuantums) >= 0 {
		return quoteQuantums
	}
	return effectiveImf
}
