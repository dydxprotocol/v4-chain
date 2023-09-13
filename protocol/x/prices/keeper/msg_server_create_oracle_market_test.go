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

func TestCreateOracleMarket(t *testing.T) {
	testMarket1 := *pricestest.GenerateMarketParamPrice(
		pricestest.WithId(1),
		pricestest.WithPair("BTC-USD"),
		pricestest.WithExponent(-8), // for both Param and Price
		pricestest.WithPriceValue(0),
	)
	testCases := map[string]struct {
		setup           func(*testing.T, sdk.Context, *keeper.Keeper)
		msg             *pricestypes.MsgCreateOracleMarket
		expectedMarkets []pricestypes.MarketParamPrice
		expectedErr     string
	}{
		"Succeeds: create new oracle market (id = 1)": {
			setup: func(t *testing.T, ctx sdk.Context, pricesKeeper *keeper.Keeper) {},
			msg: &pricestypes.MsgCreateOracleMarket{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Params:    testMarket1.Param,
			},
			expectedMarkets: []pricestypes.MarketParamPrice{testMarket1},
		},
		"Failure: empty pair": {
			setup: func(t *testing.T, ctx sdk.Context, pricesKeeper *keeper.Keeper) {},
			msg: &pricestypes.MsgCreateOracleMarket{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Params: pricestest.GenerateMarketParamPrice(
					pricestest.WithPair(""),
					pricestest.WithExponent(-8), // for both Param and Price
				).Param,
			},
			expectedMarkets: []pricestypes.MarketParamPrice{},
			expectedErr:     "Pair cannot be empty",
		},
		"Failure: typo in exchange config json": {
			setup: func(t *testing.T, ctx sdk.Context, pricesKeeper *keeper.Keeper) {},
			msg: &pricestypes.MsgCreateOracleMarket{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Params: pricestest.GenerateMarketParamPrice(
					pricestest.WithPair("BTC-USD"),
					pricestest.WithExponent(-8), // for both Param and Price
					pricestest.WithExchangeConfigJson(`{"exchanges":[{"exchangeName":"Binance"""}]}`),
				).Param,
			},
			expectedMarkets: []pricestypes.MarketParamPrice{},
			expectedErr:     "ExchangeConfigJson string is not valid",
		},
		"Failure: oracle market id already exists": {
			setup: func(t *testing.T, ctx sdk.Context, pricesKeeper *keeper.Keeper) {
				keepertest.CreateTestPriceMarkets(
					t,
					ctx,
					pricesKeeper,
					[]pricestypes.MarketParamPrice{testMarket1},
				)
			},
			msg: &pricestypes.MsgCreateOracleMarket{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Params: pricestest.GenerateMarketParamPrice(
					pricestest.WithId(1), // same id as testMarket1
				).Param,
			},
			expectedErr:     "Market params already exists",
			expectedMarkets: []pricestypes.MarketParamPrice{testMarket1},
		},
		"Failure: invalid authority": {
			setup: func(t *testing.T, ctx sdk.Context, pricesKeeper *keeper.Keeper) {},
			msg: &pricestypes.MsgCreateOracleMarket{
				Authority: "invalid",
				Params:    testMarket1.Param,
			},
			expectedMarkets: []pricestypes.MarketParamPrice{},
			expectedErr:     "invalid authority",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctx, pricesKeeper, _, _, _, mockTimeProvider := keepertest.PricesKeepers(t)
			mockTimeProvider.On("Now").Return(constants.TimeT)
			msgServer := keeper.NewMsgServerImpl(pricesKeeper)
			goCtx := sdk.WrapSDKContext(ctx)
			tc.setup(t, ctx, pricesKeeper)

			_, err := msgServer.CreateOracleMarket(goCtx, tc.msg)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
			gotAllMarketParamPrices, err := pricesKeeper.GetAllMarketParamPrices(ctx)
			require.NoError(t, err)
			require.Equal(t, tc.expectedMarkets, gotAllMarketParamPrices)
		})
	}
}
