package keeper_test

import (
	"fmt"
	"testing"

	marketmapkeeper "github.com/dydxprotocol/slinky/x/marketmap/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/lib/slinky"

	errorsmod "cosmossdk.io/errors"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/require"
)

func TestModifyMarketParam(t *testing.T) {
	ctx, keeper, _, _, mockTimeProvider, _, _ := keepertest.PricesKeepers(t)
	mockTimeProvider.On("Now").Return(constants.TimeT)
	ctx = ctx.WithTxBytes(constants.TestTxBytes)
	items := keepertest.CreateNMarkets(t, ctx, keeper, 10)
	for i, item := range items {
		// Modify each field arbitrarily and
		// verify the fields were modified in state
		newItem, err := keeper.ModifyMarketParam(
			ctx,
			types.MarketParam{
				Id:                 item.Param.Id,
				Pair:               item.Param.Pair,
				MinExchanges:       uint32(2),
				Exponent:           item.Param.Exponent,
				MinPriceChangePpm:  uint32(9_999 - i),
				ExchangeConfigJson: fmt.Sprintf(`{"id":"%v"}`, i),
			},
		)
		require.NoError(t, err)
		require.Equal(t, uint32(i), newItem.Id)
		require.Equal(t, fmt.Sprintf("%v-%v", i, i), newItem.Pair)
		require.Equal(t, item.Param.Exponent, newItem.Exponent)
		require.Equal(t, uint32(2), newItem.MinExchanges)
		require.Equal(t, uint32(9999-i), newItem.MinPriceChangePpm)
		require.Equal(t, fmt.Sprintf("%v-%v", i, i), metrics.GetMarketPairForTelemetry(item.Param.Id))
		require.Equal(t, fmt.Sprintf(`{"id":"%v"}`, i), newItem.ExchangeConfigJson)
		keepertest.AssertMarketModifyEventInIndexerBlock(t, keeper, ctx, newItem)
	}
}

func TestModifyMarketParamUpdatesCache(t *testing.T) {
	ctx, keeper, _, _, mockTimeProvider, _, _ := keepertest.PricesKeepers(t)
	mockTimeProvider.On("Now").Return(constants.TimeT)
	ctx = ctx.WithTxBytes(constants.TestTxBytes)

	id := uint32(1)
	oldParam := types.MarketParam{
		Id:                 id,
		Pair:               "foo-bar",
		MinExchanges:       uint32(2),
		Exponent:           -8,
		MinPriceChangePpm:  uint32(50),
		ExchangeConfigJson: `{"id":"1"}`,
	}
	mp, err := keepertest.CreateTestMarket(t, ctx, keeper, oldParam, types.MarketPrice{
		Id:       id,
		Exponent: -8,
		Price:    1,
	})
	require.NoError(t, err)

	// check that the existing entry exists
	cp, err := slinky.MarketPairToCurrencyPair(mp.Pair)
	require.NoError(t, err)

	cpID, found := keeper.GetIDForCurrencyPair(ctx, cp)
	require.True(t, found)
	require.Equal(t, uint64(id), cpID)

	// create new ticker in MarketMap for newParam
	// (new ticker must exist in MarketMap before MarketParam.Pair can be updated)
	newParam := oldParam
	newParam.Pair = "bar-foo"
	keepertest.CreateMarketsInMarketMapFromParams(
		t,
		ctx,
		keeper.MarketMapKeeper.(*marketmapkeeper.Keeper),
		[]types.MarketParam{newParam},
	)

	// modify the market param
	newParam, err = keeper.ModifyMarketParam(ctx, newParam)
	require.NoError(t, err)

	// check that the existing entry does not exist
	_, found = keeper.GetIDForCurrencyPair(ctx, cp)
	require.False(t, found)

	// check that the new entry exists
	cp, err = slinky.MarketPairToCurrencyPair(newParam.Pair)
	require.NoError(t, err)
	cpID, found = keeper.GetIDForCurrencyPair(ctx, cp)
	require.True(t, found)
	require.Equal(t, uint64(id), cpID)
}

func TestModifyMarketParamUpdatesMarketmap(t *testing.T) {
	ctx, keeper, _, _, mockTimeProvider, _, marketMapKeeper := keepertest.PricesKeepers(t)
	mockTimeProvider.On("Now").Return(constants.TimeT)
	ctx = ctx.WithTxBytes(constants.TestTxBytes)

	id := uint32(1)
	oldParam := types.MarketParam{
		Id:                 id,
		Pair:               "foo-bar",
		MinExchanges:       uint32(2),
		Exponent:           -8,
		MinPriceChangePpm:  uint32(50),
		ExchangeConfigJson: `{"id":"1"}`,
	}
	mp, err := keepertest.CreateTestMarket(t, ctx, keeper, oldParam, types.MarketPrice{
		Id:       id,
		Exponent: -8,
		Price:    1,
	})
	require.NoError(t, err)

	// check that the market is enabled for the existing pair
	oldCp, err := slinky.MarketPairToCurrencyPair(mp.Pair)
	require.NoError(t, err)

	oldMarket, err := marketMapKeeper.GetMarket(ctx, oldCp.String())
	require.NoError(t, err)
	require.True(t, oldMarket.Ticker.Enabled)

	// create new ticker in MarketMap for newParam
	// (new ticker must exist in MarketMap before MarketParam.Pair can be updated)
	newParam := oldParam
	newParam.Pair = "bar-foo"
	newCp, err := slinky.MarketPairToCurrencyPair(newParam.Pair)
	require.NoError(t, err)
	keepertest.CreateMarketsInMarketMapFromParams(
		t,
		ctx,
		keeper.MarketMapKeeper.(*marketmapkeeper.Keeper),
		[]types.MarketParam{newParam},
	)

	// check that the new market is disabled to start
	newMarket, err := marketMapKeeper.GetMarket(ctx, newCp.String())
	require.NoError(t, err)
	require.False(t, newMarket.Ticker.Enabled)

	// modify the market param
	newParam, err = keeper.ModifyMarketParam(ctx, newParam)
	require.NoError(t, err)

	// check that old market is disabled
	oldMarket, err = marketMapKeeper.GetMarket(ctx, oldCp.String())
	require.NoError(t, err)
	require.False(t, oldMarket.Ticker.Enabled)

	// check that the new market is enabled
	newMarket, err = marketMapKeeper.GetMarket(ctx, newCp.String())
	require.NoError(t, err)
	require.True(t, newMarket.Ticker.Enabled)
}

func TestModifyMarketParam_Errors(t *testing.T) {
	validExchangeConfigJson := `{"exchanges":[{"exchangeName":"Binance","ticker":"BTCUSDT"}]}`
	invalidUpdatePair := "nil-nil"
	invalidUpdateCurrencyPair, err := slinky.MarketPairToCurrencyPair(invalidUpdatePair)
	require.NoError(t, err)

	tests := map[string]struct {
		// Setup
		targetId           uint32
		pair               string
		minExchanges       uint32
		minPriceChangePpm  uint32
		exchangeConfigJson string

		// Expected
		expectedErr string
	}{
		"Market not found": {
			targetId:           99, // this market id does not exist
			pair:               constants.BtcUsdPair,
			minExchanges:       uint32(2),
			minPriceChangePpm:  uint32(50),
			exchangeConfigJson: validExchangeConfigJson,
			expectedErr:        errorsmod.Wrap(types.ErrMarketParamDoesNotExist, "99").Error(),
		},
		"Empty pair": {
			targetId:           0,
			pair:               "", // pair cannot be empty
			minExchanges:       uint32(2),
			minPriceChangePpm:  uint32(50),
			exchangeConfigJson: validExchangeConfigJson,
			expectedErr:        errorsmod.Wrap(types.ErrInvalidInput, constants.ErrorMsgMarketPairCannotBeEmpty).Error(),
		},
		"Pair name invalid": {
			targetId:           0,
			pair:               "test", // pair must be in format {Base}-{Quote}
			minExchanges:       uint32(2),
			minPriceChangePpm:  uint32(50),
			exchangeConfigJson: validExchangeConfigJson,
			expectedErr:        errorsmod.Wrap(types.ErrMarketPairConversionFailed, "test").Error(),
		},
		"Invalid min price change: zero": {
			targetId:           0,
			pair:               constants.BtcUsdPair,
			minExchanges:       uint32(2),
			minPriceChangePpm:  uint32(0), // must be > 0
			exchangeConfigJson: validExchangeConfigJson,
			expectedErr:        errorsmod.Wrap(types.ErrInvalidInput, constants.ErrorMsgInvalidMinPriceChange).Error(),
		},
		"Invalid min price change: ten thousand": {
			targetId:           0,
			pair:               constants.BtcUsdPair,
			minExchanges:       uint32(2),
			minPriceChangePpm:  uint32(10_000), // must be < 10,000
			exchangeConfigJson: validExchangeConfigJson,
			expectedErr:        errorsmod.Wrap(types.ErrInvalidInput, constants.ErrorMsgInvalidMinPriceChange).Error(),
		},
		"Updating pair fails due to pair already existing": {
			targetId:           0,
			pair:               "1-1",
			minExchanges:       uint32(1),
			minPriceChangePpm:  uint32(50),
			exchangeConfigJson: validExchangeConfigJson,
			expectedErr: errorsmod.Wrap(
				types.ErrMarketParamPairAlreadyExists,
				"1-1",
			).Error(),
		},
		"Updating pair fails due to no ticker in MarketMap": {
			targetId:           0,
			pair:               invalidUpdatePair,
			minExchanges:       uint32(1),
			minPriceChangePpm:  uint32(50),
			exchangeConfigJson: validExchangeConfigJson,
			expectedErr: errorsmod.Wrap(
				types.ErrTickerNotFoundInMarketMap,
				invalidUpdateCurrencyPair.String(),
			).Error(),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, keeper, _, _, mockTimeKeeper, _, _ := keepertest.PricesKeepers(t)
			mockTimeKeeper.On("Now").Return(constants.TimeT)
			ctx = ctx.WithTxBytes(constants.TestTxBytes)
			keepertest.CreateNMarkets(t, ctx, keeper, 2)
			_, err := keeper.ModifyMarketParam(
				ctx,
				types.MarketParam{
					Id:                 tc.targetId,
					Pair:               tc.pair,
					Exponent:           -int32(tc.targetId + 1), // consistent with CreateNMarkets()
					MinExchanges:       tc.minExchanges,
					MinPriceChangePpm:  tc.minPriceChangePpm,
					ExchangeConfigJson: tc.exchangeConfigJson,
				},
			)
			require.EqualError(t, err, tc.expectedErr)
		})
	}
}

func TestGetMarketParam(t *testing.T) {
	ctx, keeper, _, _, mockTimeProvider, _, _ := keepertest.PricesKeepers(t)
	mockTimeProvider.On("Now").Return(constants.TimeT)
	items := keepertest.CreateNMarkets(t, ctx, keeper, 10)
	for _, item := range items {
		rst, exists := keeper.GetMarketParam(ctx, item.Param.Id)
		require.True(t, exists)
		require.Equal(
			t,
			&item.Param,
			&rst,
		)
	}
}

func TestGetMarketParam_NotFound(t *testing.T) {
	ctx, keeper, _, _, _, _, _ := keepertest.PricesKeepers(t)
	_, exists := keeper.GetMarketParam(ctx, uint32(0))
	require.False(t, exists)
}

func TestGetAllMarketParams(t *testing.T) {
	ctx, keeper, _, _, mockTimeProvider, _, _ := keepertest.PricesKeepers(t)
	mockTimeProvider.On("Now").Return(constants.TimeT)
	items := keepertest.CreateNMarkets(t, ctx, keeper, 10)
	params := make([]types.MarketParam, len(items))
	for i, item := range items {
		params[i] = item.Param
	}
	require.ElementsMatch(
		t,
		params,
		keeper.GetAllMarketParams(ctx),
	)
}
