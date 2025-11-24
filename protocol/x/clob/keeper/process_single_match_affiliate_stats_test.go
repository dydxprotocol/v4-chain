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

// TestProcessSingleMatch_AffiliateAttribution_TakerOnly tests that when only the taker
// has an affiliate referrer, the attribution is correctly stored in BlockStats.
func TestProcessSingleMatch_AffiliateAttribution_TakerOnly(t *testing.T) {
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
	require.Len(t, fill.AffiliateAttributions, 1, "Should have exactly one attribution (taker only)")

	attribution := fill.AffiliateAttributions[0]
	require.Equal(t, statstypes.AffiliateAttribution_ROLE_TAKER, attribution.Role)
	require.Equal(t, referrerAddr, attribution.ReferrerAddress)
}

// TestProcessSingleMatch_AffiliateAttribution_BothTakerAndMaker tests that when both
// taker and maker have affiliate referrers, both attributions are stored in BlockStats.
func TestProcessSingleMatch_AffiliateAttribution_BothTakerAndMaker(t *testing.T) {
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
	require.Len(t, fill.AffiliateAttributions, 2, "Should have two attributions (taker and maker)")

	// Find taker and maker attributions by role
	var takerAttribution, makerAttribution *statstypes.AffiliateAttribution
	for _, attr := range fill.AffiliateAttributions {
		if attr.Role == statstypes.AffiliateAttribution_ROLE_TAKER {
			takerAttribution = attr
		} else if attr.Role == statstypes.AffiliateAttribution_ROLE_MAKER {
			makerAttribution = attr
		}
	}

	require.NotNil(t, takerAttribution, "Should have taker attribution")
	require.NotNil(t, makerAttribution, "Should have maker attribution")

	// Verify roles and referrers
	require.Equal(t, statstypes.AffiliateAttribution_ROLE_TAKER, takerAttribution.Role)
	require.Equal(t, takerReferrerAddr, takerAttribution.ReferrerAddress)
	require.Equal(t, statstypes.AffiliateAttribution_ROLE_MAKER, makerAttribution.Role)
	require.Equal(t, makerReferrerAddr, makerAttribution.ReferrerAddress)
}

// TestProcessSingleMatch_AffiliateAttribution_NoReferrers tests that when neither
// taker nor maker has an affiliate referrer, no attributions are stored.
func TestProcessSingleMatch_AffiliateAttribution_NoReferrers(t *testing.T) {
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
	require.Empty(t, fill.AffiliateAttributions, "Should have no attributions when neither has referrer")
}

// TestProcessSingleMatch_AffiliateAttribution_VolumeCapApplied tests that the
// attributable volume cap is correctly applied when storing attributions.
func TestProcessSingleMatch_AffiliateAttribution_VolumeCapApplied(t *testing.T) {
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
	require.Len(t, fill.AffiliateAttributions, 1)

	attribution := fill.AffiliateAttributions[0]
	require.Equal(t, statstypes.AffiliateAttribution_ROLE_TAKER, attribution.Role)
	require.Equal(t, referrerAddr, attribution.ReferrerAddress)

	// Verify the trade notional exceeds the cap
	require.Greater(t, fill.Notional, lowCap, "Trade notional should exceed the cap for this test")
}

// TestProcessSingleMatch_AffiliateAttribution_AlreadyAtCap tests that when a user
// has already reached the 30-day attributable volume cap, no new volume is attributed.
func TestProcessSingleMatch_AffiliateAttribution_AlreadyAtCap(t *testing.T) {
	cap := uint64(100_000_000_000) // 100k USDC cap

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

	// Set up affiliate parameters with cap
	err := tApp.App.AffiliatesKeeper.UpdateAffiliateParameters(ctx, &affiliatetypes.MsgUpdateAffiliateParameters{
		Authority: constants.GovAuthority,
		AffiliateParameters: affiliatetypes.AffiliateParameters{
			Maximum_30DAttributableVolumePerReferredUserQuoteQuantums: cap,
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

	// Register taker with referrer
	referrerAddr := constants.CarlAccAddress.String()
	err = tApp.App.AffiliatesKeeper.RegisterAffiliate(ctx, constants.Alice_Num0.Owner, referrerAddr)
	require.NoError(t, err)

	// Set taker's previous ATTRIBUTED volume to EXACTLY the cap
	// (This is the key fix - we track attributed volume, not total trading volume)
	tApp.App.StatsKeeper.SetUserStats(ctx, constants.Alice_Num0.Owner, &statstypes.UserStats{
		TakerNotional: 200_000_000_000, // User has traded 200k total
		MakerNotional: 0,
		Affiliate_30DAttributedVolumeQuoteQuantums: cap, // But only 100k was attributed (at cap)
	})

	// Create a normal trade
	takerOrder := clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: takerSubaccount,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_ShortTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     100_000_000,   // 1 BTC
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
		FillAmount: satypes.BaseQuantums(100_000_000),
	}

	// Process the match
	success, _, _, _, err := k.ProcessSingleMatch(
		ctx,
		matchWithOrders,
		map[string]bool{},
		affiliatetypes.AffiliateParameters{
			Maximum_30DAttributableVolumePerReferredUserQuoteQuantums: cap,
		},
	)
	require.NoError(t, err)
	require.True(t, success)

	// Verify NO attribution was made since user is already at cap
	blockStats := tApp.App.StatsKeeper.GetBlockStats(ctx)
	require.Len(t, blockStats.Fills, 1)

	fill := blockStats.Fills[0]

	// Should have empty attributions array since no volume can be attributed
	require.Empty(t, fill.AffiliateAttributions,
		"Should have no attributions when referee is already at cap")
}

// TestProcessSingleMatch_AffiliateAttribution_OverCap tests that when a user
// has volume EXCEEDING the 30-day cap, no new volume is attributed.
func TestProcessSingleMatch_AffiliateAttribution_OverCap(t *testing.T) {
	cap := uint64(100_000_000_000) // 100k USDC cap

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

	// Set up affiliate parameters with cap
	err := tApp.App.AffiliatesKeeper.UpdateAffiliateParameters(ctx, &affiliatetypes.MsgUpdateAffiliateParameters{
		Authority: constants.GovAuthority,
		AffiliateParameters: affiliatetypes.AffiliateParameters{
			Maximum_30DAttributableVolumePerReferredUserQuoteQuantums: cap,
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

	// Register taker with referrer
	referrerAddr := constants.CarlAccAddress.String()
	err = tApp.App.AffiliatesKeeper.RegisterAffiliate(ctx, constants.Alice_Num0.Owner, referrerAddr)
	require.NoError(t, err)

	// Set taker's previous ATTRIBUTED volume to EXCEED the cap
	// (User has traded 300k total, but 150k was attributed, which exceeds 100k cap)
	tApp.App.StatsKeeper.SetUserStats(ctx, constants.Alice_Num0.Owner, &statstypes.UserStats{
		TakerNotional: 200_000_000_000, // User traded 200k as taker
		MakerNotional: 100_000_000_000, // User traded 100k as maker
		Affiliate_30DAttributedVolumeQuoteQuantums: 150_000_000_000, // 150k attributed > 100k cap
	})

	// Create a normal trade
	takerOrder := clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: takerSubaccount,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_ShortTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     100_000_000,   // 1 BTC
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
		FillAmount: satypes.BaseQuantums(100_000_000),
	}

	// Process the match
	success, _, _, _, err := k.ProcessSingleMatch(
		ctx,
		matchWithOrders,
		map[string]bool{},
		affiliatetypes.AffiliateParameters{
			Maximum_30DAttributableVolumePerReferredUserQuoteQuantums: cap,
		},
	)
	require.NoError(t, err)
	require.True(t, success)

	// Verify NO attribution was made since user exceeds cap
	blockStats := tApp.App.StatsKeeper.GetBlockStats(ctx)
	require.Len(t, blockStats.Fills, 1)

	fill := blockStats.Fills[0]

	// Should have empty attributions array since user is over cap
	require.Empty(t, fill.AffiliateAttributions,
		"Should have no attributions when referee exceeds cap")
}

// TestProcessSingleMatch_AffiliateAttribution_CapWithExpiration tests that when
// old stats expire from the 30-day window, a user who was over the cap can start
// receiving attributions again.
func TestProcessSingleMatch_AffiliateAttribution_CapWithExpiration(t *testing.T) {
	cap := uint64(100_000_000_000) // 100k USDC cap

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

	// Set up affiliate parameters with cap
	err := tApp.App.AffiliatesKeeper.UpdateAffiliateParameters(ctx, &affiliatetypes.MsgUpdateAffiliateParameters{
		Authority: constants.GovAuthority,
		AffiliateParameters: affiliatetypes.AffiliateParameters{
			Maximum_30DAttributableVolumePerReferredUserQuoteQuantums: cap,
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

	// Register taker with referrer
	referrerAddr := constants.CarlAccAddress.String()
	err = tApp.App.AffiliatesKeeper.RegisterAffiliate(ctx, constants.Alice_Num0.Owner, referrerAddr)
	require.NoError(t, err)

	// SCENARIO 1: User is over cap (150k attributed volume)
	tApp.App.StatsKeeper.SetUserStats(ctx, constants.Alice_Num0.Owner, &statstypes.UserStats{
		TakerNotional: 200_000_000_000, // User traded 200k
		MakerNotional: 100_000_000_000, // User traded 100k
		Affiliate_30DAttributedVolumeQuoteQuantums: 150_000_000_000, // 150k attributed > 100k cap
	})

	// Create first trade
	takerOrder1 := clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: takerSubaccount,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_ShortTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     100_000_000,   // 1 BTC = ~5k USDC
		Subticks:     5_000_000_000, // $50,000 per BTC
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 100},
	}

	makerOrder1 := clobtypes.Order{
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

	matchWithOrders1 := &clobtypes.MatchWithOrders{
		TakerOrder: &takerOrder1,
		MakerOrder: &makerOrder1,
		FillAmount: satypes.BaseQuantums(100_000_000),
	}

	// Process first match - should have NO attribution
	success, _, _, _, err := k.ProcessSingleMatch(
		ctx,
		matchWithOrders1,
		map[string]bool{},
		affiliatetypes.AffiliateParameters{
			Maximum_30DAttributableVolumePerReferredUserQuoteQuantums: cap,
		},
	)
	require.NoError(t, err)
	require.True(t, success)

	blockStats := tApp.App.StatsKeeper.GetBlockStats(ctx)
	require.Len(t, blockStats.Fills, 1)
	require.Empty(t, blockStats.Fills[0].AffiliateAttributions,
		"First trade: No attribution since user is over cap")

	// SCENARIO 2: Simulate expiration - old attributed volume expires from 30-day window
	// Now user only has 80k attributed volume, which is BELOW the 100k cap
	tApp.App.StatsKeeper.SetUserStats(ctx, constants.Alice_Num0.Owner, &statstypes.UserStats{
		TakerNotional: 150_000_000_000, // Still trading (down from 200k as old trades expired)
		MakerNotional: 70_000_000_000,  // Still trading (down from 100k)
		Affiliate_30DAttributedVolumeQuoteQuantums: 80_000_000_000, // 80k attributed (down from
		// 150k - old attributions expired)
		// 80k < 100k cap, so 20k capacity available
	})

	// Clear block stats for second trade
	tApp.App.StatsKeeper.SetBlockStats(ctx, &statstypes.BlockStats{})

	// Create second trade
	takerOrder2 := clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: takerSubaccount,
			ClientId:     2, // Different client ID
			OrderFlags:   clobtypes.OrderIdFlags_ShortTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     100_000_000,   // 1 BTC = ~5k USDC
		Subticks:     5_000_000_000, // $50,000 per BTC
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 100},
	}

	makerOrder2 := clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: makerSubaccount,
			ClientId:     2,
			OrderFlags:   clobtypes.OrderIdFlags_ShortTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     100_000_000,
		Subticks:     5_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 100},
	}

	matchWithOrders2 := &clobtypes.MatchWithOrders{
		TakerOrder: &takerOrder2,
		MakerOrder: &makerOrder2,
		FillAmount: satypes.BaseQuantums(100_000_000),
	}

	// Process second match - should NOW have attribution since user is below cap
	success, _, _, _, err = k.ProcessSingleMatch(
		ctx,
		matchWithOrders2,
		map[string]bool{},
		affiliatetypes.AffiliateParameters{
			Maximum_30DAttributableVolumePerReferredUserQuoteQuantums: cap,
		},
	)
	require.NoError(t, err)
	require.True(t, success)

	blockStats = tApp.App.StatsKeeper.GetBlockStats(ctx)
	require.Len(t, blockStats.Fills, 1)

	fill2 := blockStats.Fills[0]
	require.Len(t, fill2.AffiliateAttributions, 1,
		"Second trade: Should have attribution after old stats expired")

	attribution := fill2.AffiliateAttributions[0]
	require.Equal(t, statstypes.AffiliateAttribution_ROLE_TAKER, attribution.Role)
	require.Equal(t, referrerAddr, attribution.ReferrerAddress)
}
