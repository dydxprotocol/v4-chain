package keeper_test

import (
	"math/big"
	"testing"

	testapp "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/app"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/keeper"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sDAIOracle/api"

	pricetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
)

func TestMsgUpdateSDAIConversionRateInitial(t *testing.T) {

	testCases := []struct {
		name             string
		input            *types.MsgUpdateSDAIConversionRate
		expectedSDAIRate string
		expErr           bool
	}{
		{
			name: "Valid input",
			input: &types.MsgUpdateSDAIConversionRate{
				Sender:              "cosmos139f7kncmglres2nf3h4hc4tade85ekfr8sulz5",
				ConversionRate:      "1",
				EthereumBlockNumber: "1",
			},
			expectedSDAIRate: "1",
			expErr:           false,
		},
		{
			name: "Invalid conversion rate doesn't match local server",
			input: &types.MsgUpdateSDAIConversionRate{
				Sender:              "cosmos139f7kncmglres2nf3h4hc4tade85ekfr8sulz5",
				ConversionRate:      "2",
				EthereumBlockNumber: "1",
			},
			expectedSDAIRate: "1",
			expErr:           true,
		},
		{
			name: "Invalid ethereum block number doesn't match local server",
			input: &types.MsgUpdateSDAIConversionRate{
				Sender:              "cosmos139f7kncmglres2nf3h4hc4tade85ekfr8sulz5",
				ConversionRate:      "1",
				EthereumBlockNumber: "2",
			},
			expectedSDAIRate: "1",
			expErr:           true,
		},
		{
			name: "Invalid conversion rate (empty)",
			input: &types.MsgUpdateSDAIConversionRate{
				Sender:              "cosmos139f7kncmglres2nf3h4hc4tade85ekfr8sulz5",
				ConversionRate:      "",
				EthereumBlockNumber: "1",
			},
			expectedSDAIRate: "",
			expErr:           true,
		},
		{
			name: "Invalid ethereum block number (empty)",
			input: &types.MsgUpdateSDAIConversionRate{
				Sender:              "cosmos139f7kncmglres2nf3h4hc4tade85ekfr8sulz5",
				ConversionRate:      "1",
				EthereumBlockNumber: "",
			},
			expectedSDAIRate: "",
			expErr:           true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.RatelimitKeeper
			ms := keeper.NewMsgServerImpl(k)

			k.SetSDAIPrice(ctx, big.NewInt(0))

			sDAIEventManager := k.GetSDAIEventManagerForTestingOnly()

			sDAIEventManager.AddsDAIEvent(&api.AddsDAIEventsRequest{
				ConversionRate:      "1",
				EthereumBlockNumber: "1",
			})

			_, err := ms.UpdateSDAIConversionRate(ctx, tc.input)
			if tc.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				sdkCtx := sdk.UnwrapSDKContext(ctx)
				price, found := k.GetSDAIPrice(sdkCtx)
				require.True(t, found)
				require.Equal(t,
					tc.expectedSDAIRate,
					price.String(),
				)
			}
		})
	}
}

func TestMsgUpdateSDAIConversionRatePostFirstEpoch(t *testing.T) {

	testCases := []struct {
		name                string
		input               *types.MsgUpdateSDAIConversionRate
		daiYieldEpochParams types.DaiYieldEpochParams
		epoch               int64
		expectedSDAIRate    string
		expErr              bool
	}{
		{
			name: "Valid input",
			input: &types.MsgUpdateSDAIConversionRate{
				Sender:              "cosmos139f7kncmglres2nf3h4hc4tade85ekfr8sulz5",
				ConversionRate:      "1",
				EthereumBlockNumber: "1",
			},
			daiYieldEpochParams: types.DaiYieldEpochParams{
				TradingDaiMinted:               "0",
				TotalTradingDaiPreMint:         "0",
				TotalTradingDaiClaimedForEpoch: "0",
				BlockNumber:                    0,
				EpochMarketPrices:              []*pricetypes.MarketPrice{},
			},
			expectedSDAIRate: "1",
			epoch:            1,
			expErr:           false,
		},
		{
			name: "Invalid epoch hasnt elapsed",
			input: &types.MsgUpdateSDAIConversionRate{
				Sender:              "cosmos139f7kncmglres2nf3h4hc4tade85ekfr8sulz5",
				ConversionRate:      "1",
				EthereumBlockNumber: "1",
			},
			daiYieldEpochParams: types.DaiYieldEpochParams{
				TradingDaiMinted:               "0",
				TotalTradingDaiPreMint:         "0",
				TotalTradingDaiClaimedForEpoch: "0",
				BlockNumber:                    100,
				EpochMarketPrices:              []*pricetypes.MarketPrice{},
			},
			expectedSDAIRate: "1",
			epoch:            1,
			expErr:           true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.RatelimitKeeper
			ms := keeper.NewMsgServerImpl(k)

			k.SetSDAIPrice(ctx, big.NewInt(0))
			k.SetCurrentDaiYieldEpochNumber(ctx, big.NewInt(tc.epoch))
			k.SetDaiYieldEpochParams(ctx, uint64(tc.epoch), tc.daiYieldEpochParams)

			ctx = ctx.WithBlockHeight(110)

			sDAIEventManager := k.GetSDAIEventManagerForTestingOnly()

			sDAIEventManager.AddsDAIEvent(&api.AddsDAIEventsRequest{
				ConversionRate:      "1",
				EthereumBlockNumber: "1",
			})

			_, err := ms.UpdateSDAIConversionRate(ctx, tc.input)
			if tc.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				sdkCtx := sdk.UnwrapSDKContext(ctx)
				price, found := k.GetSDAIPrice(sdkCtx)
				require.True(t, found)
				require.Equal(t,
					tc.expectedSDAIRate,
					price.String(),
				)
			}
		})
	}
}
