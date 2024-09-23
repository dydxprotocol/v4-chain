package keeper

import (
	"errors"
	"math/big"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
)

// getProposalPrice returns the proposed price update for the next block, which is either the smoothed price or the
// daemon price - whichever is closer to the current market price. In cases where the smoothed price and the daemon price
// are equidistant from the current market price, the smoothed price is chosen.
func getProposalPrice(smoothedPrice uint64, daemonPrice uint64, marketPrice uint64) uint64 {
	if lib.AbsDiffUint64(smoothedPrice, marketPrice) > lib.AbsDiffUint64(daemonPrice, marketPrice) {
		return daemonPrice
	}
	return smoothedPrice
}

// isAboveRequiredMinPriceChange returns true if the new price meets the required min price change
// for the market. Otherwise, returns false.
func isAboveRequiredMinSpotPriceChange(marketParamPrice types.MarketParamPrice, newPrice uint64) bool {
	minSpotChangeAmt := getMinPriceChangeAmountForSpotMarket(marketParamPrice)
	return lib.AbsDiffUint64(marketParamPrice.Price.SpotPrice, newPrice) >= minSpotChangeAmt
}

func isAboveRequiredMinPnlPriceChange(marketParamPrice types.MarketParamPrice, newPrice uint64) bool {
	minPnlChangeAmt := getMinPriceChangeAmountForPnlMarket(marketParamPrice)
	return lib.AbsDiffUint64(marketParamPrice.Price.PnlPrice, newPrice) >= minPnlChangeAmt
}

// getMinPriceChangeAmountForMarket returns the amount of price change that is needed to trigger
// a price update in accordance with the min price change parts-per-million value. This method always rounds down,
// which slightly biases towards price updates.
func getMinPriceChangeAmountForSpotMarket(marketParamPrice types.MarketParamPrice) uint64 {
	bigPrice := new(big.Int).SetUint64(marketParamPrice.Price.SpotPrice)
	// There's no need to multiply this by the market's exponent, because `Price` comparisons are
	// done without the market's exponent.
	bigMinChangeAmt := lib.BigIntMulPpm(bigPrice, marketParamPrice.Param.MinPriceChangePpm)

	if !bigMinChangeAmt.IsUint64() {
		// This means that the min change amount is greater than the max uint64. This can only
		// happen if the `MinPriceChangePpm` > 1,000,000 and there's a validation when
		// creating/modifying the `Market`.
		panic(errors.New("getMinPriceChangeAmountForMarket: min price change amount is greater than max uint64 value"))
	}

	return bigMinChangeAmt.Uint64()
}

func getMinPriceChangeAmountForPnlMarket(marketParamPrice types.MarketParamPrice) uint64 {
	bigPrice := new(big.Int).SetUint64(marketParamPrice.Price.PnlPrice)
	// There's no need to multiply this by the market's exponent, because `Price` comparisons are
	// done without the market's exponent.
	bigMinChangeAmt := lib.BigIntMulPpm(bigPrice, marketParamPrice.Param.MinPriceChangePpm)

	if !bigMinChangeAmt.IsUint64() {
		// This means that the min change amount is greater than the max uint64. This can only
		// happen if the `MinPriceChangePpm` > 1,000,000 and there's a validation when
		// creating/modifying the `Market`.
		panic(errors.New("getMinPriceChangeAmountForMarket: min price change amount is greater than max uint64 value"))
	}

	return bigMinChangeAmt.Uint64()
}

// PriceTuple labels and encapsulates the set of prices used for various price computations.
type PriceTuple struct {
	OldPrice    uint64
	DaemonPrice uint64
	NewPrice    uint64
}

// isTowardsDaemonPrice returns true if the new price is between the current price and the daemon
// price, inclusive. Otherwise, it returns false.
func isTowardsDaemonPrice(
	priceTuple PriceTuple,
) bool {
	return priceTuple.NewPrice <= lib.Max(priceTuple.OldPrice, priceTuple.DaemonPrice) &&
		priceTuple.NewPrice >= lib.Min(priceTuple.OldPrice, priceTuple.DaemonPrice)
}

// isCrossingDaemonPrice returns true if daemon price is between the current and the new price,
// noninclusive. Otherwise, returns false.
func isCrossingDaemonPrice(
	priceTuple PriceTuple,
) bool {
	return isCrossingReferencePrice(priceTuple.OldPrice, priceTuple.DaemonPrice, priceTuple.NewPrice)
}

// isCrossingOldPrice returns true if the old price is between the daemon price and the new
// price, noninclusive. Otherwise, returns false.
func isCrossingOldPrice(
	priceTuple PriceTuple,
) bool {
	return isCrossingReferencePrice(priceTuple.DaemonPrice, priceTuple.OldPrice, priceTuple.NewPrice)
}

// isCrossingReferencePrice returns true if the reference price is between the base price and the
// test price, noninclusive. Otherwise, returns false.
func isCrossingReferencePrice(
	basePrice uint64,
	referencePrice uint64,
	testPrice uint64,
) bool {
	return referencePrice < lib.Max(basePrice, testPrice) && referencePrice > lib.Min(basePrice, testPrice)
}

// computeTickSizePpm calculates the tick_size of the currency at the current price, in ppm.
// We keep the tick_size multiplied by 10^6 to reduce divisions in our calculations and avoid rounding errors.
func computeTickSizePpm(oldPrice uint64, minPriceChangePpm uint32) *big.Int {
	// tick_size = oldPrice * minPriceChangePpm / 1_000_000 ==>
	// tick_size_ppm = oldPrice * minPriceChangePpm = tick_size * 1_000_000
	return new(big.Int).Mul(
		new(big.Int).SetUint64(oldPrice),
		new(big.Int).SetUint64(uint64(minPriceChangePpm)))
}

// priceDeltaIsWithinOneTick returns true iff the price delta is within one tick, given the tick_size in ppm.
func priceDeltaIsWithinOneTick(priceDelta *big.Int, tickSizePpm *big.Int) bool {
	// To compare if a price_delta > tick_size, let's multiply by 1_000_000 and compare against the
	// tick size in ppm
	priceDeltaPpm := new(big.Int).Mul(priceDelta, new(big.Int).SetUint64(constants.OneMillion))
	return priceDeltaPpm.Cmp(tickSizePpm) <= 0
}

// newPriceMeetsSqrtCondition calculates the price acceptance condition when the new price crosses the daemon
// price and the price change from the current price to the daemon price, or old_ticks, is > 1 tick.
//
// Ticks are computed at the currency's current price.
//
// Under these conditions, price changes are valid when new_ticks <= sqrt(old_ticks)
func newPriceMeetsSqrtCondition(oldDelta *big.Int, newDelta *big.Int, tickSizePpm *big.Int) bool {
	// In order to avoid division / sqrt, which is potentially lossy, use big.Ints and refactor:
	// given that new_ticks = new_delta / tick_size, old_ticks = old_delta / tick_size
	// new_ticks < sqrt(old_ticks)                                  ==> sub in old_ticks, new_ticks
	// new_delta / tick_size <= sqrt(old_delta / tick_size)         ==>
	// new_delta * new_delta / tick_size <= old_delta               ==>
	// new_delta * new_delta <= old_delta * tick_size               ==>
	// new_delta * new_delta * 1_000_000 <= old_delta * tickSizePpm
	newDeltaSquaredPpm := new(big.Int).Mul(newDelta, newDelta)
	newDeltaSquaredPpm.Mul(newDeltaSquaredPpm, new(big.Int).SetUint64(constants.OneMillion))
	oldDeltaTimesTickSizePpm := new(big.Int).Mul(oldDelta, tickSizePpm)
	return newDeltaSquaredPpm.Cmp(oldDeltaTimesTickSizePpm) <= 0
}

// maximumAllowedPriceDelta computes the maximum allowable value of new_delta under the conditions
// that the proposed price is crossing in the daemon price, and old_ticks > 1. This method uses potentially
// lossy arithmetic and is only for logging purposes.
func maximumAllowedPriceDelta(oldDelta *big.Int, tickSizePpm *big.Int) *big.Int {
	// Compute maximum allowable new_delta, or price difference between the daemon price
	// and the proposed price:
	// max_allowed = sqrt(old_delta * tick_size_ppm / 1_000_000)
	maxAllowed := new(big.Int).Mul(oldDelta, tickSizePpm)
	maxAllowed.Div(maxAllowed, new(big.Int).SetUint64(constants.OneMillion))
	maxAllowed.Sqrt(maxAllowed)
	return maxAllowed
}
