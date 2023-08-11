package keeper

import (
	"fmt"
	"sort"
	"time"

	gometrics "github.com/armon/go-metrics"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	pricefeedmetrics "github.com/dydxprotocol/v4/daemons/pricefeed/metrics"
	"github.com/dydxprotocol/v4/lib/metrics"
	"github.com/dydxprotocol/v4/x/prices/types"
)

// GetValidMarketPriceUpdates returns a msg containing a list of "valid" price updates that should
// be included in a block. A "valid" price update means:
// 1) The index price must exist.
// 2) The index price must not be zero.
// 3) The index price must be greater than the min price change for the market.
//
// Note: the list of market price updates can be empty if there are no "valid" index prices.
func (k Keeper) GetValidMarketPriceUpdates(
	ctx sdk.Context,
	proposer sdk.AccAddress,
) *types.MsgUpdateMarketPrices {
	defer telemetry.ModuleMeasureSince(
		types.ModuleName,
		time.Now(),
		metrics.GetValidMarketPriceUpdates,
		metrics.Latency,
	)

	// 1. Get all markets from state.
	allMarkets := k.GetAllMarkets(ctx)

	// 2. Get all index prices from in-memory cache.
	allIndexPrices := k.indexPriceCache.GetValidMedianPrices(allMarkets, k.timeProvider.Now())

	// 3. Collect all "valid" price updates.
	updates := make([]*types.MsgUpdateMarketPrices_MarketPrice, 0, len(allMarkets))
	for _, market := range allMarkets {
		indexPrice, exists := allIndexPrices[market.Id]

		if !exists {
			telemetry.IncrCounterWithLabels(
				[]string{types.ModuleName, metrics.IndexPriceDoesNotExist, metrics.Count},
				1,
				[]gometrics.Label{ // To track per market, include the id as a label.
					pricefeedmetrics.GetLabelForMarketId(market.Id),
				},
			)
			ctx.Logger().Error(fmt.Sprintf("Index price for market (%v) does not exist", market.Id))
			continue
		}

		if indexPrice == 0 {
			telemetry.IncrCounterWithLabels(
				[]string{types.ModuleName, metrics.IndexPriceIsZero, metrics.Count},
				1,
				[]gometrics.Label{ // To track per market, include the id as a label.
					pricefeedmetrics.GetLabelForMarketId(market.Id),
				},
			)
			ctx.Logger().Error(fmt.Sprintf("Index price for market (%v) is zero", market.Id))
			continue
		}

		// Check if the index price is above min change required.
		if !isAboveRequiredMinPriceChange(market, indexPrice) {
			telemetry.IncrCounterWithLabels(
				[]string{types.ModuleName, metrics.IndexPriceDoesNotMeetMinPriceChange, metrics.Count},
				1,
				[]gometrics.Label{ // To track per market, include the id as a label.
					pricefeedmetrics.GetLabelForMarketId(market.Id),
				},
			)
			ctx.Logger().Info(fmt.Sprintf(
				"Index price (%v) for market (%v) does not meet min price change ppm (%v)",
				indexPrice,
				market.Id,
				market.MinPriceChangePpm,
			))
			continue
		}

		// Note: there's no need to check for price "accuracy" because the index price is used
		// for the price update.

		// Add as a "valid" price update.
		updates = append(
			updates,
			&types.MsgUpdateMarketPrices_MarketPrice{
				MarketId: market.Id,
				Price:    indexPrice,
			},
		)
	}

	// 4. Sort price updates by market id in ascending order.
	sort.Slice(updates, func(i, j int) bool { return updates[i].MarketId < updates[j].MarketId })

	return &types.MsgUpdateMarketPrices{
		Proposer:           proposer.String(),
		MarketPriceUpdates: updates,
	}
}
