package keeper

import (
	"sort"

	errorsmod "cosmossdk.io/errors"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/metrics"

	"github.com/cosmos/cosmos-sdk/store/prefix"
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
	if updatedMarketParam.Exponent != existingParam.Exponent {
		return types.MarketParam{},
			errorsmod.Wrapf(types.ErrMarketExponentCannotBeUpdated, lib.UintToString(updatedMarketParam.Id))
	}

	// Store the modified market param.
	marketParamStore := k.getMarketParamStore(ctx)
	b := k.cdc.MustMarshal(&updatedMarketParam)
	marketParamStore.Set(lib.Uint32ToKey(updatedMarketParam.Id), b)

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
