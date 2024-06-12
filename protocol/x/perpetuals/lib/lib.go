package lib

import (
	"math/big"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
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
	indexDelta := new(big.Int).Sub(perpetual.FundingIndex.BigInt(), index)

	// if indexDelta is zero, then net settlement is zero.
	if indexDelta.Sign() == 0 {
		return big.NewInt(0), perpetual.FundingIndex.BigInt()
	}

	bigNetSettlementPpm = new(big.Int).Mul(indexDelta, quantums)

	// `bigNetSettlementPpm` carries sign. `indexDelta` is the increase in `fundingIndex`, so if
	// the position is long (positive), the net settlement should be short (negative), and vice versa.
	// Thus, always negate `bigNetSettlementPpm` here.
	bigNetSettlementPpm = bigNetSettlementPpm.Neg(bigNetSettlementPpm)

	return bigNetSettlementPpm, perpetual.FundingIndex.BigInt()
}

// GetNetCollateralAndMarginRequirements returns the net collateral, initial margin requirement,
// and maintenance margin requirement in quote quantums, given the position size in base quantums.
func GetNetCollateralAndMarginRequirements(
	perpetual types.Perpetual,
	marketPrice pricestypes.MarketPrice,
	liquidityTier types.LiquidityTier,
	quantums *big.Int,
) (
	nc *big.Int,
	imr *big.Int,
	mmr *big.Int,
) {
	nc = GetNetNotionalInQuoteQuantums(
		perpetual,
		marketPrice,
		quantums,
	)
	imr, mmr = GetMarginRequirementsInQuoteQuantums(
		perpetual,
		marketPrice,
		liquidityTier,
		quantums,
	)
	return nc, imr, mmr
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
func GetMarginRequirementsInQuoteQuantums(
	perpetual types.Perpetual,
	marketPrice pricestypes.MarketPrice,
	liquidityTier types.LiquidityTier,
	bigQuantums *big.Int,
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
	)
	return bigInitialMarginQuoteQuantums, bigMaintenanceMarginQuoteQuantums
}
