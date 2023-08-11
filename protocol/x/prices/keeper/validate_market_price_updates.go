package keeper

import (
	"fmt"
	gometrics "github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	pricefeedmetrics "github.com/dydxprotocol/v4/daemons/pricefeed/metrics"
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/lib/metrics"
	"github.com/dydxprotocol/v4/x/prices/types"
	"math/big"
	"time"
)

const (
	CrossingPriceUpdateCutoffPpm = uint32(500_000) // 50%
)

// PerformStatefulPriceUpdateValidation performs stateful validations on `MsgUpdateMarketPrices`.
// Depending on the input, this func performs non-deterministic stateful validation.
func (k Keeper) PerformStatefulPriceUpdateValidation(
	ctx sdk.Context,
	marketPriceUpdates *types.MsgUpdateMarketPrices,
	performNonDeterministicValidation bool,
) error {
	var determinismMetricKeyValue string
	if performNonDeterministicValidation {
		determinismMetricKeyValue = metrics.NonDeterministic
	} else {
		determinismMetricKeyValue = metrics.Deterministic
	}
	defer telemetry.ModuleMeasureSince(
		types.ModuleName,
		time.Now(),
		metrics.StatefulPriceUpdateValidation,
		determinismMetricKeyValue,
		metrics.Latency,
	)

	markets := k.GetAllMarkets(ctx)
	if err := k.performDeterministicStatefulValidation(ctx, marketPriceUpdates, markets); err != nil {
		telemetry.IncrCounter(
			1,
			types.ModuleName,
			metrics.StatefulPriceUpdateValidation,
			metrics.Deterministic,
			metrics.Error,
		)
		return err
	}

	if performNonDeterministicValidation {
		err := k.performNonDeterministicStatefulValidation(ctx, marketPriceUpdates, markets)
		if err != nil {
			telemetry.IncrCounter(
				1,
				types.ModuleName,
				metrics.StatefulPriceUpdateValidation,
				metrics.NonDeterministic,
				metrics.Error,
			)
			return err
		}
	}

	return nil
}

// performNonDeterministicStatefulValidation performs stateful validations that are non-deterministic.
//
// Specificically, for each price update, validate the following:
//   - The index price exists.
//   - The price is "accurate". See `validatePriceAccuracy` for how "accuracy" is determined.
//
// Note: this is NOT determistic, because it relies on "index price" that is subject to each validator.
func (k Keeper) performNonDeterministicStatefulValidation(
	ctx sdk.Context,
	marketPriceUpdates *types.MsgUpdateMarketPrices,
	allMarkets []types.Market,
) error {
	idToMarket := getIdToMarket(allMarkets)
	idToIndexPrice := k.indexPriceCache.GetValidMedianPrices(allMarkets, k.timeProvider.Now())

	for _, priceUpdate := range marketPriceUpdates.GetMarketPriceUpdates() {
		// Check market exists.
		market, err := getMarket(priceUpdate.MarketId, idToMarket)
		if err != nil {
			return err
		}

		// Check index price exists.
		indexPrice, indexPriceExists := idToIndexPrice[priceUpdate.MarketId]
		if !indexPriceExists {
			// Index price not available, so the update price accuracy cannot be determined.
			telemetry.IncrCounterWithLabels(
				[]string{types.ModuleName, metrics.IndexPriceNotAvailForAccuracyCheck, metrics.Count},
				1,
				[]gometrics.Label{ // To track per market, include the id as a label.
					pricefeedmetrics.GetLabelForMarketId(market.Id),
				},
			)
			return sdkerrors.Wrapf(
				types.ErrIndexPriceNotAvailable,
				"index price for market (%d) is not available",
				priceUpdate.MarketId,
			)
		}

		// Check price is "accurate".
		if err := k.validatePriceAccuracy(market, priceUpdate, indexPrice); err != nil {
			telemetry.IncrCounterWithLabels(
				[]string{types.ModuleName, metrics.IndexPriceNotAccurate, metrics.Count},
				1,
				[]gometrics.Label{ // To track per market, include the id as a label.
					pricefeedmetrics.GetLabelForMarketId(market.Id),
				},
			)
			return err
		}
	}

	// Report missing markets in the price updates.
	missingMarketIds := k.GetMarketsMissingFromPriceUpdates(ctx, marketPriceUpdates.MarketPriceUpdates)
	if len(missingMarketIds) > 0 {
		telemetry.SetGauge(
			float32(len(missingMarketIds)),
			types.ModuleName,
			metrics.MissingPriceUpdates,
			metrics.Count,
		)
		ctx.Logger().Info(fmt.Sprintf("markets were not included in the price updates: %+v", missingMarketIds))
	}

	return nil
}

// performDeterministicStatefulValidation performs stateful validations that are deterministic.
//
// Specificically, for each price update, validate the following:
//   - The market exists.
//   - The price update is greater than the min price change.
func (k Keeper) performDeterministicStatefulValidation(
	ctx sdk.Context,
	marketPriceUpdates *types.MsgUpdateMarketPrices,
	allMarkets []types.Market,
) error {
	idToMarket := getIdToMarket(allMarkets)

	for _, priceUpdate := range marketPriceUpdates.GetMarketPriceUpdates() {
		// Check market exists.
		market, err := getMarket(priceUpdate.MarketId, idToMarket)
		if err != nil {
			return err
		}

		// Check price respects min price change.
		if !isAboveRequiredMinPriceChange(market, priceUpdate.Price) {
			return sdkerrors.Wrapf(
				types.ErrInvalidMarketPriceUpdateDeterministic,
				"update price (%d) for market (%d) does not meet min price change requirement"+
					" (%d ppm) based on the current market price (%d)",
				priceUpdate.Price,
				priceUpdate.MarketId,
				market.MinPriceChangePpm,
				market.Price,
			)
		}
	}
	return nil
}

// validatePriceAccuracy checks if the price update is "accurate".
//
// "Accurate" means either one of the following conditions must be met:
//
//   - Towards condition: the price update is between the index price and the current price, inclusive
//
//   - Crossing condition: if the proposed price crosses the index price, then the absolute difference
//     in ticks between the index price and the current price, in ticks (old_ticks), must be:
//     -- old_ticks > 1: greater than or equal to the square of the absolute difference between this node's
//     index price and the proposed price update, in ticks (new_ticks),
//     -- old_ticks <= 1: greater than or equal to the absolute difference between this node's index price
//     and the proposed price update
//
//     Note that ticks are defined as the minimum price change of the currency at the current price
func (k Keeper) validatePriceAccuracy(
	currMarket types.Market,
	priceUpdate *types.MsgUpdateMarketPrices_MarketPrice,
	indexPrice uint64,
) error {
	if isTowardsIndexPrice(currMarket.Price, priceUpdate.Price, indexPrice) {
		return nil
	}

	if !isCrossingIndexPrice(currMarket.Price, priceUpdate.Price, indexPrice) {
		return sdkerrors.Wrapf(
			types.ErrInvalidMarketPriceUpdateNonDeterministic,
			"update price (%d) for market (%d) trends in the opposite direction of the index price (%d) compared "+
				"to the current price (%d)",
			priceUpdate.Price,
			priceUpdate.MarketId,
			indexPrice,
			currMarket.Price,
		)
	}

	tickSizePpm := computeTickSizePpm(currMarket.Price, currMarket.MinPriceChangePpm)

	oldDelta := new(big.Int).SetUint64(lib.AbsDiffUint64(currMarket.Price, indexPrice))
	newDelta := new(big.Int).SetUint64(lib.AbsDiffUint64(indexPrice, priceUpdate.Price))

	// If the index price is <= 1 tick from the old price, we want to compare absolute values of old_delta
	// and new_delta to determine if the price change is valid.
	if priceDeltaIsWithinOneTick(oldDelta, tickSizePpm) {
		if newDelta.Cmp(oldDelta) > 0 {
			return sdkerrors.Wrapf(
				types.ErrInvalidMarketPriceUpdateNonDeterministic,
				"update price (%d) for market (%d) crosses the index price (%d) with current price (%d) "+
					"and deviates from index price (%d) more than minimum allowed (%d)",
				priceUpdate.Price,
				priceUpdate.MarketId,
				indexPrice,
				currMarket.Price,
				newDelta.Uint64(),
				oldDelta.Uint64(),
			)
		}
		return nil
	}

	// Update price crosses index price and old_ticks > 1: check new_ticks <= sqrt(old_ticks)
	if !newPriceMeetsSqrtCondition(oldDelta, newDelta, tickSizePpm) {
		return sdkerrors.Wrapf(
			types.ErrInvalidMarketPriceUpdateNonDeterministic,
			"update price (%d) for market (%d) crosses the index price (%d) with current price (%d) "+
				"and deviates from index price (%d) more than minimum allowed (%d)",
			priceUpdate.Price,
			priceUpdate.MarketId,
			indexPrice,
			currMarket.Price,
			newDelta.Uint64(),
			maximumAllowedPriceDelta(oldDelta, tickSizePpm),
		)
	}

	return nil
}

// GetMarketsMissingFromPriceUpdates returns a list of market ids that should have been included but
// not present in the `MsgUpdateMarketPrices`.
//
// Note: this is NOT determistic, because it relies on "index price" that is subject to each validator.
func (k Keeper) GetMarketsMissingFromPriceUpdates(
	ctx sdk.Context,
	marketPriceUpdates []*types.MsgUpdateMarketPrices_MarketPrice,
) []uint32 {
	// Gather all markets that are part of the proposed updates.
	proposedUpdatesMap := make(map[uint32]bool, len(marketPriceUpdates))
	for _, proposedUpdate := range marketPriceUpdates {
		proposedUpdatesMap[proposedUpdate.MarketId] = true
	}

	// Gather all markets that we think should be updated.
	var missingMarkets []uint32
	// Note that `GetValidMarketPriceUpdates` return value is ordered by market id.
	// This is NOT deterministic, because the returned values are based on "index price".
	allLocalUpdates := k.GetValidMarketPriceUpdates(ctx, []byte{}).MarketPriceUpdates
	for _, localUpdate := range allLocalUpdates {
		if _, exists := proposedUpdatesMap[localUpdate.MarketId]; !exists {
			missingMarkets = append(missingMarkets, localUpdate.MarketId)
		}
	}

	return missingMarkets
}

// getIdToMarket returns a map of market id to market given a slice of markets.
func getIdToMarket(markets []types.Market) map[uint32]types.Market {
	idToMarket := make(map[uint32]types.Market, len(markets))
	for _, market := range markets {
		idToMarket[market.Id] = market
	}
	return idToMarket
}

// getMarket returns a market given a market id. Returns an error if a market does not exist.
func getMarket(marketId uint32, idToMarket map[uint32]types.Market) (types.Market, error) {
	market, marketExists := idToMarket[marketId]
	if !marketExists {
		return market, sdkerrors.Wrapf(
			types.ErrInvalidMarketPriceUpdateDeterministic,
			"market (%d) does not exist",
			marketId,
		)
	}
	return market, nil
}
