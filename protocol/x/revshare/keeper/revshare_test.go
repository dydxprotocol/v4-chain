package keeper_test

import (
	"math/big"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	affiliateskeeper "github.com/dydxprotocol/v4-chain/protocol/x/affiliates/keeper"
	affiliatetypes "github.com/dydxprotocol/v4-chain/protocol/x/affiliates/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/revshare/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
	"github.com/stretchr/testify/require"
)

func TestGetSetMarketMapperRevShareDetails(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()

	k := tApp.App.RevShareKeeper
	// Set the rev share details for a market
	marketId := uint32(42)
	setDetails := types.MarketMapperRevShareDetails{
		ExpirationTs: 1735707600,
	}
	k.SetMarketMapperRevShareDetails(ctx, marketId, setDetails)

	// Get the rev share details for the market
	getDetails, err := k.GetMarketMapperRevShareDetails(ctx, marketId)
	require.NoError(t, err)
	require.Equal(t, getDetails.ExpirationTs, setDetails.ExpirationTs)

	// Set expiration ts to 0
	setDetails.ExpirationTs = 0
	k.SetMarketMapperRevShareDetails(ctx, marketId, setDetails)

	getDetails, err = k.GetMarketMapperRevShareDetails(ctx, marketId)
	require.NoError(t, err)
	require.Equal(t, getDetails.ExpirationTs, setDetails.ExpirationTs)
}

func TestGetMarketMapperRevShareDetailsFailure(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.RevShareKeeper

	// Get the rev share details for non-existent market
	_, err := k.GetMarketMapperRevShareDetails(ctx, 42)
	require.ErrorIs(t, err, types.ErrMarketMapperRevShareDetailsNotFound)
}

func TestCreateNewMarketRevShare(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()

	k := tApp.App.RevShareKeeper

	// Set base rev share params
	err := k.SetMarketMapperRevenueShareParams(
		ctx, types.MarketMapperRevenueShareParams{
			Address:         constants.AliceAccAddress.String(),
			RevenueSharePpm: 100_000, // 10%
			ValidDays:       240,
		},
	)
	require.NoError(t, err)

	// Create a new market rev share
	marketId := uint32(42)
	k.CreateNewMarketRevShare(ctx, marketId)
	require.NoError(t, err)

	// Check if the market rev share exists
	details, err := k.GetMarketMapperRevShareDetails(ctx, marketId)
	require.NoError(t, err)

	// TODO: is this blocktime call deterministic?
	expectedExpirationTs := ctx.BlockTime().Unix() + 240*24*60*60
	require.Equal(t, details.ExpirationTs, uint64(expectedExpirationTs))
}

func TestGetMarketMapperRevenueShareForMarket(t *testing.T) {
	tests := map[string]struct {
		revShareParams     types.MarketMapperRevenueShareParams
		marketId           uint32
		expirationDelta    int64
		setRevShareDetails bool

		// expected
		expectedMarketMapperAddr sdk.AccAddress
		expectedRevenueSharePpm  uint32
		expectedErr              error
	}{
		"valid market": {
			revShareParams: types.MarketMapperRevenueShareParams{
				Address:         constants.AliceAccAddress.String(),
				RevenueSharePpm: 100_000, // 10%
				ValidDays:       0,
			},
			marketId:           42,
			expirationDelta:    10,
			setRevShareDetails: true,

			expectedMarketMapperAddr: constants.AliceAccAddress,
			expectedRevenueSharePpm:  100_000,
		},
		"invalid market": {
			revShareParams: types.MarketMapperRevenueShareParams{
				Address:         constants.AliceAccAddress.String(),
				RevenueSharePpm: 100_000, // 10%
				ValidDays:       0,
			},
			marketId:           42,
			setRevShareDetails: false,

			expectedErr: types.ErrMarketMapperRevShareDetailsNotFound,
		},
		// TODO (TRA-455): investigate why tApp blocktime doesn't translate to ctx.BlockTime()
		//"expired market rev share": {
		//	revShareParams: types.MarketMapperRevenueShareParams{
		//		Address:         constants.AliceAccAddress.String(),
		//		RevenueSharePpm: 100_000, // 10%
		//		ValidDays:       0,
		//	},
		//	marketId:           42,
		//	expirationDelta:    -10,
		//	setRevShareDetails: true,
		//
		//	expectedMarketMapperAddr: nil,
		//	expectedRevenueSharePpm:  0,
		//},
	}

	for name, tc := range tests {
		t.Run(
			name, func(t *testing.T) {
				tApp := testapp.NewTestAppBuilder(t).Build()
				ctx := tApp.InitChain()
				k := tApp.App.RevShareKeeper
				tApp.AdvanceToBlock(
					2, testapp.AdvanceToBlockOptions{BlockTime: time.Now()},
				)

				// Set base rev share params
				err := k.SetMarketMapperRevenueShareParams(ctx, tc.revShareParams)
				require.NoError(t, err)

				// Set market rev share details
				if tc.setRevShareDetails {
					k.SetMarketMapperRevShareDetails(
						ctx, tc.marketId, types.MarketMapperRevShareDetails{
							ExpirationTs: uint64(ctx.BlockTime().Unix() + tc.expirationDelta),
						},
					)
				}

				// Get the revenue share for the market
				marketMapperAddr, revenueSharePpm, err := k.GetMarketMapperRevenueShareForMarket(ctx, tc.marketId)
				if tc.expectedErr != nil {
					require.ErrorIs(t, err, tc.expectedErr)
				} else {
					require.NoError(t, err)
					require.Equal(t, tc.expectedMarketMapperAddr, marketMapperAddr)
					require.Equal(t, tc.expectedRevenueSharePpm, revenueSharePpm)
				}
			},
		)
	}
}

func TestValidateRevShareSafety(t *testing.T) {
	tests := map[string]struct {
		affiliateTiers             affiliatetypes.AffiliateTiers
		revShareConfig             types.UnconditionalRevShareConfig
		marketMapperRevShareParams types.MarketMapperRevenueShareParams
		expectedValid              bool
	}{
		"valid rev share config": {
			affiliateTiers: affiliatetypes.DefaultAffiliateTiers,
			revShareConfig: types.UnconditionalRevShareConfig{
				Configs: []types.UnconditionalRevShareConfig_RecipientConfig{
					{
						Address:  constants.AliceAccAddress.String(),
						SharePpm: 100_000, // 10%
					},
				},
			},
			marketMapperRevShareParams: types.MarketMapperRevenueShareParams{
				Address:         constants.AliceAccAddress.String(),
				RevenueSharePpm: 100_000, // 10%
				ValidDays:       0,
			},
			expectedValid: true,
		},
		"invalid rev share config - sum of shares > 100%": {
			affiliateTiers: affiliatetypes.DefaultAffiliateTiers,
			revShareConfig: types.UnconditionalRevShareConfig{
				Configs: []types.UnconditionalRevShareConfig_RecipientConfig{
					{
						Address:  constants.AliceAccAddress.String(),
						SharePpm: 100_000, // 10%
					},
					{
						Address:  constants.BobAccAddress.String(),
						SharePpm: 810_000, // 81%
					},
				},
			},
			marketMapperRevShareParams: types.MarketMapperRevenueShareParams{
				Address:         constants.AliceAccAddress.String(),
				RevenueSharePpm: 100_000, // 10%
				ValidDays:       0,
			},
			expectedValid: false,
		},
		"invalid rev share config - sum of shares + highest tier share > 100%": {
			affiliateTiers: affiliatetypes.AffiliateTiers{
				Tiers: []affiliatetypes.AffiliateTiers_Tier{
					{
						ReqReferredVolumeQuoteQuantums: 0,
						ReqStakedWholeCoins:            0,
						TakerFeeSharePpm:               50_000, // 5%
					},
					{
						ReqReferredVolumeQuoteQuantums: 1_000_000_000_000, // 1 million USDC
						ReqStakedWholeCoins:            200,               // 200 whole coins
						TakerFeeSharePpm:               800_000,           // 80%
					},
				},
			},
			revShareConfig: types.UnconditionalRevShareConfig{
				Configs: []types.UnconditionalRevShareConfig_RecipientConfig{
					{
						Address:  constants.AliceAccAddress.String(),
						SharePpm: 100_000, // 10%
					},
					{
						Address:  constants.BobAccAddress.String(),
						SharePpm: 100_000, // 10%
					},
				},
			},
			marketMapperRevShareParams: types.MarketMapperRevenueShareParams{
				Address:         constants.AliceAccAddress.String(),
				RevenueSharePpm: 100_000, // 10%
				ValidDays:       0,
			},
			expectedValid: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			_ = tApp.InitChain()
			k := tApp.App.RevShareKeeper

			valid := k.ValidateRevShareSafety(tc.affiliateTiers, tc.revShareConfig, tc.marketMapperRevShareParams)
			require.Equal(t, tc.expectedValid, valid)
		})
	}
}

func TestKeeper_GetAllRevShares_Valid(t *testing.T) {
	tests := []struct {
		name                              string
		revenueSharePpmNetFees            uint32
		revenueSharePpmTakerFees          uint32
		expectedAffiliateRevShares        int
		expectedUnconditionalRevShares    int
		expectedMarketMapperRevShares     int
		monthlyRollingTakerVolumeQuantums uint64
		setup                             func(tApp *testapp.TestApp, ctx sdk.Context,
			keeper *keeper.Keeper, affiliatesKeeper *affiliateskeeper.Keeper)
	}{
		{
			name: "Valid revenue share from affiliates, unconditional " +
				"rev shares and market mapper rev share",
			revenueSharePpmNetFees:            600_000, // 60%,
			revenueSharePpmTakerFees:          150_000, // 15%
			expectedAffiliateRevShares:        1,
			expectedUnconditionalRevShares:    2,
			expectedMarketMapperRevShares:     1,
			monthlyRollingTakerVolumeQuantums: 1_000_000_000_000, // 1 million USDC
			setup: func(tApp *testapp.TestApp, ctx sdk.Context, keeper *keeper.Keeper,
				affiliatesKeeper *affiliateskeeper.Keeper) {
				err := keeper.SetMarketMapperRevenueShareParams(ctx, types.MarketMapperRevenueShareParams{
					Address:         constants.AliceAccAddress.String(),
					RevenueSharePpm: 100_000, // 10%
					ValidDays:       1,
				})
				require.NoError(t, err)

				keeper.SetUnconditionalRevShareConfigParams(ctx, types.UnconditionalRevShareConfig{
					Configs: []types.UnconditionalRevShareConfig_RecipientConfig{
						{
							Address:  constants.BobAccAddress.String(),
							SharePpm: 200_000, // 20%
						},
						{
							Address:  constants.AliceAccAddress.String(),
							SharePpm: 300_000, // 30%
						},
					},
				})

				affiliatesKeeper.UpdateAffiliateTiers(ctx, affiliatetypes.DefaultAffiliateTiers)
				err = affiliatesKeeper.RegisterAffiliate(ctx, constants.AliceAccAddress.String(),
					constants.BobAccAddress.String())
				require.NoError(t, err)
			},
		},
		{
			name:                              "Valid revenue share with 30d volume greater than max 30d referral volume",
			revenueSharePpmNetFees:            600_000, // 60%,
			revenueSharePpmTakerFees:          0,       // 0%
			expectedAffiliateRevShares:        0,
			expectedUnconditionalRevShares:    2,
			expectedMarketMapperRevShares:     1,
			monthlyRollingTakerVolumeQuantums: types.Max30dRefereeVolumeQuantums + 1,
			setup: func(tApp *testapp.TestApp, ctx sdk.Context, keeper *keeper.Keeper,
				affiliatesKeeper *affiliateskeeper.Keeper) {
				err := keeper.SetMarketMapperRevenueShareParams(ctx, types.MarketMapperRevenueShareParams{
					Address:         constants.AliceAccAddress.String(),
					RevenueSharePpm: 100_000, // 10%
					ValidDays:       1,
				})
				require.NoError(t, err)
				keeper.SetUnconditionalRevShareConfigParams(ctx, types.UnconditionalRevShareConfig{
					Configs: []types.UnconditionalRevShareConfig_RecipientConfig{
						{
							Address:  constants.BobAccAddress.String(),
							SharePpm: 200_000, // 20%
						},
						{
							Address:  constants.AliceAccAddress.String(),
							SharePpm: 300_000, // 30%
						},
					},
				})

				affiliatesKeeper.UpdateAffiliateTiers(ctx, affiliatetypes.DefaultAffiliateTiers)
				err = affiliatesKeeper.RegisterAffiliate(ctx, constants.AliceAccAddress.String(),
					constants.BobAccAddress.String())
				require.NoError(t, err)
			},
		},
		{
			name:                              "Valid revenue share with no unconditional rev shares",
			revenueSharePpmNetFees:            100_000, // 10%,
			revenueSharePpmTakerFees:          150_000, // 15%
			expectedAffiliateRevShares:        1,
			expectedUnconditionalRevShares:    0,
			expectedMarketMapperRevShares:     1,
			monthlyRollingTakerVolumeQuantums: 1_000_000_000_000, // 1 million USDC
			setup: func(tApp *testapp.TestApp, ctx sdk.Context, keeper *keeper.Keeper,
				affiliatesKeeper *affiliateskeeper.Keeper) {
				err := keeper.SetMarketMapperRevenueShareParams(ctx, types.MarketMapperRevenueShareParams{
					Address:         constants.AliceAccAddress.String(),
					RevenueSharePpm: 100_000, // 10%
					ValidDays:       1,
				})
				require.NoError(t, err)

				affiliatesKeeper.UpdateAffiliateTiers(ctx, affiliatetypes.DefaultAffiliateTiers)
				err = affiliatesKeeper.RegisterAffiliate(ctx, constants.AliceAccAddress.String(),
					constants.BobAccAddress.String())
				require.NoError(t, err)
			},
		},
		{
			name:                              "Valid revenue share with no market mapper rev share",
			revenueSharePpmNetFees:            200_000, // 20%,
			revenueSharePpmTakerFees:          150_000, // 15%
			expectedAffiliateRevShares:        1,
			expectedUnconditionalRevShares:    1,
			expectedMarketMapperRevShares:     0,
			monthlyRollingTakerVolumeQuantums: 1_000_000_000_000, // 1 million USDC
			setup: func(tApp *testapp.TestApp, ctx sdk.Context, keeper *keeper.Keeper,
				affiliatesKeeper *affiliateskeeper.Keeper) {
				keeper.SetUnconditionalRevShareConfigParams(ctx, types.UnconditionalRevShareConfig{
					Configs: []types.UnconditionalRevShareConfig_RecipientConfig{
						{
							Address:  constants.BobAccAddress.String(),
							SharePpm: 200_000, // 20%
						},
					},
				})
				affiliatesKeeper.UpdateAffiliateTiers(ctx, affiliatetypes.DefaultAffiliateTiers)
				err := affiliatesKeeper.RegisterAffiliate(ctx, constants.AliceAccAddress.String(),
					constants.BobAccAddress.String())
				require.NoError(t, err)
			},
		},
		{
			name:                              "No rev shares",
			revenueSharePpmNetFees:            0, // 0%,
			revenueSharePpmTakerFees:          0, // 0%
			expectedAffiliateRevShares:        0,
			expectedUnconditionalRevShares:    0,
			expectedMarketMapperRevShares:     0,
			monthlyRollingTakerVolumeQuantums: 1_000_000_000_000, // 1 million USDC
			setup: func(tApp *testapp.TestApp, ctx sdk.Context, keeper *keeper.Keeper,
				affiliatesKeeper *affiliateskeeper.Keeper) {
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			marketId := uint32(1)
			fill := clobtypes.CreatePerpetualFillForProcess(
				constants.AliceAccAddress.String(),
				big.NewInt(10_000_000),
				constants.BobAccAddress.String(),
				big.NewInt(2_000_000),
				big.NewInt(100000),
				marketId,
				tc.monthlyRollingTakerVolumeQuantums,
			)

			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			keeper := tApp.App.RevShareKeeper
			affiliatesKeeper := tApp.App.AffiliatesKeeper
			if tc.setup != nil {
				tc.setup(tApp, ctx, &keeper, &affiliatesKeeper)
			}

			keeper.CreateNewMarketRevShare(ctx, marketId)

			revShares, err := keeper.GetAllRevShares(ctx, fill)

			require.NoError(t, err)
			affiliateRevShares := 0
			unconditionalRevShares := 0
			marketMapperRevShares := 0
			for _, revShare := range revShares {
				switch revShare.RevShareType {
				case types.REV_SHARE_TYPE_AFFILIATE:
					affiliateRevShares++
				case types.REV_SHARE_TYPE_UNCONDITIONAL:
					unconditionalRevShares++
				case types.REV_SHARE_TYPE_MARKET_MAPPER:
					marketMapperRevShares++
				}
			}
			require.Equal(t, tc.expectedAffiliateRevShares, affiliateRevShares)
			require.Equal(t, tc.expectedUnconditionalRevShares, unconditionalRevShares)
			require.Equal(t, tc.expectedMarketMapperRevShares, marketMapperRevShares)

			if tc.expectedAffiliateRevShares > 0 || tc.expectedUnconditionalRevShares > 0 ||
				tc.expectedMarketMapperRevShares > 0 {
				totalFees := new(big.Int).Add(fill.TakerFeeQuoteQuantums(), fill.MakerFeeQuoteQuantums())
				actualShare := new(big.Int)

				expectedShareFromNetFees := lib.BigMulPpm(totalFees, lib.BigI(int64(tc.revenueSharePpmNetFees)), false)
				expectedShareFromTakerFees := lib.BigMulPpm(fill.TakerFeeQuoteQuantums(),
					lib.BigI(int64(tc.revenueSharePpmTakerFees)), false)
				expectedShare := new(big.Int).Add(expectedShareFromNetFees, expectedShareFromTakerFees)
				for _, revShare := range revShares {
					actualShare.Add(actualShare, revShare.QuoteQuantums)
				}
				require.Equal(t, expectedShare, actualShare)
			}
		})
	}
}

func TestKeeper_GetAllRevShares_Invalid(t *testing.T) {
	marketId := uint32(1)
	tests := []struct {
		name                              string
		revenueSharePpmNetFees            uint32
		revenueSharePpmTakerFees          uint32
		expectedError                     error
		monthlyRollingTakerVolumeQuantums uint64
		setup                             func(tApp *testapp.TestApp, ctx sdk.Context,
			keeper *keeper.Keeper, affiliatesKeeper *affiliateskeeper.Keeper)
	}{
		{
			name:                              "Total fees shared exceeds net fees from all sources",
			revenueSharePpmNetFees:            950_000, // 95%,
			revenueSharePpmTakerFees:          150_000, // 15%
			expectedError:                     types.ErrTotalFeesSharedExceedsNetFees,
			monthlyRollingTakerVolumeQuantums: 1_000_000_000_000, // 1 million USDC
			setup: func(tApp *testapp.TestApp, ctx sdk.Context, keeper *keeper.Keeper,
				affiliatesKeeper *affiliateskeeper.Keeper) {
				err := keeper.SetMarketMapperRevenueShareParams(ctx, types.MarketMapperRevenueShareParams{
					Address:         constants.AliceAccAddress.String(),
					RevenueSharePpm: 800_000, // 80%
					ValidDays:       1,
				})
				require.NoError(t, err)

				keeper.SetUnconditionalRevShareConfigParams(ctx, types.UnconditionalRevShareConfig{
					Configs: []types.UnconditionalRevShareConfig_RecipientConfig{
						{
							Address:  constants.BobAccAddress.String(),
							SharePpm: 150_000, // 15%
						},
					},
				})
				affiliatesKeeper.UpdateAffiliateTiers(ctx, affiliatetypes.DefaultAffiliateTiers)
				err = affiliatesKeeper.RegisterAffiliate(ctx, constants.AliceAccAddress.String(),
					constants.BobAccAddress.String())
				require.NoError(t, err)
			},
		},
		{
			name:                              "Total fees shared exceeds net fees from market mapper and affiliates",
			revenueSharePpmNetFees:            950_000, // 95%,
			revenueSharePpmTakerFees:          150_000, // 15%
			expectedError:                     types.ErrTotalFeesSharedExceedsNetFees,
			monthlyRollingTakerVolumeQuantums: 1_000_000_000_000, // 1 million USDC
			setup: func(tApp *testapp.TestApp, ctx sdk.Context, keeper *keeper.Keeper,
				affiliatesKeeper *affiliateskeeper.Keeper) {
				err := keeper.SetMarketMapperRevenueShareParams(ctx, types.MarketMapperRevenueShareParams{
					Address:         constants.AliceAccAddress.String(),
					RevenueSharePpm: 950_000, // 95%
					ValidDays:       1,
				})
				require.NoError(t, err)

				affiliatesKeeper.UpdateAffiliateTiers(ctx, affiliatetypes.DefaultAffiliateTiers)
				err = affiliatesKeeper.RegisterAffiliate(ctx, constants.AliceAccAddress.String(),
					constants.BobAccAddress.String())
				require.NoError(t, err)
			},
		},
		{
			name:                              "Total fees shared exceeds net fees from affiliates and unconditional rev shares",
			revenueSharePpmNetFees:            950_000, // 95%,
			revenueSharePpmTakerFees:          150_000, // 15%
			expectedError:                     types.ErrTotalFeesSharedExceedsNetFees,
			monthlyRollingTakerVolumeQuantums: 1_000_000_000_000, // 1 million USDC
			setup: func(tApp *testapp.TestApp, ctx sdk.Context, keeper *keeper.Keeper,
				affiliatesKeeper *affiliateskeeper.Keeper) {
				keeper.SetUnconditionalRevShareConfigParams(ctx, types.UnconditionalRevShareConfig{
					Configs: []types.UnconditionalRevShareConfig_RecipientConfig{
						{
							Address:  constants.BobAccAddress.String(),
							SharePpm: 950_000, // 95%
						},
					},
				})

				affiliatesKeeper.UpdateAffiliateTiers(ctx, affiliatetypes.DefaultAffiliateTiers)
				err := affiliatesKeeper.RegisterAffiliate(ctx, constants.AliceAccAddress.String(), constants.BobAccAddress.String())
				require.NoError(t, err)
			},
		},
		{
			name:                              "Total fees shared exceeds net fees - no affiliate rev shares",
			revenueSharePpmNetFees:            1_150_000, // 115%,
			revenueSharePpmTakerFees:          0,         // 0%
			expectedError:                     types.ErrTotalFeesSharedExceedsNetFees,
			monthlyRollingTakerVolumeQuantums: 1_000_000_000_000, // 1 million USDC
			setup: func(tApp *testapp.TestApp, ctx sdk.Context, keeper *keeper.Keeper,
				affiliatesKeeper *affiliateskeeper.Keeper) {
				err := keeper.SetMarketMapperRevenueShareParams(ctx, types.MarketMapperRevenueShareParams{
					Address:         constants.AliceAccAddress.String(),
					RevenueSharePpm: 950_000, // 95%
					ValidDays:       1,
				})
				require.NoError(t, err)

				keeper.SetUnconditionalRevShareConfigParams(ctx, types.UnconditionalRevShareConfig{
					Configs: []types.UnconditionalRevShareConfig_RecipientConfig{
						{
							Address:  constants.BobAccAddress.String(),
							SharePpm: 200_000, // 20%
						},
					},
				})
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			fill := clobtypes.CreatePerpetualFillForProcess(
				constants.AliceAccAddress.String(),
				big.NewInt(10_000_000),
				constants.BobAccAddress.String(),
				big.NewInt(2_000_000),
				big.NewInt(100000),
				uint32(1),
				tc.monthlyRollingTakerVolumeQuantums,
			)
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			keeper := tApp.App.RevShareKeeper
			affiliatesKeeper := tApp.App.AffiliatesKeeper
			if tc.setup != nil {
				tc.setup(tApp, ctx, &keeper, &affiliatesKeeper)
			}

			keeper.CreateNewMarketRevShare(ctx, marketId)

			_, err := keeper.GetAllRevShares(ctx, fill)

			require.ErrorIs(t, err, tc.expectedError)
		})
	}
}
