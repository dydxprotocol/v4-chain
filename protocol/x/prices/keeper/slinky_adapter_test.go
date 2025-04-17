package keeper_test

import (
	"fmt"
	"testing"

	oracletypes "github.com/dydxprotocol/slinky/pkg/types"
	"github.com/stretchr/testify/require"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

func TestCurrencyPairIDStoreFunctions(t *testing.T) {
	ctx, keeper, _, _, _, _, _ := keepertest.PricesKeepers(t)

	currencyPair := oracletypes.CurrencyPair{
		Base:  "BTC",
		Quote: "USD",
	}

	// Add the currency pair ID to the store
	marketID := uint32(1)
	keeper.AddCurrencyPairIDToStore(ctx, marketID, currencyPair)

	// Retrieve the currency pair ID from the store
	storedMarketID, found := keeper.GetCurrencyPairIDFromStore(ctx, currencyPair)

	require.True(t, found)
	require.Equal(t, uint64(marketID), storedMarketID)

	// Remove the currency pair ID from the store
	keeper.RemoveCurrencyPairFromStore(ctx, currencyPair)

	_, found = keeper.GetCurrencyPairIDFromStore(ctx, currencyPair)
	require.False(t, found)
}

func TestGetCurrencyPairFromID(t *testing.T) {
	ctx, keeper, _, _, mockTimeProvider, _, _ := keepertest.PricesKeepers(t)
	mockTimeProvider.On("Now").Return(constants.TimeT)

	marketNumber := 10
	items := keepertest.CreateNMarkets(t, ctx, keeper, marketNumber)
	marketParams := keeper.GetAllMarketParams(ctx)
	require.Equal(t, len(marketParams), marketNumber)
	for _, mpp := range items {
		mpId := mpp.Param.Id
		_, found := keeper.GetCurrencyPairFromID(ctx, uint64(mpId))
		require.True(t, found)
	}
	_, found := keeper.GetCurrencyPairFromID(ctx, uint64(marketNumber+1))
	require.False(t, found)
}

func TestIDForCurrencyPair(t *testing.T) {
	ctx, keeper, _, _, mockTimeProvider, _, _ := keepertest.PricesKeepers(t)
	mockTimeProvider.On("Now").Return(constants.TimeT)

	marketNumber := 10
	_ = keepertest.CreateNMarkets(t, ctx, keeper, marketNumber)
	marketParams := keeper.GetAllMarketParams(ctx)
	require.Equal(t, len(marketParams), marketNumber)
	for i := 0; i < marketNumber; i++ {
		pair := oracletypes.CurrencyPair{
			Base:  fmt.Sprint(i),
			Quote: fmt.Sprint(i),
		}
		id, found := keeper.GetIDForCurrencyPair(ctx, pair)
		require.True(t, found)
		require.Equal(t, uint64(i), id)
	}
	_, found := keeper.GetIDForCurrencyPair(ctx, oracletypes.CurrencyPair{
		Base:  fmt.Sprint(marketNumber + 1),
		Quote: fmt.Sprint(marketNumber + 1),
	})
	require.False(t, found)
}

func TestGetPriceForCurrencyPair(t *testing.T) {
	ctx, keeper, _, _, mockTimeProvider, _, _ := keepertest.PricesKeepers(t)
	mockTimeProvider.On("Now").Return(constants.TimeT)

	marketNumber := 10
	items := keepertest.CreateNMarkets(t, ctx, keeper, marketNumber)
	marketParams := keeper.GetAllMarketParams(ctx)
	require.Equal(t, len(marketParams), marketNumber)
	for i := 0; i < marketNumber; i++ {
		pair := oracletypes.CurrencyPair{
			Base:  fmt.Sprint(i),
			Quote: fmt.Sprint(i),
		}
		price, err := keeper.GetPriceForCurrencyPair(ctx, pair)
		require.NoError(t, err)
		require.Equal(t, items[i].Price.Price, price.Price.Uint64())
	}
	_, err := keeper.GetPriceForCurrencyPair(ctx, oracletypes.CurrencyPair{
		Base:  fmt.Sprint(marketNumber + 1),
		Quote: fmt.Sprint(marketNumber + 1),
	})
	require.Error(t, err)
}

func TestBadMarketData(t *testing.T) {
	ctx, keeper, _, _, mockTimeProvider, _, _ := keepertest.PricesKeepers(t)
	mockTimeProvider.On("Now").Return(constants.TimeT)

	_, err := keeper.CreateMarket(
		ctx,
		types.MarketParam{
			Id:                uint32(0),
			Pair:              "00",
			MinPriceChangePpm: 1,
		},
		types.MarketPrice{})
	require.Error(t, err)

	_, found := keeper.GetCurrencyPairFromID(ctx, uint64(0))
	require.False(t, found)

	_, found = keeper.GetIDForCurrencyPair(ctx, oracletypes.CurrencyPair{})
	require.False(t, found)

	_, err = keeper.GetPriceForCurrencyPair(ctx, oracletypes.CurrencyPair{})
	require.Error(t, err)
}

func TestGetNumCurrencyPairs(t *testing.T) {
	ctx, keeper, _, _, mockTimeProvider, _, _ := keepertest.PricesKeepers(t)
	mockTimeProvider.On("Now").Return(constants.TimeT)

	marketNumber := 10
	_ = keepertest.CreateNMarkets(t, ctx, keeper, marketNumber)
	cpCounter, err := keeper.GetNumCurrencyPairs(ctx)
	require.NoError(t, err)
	require.Equal(t, uint64(marketNumber), cpCounter)
}
