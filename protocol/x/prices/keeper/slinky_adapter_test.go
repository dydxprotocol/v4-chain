package keeper_test

import (
	"fmt"
	"github.com/dydxprotocol/v4-chain/protocol/lib/slinky"
	"testing"

	oracletypes "github.com/skip-mev/slinky/pkg/types"
	"github.com/stretchr/testify/require"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

func TestGetCurrencyPairFromID(t *testing.T) {
	ctx, keeper, _, _, _, mockTimeProvider := keepertest.PricesKeepers(t)
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
	require.True(t, !found)
}

func TestIDForCurrencyPair(t *testing.T) {
	ctx, keeper, _, _, _, mockTimeProvider := keepertest.PricesKeepers(t)
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
	require.True(t, !found)
}

func TestGetPriceForCurrencyPair(t *testing.T) {
	ctx, keeper, _, _, _, mockTimeProvider := keepertest.PricesKeepers(t)
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
	ctx, keeper, _, _, _, mockTimeProvider := keepertest.PricesKeepers(t)
	mockTimeProvider.On("Now").Return(constants.TimeT)

	_, err := keeper.CreateMarket(
		ctx,
		types.MarketParam{
			Id:                 uint32(0),
			Pair:               "00",
			MinExchanges:       1,
			MinPriceChangePpm:  1,
			ExchangeConfigJson: "{}",
		},
		types.MarketPrice{})
	require.NoError(t, err)

	_, found := keeper.GetCurrencyPairFromID(ctx, uint64(0))
	require.False(t, found)

	_, found = keeper.GetIDForCurrencyPair(ctx, oracletypes.CurrencyPair{})
	require.False(t, found)

	_, err = keeper.GetPriceForCurrencyPair(ctx, oracletypes.CurrencyPair{})
	require.Error(t, err)
}

func TestMarketPairToCurrencyPair(t *testing.T) {
	testCases := []struct {
		marketPair        string
		currencyPairBase  string
		currencyPairQuote string
		shouldFail        bool
	}{
		{"0-0", "0", "0", false},
		{"0/0", "", "", true},
		{"00", "", "", true},
	}
	for _, tc := range testCases {
		t.Run(tc.marketPair, func(t *testing.T) {
			cp, err := slinky.MarketPairToCurrencyPair(tc.marketPair)
			if tc.shouldFail {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.currencyPairBase, cp.Base)
				require.Equal(t, tc.currencyPairQuote, cp.Quote)
			}
		})
	}
}
