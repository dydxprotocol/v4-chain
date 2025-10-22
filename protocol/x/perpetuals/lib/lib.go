package lib

import (
	"math/big"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/margin"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

// GetSettlementPpm returns the net settlement amount ppm (in quote quantums) given
// the perpetual and position size (in base quantums).
func GetSettlementPpmWithPerpetual(
	perpetual types.Perpetual,
	quantums *big.Int,
	index *big.Int,
) (
	bigNetSettlementPpm *big.Int,
	newFundingIndex *big.Int,
) {
	fundingIndex := perpetual.FundingIndex.BigInt()

	// If no change in funding, return 0.
	if fundingIndex.Cmp(index) == 0 {
		return big.NewInt(0), fundingIndex
	}

	// The settlement is a signed value.
	// If the index delta is positive and the quantums is positive (long), then settlement is negative.
	// Thus, always negate the value of the multiplication of the index delta and the quantums.
	result := new(big.Int).Sub(fundingIndex, index)
	result = result.Mul(result, quantums)
	result = result.Neg(result)

	return result, fundingIndex
}

// GetPositionNetNotionalValueAndMarginRequirements returns the net collateral, initial margin requirement,
// and maintenance margin requirement in quote quantums, given the position size in base quantums.
func GetPositionNetNotionalValueAndMarginRequirements(
	perpetual types.Perpetual,
	marketPrice pricestypes.MarketPrice,
	liquidityTier types.LiquidityTier,
	quantums *big.Int,
	custom_imf_ppm uint32,
) (
	risk margin.Risk,
) {
	nc := GetNetNotionalInQuoteQuantums(
		perpetual,
		marketPrice,
		quantums,
	)
	imr, mmr := GetMarginRequirementsInQuoteQuantums(
		perpetual,
		marketPrice,
		liquidityTier,
		quantums,
		custom_imf_ppm,
	)
	return margin.Risk{
		NC:  nc,
		IMR: imr,
		MMR: mmr,
	}
}

// GetNetCollateralAndMarginRequirements returns the net collateral, initial margin requirement,
// and maintenance margin requirement in quote quantums, given the position size in base quantums.
func GetNetCollateralAndMarginRequirements(
	perpetual types.Perpetual,
	marketPrice pricestypes.MarketPrice,
	liquidityTier types.LiquidityTier,
	quantums *big.Int,
	quoteBalance *big.Int,
	custom_imf_ppm uint32, // 0 means use default liquidity tier margins
) (
	risk margin.Risk,
) {
	risk = GetPositionNetNotionalValueAndMarginRequirements(
		perpetual,
		marketPrice,
		liquidityTier,
		quantums,
		custom_imf_ppm,
	)
	risk.NC.Add(risk.NC, quoteBalance)
	return risk
}

// GetNetNotionalInQuoteQuantums returns the net notional in quote quantums, which can be
// represented by the following equation:
//
// `quantums / 10^baseAtomicResolution * marketPrice * 10^marketExponent * 10^quoteAtomicResolution`.
// Note that longs are positive, and shorts are negative.
func GetNetNotionalInQuoteQuantums(
	perpetual types.Perpetual,
	marketPrice pricestypes.MarketPrice,
	bigQuantums *big.Int,
) (
	bigNetNotionalQuoteQuantums *big.Int,
) {
	bigQuoteQuantums := lib.BaseToQuoteQuantums(
		bigQuantums,
		perpetual.Params.AtomicResolution,
		marketPrice.Price,
		marketPrice.Exponent,
	)

	return bigQuoteQuantums
}

// GetMarginRequirementsInQuoteQuantums returns initial and maintenance margin requirements
// in quote quantums, given the position size in base quantums.
// If leverage > 0, scales the margin requirements based on maxLeverage/userLeverage ratio.
func GetMarginRequirementsInQuoteQuantums(
	perpetual types.Perpetual,
	marketPrice pricestypes.MarketPrice,
	liquidityTier types.LiquidityTier,
	bigQuantums *big.Int,
	custom_imf_ppm uint32, // 0 means use default liquidity tier margins
) (
	bigInitialMarginQuoteQuantums *big.Int,
	bigMaintenanceMarginQuoteQuantums *big.Int,
) {
	// Always consider the magnitude of the position regardless of whether it is long/short.
	bigAbsQuantums := new(big.Int).Abs(bigQuantums)

	// Calculate the notional value of the position in quote quantums.
	bigQuoteQuantums := lib.BaseToQuoteQuantums(
		bigAbsQuantums,
		perpetual.Params.AtomicResolution,
		marketPrice.Price,
		marketPrice.Exponent,
	)
	// Calculate the perpetual's open interest in quote quantums.
	openInterestQuoteQuantums := lib.BaseToQuoteQuantums(
		perpetual.OpenInterest.BigInt(), // OpenInterest is represented as base quantums.
		perpetual.Params.AtomicResolution,
		marketPrice.Price,
		marketPrice.Exponent,
	)

	// Initial margin requirement quote quantums = size in quote quantums * initial margin PPM.
	bigBaseInitialMarginQuoteQuantums := liquidityTier.GetInitialMarginQuoteQuantums(
		bigQuoteQuantums,
		big.NewInt(0), // pass in 0 as open interest to get base IMR.
		big.NewInt(0), // pass in 0 to use the base IMR
	)
	// Maintenance margin requirement quote quantums = IM in quote quantums * maintenance fraction PPM.
	bigMaintenanceMarginQuoteQuantums = lib.BigMulPpm(
		bigBaseInitialMarginQuoteQuantums,
		lib.BigU(liquidityTier.MaintenanceFractionPpm),
		true,
	)

	bigInitialMarginQuoteQuantums = liquidityTier.GetInitialMarginQuoteQuantums(
		bigQuoteQuantums,
		openInterestQuoteQuantums, // pass in current OI to get scaled IMR.
		lib.BigU(custom_imf_ppm),
	)

	return bigInitialMarginQuoteQuantums, bigMaintenanceMarginQuoteQuantums
}
