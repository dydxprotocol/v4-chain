package keeper_test

import (
	"math/big"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

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
		revShareConfig             types.UnconditionalRevShareConfig
		marketMapperRevShareParams types.MarketMapperRevenueShareParams
		lowestTakerFee             int32
		lowestMakerFee             int32
		expectedValid              bool
	}{
		"valid rev share config": {
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
			lowestTakerFee: 350,
			lowestMakerFee: -110,
			expectedValid:  true,
		},
		"rev share safety violation: unconditional + marketmapper > 100%": {
			revShareConfig: types.UnconditionalRevShareConfig{
				Configs: []types.UnconditionalRevShareConfig_RecipientConfig{
					{
						Address:  constants.AliceAccAddress.String(),
						SharePpm: 600_000, // 60%
					},
				},
			},
			marketMapperRevShareParams: types.MarketMapperRevenueShareParams{
				Address:         constants.AliceAccAddress.String(),
				RevenueSharePpm: 500_000, // 50%
				ValidDays:       0,
			},
			lowestTakerFee: 350,
			lowestMakerFee: -110,
			expectedValid:  false,
		},
		"rev share safety violation: affiliate and fees": {
			lowestTakerFee: 350,
			lowestMakerFee: -220,
			expectedValid:  false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.RevShareKeeper

			valid := k.ValidateRevShareSafety(
				ctx,
				tc.revShareConfig,
				tc.marketMapperRevShareParams,
				tc.lowestTakerFee,
				tc.lowestMakerFee,
			)
			require.Equal(t, tc.expectedValid, valid)
		})
	}
}

func TestKeeper_GetAllRevShares_Valid(t *testing.T) {
	perpetualId := uint32(1)
	marketId := uint32(1)
	tests := []struct {
		name                     string
		fill                     clobtypes.FillForProcess
		expectedRevSharesForFill types.RevSharesForFill
		setup                    func(tApp *testapp.TestApp, ctx sdk.Context,
			keeper *keeper.Keeper, affiliatesKeeper *affiliateskeeper.Keeper)
	}{
		{
			name: "Valid revenue share from affiliates, unconditional and market mapper",
			expectedRevSharesForFill: types.RevSharesForFill{
				AllRevShares: []types.RevShare{
					{
						// BobAddress has 500000000000000000000000 coins staked
						// via genesis. Which puts Bob at the highest tier
						// and gives him a 15% taker fee share.
						Recipient:         constants.BobAccAddress.String(),
						RevShareFeeSource: types.REV_SHARE_FEE_SOURCE_TAKER_FEE,
						RevShareType:      types.REV_SHARE_TYPE_AFFILIATE,
						QuoteQuantums:     big.NewInt(1_500_000), // 15 % of 10 million taker fee quote quantums
						RevSharePpm:       150_000,
					},
					{
						Recipient:         constants.BobAccAddress.String(),
						RevShareFeeSource: types.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE,
						RevShareType:      types.REV_SHARE_TYPE_UNCONDITIONAL,
						QuoteQuantums:     big.NewInt(2_100_000), // (10 + 2 - 1.5) * 20%
						RevSharePpm:       200_000,
					},
					{
						Recipient:         constants.AliceAccAddress.String(),
						RevShareFeeSource: types.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE,
						RevShareType:      types.REV_SHARE_TYPE_UNCONDITIONAL,
						QuoteQuantums:     big.NewInt(3_150_000), // (10 + 2 - 1.5) * 30%
						RevSharePpm:       300_000,
					},
					{
						Recipient:         constants.AliceAccAddress.String(),
						RevShareFeeSource: types.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE,
						RevShareType:      types.REV_SHARE_TYPE_MARKET_MAPPER,
						QuoteQuantums:     big.NewInt(1_050_000), // (10 + 2 - 1.5) * 10%
						RevSharePpm:       100_000,
					},
				},
				AffiliateRevShare: &types.RevShare{
					Recipient:         constants.BobAccAddress.String(),
					RevShareFeeSource: types.REV_SHARE_FEE_SOURCE_TAKER_FEE,
					RevShareType:      types.REV_SHARE_TYPE_AFFILIATE,
					QuoteQuantums:     big.NewInt(1_500_000), // 15 % of 10 million taker fee quote quantums
					RevSharePpm:       150_000,
				},
				FeeSourceToQuoteQuantums: map[types.RevShareFeeSource]*big.Int{
					types.REV_SHARE_FEE_SOURCE_TAKER_FEE: big.NewInt(1_500_000), // affiliate rev share fees
					// unconditional + market mapper rev shares fees
					types.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE: big.NewInt(6_300_000),
				},
				FeeSourceToRevSharePpm: map[types.RevShareFeeSource]uint32{
					types.REV_SHARE_FEE_SOURCE_TAKER_FEE:            150_000, // affiliate rev share fee ppm
					types.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE: 600_000, // unconditional + market mapper rev share fee ppm
				},
			},
			fill: clobtypes.FillForProcess{
				TakerAddr:                         constants.AliceAccAddress.String(),
				TakerFeeQuoteQuantums:             big.NewInt(10_000_000),
				MakerAddr:                         constants.BobAccAddress.String(),
				MakerFeeQuoteQuantums:             big.NewInt(2_000_000),
				FillQuoteQuantums:                 big.NewInt(100_000_000_000),
				ProductId:                         perpetualId,
				MonthlyRollingTakerVolumeQuantums: 1_000_000_000_000,
				MarketId:                          marketId,
			},
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

				err = affiliatesKeeper.UpdateAffiliateTiers(ctx, affiliatetypes.DefaultAffiliateTiers)
				require.NoError(t, err)
				err = affiliatesKeeper.RegisterAffiliate(ctx, constants.AliceAccAddress.String(),
					constants.BobAccAddress.String())
				require.NoError(t, err)
			},
		},
		{
			name: "Valid rev-share from affiliates, negative maker fee and unconditional and market mapper",
			expectedRevSharesForFill: types.RevSharesForFill{
				AllRevShares: []types.RevShare{
					{
						Recipient:         constants.BobAccAddress.String(),
						RevShareFeeSource: types.REV_SHARE_FEE_SOURCE_TAKER_FEE,
						RevShareType:      types.REV_SHARE_TYPE_AFFILIATE,
						QuoteQuantums:     big.NewInt(1_500_000),
						RevSharePpm:       150_000,
					},
					{
						Recipient:         constants.BobAccAddress.String(),
						RevShareFeeSource: types.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE,
						RevShareType:      types.REV_SHARE_TYPE_UNCONDITIONAL,
						QuoteQuantums:     big.NewInt(1_300_000),
						RevSharePpm:       200_000,
					},
					{
						Recipient:         constants.AliceAccAddress.String(),
						RevShareFeeSource: types.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE,
						RevShareType:      types.REV_SHARE_TYPE_UNCONDITIONAL,
						QuoteQuantums:     big.NewInt(1_950_000),
						RevSharePpm:       300_000,
					},
					{
						Recipient:         constants.AliceAccAddress.String(),
						RevShareFeeSource: types.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE,
						RevShareType:      types.REV_SHARE_TYPE_MARKET_MAPPER,
						QuoteQuantums:     big.NewInt(650_000),
						RevSharePpm:       100_000,
					},
				},
				FeeSourceToQuoteQuantums: map[types.RevShareFeeSource]*big.Int{
					types.REV_SHARE_FEE_SOURCE_TAKER_FEE:            big.NewInt(1_500_000),
					types.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE: big.NewInt(3_900_000),
				},
				FeeSourceToRevSharePpm: map[types.RevShareFeeSource]uint32{
					types.REV_SHARE_FEE_SOURCE_TAKER_FEE:            150_000,
					types.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE: 600_000,
				},
				AffiliateRevShare: &types.RevShare{
					Recipient:         constants.BobAccAddress.String(),
					RevShareFeeSource: types.REV_SHARE_FEE_SOURCE_TAKER_FEE,
					RevShareType:      types.REV_SHARE_TYPE_AFFILIATE,
					QuoteQuantums:     big.NewInt(1_500_000),
					RevSharePpm:       150_000,
				},
			},
			fill: clobtypes.FillForProcess{
				TakerAddr:                         constants.AliceAccAddress.String(),
				TakerFeeQuoteQuantums:             big.NewInt(10_000_000),
				MakerAddr:                         constants.BobAccAddress.String(),
				MakerFeeQuoteQuantums:             big.NewInt(-2_000_000),
				FillQuoteQuantums:                 big.NewInt(100_000_000_000),
				ProductId:                         perpetualId,
				MarketId:                          marketId,
				MonthlyRollingTakerVolumeQuantums: 1_000_000_000_000,
			},
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

				err = affiliatesKeeper.UpdateAffiliateTiers(ctx, affiliatetypes.DefaultAffiliateTiers)
				require.NoError(t, err)
				err = affiliatesKeeper.RegisterAffiliate(ctx, constants.AliceAccAddress.String(),
					constants.BobAccAddress.String())
				require.NoError(t, err)
			},
		},
		{
			name: "Valid revenue share with 30d volume greater than max 30d referral volume",
			fill: clobtypes.FillForProcess{
				TakerAddr:                         constants.AliceAccAddress.String(),
				TakerFeeQuoteQuantums:             big.NewInt(10_000_000),
				MakerAddr:                         constants.BobAccAddress.String(),
				MakerFeeQuoteQuantums:             big.NewInt(2_000_000),
				FillQuoteQuantums:                 big.NewInt(100_000_000_000),
				ProductId:                         perpetualId,
				MonthlyRollingTakerVolumeQuantums: types.MaxReferee30dVolumeForAffiliateShareQuantums + 1,
				MarketId:                          marketId,
			},
			expectedRevSharesForFill: types.RevSharesForFill{
				AllRevShares: []types.RevShare{
					{
						Recipient:         constants.BobAccAddress.String(),
						RevShareFeeSource: types.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE,
						RevShareType:      types.REV_SHARE_TYPE_UNCONDITIONAL,
						QuoteQuantums:     big.NewInt(2_400_000),
						RevSharePpm:       200_000,
					},
					{
						Recipient:         constants.AliceAccAddress.String(),
						RevShareFeeSource: types.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE,
						RevShareType:      types.REV_SHARE_TYPE_UNCONDITIONAL,
						QuoteQuantums:     big.NewInt(3_600_000),
						RevSharePpm:       300_000,
					},
					{
						Recipient:         constants.AliceAccAddress.String(),
						RevShareFeeSource: types.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE,
						RevShareType:      types.REV_SHARE_TYPE_MARKET_MAPPER,
						QuoteQuantums:     big.NewInt(1_200_000),
						RevSharePpm:       100_000,
					},
				},
				AffiliateRevShare: nil,
				FeeSourceToQuoteQuantums: map[types.RevShareFeeSource]*big.Int{
					types.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE: big.NewInt(7_200_000),
					types.REV_SHARE_FEE_SOURCE_TAKER_FEE:            big.NewInt(0),
				},
				FeeSourceToRevSharePpm: map[types.RevShareFeeSource]uint32{
					types.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE: 600_000,
					types.REV_SHARE_FEE_SOURCE_TAKER_FEE:            0,
				},
			},
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

				err = affiliatesKeeper.UpdateAffiliateTiers(ctx, affiliatetypes.DefaultAffiliateTiers)
				require.NoError(t, err)
				err = affiliatesKeeper.RegisterAffiliate(ctx, constants.AliceAccAddress.String(),
					constants.BobAccAddress.String())
				require.NoError(t, err)
			},
		},
		{
			name: "Valid revenue share with no unconditional rev shares",
			expectedRevSharesForFill: types.RevSharesForFill{
				AllRevShares: []types.RevShare{
					{
						Recipient:         constants.BobAccAddress.String(),
						RevShareFeeSource: types.REV_SHARE_FEE_SOURCE_TAKER_FEE,
						RevShareType:      types.REV_SHARE_TYPE_AFFILIATE,
						QuoteQuantums:     big.NewInt(1_500_000),
						RevSharePpm:       150_000, // 15%
					},
					{
						Recipient:         constants.AliceAccAddress.String(),
						RevShareFeeSource: types.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE,
						RevShareType:      types.REV_SHARE_TYPE_MARKET_MAPPER,
						QuoteQuantums:     big.NewInt(1_050_000),
						RevSharePpm:       100_000, // 10%
					},
				},
				FeeSourceToQuoteQuantums: map[types.RevShareFeeSource]*big.Int{
					types.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE: big.NewInt(1_050_000),
					types.REV_SHARE_FEE_SOURCE_TAKER_FEE:            big.NewInt(1_500_000),
				},
				FeeSourceToRevSharePpm: map[types.RevShareFeeSource]uint32{
					types.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE: 100_000, // 10%
					types.REV_SHARE_FEE_SOURCE_TAKER_FEE:            150_000, // 15%
				},
				AffiliateRevShare: &types.RevShare{
					Recipient:         constants.BobAccAddress.String(),
					RevShareFeeSource: types.REV_SHARE_FEE_SOURCE_TAKER_FEE,
					RevShareType:      types.REV_SHARE_TYPE_AFFILIATE,
					QuoteQuantums:     big.NewInt(1_500_000),
					RevSharePpm:       150_000, // 15%
				},
			},
			fill: clobtypes.FillForProcess{
				TakerAddr:                         constants.AliceAccAddress.String(),
				TakerFeeQuoteQuantums:             big.NewInt(10_000_000),
				MakerAddr:                         constants.BobAccAddress.String(),
				MakerFeeQuoteQuantums:             big.NewInt(2_000_000),
				FillQuoteQuantums:                 big.NewInt(100_000_000_000),
				ProductId:                         perpetualId,
				MarketId:                          marketId,
				MonthlyRollingTakerVolumeQuantums: 1_000_000_000_000,
			},
			setup: func(tApp *testapp.TestApp, ctx sdk.Context, keeper *keeper.Keeper,
				affiliatesKeeper *affiliateskeeper.Keeper) {
				err := keeper.SetMarketMapperRevenueShareParams(ctx, types.MarketMapperRevenueShareParams{
					Address:         constants.AliceAccAddress.String(),
					RevenueSharePpm: 100_000, // 10%
					ValidDays:       1,
				})
				require.NoError(t, err)

				err = affiliatesKeeper.UpdateAffiliateTiers(ctx, affiliatetypes.DefaultAffiliateTiers)
				require.NoError(t, err)
				err = affiliatesKeeper.RegisterAffiliate(ctx, constants.AliceAccAddress.String(),
					constants.BobAccAddress.String())
				require.NoError(t, err)
			},
		},
		{
			name: "Valid revenue share with no market mapper rev share",
			expectedRevSharesForFill: types.RevSharesForFill{
				AllRevShares: []types.RevShare{
					{
						Recipient:         constants.BobAccAddress.String(),
						RevShareFeeSource: types.REV_SHARE_FEE_SOURCE_TAKER_FEE,
						RevShareType:      types.REV_SHARE_TYPE_AFFILIATE,
						QuoteQuantums:     big.NewInt(1_500_000),
						RevSharePpm:       150_000, // 15%
					},
					{
						Recipient:         constants.BobAccAddress.String(),
						RevShareFeeSource: types.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE,
						RevShareType:      types.REV_SHARE_TYPE_UNCONDITIONAL,
						QuoteQuantums:     big.NewInt(2_100_000),
						RevSharePpm:       200_000, // 20%
					},
				},
				FeeSourceToQuoteQuantums: map[types.RevShareFeeSource]*big.Int{
					types.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE: big.NewInt(2_100_000),
					types.REV_SHARE_FEE_SOURCE_TAKER_FEE:            big.NewInt(1_500_000),
				},
				FeeSourceToRevSharePpm: map[types.RevShareFeeSource]uint32{
					types.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE: 200_000, // 20%
					types.REV_SHARE_FEE_SOURCE_TAKER_FEE:            150_000, // 15%
				},
				AffiliateRevShare: &types.RevShare{
					Recipient:         constants.BobAccAddress.String(),
					RevShareFeeSource: types.REV_SHARE_FEE_SOURCE_TAKER_FEE,
					RevShareType:      types.REV_SHARE_TYPE_AFFILIATE,
					QuoteQuantums:     big.NewInt(1_500_000),
					RevSharePpm:       150_000, // 15%
				},
			},
			fill: clobtypes.FillForProcess{
				TakerAddr:                         constants.AliceAccAddress.String(),
				TakerFeeQuoteQuantums:             big.NewInt(10_000_000),
				MakerAddr:                         constants.BobAccAddress.String(),
				MakerFeeQuoteQuantums:             big.NewInt(2_000_000),
				FillQuoteQuantums:                 big.NewInt(100_000_000_000),
				ProductId:                         perpetualId,
				MarketId:                          marketId,
				MonthlyRollingTakerVolumeQuantums: 1_000_000_000_000,
			},
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
				err := affiliatesKeeper.UpdateAffiliateTiers(ctx, affiliatetypes.DefaultAffiliateTiers)
				require.NoError(t, err)
				err = affiliatesKeeper.RegisterAffiliate(ctx, constants.AliceAccAddress.String(),
					constants.BobAccAddress.String())
				require.NoError(t, err)
			},
		},
		{
			name: "Valid revenue share with whitelisted affiliate and no unconditional rev shares",
			expectedRevSharesForFill: types.RevSharesForFill{
				AllRevShares: []types.RevShare{
					{
						Recipient:         constants.BobAccAddress.String(),
						RevShareFeeSource: types.REV_SHARE_FEE_SOURCE_TAKER_FEE,
						RevShareType:      types.REV_SHARE_TYPE_AFFILIATE,
						QuoteQuantums:     big.NewInt(2_500_000),
						RevSharePpm:       250_000, // 25%
					},
					{
						Recipient:         constants.AliceAccAddress.String(),
						RevShareFeeSource: types.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE,
						RevShareType:      types.REV_SHARE_TYPE_MARKET_MAPPER,
						QuoteQuantums:     big.NewInt(950_000),
						RevSharePpm:       100_000, // 10%
					},
				},
				AffiliateRevShare: &types.RevShare{
					Recipient:         constants.BobAccAddress.String(),
					RevShareFeeSource: types.REV_SHARE_FEE_SOURCE_TAKER_FEE,
					RevShareType:      types.REV_SHARE_TYPE_AFFILIATE,
					QuoteQuantums:     big.NewInt(2_500_000),
					RevSharePpm:       250_000, // 25%
				},
				FeeSourceToQuoteQuantums: map[types.RevShareFeeSource]*big.Int{
					types.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE: big.NewInt(950_000),
					types.REV_SHARE_FEE_SOURCE_TAKER_FEE:            big.NewInt(2_500_000),
				},
				FeeSourceToRevSharePpm: map[types.RevShareFeeSource]uint32{
					types.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE: 100_000, // 10%
					types.REV_SHARE_FEE_SOURCE_TAKER_FEE:            250_000, // 25%
				},
			},
			fill: clobtypes.FillForProcess{
				TakerAddr:                         constants.AliceAccAddress.String(),
				TakerFeeQuoteQuantums:             big.NewInt(10_000_000),
				MakerAddr:                         constants.BobAccAddress.String(),
				MakerFeeQuoteQuantums:             big.NewInt(2_000_000),
				FillQuoteQuantums:                 big.NewInt(100_000_000_000),
				ProductId:                         marketId,
				MonthlyRollingTakerVolumeQuantums: 1_000_000_000_000,
			},
			setup: func(tApp *testapp.TestApp, ctx sdk.Context, keeper *keeper.Keeper,
				affiliatesKeeper *affiliateskeeper.Keeper) {
				err := keeper.SetMarketMapperRevenueShareParams(ctx, types.MarketMapperRevenueShareParams{
					Address:         constants.AliceAccAddress.String(),
					RevenueSharePpm: 100_000, // 10%
					ValidDays:       1,
				})
				require.NoError(t, err)

				err = affiliatesKeeper.UpdateAffiliateTiers(ctx, affiliatetypes.DefaultAffiliateTiers)
				require.NoError(t, err)
				err = affiliatesKeeper.RegisterAffiliate(ctx, constants.AliceAccAddress.String(),
					constants.BobAccAddress.String())
				require.NoError(t, err)
				err = affiliatesKeeper.SetAffiliateWhitelist(ctx, affiliatetypes.AffiliateWhitelist{
					Tiers: []affiliatetypes.AffiliateWhitelist_Tier{
						{
							Addresses:        []string{constants.BobAccAddress.String()},
							TakerFeeSharePpm: 250_000, // 25%
						},
					},
				})
				require.NoError(t, err)
			},
		},
		{
			name: "No rev shares",
			expectedRevSharesForFill: types.RevSharesForFill{
				AllRevShares:      []types.RevShare{},
				AffiliateRevShare: nil,
				FeeSourceToQuoteQuantums: map[types.RevShareFeeSource]*big.Int{
					types.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE: big.NewInt(0),
					types.REV_SHARE_FEE_SOURCE_TAKER_FEE:            big.NewInt(0),
				},
				FeeSourceToRevSharePpm: map[types.RevShareFeeSource]uint32{
					types.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE: 0,
					types.REV_SHARE_FEE_SOURCE_TAKER_FEE:            0,
				},
			},
			fill: clobtypes.FillForProcess{
				TakerAddr:                         constants.AliceAccAddress.String(),
				TakerFeeQuoteQuantums:             big.NewInt(10_000_000),
				MakerAddr:                         constants.BobAccAddress.String(),
				MakerFeeQuoteQuantums:             big.NewInt(2_000_000),
				FillQuoteQuantums:                 big.NewInt(100_000_000_000),
				ProductId:                         perpetualId,
				MarketId:                          marketId,
				MonthlyRollingTakerVolumeQuantums: 1_000_000_000_000,
			},
			setup: func(tApp *testapp.TestApp, ctx sdk.Context, keeper *keeper.Keeper,
				affiliatesKeeper *affiliateskeeper.Keeper) {
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup

			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			keeper := tApp.App.RevShareKeeper
			affiliatesKeeper := tApp.App.AffiliatesKeeper
			if tc.setup != nil {
				tc.setup(tApp, ctx, &keeper, &affiliatesKeeper)
			}

			keeper.CreateNewMarketRevShare(ctx, marketId)
			affiliateWhitelistMap, err := affiliatesKeeper.GetAffiliateWhitelistMap(ctx)
			require.NoError(t, err)

			revSharesForFill, err := keeper.GetAllRevShares(ctx, tc.fill, affiliateWhitelistMap)

			require.NoError(t, err)
			require.Equal(t, tc.expectedRevSharesForFill, revSharesForFill)
		})
	}
}

func TestKeeper_GetAllRevShares_Invalid(t *testing.T) {
	perpetualId := uint32(1)
	marketId := uint32(1)
	tests := []struct {
		name                              string
		revenueSharePpmNetFees            uint32
		revenueSharePpmTakerFees          uint32
		monthlyRollingTakerVolumeQuantums uint64
		expectedError                     error
		setup                             func(tApp *testapp.TestApp, ctx sdk.Context,
			keeper *keeper.Keeper, affiliatesKeeper *affiliateskeeper.Keeper)
	}{
		{
			name:                              "Total fees shared exceeds net fees from all sources",
			revenueSharePpmNetFees:            950_000,           // 95%,
			revenueSharePpmTakerFees:          150_000,           // 15%
			monthlyRollingTakerVolumeQuantums: 1_000_000_000_000, // 1 million USDC
			expectedError:                     types.ErrTotalFeesSharedExceedsNetFees,
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
							SharePpm: 250_000, // 25%
						},
					},
				})
				err = affiliatesKeeper.UpdateAffiliateTiers(ctx, affiliatetypes.DefaultAffiliateTiers)
				require.NoError(t, err)
				err = affiliatesKeeper.RegisterAffiliate(ctx, constants.AliceAccAddress.String(),
					constants.BobAccAddress.String())
				require.NoError(t, err)
			},
		},
		{
			name:                              "Total fees shared exceeds net fees - no affiliate rev shares",
			revenueSharePpmNetFees:            1_150_000,         // 115%,
			revenueSharePpmTakerFees:          0,                 // 0%
			monthlyRollingTakerVolumeQuantums: 1_000_000_000_000, // 1 million USDC
			expectedError:                     types.ErrTotalFeesSharedExceedsNetFees,
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
			fill := clobtypes.FillForProcess{
				TakerAddr:                         constants.AliceAccAddress.String(),
				TakerFeeQuoteQuantums:             big.NewInt(10_000_000),
				MakerAddr:                         constants.BobAccAddress.String(),
				MakerFeeQuoteQuantums:             big.NewInt(2_000_000),
				FillQuoteQuantums:                 big.NewInt(100_000_000_000),
				ProductId:                         perpetualId,
				MarketId:                          marketId,
				MonthlyRollingTakerVolumeQuantums: 1_000_000_000_000,
			}
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			keeper := tApp.App.RevShareKeeper
			affiliatesKeeper := tApp.App.AffiliatesKeeper
			if tc.setup != nil {
				tc.setup(tApp, ctx, &keeper, &affiliatesKeeper)
			}

			keeper.CreateNewMarketRevShare(ctx, marketId)

			_, err := keeper.GetAllRevShares(ctx, fill, map[string]uint32{})

			require.ErrorIs(t, err, tc.expectedError)
		})
	}
}
