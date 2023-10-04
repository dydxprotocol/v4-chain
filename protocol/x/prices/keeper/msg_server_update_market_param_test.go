package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
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
		"Succeeds: update all parameters except exponent": {
			msg: &pricestypes.MsgUpdateMarketParam{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				MarketParam: pricestypes.MarketParam{
					Id:                 testMarketParam.Id,
					Pair:               "PIKACHU-XXX",
					Exponent:           testMarketParam.Exponent,
					MinExchanges:       72,
					MinPriceChangePpm:  2_023,
					ExchangeConfigJson: `{"exchanges":[{"exchangeName":"XYZ","ticker":"PIKACHU"}]}`,
				},
			},
		},
		"Succeeds: update min price change ppm only": {
			msg: &pricestypes.MsgUpdateMarketParam{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
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
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
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
		"Failure: update to 0 min exchanges": {
			msg: &pricestypes.MsgUpdateMarketParam{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				MarketParam: pricestypes.MarketParam{
					Id:                 testMarketParam.Id,
					Pair:               testMarketParam.Pair,
					Exponent:           testMarketParam.Exponent,
					MinExchanges:       0, // invalid
					MinPriceChangePpm:  testMarketParam.MinPriceChangePpm,
					ExchangeConfigJson: testMarketParam.ExchangeConfigJson,
				},
			},
			expectedErr: "Min exchanges must be greater than zero",
		},
		"Failure: update to 0 min price change ppm": {
			msg: &pricestypes.MsgUpdateMarketParam{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
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
		"Failure: update to invalid exchange config json": {
			msg: &pricestypes.MsgUpdateMarketParam{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				MarketParam: pricestypes.MarketParam{
					Id:                 testMarketParam.Id,
					Pair:               testMarketParam.Pair,
					Exponent:           testMarketParam.Exponent,
					MinExchanges:       testMarketParam.MinExchanges,
					MinPriceChangePpm:  testMarketParam.MinPriceChangePpm,
					ExchangeConfigJson: `{{"exchanges":[{"exchangeName":"XYZ","ticker":"PIKACHU"}]}`, // invalid json
				},
			},
			expectedErr: "Invalid input",
		},
		"Failure: update market exponent": {
			msg: &pricestypes.MsgUpdateMarketParam{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				MarketParam: pricestypes.MarketParam{
					Id:                 testMarketParam.Id,
					Pair:               testMarketParam.Pair,
					Exponent:           testMarketParam.Exponent + 1, // cannot be updated
					MinExchanges:       testMarketParam.MinExchanges,
					MinPriceChangePpm:  testMarketParam.MinPriceChangePpm,
					ExchangeConfigJson: "{}",
				},
			},
			expectedErr: "Market exponent cannot be updated",
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
			ctx, pricesKeeper, _, _, _, mockTimeProvider := keepertest.PricesKeepers(t)
			mockTimeProvider.On("Now").Return(constants.TimeT)
			msgServer := keeper.NewMsgServerImpl(pricesKeeper)
			goCtx := sdk.WrapSDKContext(ctx)
			initialMarketParam, err := pricesKeeper.CreateMarket(ctx, testMarketParam, testMarketPrice)
			require.NoError(t, err)

			_, err = msgServer.UpdateMarketParam(goCtx, tc.msg)
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
			}
		})
	}
}
