package keeper

import (
	"errors"
	"fmt"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// UpdateSmoothedPrices updates the internal map of smoothed prices for all markets.
// The smoothing is calculated with Basic Exponential Smoothing, see
// https://en.wikipedia.org/wiki/Exponential_smoothing
// If there is no valid daemon price for a market at this time, the smoothed price does not change.
func (k Keeper) UpdateSmoothedSpotPrices(
	ctx sdk.Context,
	linearInterpolateFunc func(v0 uint64, v1 uint64, ppm uint32) (uint64, error),
) error {
	allMarketParams := k.GetAllMarketParams(ctx)
	daemonPrices := k.DaemonPriceCache.GetValidMedianPrices(allMarketParams, k.timeProvider.Now())

	// Track errors for each market.
	updateErrors := make([]error, 0)

	// Iterate through allMarketParams instead of daemonPrices to ensure that we generate deterministic error messages
	// in the case of failed updates.
	for _, marketParam := range allMarketParams {
		daemonPrice, exists := daemonPrices[marketParam.Id]
		if !exists {
			continue
		}

		smoothedPrice, ok := k.marketToSmoothedPrices.GetSmoothedSpotPrice(marketParam.Id)
		if !ok {
			smoothedPrice = daemonPrice
		}
		update, err := linearInterpolateFunc(
			smoothedPrice,
			daemonPrice,
			types.PriceSmoothingPpm,
		)
		if err != nil {
			updateErrors = append(
				updateErrors,
				fmt.Errorf("Error updating smoothed price for market %v: %w", marketParam.Id, err),
			)
			continue
		}

		k.marketToSmoothedPrices.PushSmoothedSpotPrice(marketParam.Id, update)
	}

	return errors.Join(updateErrors...)
}

func (k Keeper) GetSmoothedSpotPrice(
	markedId uint32,
) (uint64, bool) {
	return k.marketToSmoothedPrices.GetSmoothedSpotPrice(markedId)
}
