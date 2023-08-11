package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	indexerevents "github.com/dydxprotocol/v4/indexer/events"
	"github.com/dydxprotocol/v4/indexer/indexer_manager"
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/x/prices/types"
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
			sdkerrors.Wrapf(types.ErrMarketExponentCannotBeUpdated, lib.Uint32ToString(marketParam.Id))
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
		return types.MarketParam{}, sdkerrors.Wrap(types.ErrMarketParamDoesNotExist, lib.Uint32ToString(id))
	}

	var market = types.MarketParam{}
	k.cdc.MustUnmarshal(b, &market)
	return market, nil
}

// GetAllMarketParams returns all market params.
func (k Keeper) GetAllMarketParams(ctx sdk.Context) []types.MarketParam {
	num := k.GetNumMarkets(ctx)
	marketParams := make([]types.MarketParam, num)

	for i := uint32(0); i < num; i++ {
		marketParam, err := k.GetMarketParam(ctx, i)
		if err != nil {
			panic(err)
		}

		marketParams[i] = marketParam
	}

	return marketParams
}
