package keeper_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/slinky"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	pricestest "github.com/dydxprotocol/v4-chain/protocol/testutil/prices"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/keeper"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/require"
)

func TestUpdateMarketParam(t *testing.T) {
	testMarket := *pricestest.GenerateMarketParamPrice(
		pricestest.WithId(1),
		pricestest.WithPair("BTC-USD"),
		pricestest.WithExponent(-8),
		pricestest.WithPriceValue(0),
	)
	testMarketParam, testMarketPrice := testMarket.Param, testMarket.Price

	tests := map[string]struct {
		msg         *pricestypes.MsgUpdateMarketParam
		expectedErr string
	}{
		"Succeeds: update all parameters except exponent and pair": {
			msg: &pricestypes.MsgUpdateMarketParam{
				Authority: lib.GovModuleAddress.String(),
				MarketParam: pricestypes.MarketParam{
					Id:                 testMarketParam.Id,
					Pair:               testMarketParam.Pair,
					Exponent:           testMarketParam.Exponent,
					MinExchanges:       72,
					MinPriceChangePpm:  2_023,
					ExchangeConfigJson: `{"exchanges":[{"exchangeName":"XYZ","ticker":"PIKACHU"}]}`,
				},
			},
		},
		"Succeeds: update pair name": {
			msg: &pricestypes.MsgUpdateMarketParam{
				Authority: lib.GovModuleAddress.String(),
				MarketParam: pricestypes.MarketParam{
					Id:                 testMarketParam.Id,
					Pair:               "NEWMARKET-USD",
					Exponent:           testMarketParam.Exponent,
					MinExchanges:       72,
					MinPriceChangePpm:  2_023,
					ExchangeConfigJson: `{"exchanges":[{"exchangeName":"XYZ","ticker":"PIKACHU"}]}`,
				},
			},
		},
		"Succeeds: update min price change ppm only": {
			msg: &pricestypes.MsgUpdateMarketParam{
				Authority: lib.GovModuleAddress.String(),
				MarketParam: pricestypes.MarketParam{
					Id:                 testMarketParam.Id,
					Pair:               testMarketParam.Pair,
					Exponent:           testMarketParam.Exponent,
					MinExchanges:       testMarketParam.MinExchanges,
					MinPriceChangePpm:  4_321,
					ExchangeConfigJson: testMarketParam.ExchangeConfigJson,
				},
			},
		},
		"Failure: update to an empty pair": {
			msg: &pricestypes.MsgUpdateMarketParam{
				Authority: lib.GovModuleAddress.String(),
				MarketParam: pricestypes.MarketParam{
					Id:                 testMarketParam.Id,
					Pair:               "", // invalid
					Exponent:           testMarketParam.Exponent,
					MinExchanges:       testMarketParam.MinExchanges,
					MinPriceChangePpm:  testMarketParam.MinPriceChangePpm,
					ExchangeConfigJson: testMarketParam.ExchangeConfigJson,
				},
			},
			expectedErr: "Pair cannot be empty",
		},
		"Failure: update to 0 min price change ppm": {
			msg: &pricestypes.MsgUpdateMarketParam{
				Authority: lib.GovModuleAddress.String(),
				MarketParam: pricestypes.MarketParam{
					Id:                 testMarketParam.Id,
					Pair:               testMarketParam.Pair,
					Exponent:           testMarketParam.Exponent,
					MinExchanges:       testMarketParam.MinExchanges,
					MinPriceChangePpm:  0, // invalid
					ExchangeConfigJson: testMarketParam.ExchangeConfigJson,
				},
			},
			expectedErr: "Invalid input",
		},
		"Failure: new pair name does not exist in marketmap": {
			msg: &pricestypes.MsgUpdateMarketParam{
				Authority: lib.GovModuleAddress.String(),
				MarketParam: pricestypes.MarketParam{
					Id:                 testMarketParam.Id,
					Pair:               "nonexistent-pair",
					Exponent:           testMarketParam.Exponent,
					MinExchanges:       testMarketParam.MinExchanges,
					MinPriceChangePpm:  testMarketParam.MinPriceChangePpm,
					ExchangeConfigJson: "{}",
				},
			},
			expectedErr: "NONEXISTENT/PAIR: Ticker not found in market map",
		},
		"Failure: empty authority": {
			msg: &pricestypes.MsgUpdateMarketParam{
				Authority: "",
				MarketParam: pricestypes.MarketParam{
					Id: testMarketParam.Id,
				},
			},
			expectedErr: "invalid authority",
		},
		"Failure: non-gov authority": {
			msg: &pricestypes.MsgUpdateMarketParam{
				Authority: constants.BobAccAddress.String(),
				MarketParam: pricestypes.MarketParam{
					Id: testMarketParam.Id,
				},
			},
			expectedErr: "invalid authority",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, pricesKeeper, _, _, mockTimeProvider, _, marketMapKeeper := keepertest.PricesKeepers(t)
			mockTimeProvider.On("Now").Return(constants.TimeT)
			msgServer := keeper.NewMsgServerImpl(pricesKeeper)
			initialMarketParam, err := keepertest.CreateTestMarket(t, ctx, pricesKeeper, testMarketParam, testMarketPrice)
			require.NoError(t, err)

			// Create new pair in marketmap if test is expected to succeed
			if (initialMarketParam.Pair != tc.msg.MarketParam.Pair) && tc.expectedErr == "" {
				keepertest.CreateMarketsInMarketMapFromParams(
					t,
					ctx,
					marketMapKeeper,
					[]pricestypes.MarketParam{tc.msg.MarketParam},
				)
			}

			_, err = msgServer.UpdateMarketParam(ctx, tc.msg)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
				// Verify that market param was not updated.
				marketParam, exists := pricesKeeper.GetMarketParam(ctx, tc.msg.MarketParam.Id)
				require.True(t, exists)
				require.Equal(t, initialMarketParam, marketParam)
			} else {
				require.NoError(t, err)
				// Verify that market param was updated.
				updatedMarketParam, exists := pricesKeeper.GetMarketParam(ctx, tc.msg.MarketParam.Id)
				require.True(t, exists)
				require.Equal(t, tc.msg.MarketParam, updatedMarketParam)

				// If pair name changed, verify that old pair is disabled in the marketmap and new pair is enabled
				if initialMarketParam.Pair != updatedMarketParam.Pair {
					oldCp, err := slinky.MarketPairToCurrencyPair(initialMarketParam.Pair)
					require.NoError(t, err)
					oldMarket, err := marketMapKeeper.GetMarket(ctx, oldCp.String())
					require.NoError(t, err)
					require.False(t, oldMarket.Ticker.Enabled)

					newCp, err := slinky.MarketPairToCurrencyPair(updatedMarketParam.Pair)
					require.NoError(t, err)
					market, err := marketMapKeeper.GetMarket(ctx, newCp.String())
					require.NoError(t, err)
					require.True(t, market.Ticker.Enabled)
				}
			}
		})
	}
}
