package keeper

import (
	"fmt"
	"math/big"
	"time"

	errorsmod "cosmossdk.io/errors"
	errorlib "github.com/dydxprotocol/v4-chain/protocol/lib/error"
	"github.com/dydxprotocol/v4-chain/protocol/lib/log"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	pricefeedmetrics "github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	gometrics "github.com/hashicorp/go-metrics"
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

	marketParamPrices, err := k.GetAllMarketParamPrices(ctx)
	if err != nil {
		telemetry.IncrCounter(
			1,
			types.ModuleName,
			metrics.StatefulPriceUpdateValidation,
			metrics.Error,
		)
		return errorlib.WrapErrorWithSourceModuleContext(
			errorsmod.Wrap(err, "failed to get all market param prices"),
			types.ModuleName,
		)
	}

	if err := k.performDeterministicStatefulValidation(ctx, marketPriceUpdates, marketParamPrices); err != nil {
		telemetry.IncrCounter(
			1,
			types.ModuleName,
			metrics.StatefulPriceUpdateValidation,
			metrics.Deterministic,
			metrics.Error,
		)
		return errorlib.WrapErrorWithSourceModuleContext(err, types.ModuleName)
	}

	if performNonDeterministicValidation {
		err := k.performNonDeterministicStatefulValidation(ctx, marketPriceUpdates, marketParamPrices)
		if err != nil {
			telemetry.IncrCounter(
				1,
				types.ModuleName,
				metrics.StatefulPriceUpdateValidation,
				metrics.NonDeterministic,
				metrics.Error,
			)
			return errorlib.WrapErrorWithSourceModuleContext(err, types.ModuleName)
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
	allMarketParamPrices []types.MarketParamPrice,
) error {
	idToMarket := getIdToMarketParamPrice(allMarketParamPrices)
	allMarketParams := make([]types.MarketParam, len(allMarketParamPrices))
	for i, marketParamPrice := range allMarketParamPrices {
		allMarketParams[i] = marketParamPrice.Param
	}

	idToIndexPrice := k.indexPriceCache.GetValidMedianPrices(
		k.Logger(ctx),
		allMarketParams,
		k.timeProvider.Now(),
	)

	for _, priceUpdate := range marketPriceUpdates.GetMarketPriceUpdates() {
		// Check market exists.
		marketParamPrice, err := getMarketParamPrice(priceUpdate.MarketId, idToMarket)
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
					pricefeedmetrics.GetLabelForMarketId(marketParamPrice.Param.Id),
				},
			)
			return errorsmod.Wrapf(
				types.ErrIndexPriceNotAvailable,
				"index price for market (%d) is not available",
				priceUpdate.MarketId,
			)
		}

		// Check price is "accurate".
		if err := k.validatePriceAccuracy(marketParamPrice, priceUpdate, indexPrice); err != nil {
			telemetry.IncrCounterWithLabels(
				[]string{types.ModuleName, metrics.IndexPriceNotAccurate, metrics.Count},
				1,
				[]gometrics.Label{ // To track per market, include the id as a label.
					pricefeedmetrics.GetLabelForMarketId(marketParamPrice.Param.Id),
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
		log.InfoLog(
			ctx,
			fmt.Sprintf("markets were not included in the price updates: %+v", missingMarketIds),
		)
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
	allMarketParamPrices []types.MarketParamPrice,
) error {
	idToMarketParamPrice := getIdToMarketParamPrice(allMarketParamPrices)

	for _, priceUpdate := range marketPriceUpdates.GetMarketPriceUpdates() {
		// Check market exists.
		marketParamPrice, err := getMarketParamPrice(priceUpdate.MarketId, idToMarketParamPrice)
		if err != nil {
			return err
		}

		// Check price respects min price change.
		if !isAboveRequiredMinPriceChange(marketParamPrice, priceUpdate.Price) {
			return errorsmod.Wrapf(
				types.ErrInvalidMarketPriceUpdateDeterministic,
				"update price (%d) for market (%d) does not meet min price change requirement"+
					" (%d ppm) based on the current market price (%d)",
				priceUpdate.Price,
				priceUpdate.MarketId,
				marketParamPrice.Param.MinPriceChangePpm,
				marketParamPrice.Price.Price,
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
	currMarketParamPrice types.MarketParamPrice,
	priceUpdate *types.MsgUpdateMarketPrices_MarketPrice,
	indexPrice uint64,
) error {
	if isTowardsIndexPrice(PriceTuple{
		OldPrice:   currMarketParamPrice.Price.Price,
		IndexPrice: indexPrice,
		NewPrice:   priceUpdate.Price,
	}) {
		return nil
	}

	if !isCrossingIndexPrice(PriceTuple{
		OldPrice:   currMarketParamPrice.Price.Price,
		IndexPrice: indexPrice,
		NewPrice:   priceUpdate.Price,
	}) {
		return errorsmod.Wrapf(
			types.ErrInvalidMarketPriceUpdateNonDeterministic,
			"update price (%d) for market (%d) trends in the opposite direction of the index price (%d) compared "+
				"to the current price (%d)",
			priceUpdate.Price,
			priceUpdate.MarketId,
			indexPrice,
			currMarketParamPrice.Price.Price,
		)
	}

	tickSizePpm := computeTickSizePpm(currMarketParamPrice.Price.Price, currMarketParamPrice.Param.MinPriceChangePpm)

	oldDelta := new(big.Int).SetUint64(lib.AbsDiffUint64(currMarketParamPrice.Price.Price, indexPrice))
	newDelta := new(big.Int).SetUint64(lib.AbsDiffUint64(indexPrice, priceUpdate.Price))

	// If the index price is <= 1 tick from the old price, we want to compare absolute values of old_delta
	// and new_delta to determine if the price change is valid.
	if priceDeltaIsWithinOneTick(oldDelta, tickSizePpm) {
		if newDelta.Cmp(oldDelta) > 0 {
			return errorsmod.Wrapf(
				types.ErrInvalidMarketPriceUpdateNonDeterministic,
				"update price (%d) for market (%d) crosses the index price (%d) with current price (%d) "+
					"and deviates from index price (%d) more than minimum allowed (%d)",
				priceUpdate.Price,
				priceUpdate.MarketId,
				indexPrice,
				currMarketParamPrice.Price.Price,
				newDelta.Uint64(),
				oldDelta.Uint64(),
			)
		}
		return nil
	}

	// Update price crosses index price and old_ticks > 1: check new_ticks <= sqrt(old_ticks)
	if !newPriceMeetsSqrtCondition(oldDelta, newDelta, tickSizePpm) {
		return errorsmod.Wrapf(
			types.ErrInvalidMarketPriceUpdateNonDeterministic,
			"update price (%d) for market (%d) crosses the index price (%d) with current price (%d) "+
				"and deviates from index price (%d) more than minimum allowed (%d)",
			priceUpdate.Price,
			priceUpdate.MarketId,
			indexPrice,
			currMarketParamPrice.Price.Price,
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
	proposedUpdatesMap := make(map[uint32]struct{}, len(marketPriceUpdates))
	for _, proposedUpdate := range marketPriceUpdates {
		proposedUpdatesMap[proposedUpdate.MarketId] = struct{}{}
	}

	// Gather all markets that we think should be updated.
	var missingMarkets []uint32
	// Note that `GetValidMarketPriceUpdates` return value is ordered by market id.
	// This is NOT deterministic, because the returned values are based on "index price".
	allLocalUpdates := k.GetValidMarketPriceUpdates(ctx).MarketPriceUpdates
	for _, localUpdate := range allLocalUpdates {
		if _, exists := proposedUpdatesMap[localUpdate.MarketId]; !exists {
			missingMarkets = append(missingMarkets, localUpdate.MarketId)
		}
	}

	return missingMarkets
}

// getIdToMarketParamPrice returns a map of market id to market param price given a slice of market param prices.
func getIdToMarketParamPrice(marketParamPrices []types.MarketParamPrice) map[uint32]types.MarketParamPrice {
	idToMarketParamPrice := make(map[uint32]types.MarketParamPrice, len(marketParamPrices))
	for _, marketParamPrice := range marketParamPrices {
		idToMarketParamPrice[marketParamPrice.Param.Id] = marketParamPrice
	}
	return idToMarketParamPrice
}

// getMarketParamPrice returns a market param price given a market id. Returns an error if a market
// param price does not exist.
func getMarketParamPrice(
	marketId uint32,
	idToMarketParamPrice map[uint32]types.MarketParamPrice,
) (
	types.MarketParamPrice,
	error,
) {
	marketParamPrice, marketParamPriceExists := idToMarketParamPrice[marketId]
	if !marketParamPriceExists {
		return marketParamPrice, errorsmod.Wrapf(
			types.ErrInvalidMarketPriceUpdateDeterministic,
			"market param price (%d) does not exist",
			marketId,
		)
	}
	return marketParamPrice, nil
}
