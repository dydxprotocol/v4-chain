package keeper

import (
	"fmt"
	"sort"
	"time"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/pricefeed/client/constants"

	gometrics "github.com/hashicorp/go-metrics"

	pricefeedmetrics "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/pricefeed/metrics"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib/log"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib/metrics"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetValidMarketPriceUpdates returns a msg containing a list of "valid" price updates that should
// be included in a block. A "valid" price update means:
// 1) All values used to compute the valid price must exist and be valid:
// - the daemon price exists and is nonzero.
// - the smoothed price exists and is nonzero.
// 2) The smoothed price and the daemon price must be on the same side compared to the oracle (spot) price.
// 3) The proposed price is either the daemon price or the smoothed price, depending on which is closer to the
// oracle (spot) price.
// 4) The proposed price meets the minimum price change ppm requirement.
// Note: the list of market price updates can be empty if there are no "valid" daemon prices, smoothed prices, and/or
// proposed prices for any market.
func (k Keeper) GetValidMarketSpotPriceUpdates(
	ctx sdk.Context,
) []*types.MarketSpotPriceUpdate {
	defer telemetry.ModuleMeasureSince(
		types.ModuleName,
		time.Now(),
		metrics.GetValidMarketPriceUpdates,
		metrics.Latency,
	)

	// 1. Get all markets from state.
	allMarketParamPrices, err := k.GetAllMarketParamPrices(ctx)
	if err != nil {
		log.ErrorLogWithError(
			ctx,
			"error getting all market param prices",
			err,
		)
		// If there is an error, return an empty list of price updates. We don't want to introduce
		// liveness issues due to an error in market state.
	}

	allMarketParams := make([]types.MarketParam, len(allMarketParamPrices))
	for i, marketParamPrice := range allMarketParamPrices {
		allMarketParams[i] = marketParamPrice.Param
	}

	// 2. Get all daemon prices from in-memory cache.
	allDaemonPrices := k.daemonPriceCache.GetValidMedianPrices(allMarketParams, k.timeProvider.Now())

	// 3. Collect all "valid" price updates.
	updates := make([]*types.MarketSpotPriceUpdate, 0, len(allMarketParamPrices))
	for _, marketParamPrice := range allMarketParamPrices {
		marketId := marketParamPrice.Param.Id
		daemonPrice, daemonPriceExists := allDaemonPrices[marketId]

		marketMetricsLabel := pricefeedmetrics.GetLabelForMarketId(marketId)

		// Skip proposal logic in the event of invalid inputs, which is only likely to occur around network genesis.
		if !daemonPriceExists {
			metrics.IncrCountMetricWithLabels(types.ModuleName, metrics.DaemonPriceDoesNotExist, marketMetricsLabel)
			// Conditionally log missing daemon prices at least 20s after genesis/restart/market creation. We expect that
			// there will be a delay in populating daemon prices after network genesis or a network restart, or when a
			// market is created, it takes the daemon some time to warm up.
			if !k.IsRecentlyAvailable(ctx, marketId) {
				log.ErrorLog(
					ctx,
					"daemon price for market does not exist",
					constants.MarketIdLogKey,
					marketId,
				)
			}
			continue
		}

		// daemon prices of 0 are unexpected. In this scenario, we skip the proposal logic for the market and report an
		// error.
		if daemonPrice == 0 {
			metrics.IncrCountMetricWithLabels(types.ModuleName, metrics.DaemonPriceIsZero, marketMetricsLabel)
			log.ErrorLog(
				ctx,
				"Unexpected error: daemon price for market is zero",
				constants.MarketIdLogKey,
				marketId,
			)
			continue
		}

		historicalSmoothedPrices := k.marketToSmoothedPrices.GetHistoricalSmoothedSpotPrices(marketId)
		// We generally expect to have a smoothed price history for each market, except during the first few blocks
		// after network genesis or a network restart. In this scenario, we use the daemon price as the smoothed price.
		if len(historicalSmoothedPrices) == 0 {
			// Conditionally log missing smoothed prices at least 20s after genesis/restart/market creation. We expect
			// that there will be a delay in populating historical smoothed prices after network genesis or a network
			// restart, or when a market is created, because they depend on present daemon prices, and it takes the
			// daemon some time to warm up.
			if !k.IsRecentlyAvailable(ctx, marketId) {
				log.ErrorLog(
					ctx,
					"Smoothed price for market does not exist",
					constants.MarketIdLogKey,
					marketId,
				)
			}
			historicalSmoothedPrices = []uint64{daemonPrice}
		}
		smoothedPrice := historicalSmoothedPrices[0]

		proposalPrice := getProposalPrice(smoothedPrice, daemonPrice, marketParamPrice.Price.SpotPrice)

		shouldPropose, reasons := shouldProposePrice(
			proposalPrice,
			marketParamPrice,
			daemonPrice,
			historicalSmoothedPrices,
		)

		// If the daemon price would have updated, track how the proposal price changes the update
		// decision / amount.
		if isAboveRequiredMinSpotPriceChange(marketParamPrice, daemonPrice) {
			logPriceUpdateBehavior(
				ctx,
				marketParamPrice,
				proposalPrice,
				daemonPrice,
				marketMetricsLabel,
				shouldPropose,
				reasons,
			)
		}

		if shouldPropose {
			// Add as a "valid" price update.
			updates = append(
				updates,
				&types.MarketSpotPriceUpdate{
					MarketId:  marketId,
					SpotPrice: proposalPrice,
				},
			)
		}
	}

	// 4. Sort price updates by market id in ascending order.
	sort.Slice(updates, func(i, j int) bool { return updates[i].MarketId < updates[j].MarketId })

	return updates
}

func logPriceUpdateBehavior(
	ctx sdk.Context,
	marketParamPrice types.MarketParamPrice,
	proposalPrice uint64,
	daemonPrice uint64,
	marketMetricsLabel gometrics.Label,
	shouldPropose bool,
	reasons []proposeCancellationReason,
) {
	loggingVerb := "proposed"
	if !shouldPropose {
		loggingVerb = "not proposed"

		// Convert reasons map to a slice of metrics labels and include market label.
		labels := make([]gometrics.Label, 0, len(reasons)+1)
		labels = append(labels, marketMetricsLabel)

		for _, reason := range reasons {
			labels = append(labels, metrics.GetLabelForBoolValue(reason.Reason, reason.Value))
		}

		metrics.IncrCountMetricWithLabels(
			types.ModuleName,
			metrics.ProposedPriceChangesPriceUpdateDecision,
			labels...,
		)
	}
	log.InfoLog(
		ctx,
		fmt.Sprintf(
			"Proposal price (%v) %v for market (%v), daemon price (%v), oracle price (%v), min price change (%v)",
			proposalPrice,
			loggingVerb,
			marketParamPrice.Param.Id,
			daemonPrice,
			marketParamPrice.Price.SpotPrice,
			getMinPriceChangeAmountForSpotMarket(marketParamPrice),
		),
	)
}

type proposeCancellationReason struct {
	Reason string
	Value  bool
}

// shouldProposePrice determines if a price should be proposed for a market. It returns the result as well as a list of
// possible reasons why the price should not be proposed that can be used to create metrics labels.
// All of the logic for determining if a price should be proposed was consolidated here to prevent inconsistencies
// between proposal behavior and logging/metrics.
func shouldProposePrice(
	proposalPrice uint64,
	marketParamPrice types.MarketParamPrice,
	daemonPrice uint64,
	historicalSmoothedPrices []uint64,
) (
	shouldPropose bool,
	reasons []proposeCancellationReason,
) {
	reasons = make([]proposeCancellationReason, 0, 4)
	shouldPropose = true

	// If any smoothed price crosses the old price compared to the daemon price, do not update.
	reasons = append(
		reasons,
		proposeCancellationReason{
			Reason: metrics.RecentSmoothedPriceCrossesOraclePrice,
			Value:  false,
		},
	)
	for _, smoothedPrice := range historicalSmoothedPrices {
		if isCrossingOldPrice(PriceTuple{
			OldPrice:    marketParamPrice.Price.SpotPrice,
			DaemonPrice: daemonPrice,
			NewPrice:    smoothedPrice,
		}) {
			shouldPropose = false
			reasons[len(reasons)-1].Value = true
			break
		}
	}

	// If the proposal price crosses the old price compared to the daemon price, do not update.
	reasons = append(
		reasons,
		proposeCancellationReason{
			Reason: metrics.ProposedPriceCrossesOraclePrice,
		},
	)
	if isCrossingOldPrice(PriceTuple{
		OldPrice:    marketParamPrice.Price.SpotPrice,
		DaemonPrice: daemonPrice,
		NewPrice:    proposalPrice,
	}) {
		shouldPropose = false
		reasons[len(reasons)-1].Value = true
	}

	// If any smoothed price does not meet the min price change, do not update.
	reasons = append(
		reasons,
		proposeCancellationReason{
			Reason: metrics.RecentSmoothedPriceDoesNotMeetMinPriceChange,
			Value:  false,
		},
	)
	for _, smoothedPrice := range historicalSmoothedPrices {
		if !isAboveRequiredMinSpotPriceChange(marketParamPrice, smoothedPrice) {
			shouldPropose = false
			reasons[len(reasons)-1].Value = true
			break
		}
	}

	// If the proposal price does not meet the min price change, do not update.
	reasons = append(
		reasons,
		proposeCancellationReason{
			Reason: metrics.ProposedPriceDoesNotMeetMinPriceChange,
		},
	)
	if !isAboveRequiredMinSpotPriceChange(marketParamPrice, proposalPrice) {
		shouldPropose = false
		reasons[len(reasons)-1].Value = true
	}

	return shouldPropose, reasons
}
