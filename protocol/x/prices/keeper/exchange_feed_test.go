package keeper_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/testutil/constants"
	keepertest "github.com/dydxprotocol/v4/testutil/keeper"
	"github.com/dydxprotocol/v4/testutil/nullify"
	"github.com/dydxprotocol/v4/x/prices/keeper"
	"github.com/dydxprotocol/v4/x/prices/types"
	"github.com/stretchr/testify/require"
)

func createNExchangeFeeds(t *testing.T, keeper *keeper.Keeper, ctx sdk.Context, n int) []types.ExchangeFeed {
	items := make([]types.ExchangeFeed, n)
	for i := range items {
		items[i].Id = uint32(i)
		items[i].Memo = fmt.Sprintf("%v", i)
		items[i].Name = fmt.Sprintf("%v", i)

		_, err := keeper.CreateExchangeFeed(ctx, items[i].Name, items[i].Memo)

		require.NoError(t, err)
	}
	return items
}

func TestCreateExchangeFeed(t *testing.T) {
	ctx, keeper, _, _, _, _ := keepertest.PricesKeepers(t)
	created, err := keeper.CreateExchangeFeed(ctx, constants.CoinbaseExchangeName, "bar")

	require.NoError(t, err)
	require.Equal(t, uint32(0), created.Id)
	require.Equal(t, constants.CoinbaseExchangeName, created.Name)
	require.Equal(t, "bar", created.Memo)
}

func TestCreateExchangeFeed_Errors(t *testing.T) {
	tests := map[string]struct {
		// Setup
		name string
		memo string

		// Expected
		expectedErr string
	}{
		"Empty name": {
			name:        "", // name cannot be empty
			memo:        "valid memo",
			expectedErr: sdkerrors.Wrap(types.ErrInvalidInput, "Name cannot be empty").Error(),
		},
		"Empty memo": {
			name:        "dee-why-dee-ex",
			memo:        "", // memo cannot be empty
			expectedErr: sdkerrors.Wrap(types.ErrInvalidInput, "Memo cannot be empty").Error(),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, keeper, _, _, _, _ := keepertest.PricesKeepers(t)
			_, err := keeper.CreateExchangeFeed(ctx, tc.name, tc.memo)
			require.EqualError(t, err, tc.expectedErr)
			require.ErrorIs(t, err, types.ErrInvalidInput)
		})
	}
}

func TestModifyExchangeFeed(t *testing.T) {
	ctx, keeper, _, _, _, _ := keepertest.PricesKeepers(t)
	exchange, exchangeErr := keeper.CreateExchangeFeed(ctx, constants.CoinbaseExchangeName, "bar")
	updated, err := keeper.ModifyExchangeFeed(ctx, exchange.Id, "foo")

	require.NoError(t, exchangeErr)
	require.NoError(t, err)
	require.Equal(t, constants.CoinbaseExchangeName, updated.Name)
	require.Equal(t, "foo", updated.Memo)
}

func TestModifyExchangeFeed_InvalidInput(t *testing.T) {
	ctx, keeper, _, _, _, _ := keepertest.PricesKeepers(t)
	exchange, exchangeErr := keeper.CreateExchangeFeed(ctx, constants.CoinbaseExchangeName, "bar")
	_, err := keeper.ModifyExchangeFeed(ctx, exchange.Id, "")

	require.NoError(t, exchangeErr)
	require.EqualError(t, err, sdkerrors.Wrap(types.ErrInvalidInput, "Memo cannot be empty").Error())
	require.ErrorIs(t, err, types.ErrInvalidInput)
}

func TestModifyExchangeFeed_NotFound(t *testing.T) {
	ctx, keeper, _, _, _, _ := keepertest.PricesKeepers(t)
	_, err := keeper.ModifyExchangeFeed(ctx, 0, "foo")

	require.EqualError(t, err, sdkerrors.Wrap(types.ErrExchangeFeedDoesNotExist, "0").Error())
	require.ErrorIs(t, err, types.ErrExchangeFeedDoesNotExist)
}

func TestGetExchangeFeed(t *testing.T) {
	ctx, keeper, _, _, _, _ := keepertest.PricesKeepers(t)
	exchange, exchangeErr := keeper.CreateExchangeFeed(ctx, constants.CoinbaseExchangeName, "bar")
	retrieved, err := keeper.GetExchangeFeed(ctx, exchange.Id)

	require.NoError(t, exchangeErr)
	require.NoError(t, err)
	require.Equal(t, uint32(0), retrieved.Id)
	require.Equal(t, constants.CoinbaseExchangeName, retrieved.Name)
	require.Equal(t, "bar", retrieved.Memo)
}

func TestGetExchangeFeed_NotFound(t *testing.T) {
	ctx, keeper, _, _, _, _ := keepertest.PricesKeepers(t)
	_, err := keeper.GetExchangeFeed(ctx, 0)

	require.EqualError(t, err, sdkerrors.Wrap(types.ErrExchangeFeedDoesNotExist, "0").Error())
	require.ErrorIs(t, err, types.ErrExchangeFeedDoesNotExist)
}

func TestGetAllExchangeFeeds(t *testing.T) {
	ctx, keeper, _, _, _, _ := keepertest.PricesKeepers(t)
	items := createNExchangeFeeds(t, keeper, ctx, 10)
	require.ElementsMatch(
		t,
		nullify.Fill(items), //nolint:staticcheck
		nullify.Fill(keeper.GetAllExchangeFeeds(ctx)), //nolint:staticcheck
	)
}

func TestGetAllExchangeFeeds_MissingExchange(t *testing.T) {
	ctx, keeper, storeKey, _, _, _ := keepertest.PricesKeepers(t)

	// Write some bad data to the store
	store := ctx.KVStore(storeKey)
	store.Set(types.KeyPrefix(types.NumExchangeFeedsKey), lib.Uint32ToBytes(20))

	// Expect a panic
	require.PanicsWithError(
		t,
		"0: ExchangeFeed does not exist",
		func() { keeper.GetAllExchangeFeeds(ctx) },
	)
}

func TestGetNumExchangeFeeds(t *testing.T) {
	ctx, keeper, _, _, _, _ := keepertest.PricesKeepers(t)
	require.Equal(t, uint32(0), keeper.GetNumExchangeFeeds(ctx))

	createNExchangeFeeds(t, keeper, ctx, 10)
	require.Equal(t, uint32(10), keeper.GetNumExchangeFeeds(ctx))
}
