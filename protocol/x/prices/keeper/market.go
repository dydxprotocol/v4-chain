package keeper

import (
	"fmt"

	gogotypes "github.com/cosmos/gogoproto/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/metrics"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/slinky"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

// CreateMarket creates a new market param in the store along with a new market price
// for that market param.
// This is the only path to creating new MarketPrices, so if we have a param
// defined for a market, we should expect to see a price defined, and vice versa.
func (k Keeper) CreateMarket(
	ctx sdk.Context,
	marketParam types.MarketParam,
	marketPrice types.MarketPrice,
) (types.MarketParam, error) {
	if _, exists := k.GetMarketParam(ctx, marketParam.Id); exists {
		return types.MarketParam{}, errorsmod.Wrapf(
			types.ErrMarketParamAlreadyExists,
			"market param with id %d already exists",
			marketParam.Id,
		)
	}

	// Validate input.
	if err := marketParam.Validate(); err != nil {
		return types.MarketParam{}, err
	}
	if err := marketPrice.ValidateFromParam(marketParam); err != nil {
		return types.MarketParam{}, err
	}
	// Stateful Validation
	for _, market := range k.GetAllMarketParams(ctx) {
		if market.Pair == marketParam.Pair {
			return types.MarketParam{}, errorsmod.Wrap(
				types.ErrMarketParamPairAlreadyExists,
				marketParam.Pair,
			)
		}
	}
	// check that the market exists in market map
	currencyPair, err := slinky.MarketPairToCurrencyPair(marketParam.Pair)
	if err != nil {
		return types.MarketParam{}, errorsmod.Wrap(
			types.ErrMarketPairConversionFailed,
			marketParam.Pair,
		)
	}
	currencyPairStr := currencyPair.String()
	marketMapDetails, err := k.MarketMapKeeper.GetMarket(ctx, currencyPairStr)
	if err != nil {
		return types.MarketParam{}, errorsmod.Wrap(
			types.ErrTickerNotFoundInMarketMap,
			currencyPairStr,
		)
	}

	// Check that the exponent of market price is the negation of the decimals value in the market map
	if marketPrice.Exponent != int32(marketMapDetails.Ticker.Decimals)*-1 {
		return types.MarketParam{}, errorsmod.Wrap(
			types.ErrInvalidMarketPriceExponent,
			currencyPairStr,
		)
	}

	paramBytes := k.cdc.MustMarshal(&marketParam)
	priceBytes := k.cdc.MustMarshal(&marketPrice)

	marketParamStore := k.getMarketParamStore(ctx)
	marketParamStore.Set(lib.Uint32ToKey(marketParam.Id), paramBytes)

	marketPriceStore := k.getMarketPriceStore(ctx)
	marketPriceStore.Set(lib.Uint32ToKey(marketPrice.Id), priceBytes)

	// add the pair to the currency-pair-id cache
	k.AddCurrencyPairIDToStore(ctx, marketParam.Id, currencyPair)

	// Generate indexer event.
	k.GetIndexerEventManager().AddTxnEvent(
		ctx,
		indexerevents.SubtypeMarket,
		indexerevents.MarketEventVersion,
		indexer_manager.GetBytes(
			indexerevents.NewMarketCreateEvent(
				marketParam.Id,
				marketParam.Pair,
				marketParam.MinPriceChangePpm,
				// The exponent in market price is the source of truth, the exponent of the param is deprecated as of v7.1.x
				marketPrice.Exponent,
			),
		),
	)

	metrics.SetMarketPairForTelemetry(marketParam.Id, marketParam.Pair)

	// create a new market rev share
	k.RevShareKeeper.CreateNewMarketRevShare(ctx, marketParam.Id)

	// enable the market in the market map
	err = k.MarketMapKeeper.EnableMarket(ctx, currencyPairStr)
	if err != nil {
		k.Logger(ctx).Error(
			"failed to enable market in market map",
			"market ticker",
			currencyPairStr,
			"err",
			err,
		)
	}
	return marketParam, nil
}

// Get the exponent for a market as the negation of the decimals value in the market map
func (k Keeper) GetExponent(ctx sdk.Context, ticker string) (int32, error) {
	currencyPair, err := slinky.MarketPairToCurrencyPair(ticker)
	if err != nil {
		k.Logger(ctx).Error("Could not convert market pair to currency pair", "error", err)
		return 0, err
	}

	marketMapDetails, err := k.MarketMapKeeper.GetMarket(ctx, currencyPair.String())
	if err != nil {
		return 0, errorsmod.Wrap(
			types.ErrTickerNotFoundInMarketMap,
			ticker,
		)
	}
	return int32(marketMapDetails.Ticker.Decimals) * -1, nil
}

// GetAllMarketParamPrices returns a slice of MarketParam, MarketPrice tuples for all markets.
func (k Keeper) GetAllMarketParamPrices(ctx sdk.Context) ([]types.MarketParamPrice, error) {
	marketParams := k.GetAllMarketParams(ctx)
	marketPrices := k.GetAllMarketPrices(ctx)

	if len(marketParams) != len(marketPrices) {
		return nil, errorsmod.Wrap(types.ErrMarketPricesAndParamsDontMatch, "market param and price lengths do not match")
	}

	marketParamPrices := make([]types.MarketParamPrice, len(marketParams))
	for i, param := range marketParams {
		marketParamPrices[i].Param = param
		price := marketPrices[i]
		if param.Id != price.Id {
			return nil, errorsmod.Wrap(types.ErrMarketPricesAndParamsDontMatch,
				fmt.Sprintf("market param and price ids do not match: %d != %d", param.Id, price.Id))
		}
		marketParamPrices[i].Price = price
	}
	return marketParamPrices, nil
}

// GetNextMarketID returns the next market id to be used from the module store
func (k Keeper) GetNextMarketID(ctx sdk.Context) uint32 {
	store := ctx.KVStore(k.storeKey)
	b := store.Get([]byte(types.NextMarketIDKey))
	var result gogotypes.UInt32Value
	k.cdc.MustUnmarshal(b, &result)
	return result.Value
}

// SetNextMarketID sets the next market id to be used
func (k Keeper) SetNextMarketID(ctx sdk.Context, nextID uint32) {
	store := ctx.KVStore(k.storeKey)
	value := gogotypes.UInt32Value{Value: nextID}
	store.Set([]byte(types.NextMarketIDKey), k.cdc.MustMarshal(&value))
}

// AcquireNextMarketID returns the next market id to be used and increments the next market id
func (k Keeper) AcquireNextMarketID(ctx sdk.Context) uint32 {
	nextID := k.GetNextMarketID(ctx)
	// if market id already exists, increment until we find one that doesn't
	maxAttempts, attempts := 1000, 0
	for {
		_, exists := k.GetMarketParam(ctx, nextID)
		if !exists {
			break
		}
		nextID++

		// panic if we've tried too many times and are stuck in a loop
		attempts++
		if attempts >= maxAttempts {
			panic("Exceeded maximum attempts to find a unique market id")
		}
	}

	k.SetNextMarketID(ctx, nextID+1)
	return nextID
}
