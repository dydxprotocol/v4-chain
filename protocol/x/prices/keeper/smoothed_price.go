package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

// UpdateSmoothedPrices updates the internal map of smoothed prices for all markets.
// The smoothing is calculated with Basic Exponential Smoothing, see
// https://en.wikipedia.org/wiki/Exponential_smoothing
// If there is no valid index price for a market at this time, the smoothed price does not change.
func (k Keeper) UpdateSmoothedPrices(ctx sdk.Context) error {
	allMarketParams := k.GetAllMarketParams(ctx)
	indexPrices := k.indexPriceCache.GetValidMedianPrices(allMarketParams, k.timeProvider.Now())

	for market, indexPrice := range indexPrices {
		smoothedPrice, ok := k.marketToSmoothedPrices.GetSmoothedPrice(market)
		if !ok {
			smoothedPrice = indexPrice
		}
		update, err := lib.Uint64LinearInterpolate(
			smoothedPrice,
			indexPrice,
			types.PriceSmoothingPpm,
		)
		if err != nil {
			return err
		}
		k.marketToSmoothedPrices.PushSmoothedPrice(market, update)
	}
	return nil
}
