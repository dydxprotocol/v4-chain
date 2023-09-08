package keeper_test

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/require"
)

func TestModifyMarketParam(t *testing.T) {
	ctx, keeper, _, _, _, _ := keepertest.PricesKeepers(t)
	ctx = ctx.WithTxBytes(constants.TestTxBytes)
	items := keepertest.CreateNMarkets(t, ctx, keeper, 10)
	for i, item := range items {
		// Modify each field arbitrarily and
		// verify the fields were modified in state
		newItem, err := keeper.ModifyMarketParam(
			ctx,
			types.MarketParam{
				Id:                 item.Param.Id,
				Pair:               fmt.Sprintf("foo_%v", i),
				MinExchanges:       uint32(2),
				Exponent:           item.Param.Exponent,
				MinPriceChangePpm:  uint32(9_999 - i),
				ExchangeConfigJson: fmt.Sprintf("config_%v", i),
			},
		)
		require.NoError(t, err)
		require.Equal(t, uint32(i), newItem.Id)
		require.Equal(t, fmt.Sprintf("foo_%v", i), newItem.Pair)
		require.Equal(t, item.Param.Exponent, newItem.Exponent)
		require.Equal(t, uint32(2), newItem.MinExchanges)
		require.Equal(t, uint32(9999-i), newItem.MinPriceChangePpm)
		require.Equal(t, fmt.Sprintf("config_%v", i), newItem.ExchangeConfigJson)
		keepertest.AssertMarketModifyEventInIndexerBlock(t, keeper, ctx, newItem)
	}
}

func TestModifyMarketParam_Errors(t *testing.T) {
	tests := map[string]struct {
		// Setup
		targetId          uint32
		pair              string
		minExchanges      uint32
		minPriceChangePpm uint32

		// Expected
		expectedErr string
	}{
		"Market not found": {
			targetId:          99, // this market id does not exist
			pair:              constants.BtcUsdPair,
			minExchanges:      uint32(2),
			minPriceChangePpm: uint32(50),
			expectedErr:       errorsmod.Wrap(types.ErrMarketParamDoesNotExist, "99").Error(),
		},
		"Empty pair": {
			targetId:          0,
			pair:              "", // pair cannot be empty
			minExchanges:      uint32(2),
			minPriceChangePpm: uint32(50),
			expectedErr:       errorsmod.Wrap(types.ErrInvalidInput, constants.ErrorMsgMarketPairCannotBeEmpty).Error(),
		},
		"Invalid min price change: zero": {
			targetId:          0,
			pair:              constants.BtcUsdPair,
			minExchanges:      uint32(2),
			minPriceChangePpm: uint32(0), // must be > 0
			expectedErr:       errorsmod.Wrap(types.ErrInvalidInput, constants.ErrorMsgInvalidMinPriceChange).Error(),
		},
		"Invalid min price change: ten thousand": {
			targetId:          0,
			pair:              constants.BtcUsdPair,
			minExchanges:      uint32(2),
			minPriceChangePpm: uint32(10_000), // must be < 10,000
			expectedErr:       errorsmod.Wrap(types.ErrInvalidInput, constants.ErrorMsgInvalidMinPriceChange).Error(),
		},
		"Min exchanges cannot be zero": {
			pair:              constants.BtcUsdPair,
			minExchanges:      uint32(0), // cannot be zero
			minPriceChangePpm: uint32(50),
			expectedErr:       types.ErrZeroMinExchanges.Error(),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, keeper, _, _, _, _ := keepertest.PricesKeepers(t)
			ctx = ctx.WithTxBytes(constants.TestTxBytes)
			keepertest.CreateNMarkets(t, ctx, keeper, 1)
			_, err := keeper.ModifyMarketParam(
				ctx,
				types.MarketParam{
					Id:                tc.targetId,
					Pair:              tc.pair,
					MinExchanges:      tc.minExchanges,
					MinPriceChangePpm: tc.minPriceChangePpm,
				},
			)
			require.EqualError(t, err, tc.expectedErr)
		})
	}
}

func TestGetMarketParam(t *testing.T) {
	ctx, keeper, _, _, _, _ := keepertest.PricesKeepers(t)
	items := keepertest.CreateNMarkets(t, ctx, keeper, 10)
	for _, item := range items {
		rst, err := keeper.GetMarketParam(ctx, item.Param.Id)
		require.NoError(t, err)
		require.Equal(
			t,
			&item.Param,
			&rst,
		)
	}
}

func TestGetMarketParam_NotFound(t *testing.T) {
	ctx, keeper, _, _, _, _ := keepertest.PricesKeepers(t)
	_, err := keeper.GetMarketParam(ctx, uint32(0))
	require.EqualError(t, err, "0: Market param does not exist")
}

func TestGetAllMarketParams(t *testing.T) {
	ctx, keeper, _, _, _, _ := keepertest.PricesKeepers(t)
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
