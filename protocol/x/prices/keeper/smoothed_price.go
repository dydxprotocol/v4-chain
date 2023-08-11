package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/x/prices/types"
)

// UpdateSmoothedPrices updates the internal map of smoothed prices for all markets.
// The smoothing is calculated with Basic Exponential Smoothing, see
// https://en.wikipedia.org/wiki/Exponential_smoothing
// If there is no valid index price for a market at this time, the smoothed price does not change.
func (k Keeper) UpdateSmoothedPrices(ctx sdk.Context) error {
	allMarkets := k.GetAllMarkets(ctx)
	indexPrices := k.indexPriceCache.GetValidMedianPrices(allMarkets, k.timeProvider.Now())

	for market, indexPrice := range indexPrices {
		smoothed, ok := k.marketToSmoothedPrices[market]
		if !ok {
			smoothed = indexPrice
		}
		update, err := lib.Uint64LinearInterpolate(
			smoothed,
			indexPrice,
			types.PriceSmoothingPpm,
		)
		if err != nil {
			return err
		}
		k.marketToSmoothedPrices[market] = update
	}
	return nil
}
