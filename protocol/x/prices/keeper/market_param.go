package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"sort"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

// newMarketParamStore creates a new prefix store for MarketParams.
func (k Keeper) newMarketParamStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MarketParamKeyPrefix))
}

// ModifyMarketParam modifies an existing market param in the store.
func (k Keeper) ModifyMarketParam(
	ctx sdk.Context,
	marketParam types.MarketParam,
) (types.MarketParam, error) {
	// Validate input.
	if err := marketParam.Validate(); err != nil {
		return types.MarketParam{}, err
	}

	// Get existing market param.
	existingParam, err := k.GetMarketParam(ctx, marketParam.Id)
	if err != nil {
		return types.MarketParam{}, err
	}

	// Validate update is permitted.
	if marketParam.Exponent != existingParam.Exponent {
		return types.MarketParam{},
			errorsmod.Wrapf(types.ErrMarketExponentCannotBeUpdated, lib.Uint32ToString(marketParam.Id))
	}

	// Store the modified market param.
	marketParamStore := k.newMarketParamStore(ctx)
	b := k.cdc.MustMarshal(&marketParam)
	marketParamStore.Set(types.MarketKey(marketParam.Id), b)

	k.GetIndexerEventManager().AddTxnEvent(
		ctx,
		indexerevents.SubtypeMarket,
		indexer_manager.GetB64EncodedEventMessage(
			indexerevents.NewMarketModifyEvent(
				marketParam.Id,
				marketParam.Pair,
				marketParam.MinPriceChangePpm,
			),
		),
	)

	return marketParam, nil
}

// GetMarketParam returns a market param from its id.
func (k Keeper) GetMarketParam(
	ctx sdk.Context,
	id uint32,
) (types.MarketParam, error) {
	marketParamStore := k.newMarketParamStore(ctx)
	b := marketParamStore.Get(types.MarketKey(id))
	if b == nil {
		return types.MarketParam{}, errorsmod.Wrap(types.ErrMarketParamDoesNotExist, lib.Uint32ToString(id))
	}

	var market = types.MarketParam{}
	k.cdc.MustUnmarshal(b, &market)
	return market, nil
}

// GetAllMarketParams returns all market params.
func (k Keeper) GetAllMarketParams(ctx sdk.Context) []types.MarketParam {
	marketParamStore := k.newMarketParamStore(ctx)

	marketParams := make([]types.MarketParam, 0)

	iterator := marketParamStore.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		marketParam := types.MarketParam{}
		k.cdc.MustUnmarshal(iterator.Value(), &marketParam)
		marketParams = append(marketParams, marketParam)
	}

	// Sort the market params to return them in ascending order based on Id.
	sort.Slice(marketParams, func(i, j int) bool {
		return marketParams[i].Id < marketParams[j].Id
	})

	return marketParams
}
