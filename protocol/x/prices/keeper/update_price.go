package keeper

import (
	"fmt"
	"sort"
	"time"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/constants"

	gometrics "github.com/hashicorp/go-metrics"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	pricefeedmetrics "github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/lib/log"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
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

	// 2. Get all index prices from in-memory cache.
	allIndexPrices := k.indexPriceCache.GetValidMedianPrices(
		k.Logger(ctx),
		allMarketParams,
		k.timeProvider.Now(),
	)

	// 3. Collect all "valid" price updates.
	updates := make([]*types.MsgUpdateMarketPrices_MarketPrice, 0, len(allMarketParamPrices))
	nonExistentMarkets := []uint32{}
	for _, marketParamPrice := range allMarketParamPrices {
		marketId := marketParamPrice.Param.Id
		indexPrice, indexPriceExists := allIndexPrices[marketId]

		marketMetricsLabel := pricefeedmetrics.GetLabelForMarketId(marketId)

		// Skip proposal logic in the event of invalid inputs, which is only likely to occur around network genesis.
		if !indexPriceExists {
			metrics.IncrCountMetricWithLabels(types.ModuleName, metrics.IndexPriceDoesNotExist, marketMetricsLabel)
			nonExistentMarkets = append(nonExistentMarkets, marketId)
			continue
		}

		// Index prices of 0 are unexpected. In this scenario, we skip the proposal logic for the market and report an
		// error.
		if indexPrice == 0 {
			metrics.IncrCountMetricWithLabels(types.ModuleName, metrics.IndexPriceIsZero, marketMetricsLabel)
			log.ErrorLog(
				ctx,
				"Unexpected error: index price for market is zero",
				constants.MarketIdLogKey,
				marketId,
			)
			continue
		}

		shouldPropose, reasons := shouldProposePrice(
			indexPrice,
			marketParamPrice,
		)

		// If the index price would have updated, track how the proposal price changes the update
		// decision / amount.
		if isAboveRequiredMinPriceChange(marketParamPrice, indexPrice) {
			logPriceUpdateBehavior(
				ctx,
				marketParamPrice,
				indexPrice,
				marketMetricsLabel,
				shouldPropose,
				reasons,
			)
		}

		if shouldPropose {
			// Add as a "valid" price update.
			updates = append(
				updates,
				&types.MsgUpdateMarketPrices_MarketPrice{
					MarketId: marketId,
					Price:    indexPrice,
				},
			)
		}
	}

	if len(nonExistentMarkets) > 0 {
		ctx.Logger().Warn(
			fmt.Sprintf("Index price for markets does not exist, marketIds: %v", nonExistentMarkets),
		)
	}

	// 4. Sort price updates by market id in ascending order.
	sort.Slice(updates, func(i, j int) bool { return updates[i].MarketId < updates[j].MarketId })

	return &types.MsgUpdateMarketPrices{
		MarketPriceUpdates: updates,
	}
}

func logPriceUpdateBehavior(
	ctx sdk.Context,
	marketParamPrice types.MarketParamPrice,
	proposalPrice uint64,
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
			"Proposal price (%v) %v for market (%v), oracle price (%v), min price change (%v)",
			proposalPrice,
			loggingVerb,
			marketParamPrice.Param.Id,
			marketParamPrice.Price.Price,
			getMinPriceChangeAmountForMarket(marketParamPrice),
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
) (
	shouldPropose bool,
	reasons []proposeCancellationReason,
) {
	reasons = make([]proposeCancellationReason, 0, 4)
	shouldPropose = true

	// If the proposal price does not meet the min price change, do not update.
	reasons = append(
		reasons,
		proposeCancellationReason{
			Reason: metrics.ProposedPriceDoesNotMeetMinPriceChange,
		},
	)
	if !isAboveRequiredMinPriceChange(marketParamPrice, proposalPrice) {
		shouldPropose = false
		reasons[len(reasons)-1].Value = true
	}

	return shouldPropose, reasons
}
