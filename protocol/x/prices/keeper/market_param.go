package keeper

import (
	"sort"

	errorsmod "cosmossdk.io/errors"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/lib/slinky"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

// getMarketParamStore returns a prefix store for MarketParams.
func (k Keeper) getMarketParamStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.MarketParamKeyPrefix))
}

// ModifyMarketParam modifies an existing market param in the store.
func (k Keeper) ModifyMarketParam(
	ctx sdk.Context,
	updatedMarketParam types.MarketParam,
) (types.MarketParam, error) {
	// Validate input.
	if err := updatedMarketParam.Validate(); err != nil {
		return types.MarketParam{}, err
	}

	// Get existing market param.
	existingParam, exists := k.GetMarketParam(ctx, updatedMarketParam.Id)
	if !exists {
		return types.MarketParam{}, errorsmod.Wrap(
			types.ErrMarketParamDoesNotExist,
			lib.UintToString(updatedMarketParam.Id),
		)
	}

	// Validate update is permitted.
	for _, market := range k.GetAllMarketParams(ctx) {
		if market.Pair == updatedMarketParam.Pair && market.Id != updatedMarketParam.Id {
			return types.MarketParam{}, errorsmod.Wrap(types.ErrMarketParamPairAlreadyExists, updatedMarketParam.Pair)
		}
	}

	// Validate that modified market param has a corresponding ticker in MarketMap
	cp, err := slinky.MarketPairToCurrencyPair(updatedMarketParam.Pair)
	if err != nil {
		return types.MarketParam{}, errorsmod.Wrap(types.ErrMarketPairConversionFailed, updatedMarketParam.Pair)
	}
	if _, err := k.MarketMapKeeper.GetMarket(ctx, cp.String()); err != nil {
		return types.MarketParam{}, errorsmod.Wrap(types.ErrTickerNotFoundInMarketMap, cp.String())
	}

	// Store the modified market param.
	marketParamStore := k.getMarketParamStore(ctx)
	b := k.cdc.MustMarshal(&updatedMarketParam)
	marketParamStore.Set(lib.Uint32ToKey(updatedMarketParam.Id), b)

	// if the market pair has been changed, we need to update the in-memory market pair cache
	if existingParam.Pair != updatedMarketParam.Pair {
		// remove the old cache entry and disable the old market
		oldCurrencyPair, err := slinky.MarketPairToCurrencyPair(existingParam.Pair)
		if err != nil {
			return types.MarketParam{}, errorsmod.Wrap(
				types.ErrMarketPairConversionFailed,
				existingParam.Pair,
			)
		}

		k.RemoveCurrencyPairFromStore(ctx, oldCurrencyPair)
		if err = k.MarketMapKeeper.DisableMarket(ctx, oldCurrencyPair.String()); err != nil {
			return types.MarketParam{}, errorsmod.Wrap(
				types.ErrMarketCouldNotBeDisabled,
				existingParam.Pair,
			)
		}

		// add the new cache entry and enable the new market
		newCurrencyPair, err := slinky.MarketPairToCurrencyPair(updatedMarketParam.Pair)
		if err != nil {
			return types.MarketParam{}, errorsmod.Wrap(
				types.ErrMarketPairConversionFailed,
				updatedMarketParam.Pair,
			)
		}

		k.AddCurrencyPairIDToStore(ctx, updatedMarketParam.Id, newCurrencyPair)
		if err = k.MarketMapKeeper.EnableMarket(ctx, newCurrencyPair.String()); err != nil {
			return types.MarketParam{}, errorsmod.Wrap(
				types.ErrMarketCouldNotBeEnabled,
				existingParam.Pair,
			)
		}
	}

	// Generate indexer event.
	k.GetIndexerEventManager().AddTxnEvent(
		ctx,
		indexerevents.SubtypeMarket,
		indexerevents.MarketEventVersion,
		indexer_manager.GetBytes(
			indexerevents.NewMarketModifyEvent(
				updatedMarketParam.Id,
				updatedMarketParam.Pair,
				updatedMarketParam.MinPriceChangePpm,
			),
		),
	)

	// Update the in-memory market pair map for labelling metrics.
	metrics.SetMarketPairForTelemetry(updatedMarketParam.Id, updatedMarketParam.Pair)

	return updatedMarketParam, nil
}

// GetMarketParam returns a market param from its id.
func (k Keeper) GetMarketParam(
	ctx sdk.Context,
	id uint32,
) (
	market types.MarketParam,
	exists bool,
) {
	marketParamStore := k.getMarketParamStore(ctx)
	b := marketParamStore.Get(lib.Uint32ToKey(id))
	if b == nil {
		return types.MarketParam{}, false
	}

	k.cdc.MustUnmarshal(b, &market)
	return market, true
}

// GetAllMarketParams returns all market params.
func (k Keeper) GetAllMarketParams(ctx sdk.Context) []types.MarketParam {
	marketParamStore := k.getMarketParamStore(ctx)

	allMarketParams := make([]types.MarketParam, 0)

	iterator := marketParamStore.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		marketParam := types.MarketParam{}
		k.cdc.MustUnmarshal(iterator.Value(), &marketParam)
		allMarketParams = append(allMarketParams, marketParam)
	}

	// Sort the market params to return them in ascending order based on Id.
	sort.Slice(allMarketParams, func(i, j int) bool {
		return allMarketParams[i].Id < allMarketParams[j].Id
	})

	return allMarketParams
}
