package keeper_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	affiliatetypes "github.com/dydxprotocol/v4-chain/protocol/x/affiliates/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	statstypes "github.com/dydxprotocol/v4-chain/protocol/x/stats/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

// TestProcessSingleMatch_AffiliateRevenueAttribution_TakerOnly tests that when only the taker
// has an affiliate referrer, the attribution is correctly stored in BlockStats.
func TestProcessSingleMatch_AffiliateRevenueAttribution_TakerOnly(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()

	ctx := tApp.InitChain()
	k := tApp.App.ClobKeeper

	// Create taker and maker subaccounts with sufficient collateral
	takerSubaccount := constants.Alice_Num0
	makerSubaccount := constants.Bob_Num0

	tApp.App.SubaccountsKeeper.SetSubaccount(ctx, satypes.Subaccount{
		Id: &takerSubaccount,
		AssetPositions: []*satypes.AssetPosition{
			{
				AssetId:  0,
				Quantums: dtypes.NewInt(1_000_000_000_000), // 1M USDC
			},
		},
	})

	tApp.App.SubaccountsKeeper.SetSubaccount(ctx, satypes.Subaccount{
		Id: &makerSubaccount,
		AssetPositions: []*satypes.AssetPosition{
			{
				AssetId:  0,
				Quantums: dtypes.NewInt(1_000_000_000_000), // 1M USDC
			},
		},
	})

	// Set up affiliate tiers
	err := tApp.App.AffiliatesKeeper.UpdateAffiliateTiers(ctx, affiliatetypes.AffiliateTiers{
		Tiers: []affiliatetypes.AffiliateTiers_Tier{
			{
				ReqReferredVolumeQuoteQuantums: 0,
				ReqStakedWholeCoins:            0,
				TakerFeeSharePpm:               100_000, // 10%
			},
		},
	})
	require.NoError(t, err)

	// Register taker with an affiliate referrer (use Carl as referrer)
	referrerAddr := constants.CarlAccAddress.String()
	err = tApp.App.AffiliatesKeeper.RegisterAffiliate(ctx, constants.Alice_Num0.Owner, referrerAddr)
	require.NoError(t, err)

	// Create orders
	takerOrder := clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: takerSubaccount,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_ShortTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     100_000_000,   // 1 BTC (in base quantums)
		Subticks:     5_000_000_000, // $50,000 per BTC
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 100},
	}

	makerOrder := clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: makerSubaccount,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_ShortTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     100_000_000,
		Subticks:     5_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 100},
	}

	matchWithOrders := &clobtypes.MatchWithOrders{
		TakerOrder: &takerOrder,
		MakerOrder: &makerOrder,
		FillAmount: satypes.BaseQuantums(100_000_000), // Fill full amount
	}

	// Process the match
	success, _, _, _, err := k.ProcessSingleMatch(
		ctx,
		matchWithOrders,
		map[string]bool{},
		affiliatetypes.AffiliateParameters{
			Maximum_30DAttributableVolumePerReferredUserQuoteQuantums: 1_000_000_000_000,
		},
	)
	require.NoError(t, err)
	require.True(t, success)

	// Verify BlockStats contains the affiliate revenue attribution for taker only
	blockStats := tApp.App.StatsKeeper.GetBlockStats(ctx)
	require.Len(t, blockStats.Fills, 1)

	fill := blockStats.Fills[0]
	require.Equal(t, constants.Alice_Num0.Owner, fill.Taker)
	require.Equal(t, constants.Bob_Num0.Owner, fill.Maker)

	// Verify affiliate revenue attributions array
	require.Len(t, fill.AffiliateRevenueAttributions, 1, "Should have exactly one attribution (taker only)")

	attribution := fill.AffiliateRevenueAttributions[0]
	require.Equal(t, referrerAddr, attribution.ReferrerAddress)
	require.Greater(t, attribution.ReferredVolumeQuoteQuantums, uint64(0), "Should have non-zero attributed volume")

	// The notional should match the attributed volume (no cap hit)
	require.Equal(t, fill.Notional, attribution.ReferredVolumeQuoteQuantums)
}

// TestProcessSingleMatch_AffiliateRevenueAttribution_BothTakerAndMaker tests that when both
// taker and maker have affiliate referrers, both attributions are stored in BlockStats.
func TestProcessSingleMatch_AffiliateRevenueAttribution_BothTakerAndMaker(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()

	ctx := tApp.InitChain()
	k := tApp.App.ClobKeeper

	// Create subaccounts
	takerSubaccount := constants.Alice_Num0
	makerSubaccount := constants.Bob_Num0

	tApp.App.SubaccountsKeeper.SetSubaccount(ctx, satypes.Subaccount{
		Id: &takerSubaccount,
		AssetPositions: []*satypes.AssetPosition{
			{
				AssetId:  0,
				Quantums: dtypes.NewInt(1_000_000_000_000),
			},
		},
	})

	tApp.App.SubaccountsKeeper.SetSubaccount(ctx, satypes.Subaccount{
		Id: &makerSubaccount,
		AssetPositions: []*satypes.AssetPosition{
			{
				AssetId:  0,
				Quantums: dtypes.NewInt(1_000_000_000_000),
			},
		},
	})

	// Set up affiliate tiers
	err := tApp.App.AffiliatesKeeper.UpdateAffiliateTiers(ctx, affiliatetypes.AffiliateTiers{
		Tiers: []affiliatetypes.AffiliateTiers_Tier{
			{
				ReqReferredVolumeQuoteQuantums: 0,
				ReqStakedWholeCoins:            0,
				TakerFeeSharePpm:               100_000, // 10%
			},
		},
	})
	require.NoError(t, err)

	// Register BOTH taker and maker with different affiliate referrers
	takerReferrerAddr := constants.CarlAccAddress.String()
	makerReferrerAddr := constants.DaveAccAddress.String()

	err = tApp.App.AffiliatesKeeper.RegisterAffiliate(ctx, constants.Alice_Num0.Owner, takerReferrerAddr)
	require.NoError(t, err)

	err = tApp.App.AffiliatesKeeper.RegisterAffiliate(ctx, constants.Bob_Num0.Owner, makerReferrerAddr)
	require.NoError(t, err)

	// Create orders
	takerOrder := clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: takerSubaccount,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_ShortTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     100_000_000,
		Subticks:     5_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 100},
	}

	makerOrder := clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: makerSubaccount,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_ShortTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     100_000_000,
		Subticks:     5_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 100},
	}

	matchWithOrders := &clobtypes.MatchWithOrders{
		TakerOrder: &takerOrder,
		MakerOrder: &makerOrder,
		FillAmount: satypes.BaseQuantums(100_000_000),
	}

	// Process the match
	success, _, _, _, err := k.ProcessSingleMatch(
		ctx,
		matchWithOrders,
		map[string]bool{},
		affiliatetypes.AffiliateParameters{
			Maximum_30DAttributableVolumePerReferredUserQuoteQuantums: 1_000_000_000_000,
		},
	)
	require.NoError(t, err)
	require.True(t, success)

	// Verify BlockStats contains BOTH affiliate revenue attributions
	blockStats := tApp.App.StatsKeeper.GetBlockStats(ctx)
	require.Len(t, blockStats.Fills, 1)

	fill := blockStats.Fills[0]
	require.Equal(t, constants.Alice_Num0.Owner, fill.Taker)
	require.Equal(t, constants.Bob_Num0.Owner, fill.Maker)

	// Verify we have TWO attributions
	require.Len(t, fill.AffiliateRevenueAttributions, 2, "Should have two attributions (taker and maker)")

	// Find taker and maker attributions (order may vary)
	var takerAttribution, makerAttribution *statstypes.AffiliateRevenueAttribution
	for _, attr := range fill.AffiliateRevenueAttributions {
		if attr.ReferrerAddress == takerReferrerAddr {
			takerAttribution = attr
		} else if attr.ReferrerAddress == makerReferrerAddr {
			makerAttribution = attr
		}
	}

	require.NotNil(t, takerAttribution, "Should have taker attribution")
	require.NotNil(t, makerAttribution, "Should have maker attribution")

	// Both should have the same notional volume attributed
	require.Equal(t, fill.Notional, takerAttribution.ReferredVolumeQuoteQuantums)
	require.Equal(t, fill.Notional, makerAttribution.ReferredVolumeQuoteQuantums)
	require.Greater(t, takerAttribution.ReferredVolumeQuoteQuantums, uint64(0))
	require.Greater(t, makerAttribution.ReferredVolumeQuoteQuantums, uint64(0))
}

// TestProcessSingleMatch_AffiliateRevenueAttribution_NoReferrers tests that when neither
// taker nor maker has an affiliate referrer, no attributions are stored.
func TestProcessSingleMatch_AffiliateRevenueAttribution_NoReferrers(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.ClobKeeper

	// Create subaccounts (no affiliate registrations)
	takerSubaccount := constants.Alice_Num0
	makerSubaccount := constants.Bob_Num0

	tApp.App.SubaccountsKeeper.SetSubaccount(ctx, satypes.Subaccount{
		Id: &takerSubaccount,
		AssetPositions: []*satypes.AssetPosition{
			{
				AssetId:  0,
				Quantums: dtypes.NewInt(1_000_000_000_000),
			},
		},
	})

	tApp.App.SubaccountsKeeper.SetSubaccount(ctx, satypes.Subaccount{
		Id: &makerSubaccount,
		AssetPositions: []*satypes.AssetPosition{
			{
				AssetId:  0,
				Quantums: dtypes.NewInt(1_000_000_000_000),
			},
		},
	})

	// Create orders
	takerOrder := clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: takerSubaccount,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_ShortTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     100_000_000,
		Subticks:     5_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 100},
	}

	makerOrder := clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: makerSubaccount,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_ShortTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     100_000_000,
		Subticks:     5_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 100},
	}

	matchWithOrders := &clobtypes.MatchWithOrders{
		TakerOrder: &takerOrder,
		MakerOrder: &makerOrder,
		FillAmount: satypes.BaseQuantums(100_000_000),
	}

	// Process the match
	success, _, _, _, err := k.ProcessSingleMatch(
		ctx,
		matchWithOrders,
		map[string]bool{},
		affiliatetypes.AffiliateParameters{},
	)
	require.NoError(t, err)
	require.True(t, success)

	// Verify BlockStats has no affiliate revenue attributions
	blockStats := tApp.App.StatsKeeper.GetBlockStats(ctx)
	require.Len(t, blockStats.Fills, 1)

	fill := blockStats.Fills[0]
	require.Empty(t, fill.AffiliateRevenueAttributions, "Should have no attributions when neither has referrer")
}

// TestProcessSingleMatch_AffiliateRevenueAttribution_VolumeCapApplied tests that the
// attributable volume cap is correctly applied when storing attributions.
func TestProcessSingleMatch_AffiliateRevenueAttribution_VolumeCapApplied(t *testing.T) {
	lowCap := uint64(100_000_000_000) // Cap at 100k USDC

	tApp := testapp.NewTestAppBuilder(t).Build()

	ctx := tApp.InitChain()
	k := tApp.App.ClobKeeper

	// Create subaccounts
	takerSubaccount := constants.Alice_Num0
	makerSubaccount := constants.Bob_Num0

	tApp.App.SubaccountsKeeper.SetSubaccount(ctx, satypes.Subaccount{
		Id: &takerSubaccount,
		AssetPositions: []*satypes.AssetPosition{
			{
				AssetId:  0,
				Quantums: dtypes.NewInt(1_000_000_000_000),
			},
		},
	})

	tApp.App.SubaccountsKeeper.SetSubaccount(ctx, satypes.Subaccount{
		Id: &makerSubaccount,
		AssetPositions: []*satypes.AssetPosition{
			{
				AssetId:  0,
				Quantums: dtypes.NewInt(1_000_000_000_000),
			},
		},
	})

	// Set up affiliate parameters with low cap
	err := tApp.App.AffiliatesKeeper.UpdateAffiliateParameters(ctx, &affiliatetypes.MsgUpdateAffiliateParameters{
		Authority: constants.GovAuthority,
		AffiliateParameters: affiliatetypes.AffiliateParameters{
			Maximum_30DAttributableVolumePerReferredUserQuoteQuantums: lowCap,
		},
	})
	require.NoError(t, err)

	// Set up affiliate tiers
	err = tApp.App.AffiliatesKeeper.UpdateAffiliateTiers(ctx, affiliatetypes.AffiliateTiers{
		Tiers: []affiliatetypes.AffiliateTiers_Tier{
			{
				ReqReferredVolumeQuoteQuantums: 0,
				ReqStakedWholeCoins:            0,
				TakerFeeSharePpm:               100_000, // 10%
			},
		},
	})
	require.NoError(t, err)

	// Register taker with referrer (use Carl as referrer)
	referrerAddr := constants.CarlAccAddress.String()
	err = tApp.App.AffiliatesKeeper.RegisterAffiliate(ctx, constants.Alice_Num0.Owner, referrerAddr)
	require.NoError(t, err)

	// Create a large trade that exceeds the cap (400k USDC notional)
	// Need to create a very large trade to exceed 100k cap
	takerOrder := clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: takerSubaccount,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_ShortTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     8_000_000_000, // 80 BTC to get large enough notional
		Subticks:     5_000_000_000, // $50,000 per BTC
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 100},
	}

	makerOrder := clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: makerSubaccount,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_ShortTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     8_000_000_000,
		Subticks:     5_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 100},
	}

	matchWithOrders := &clobtypes.MatchWithOrders{
		TakerOrder: &takerOrder,
		MakerOrder: &makerOrder,
		FillAmount: satypes.BaseQuantums(8_000_000_000),
	}

	// Process the match
	success, _, _, _, err := k.ProcessSingleMatch(
		ctx,
		matchWithOrders,
		map[string]bool{},
		affiliatetypes.AffiliateParameters{
			Maximum_30DAttributableVolumePerReferredUserQuoteQuantums: lowCap,
		},
	)
	require.NoError(t, err)
	require.True(t, success)

	// Verify the attributed volume is capped
	blockStats := tApp.App.StatsKeeper.GetBlockStats(ctx)
	require.Len(t, blockStats.Fills, 1)

	fill := blockStats.Fills[0]
	require.Len(t, fill.AffiliateRevenueAttributions, 1)

	attribution := fill.AffiliateRevenueAttributions[0]
	require.Equal(t, referrerAddr, attribution.ReferrerAddress)

	// Verify the trade notional exceeds the cap
	require.Greater(t, fill.Notional, lowCap, "Trade notional should exceed the cap for this test")

	// The attributed volume should be CAPPED at lowCap, not the full notional
	require.LessOrEqual(t, attribution.ReferredVolumeQuoteQuantums, lowCap,
		"Attributed volume should not exceed the cap")
	require.Less(t, attribution.ReferredVolumeQuoteQuantums, fill.Notional,
		"Attributed volume should be less than full notional due to cap")
}
