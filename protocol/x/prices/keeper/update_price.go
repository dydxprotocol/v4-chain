package keeper

import (
	"fmt"
	gometrics "github.com/armon/go-metrics"
	"sort"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	pricefeedmetrics "github.com/dydxprotocol/v4/daemons/pricefeed/metrics"
	"github.com/dydxprotocol/v4/lib/metrics"
	"github.com/dydxprotocol/v4/x/prices/types"
)

// GetValidMarketPriceUpdates returns a msg containing a list of "valid" price updates that should
// be included in a block. A "valid" price update means:
// 1) All values used to compute the valid price must exist and be valid:
// - the index price exists and is nonzero.
// - the smoothed price exists and is nonzero.
// 2) The smoothed price and the index price must be on the same side compared to the oracle price.
// 3) The proposed price is either the index price or the smoothed price, depending on which is closer to the
// oracle price.
// 4) The proposed price meets the minimum price change ppm requirement.
// Note: the list of market price updates can be empty if there are no "valid" index prices, smoothed prices, and/or
// proposed prices for any market.
func (k Keeper) GetValidMarketPriceUpdates(
	ctx sdk.Context,
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
		indexPrice, indexPriceExists := allIndexPrices[market.Id]
		smoothedPrice, smoothedPriceExists := k.marketToSmoothedPrices[market.Id]

		marketMetricsLabel := pricefeedmetrics.GetLabelForMarketId(market.Id)

		// Skip proposal logic in the event of invalid inputs, which is only likely to occur around network genesis.
		if !indexPriceExists {
			metrics.IncrCountMetricWithLabels(types.ModuleName, metrics.IndexPriceDoesNotExist, marketMetricsLabel)
			ctx.Logger().Error(fmt.Sprintf("Index price for market (%v) does not exist", market.Id))
			continue
		}

		if indexPrice == 0 {
			metrics.IncrCountMetricWithLabels(types.ModuleName, metrics.IndexPriceIsZero, marketMetricsLabel)
			ctx.Logger().Error(fmt.Sprintf("Index price for market (%v) is zero", market.Id))
			continue
		}

		if !smoothedPriceExists || smoothedPrice == 0 {
			if !smoothedPriceExists {
				ctx.Logger().Error(fmt.Sprintf("Smoothed price for market (%v) does not exist", market.Id))
			} else { // smoothedPrice == 0
				ctx.Logger().Error(fmt.Sprintf("Smoothed price for market (%v) is zero", market.Id))
			}
			// This could happen before the first block is proposed, because we update smoothed prices in the
			// propose handler. Otherwise, we would never expect to need to fall back to the index price again.
			smoothedPrice = indexPrice
		}

		proposalPrice := getProposalPrice(smoothedPrice, indexPrice, market.Price)

		// If the index price would have updated, track how the proposal price changes the update
		// decision / amount.
		if isAboveRequiredMinPriceChange(market, indexPrice) {
			logPriceUpdateBehavior(ctx, market, proposalPrice, smoothedPrice, indexPrice, marketMetricsLabel)
		}

		if isAboveRequiredMinPriceChange(market, proposalPrice) &&
			!isCrossingOldPrice(PriceTuple{
				OldPrice:   market.Price,
				IndexPrice: indexPrice,
				NewPrice:   smoothedPrice,
			}) {
			// Add as a "valid" price update.
			updates = append(
				updates,
				&types.MsgUpdateMarketPrices_MarketPrice{
					MarketId: market.Id,
					Price:    proposalPrice,
				},
			)
		}
	}

	// 4. Sort price updates by market id in ascending order.
	sort.Slice(updates, func(i, j int) bool { return updates[i].MarketId < updates[j].MarketId })

	return &types.MsgUpdateMarketPrices{
		MarketPriceUpdates: updates,
	}
}

func logPriceUpdateBehavior(
	ctx sdk.Context,
	market types.Market,
	proposalPrice uint64,
	smoothedPrice uint64,
	indexPrice uint64,
	marketMetricsLabel gometrics.Label,
) {
	loggingVerb := "proposed"
	proposalPriceIsAboveMinChange := isAboveRequiredMinPriceChange(market, proposalPrice)
	proposalPriceCrossesOldPrice := isCrossingOldPrice(PriceTuple{
		OldPrice:   market.Price,
		IndexPrice: indexPrice,
		NewPrice:   smoothedPrice,
	})

	if proposalPriceIsAboveMinChange && !proposalPriceCrossesOldPrice {
		loggingVerb = "not proposed"

		minPriceChangeLabel := metrics.NewBinaryStringLabel(
			metrics.DoesNotMeetMinPriceChange,
			!proposalPriceIsAboveMinChange,
		)
		trendsAwayFromIndexPriceLabel := metrics.NewBinaryStringLabel(
			metrics.ProposedPriceCrossesOldPrice,
			proposalPriceCrossesOldPrice,
		)
		metrics.IncrCountMetricWithLabels(
			types.ModuleName,
			metrics.ProposedPriceChangesPriceUpdateDecision,
			marketMetricsLabel,
			minPriceChangeLabel,
			trendsAwayFromIndexPriceLabel,
		)
	}

	ctx.Logger().Info(fmt.Sprintf(
		"Proposal price (%v) %v for market (%v), index price (%v), oracle price (%v), min price change (%v)",
		proposalPrice,
		loggingVerb,
		market.Id,
		indexPrice,
		market.Price,
		getMinPriceChangeAmountForMarket(market),
	))
}
