package keeper

import (
	"time"

	errorsmod "cosmossdk.io/errors"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib/metrics"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	CrossingPriceUpdateCutoffPpm = uint32(500_000) // 50%
)

// PerformStatefulPriceUpdateValidation performs stateful validations on `UpdateMarketPrices`.
// Depending on the input, this func performs non-deterministic stateful validation.
func (k Keeper) PerformStatefulPriceUpdateValidation(
	ctx sdk.Context,
	marketPriceUpdate *types.MarketPriceUpdate,
) (isSpotValid bool, isPnlValid bool) {
	defer telemetry.ModuleMeasureSince(
		types.ModuleName,
		time.Now(),
		metrics.StatefulPriceUpdateValidation,
		metrics.Deterministic,
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
		return false, false
	}

	return k.performDeterministicStatefulValidation(marketPriceUpdate, marketParamPrices)
}

// performDeterministicStatefulValidation performs stateful validations that are deterministic.
//
// Specificically, for each price update, validate the following:
//   - The market exists.
//   - The price update is greater than the min price change.
func (k Keeper) performDeterministicStatefulValidation(
	marketPriceUpdate *types.MarketPriceUpdate,
	allMarketParamPrices []types.MarketParamPrice,
) (isSpotValid bool, isPnlValid bool) {
	// TODO: clean up, we don't need to loop to get ID - change this
	idToMarketParamPrice := getIdToMarketParamPrice(allMarketParamPrices)

	// Check market exists.
	marketParamPrice, err := getMarketParamPrice(marketPriceUpdate.MarketId, idToMarketParamPrice)
	if err != nil {
		return false, false
	}

	isSpotValid = isAboveRequiredMinSpotPriceChange(marketParamPrice, marketPriceUpdate.SpotPrice)
	isPnlValid = isAboveRequiredMinPnlPriceChange(marketParamPrice, marketPriceUpdate.PnlPrice)

	return isSpotValid, isPnlValid
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
