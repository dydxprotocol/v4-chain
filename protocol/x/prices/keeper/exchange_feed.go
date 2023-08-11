package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/x/prices/types"
)

// CreateExchangeFeed creates a new ExchangeFeed.
func (k Keeper) CreateExchangeFeed(
	ctx sdk.Context,
	name string,
	memo string,
) (types.ExchangeFeed, error) {
	// Validate input.
	if err := validateExchangeName(name); err != nil {
		return types.ExchangeFeed{}, err
	}
	if err := validateExchangeMemo(memo); err != nil {
		return types.ExchangeFeed{}, err
	}

	// Get the `nextId`.
	nextId := k.GetNumExchangeFeeds(ctx)

	// Create the new feed.
	feed := types.ExchangeFeed{
		Id:   nextId,
		Name: name,
		Memo: memo,
	}

	// Store the new feed.
	b := k.cdc.MustMarshal(&feed)
	feedStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ExchangeFeedKeyPrefix))
	feedStore.Set(types.ExchangeFeedKey(nextId), b)

	// Store the new `numExchangeFeeds`.
	k.setNumExchangeFeeds(ctx, nextId+1)

	return feed, nil
}

// ModifyExchangeFeed modifies an existing ExchangeFeed.
func (k Keeper) ModifyExchangeFeed(
	ctx sdk.Context,
	id uint32,
	memo string,
) (types.ExchangeFeed, error) {
	// Validate input.
	// No validation for the name field since it's not modifiable.
	if err := validateExchangeMemo(memo); err != nil {
		return types.ExchangeFeed{}, err
	}

	// Get necessary store.
	feedStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ExchangeFeedKeyPrefix))

	// Get feed.
	feed, err := k.GetExchangeFeed(ctx, id)
	if err != nil {
		return types.ExchangeFeed{}, err
	}

	// Modify exchange.
	feed.Memo = memo

	// Store the modified feed.
	b := k.cdc.MustMarshal(&feed)
	feedStore.Set(types.ExchangeFeedKey(id), b)

	return feed, nil
}

// GetExchangeFeed returns a feed from its id.
func (k Keeper) GetExchangeFeed(
	ctx sdk.Context,
	id uint32,
) (types.ExchangeFeed, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ExchangeFeedKeyPrefix))

	b := store.Get(types.ExchangeFeedKey(id))
	if b == nil {
		return types.ExchangeFeed{}, sdkerrors.Wrap(types.ErrExchangeFeedDoesNotExist, lib.Uint32ToString(id))
	}

	var exchange = types.ExchangeFeed{}
	k.cdc.MustUnmarshal(b, &exchange)
	return exchange, nil
}

// GetAllExchangeFeed returns all exchange feeds.
func (k Keeper) GetAllExchangeFeeds(ctx sdk.Context) []types.ExchangeFeed {
	num := k.GetNumExchangeFeeds(ctx)
	feeds := make([]types.ExchangeFeed, num)

	for i := uint32(0); i < num; i++ {
		feed, err := k.GetExchangeFeed(ctx, i)
		if err != nil {
			panic(err)
		}

		feeds[i] = feed
	}

	return feeds
}

// GetNumExchangeFeeds returns the total number of exchange feeds.
func (k Keeper) GetNumExchangeFeeds(
	ctx sdk.Context,
) uint32 {
	store := ctx.KVStore(k.storeKey)
	var numExchangeFeedsBytes []byte = store.Get(types.KeyPrefix(types.NumExchangeFeedsKey))
	return lib.BytesToUint32(numExchangeFeedsBytes)
}

// setNumExchangeFeeds sets the number of exchanges with a new value.
func (k Keeper) setNumExchangeFeeds(
	ctx sdk.Context,
	newValue uint32,
) {
	// Get necessary stores.
	store := ctx.KVStore(k.storeKey)

	// Set `numExchangeFeeds`.
	store.Set(types.KeyPrefix(types.NumExchangeFeedsKey), lib.Uint32ToBytes(newValue))
}

// validateExchangeName validates exchange name input field.
func validateExchangeName(name string) error {
	if name == "" {
		return sdkerrors.Wrap(types.ErrInvalidInput, "Name cannot be empty")
	}
	return nil
}

// validateExchangeMemo validates exchange memo input field.
func validateExchangeMemo(memo string) error {
	if memo == "" {
		return sdkerrors.Wrap(types.ErrInvalidInput, "Memo cannot be empty")
	}
	return nil
}
