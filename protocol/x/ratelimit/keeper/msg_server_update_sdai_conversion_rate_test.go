package keeper_test

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	testapp "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/app"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/keeper"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sDAIOracle/api"
	sdaiserver "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/sDAIOracle"
)

func TestMsgUpdateSDAIConversionRate_Initial(t *testing.T) {

	testCases := []struct {
		name             string
		input            *types.MsgUpdateSDAIConversionRate
		expectedSDAIRate string
		expErr           bool
	}{
		{
			name: "Valid input: basic",
			input: &types.MsgUpdateSDAIConversionRate{
				Sender:         "cosmos139f7kncmglres2nf3h4hc4tade85ekfr8sulz5",
				ConversionRate: "1",
			},
			expectedSDAIRate: "1",
			expErr:           false,
		},
		{
			name: "Invalid conversion rate doesn't match local server",
			input: &types.MsgUpdateSDAIConversionRate{
				Sender:         "cosmos139f7kncmglres2nf3h4hc4tade85ekfr8sulz5",
				ConversionRate: "2",
			},
			expectedSDAIRate: "",
			expErr:           true,
		},
		{
			name: "Invalid conversion rate (empty)",
			input: &types.MsgUpdateSDAIConversionRate{
				Sender:         "cosmos139f7kncmglres2nf3h4hc4tade85ekfr8sulz5",
				ConversionRate: "",
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

			sDAIEventManager := k.GetSDAIEventManagerForTestingOnly()

			sDAIEventManager.AddsDAIEvent(&api.AddsDAIEventsRequest{
				ConversionRate: "1",
			})

			_, err := ms.UpdateSDAIConversionRate(ctx, tc.input)
			sdkCtx := sdk.UnwrapSDKContext(ctx)
			conversionRate, conversionRateFound := k.GetSDAIPrice(sdkCtx)
			assetYieldIndex, assetYieldIndexFound := tApp.App.RatelimitKeeper.GetAssetYieldIndex(ctx)

			if tc.expErr {
				require.Error(t, err)
				require.False(t, conversionRateFound)
				require.Nil(t, conversionRate)
				require.False(t, assetYieldIndexFound)
				require.Nil(t, assetYieldIndex)
			} else {
				require.NoError(t, err)
				require.True(t, conversionRateFound)
				require.Equal(t, tc.expectedSDAIRate, conversionRate.String())
				require.True(t, assetYieldIndexFound)
				require.Equal(t, 0, big.NewRat(0, 1).Cmp(assetYieldIndex))
			}
		})
	}
}

func TestMsgUpdateSDAIConversionRate_PostFirstEpoch(t *testing.T) {

	testCases := []struct {
		name                    string
		input                   *types.MsgUpdateSDAIConversionRate
		epoch                   uint64
		expectedSDAIRate        string
		expectedAssetYieldIndex string
		setup                   func(sdk.Context, *testapp.TestApp, *sdaiserver.SDAIEventManager)
		expErr                  bool
	}{
		{
			name: "Valid input: basic",
			input: &types.MsgUpdateSDAIConversionRate{
				Sender:         "cosmos139f7kncmglres2nf3h4hc4tade85ekfr8sulz5",
				ConversionRate: "1",
			},
			expectedSDAIRate:        "1",
			expectedAssetYieldIndex: "0/1",
			epoch:                   uint64(1),
			setup: func(ctx sdk.Context, tApp *testapp.TestApp, sDAIEventManager *sdaiserver.SDAIEventManager) {
				sDAIEventManager.AddsDAIEvent(&api.AddsDAIEventsRequest{
					ConversionRate: "1",
				})
			},
			expErr: false,
		},
		{
			name: "Valid input: multiple conversion rates from daemon",
			input: &types.MsgUpdateSDAIConversionRate{
				Sender:         "cosmos139f7kncmglres2nf3h4hc4tade85ekfr8sulz5",
				ConversionRate: "10",
			},
			expectedSDAIRate:        "10",
			expectedAssetYieldIndex: "0/1",
			epoch:                   uint64(1),
			setup: func(ctx sdk.Context, tApp *testapp.TestApp, sDAIEventManager *sdaiserver.SDAIEventManager) {
				for i := 1; i <= 10; i++ {
					sDAIEventManager.AddsDAIEvent(&api.AddsDAIEventsRequest{
						ConversionRate: fmt.Sprintf("%d", i),
					})
				}
			},
			expErr: false,
		},
		{
			name: "Valid input with minted DAI",
			input: &types.MsgUpdateSDAIConversionRate{
				Sender:         "cosmos139f7kncmglres2nf3h4hc4tade85ekfr8sulz5",
				ConversionRate: "1" + strings.Repeat("0", 27),
			},
			expectedSDAIRate:        "1" + strings.Repeat("0", 27),
			expectedAssetYieldIndex: "1/1",
			epoch:                   uint64(1),
			setup: func(ctx sdk.Context, tApp *testapp.TestApp, sDAIEventManager *sdaiserver.SDAIEventManager) {
				for i := 0; i < 8; i++ {
					sDAIEventManager.AddsDAIEvent(&api.AddsDAIEventsRequest{
						ConversionRate: "1" + strings.Repeat("0", 26) + fmt.Sprintf("%d", i),
					})
				}
				burnAllCoinsOfDenom(t, ctx, tApp, types.TDaiDenom)
				burnAllCoinsOfDenom(t, ctx, tApp, types.SDaiDenom)

				sDaiCoins := sdk.NewCoins(sdk.NewCoin(types.SDaiDenom, sdkmath.NewInt(2000000000000)))
				err := tApp.App.BankKeeper.MintCoins(
					ctx,
					types.TDaiPoolAccount,
					sDaiCoins,
				)
				require.NoError(t, err)
				err = tApp.App.BankKeeper.SendCoinsFromModuleToModule(
					ctx,
					types.TDaiPoolAccount,
					types.SDaiPoolAccount,
					sDaiCoins,
				)
				require.NoError(t, err)
				tDaiCoins := sdk.NewCoins(sdk.NewCoin(types.TDaiDenom, sdkmath.NewInt(1)))
				err = tApp.App.BankKeeper.MintCoins(
					ctx,
					types.TDaiPoolAccount,
					tDaiCoins,
				)
				require.NoError(t, err)
			},
			expErr: false,
		},
		{
			name: "Invalid: no conversion rate match with daemon",
			input: &types.MsgUpdateSDAIConversionRate{
				Sender:         "cosmos139f7kncmglres2nf3h4hc4tade85ekfr8sulz5",
				ConversionRate: "11",
			},
			expectedSDAIRate:        "0",
			expectedAssetYieldIndex: "0/1",
			epoch:                   uint64(1),
			setup: func(ctx sdk.Context, tApp *testapp.TestApp, sDAIEventManager *sdaiserver.SDAIEventManager) {
				for i := 1; i <= 10; i++ {
					sDAIEventManager.AddsDAIEvent(&api.AddsDAIEventsRequest{
						ConversionRate: fmt.Sprintf("%d", i),
					})
				}
			},
			expErr: true,
		},
		{
			name: "Invalid: updating conversion rate fails due to minting",
			input: &types.MsgUpdateSDAIConversionRate{
				Sender:         "cosmos139f7kncmglres2nf3h4hc4tade85ekfr8sulz5",
				ConversionRate: "10",
			},
			expectedSDAIRate:        "10",
			expectedAssetYieldIndex: "0/1",
			epoch:                   uint64(1),
			setup: func(ctx sdk.Context, tApp *testapp.TestApp, sDAIEventManager *sdaiserver.SDAIEventManager) {
				for i := 1; i <= 10; i++ {
					sDAIEventManager.AddsDAIEvent(&api.AddsDAIEventsRequest{
						ConversionRate: fmt.Sprintf("%d", i),
					})
				}
				burnAllCoinsOfDenom(t, ctx, tApp, types.TDaiDenom)
				burnAllCoinsOfDenom(t, ctx, tApp, types.SDaiDenom)

				coins := sdk.NewCoins(sdk.NewCoin(types.SDaiDenom, sdkmath.NewInt(100)))
				err := tApp.App.BankKeeper.MintCoins(
					ctx,
					types.TDaiPoolAccount,
					coins,
				)
				require.NoError(t, err)
				err = tApp.App.BankKeeper.SendCoinsFromModuleToModule(
					ctx,
					types.TDaiPoolAccount,
					types.SDaiPoolAccount,
					coins,
				)
				require.NoError(t, err)
			},
			expErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.RatelimitKeeper
			ms := keeper.NewMsgServerImpl(k)

			k.SetSDAIPrice(ctx, big.NewInt(0))
			k.SetAssetYieldIndex(ctx, big.NewRat(0, 1))

			ctx = ctx.WithBlockHeight(110)
			sDAIEventManager := k.GetSDAIEventManagerForTestingOnly()

			tc.setup(ctx, tApp, sDAIEventManager)

			_, err := ms.UpdateSDAIConversionRate(ctx, tc.input)
			sdkCtx := sdk.UnwrapSDKContext(ctx)
			if tc.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			conversionRate, conversionRateFound := k.GetSDAIPrice(sdkCtx)
			assetYieldIndex, assetYieldIndexFound := tApp.App.RatelimitKeeper.GetAssetYieldIndex(ctx)
			require.True(t, conversionRateFound)
			require.Equal(t, tc.expectedSDAIRate, conversionRate.String())
			require.True(t, assetYieldIndexFound)
			require.Equal(t, tc.expectedAssetYieldIndex, assetYieldIndex.String())
		})
	}
}

func TestMsgUpdateSDAIConversionRate_PerformsAllStateChanges(t *testing.T) {

	testCases := []struct {
		name                     string
		input                    *types.MsgUpdateSDAIConversionRate
		epoch                    uint64
		expectedSDAIRate         string
		expectedAssetYieldIndex  string
		expectedPerpYieldIndexes []string
		setup                    func(sdk.Context, *testapp.TestApp, *sdaiserver.SDAIEventManager)
		expErr                   bool
	}{
		{
			name: "Valid input with minted DAI",
			input: &types.MsgUpdateSDAIConversionRate{
				Sender:         "cosmos139f7kncmglres2nf3h4hc4tade85ekfr8sulz5",
				ConversionRate: "1" + strings.Repeat("0", 27),
			},
			expectedSDAIRate:         "1" + strings.Repeat("0", 27),
			expectedAssetYieldIndex:  "1/1",
			expectedPerpYieldIndexes: []string{"2/1", "1/1"},
			epoch:                    uint64(1),
			setup: func(ctx sdk.Context, tApp *testapp.TestApp, sDAIEventManager *sdaiserver.SDAIEventManager) {
				for i := 0; i < 8; i++ {
					sDAIEventManager.AddsDAIEvent(&api.AddsDAIEventsRequest{
						ConversionRate: "1" + strings.Repeat("0", 26) + fmt.Sprintf("%d", i),
					})
				}
				burnAllCoinsOfDenom(t, ctx, tApp, types.TDaiDenom)
				burnAllCoinsOfDenom(t, ctx, tApp, types.SDaiDenom)

				sDaiCoins := sdk.NewCoins(sdk.NewCoin(types.SDaiDenom, sdkmath.NewInt(2000000000000)))
				err := tApp.App.BankKeeper.MintCoins(
					ctx,
					types.TDaiPoolAccount,
					sDaiCoins,
				)
				require.NoError(t, err)
				err = tApp.App.BankKeeper.SendCoinsFromModuleToModule(
					ctx,
					types.TDaiPoolAccount,
					types.SDaiPoolAccount,
					sDaiCoins,
				)
				require.NoError(t, err)
				tDaiCoins := sdk.NewCoins(sdk.NewCoin(types.TDaiDenom, sdkmath.NewInt(1)))
				err = tApp.App.BankKeeper.MintCoins(
					ctx,
					types.TDaiPoolAccount,
					tDaiCoins,
				)
				require.NoError(t, err)
			},
			expErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.RatelimitKeeper
			ms := keeper.NewMsgServerImpl(k)

			k.SetSDAIPrice(ctx, big.NewInt(0))
			k.SetAssetYieldIndex(ctx, big.NewRat(0, 1))

			ctx = ctx.WithBlockHeight(110)
			sDAIEventManager := k.GetSDAIEventManagerForTestingOnly()

			tc.setup(ctx, tApp, sDAIEventManager)

			_, err := ms.UpdateSDAIConversionRate(ctx, tc.input)
			sdkCtx := sdk.UnwrapSDKContext(ctx)

			if tc.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			conversionRate, conversionRateFound := k.GetSDAIPrice(sdkCtx)
			assetYieldIndex, assetYieldIndexFound := tApp.App.RatelimitKeeper.GetAssetYieldIndex(ctx)
			require.True(t, conversionRateFound)
			require.Equal(t, tc.expectedSDAIRate, conversionRate.String())
			require.True(t, assetYieldIndexFound)
			require.Equal(t, tc.expectedAssetYieldIndex, assetYieldIndex.String())

			allPerps := tApp.App.PerpetualsKeeper.GetAllPerpetuals(ctx)
			for i, perp := range allPerps {
				require.Equal(t, tc.expectedPerpYieldIndexes[i], perp.YieldIndex)
			}
		})
	}
}
