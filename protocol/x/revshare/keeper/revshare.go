package keeper

import (
	"fmt"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
)

// Function to serialize market mapper revenue share params and store in the module store
func (k Keeper) SetMarketMapperRevenueShareParams(
	ctx sdk.Context,
	params types.MarketMapperRevenueShareParams,
) (err error) {
	// Validate the params
	if err := params.Validate(); err != nil {
		return err
	}

	// Store the params in the module store
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&params)
	store.Set([]byte(types.MarketMapperRevenueShareParamsKey), b)

	return nil
}

// Function to get market mapper revenue share params from the module store
func (k Keeper) GetMarketMapperRevenueShareParams(
	ctx sdk.Context,
) (params types.MarketMapperRevenueShareParams) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get([]byte(types.MarketMapperRevenueShareParamsKey))
	k.cdc.MustUnmarshal(b, &params)
	return params
}

// Function to serialize market mapper rev share details for a market
// and store in the module store
func (k Keeper) SetMarketMapperRevShareDetails(
	ctx sdk.Context,
	marketId uint32,
	params types.MarketMapperRevShareDetails,
) {
	// Store the rev share details for provided market in module store
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.MarketMapperRevSharePrefix))
	b := k.cdc.MustMarshal(&params)
	store.Set(lib.Uint32ToKey(marketId), b)
	return
}

// Function to retrieve marketmapper revshare details for a market from module store
func (k Keeper) GetMarketMapperRevShareDetails(
	ctx sdk.Context,
	marketId uint32,
) (params types.MarketMapperRevShareDetails, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.MarketMapperRevSharePrefix))
	b := store.Get(lib.Uint32ToKey(marketId))
	if b == nil {
		return params, fmt.Errorf("MarketMapperRevShareDetails not found for marketId: %d", marketId)
	}
	k.cdc.MustUnmarshal(b, &params)
	return params, nil
}

// Function to perform all market creation actions for the revshare module
func (k Keeper) CreateNewMarketRevShare(ctx sdk.Context, marketId uint32) {
	revShareParams := k.GetMarketMapperRevenueShareParams(ctx)

	validDurationSeconds := int64(revShareParams.ValidDays * 24 * 60 * 60)

	// set the rev share details for the market
	details := types.MarketMapperRevShareDetails{
		ExpirationTs: uint64(ctx.BlockTime().Unix() + validDurationSeconds),
	}
	k.SetMarketMapperRevShareDetails(ctx, marketId, details)
}
