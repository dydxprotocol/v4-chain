package keeper_test

import (
	"math/big"
	"testing"
	"time"

	"cosmossdk.io/math"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	epochstypes "github.com/dydxprotocol/v4-chain/protocol/x/epochs/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/stats/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()

	logger := tApp.App.StatsKeeper.Logger(ctx)
	require.NotNil(t, logger)
}

type recordFillArgs struct {
	taker                 string
	maker                 string
	notional              *big.Int
	affiliateFee          *big.Int
	affiliateAttributions []*types.AffiliateAttribution
}

func TestRecordFill(t *testing.T) {
	tests := map[string]struct {
		args               []recordFillArgs
		expectedBlockStats *types.BlockStats
	}{
		"no fills": {
			[]recordFillArgs{},
			&types.BlockStats{
				Fills: nil,
			},
		},
		"single fill": {
			[]recordFillArgs{
				{
					taker:        "taker",
					maker:        "maker",
					notional:     new(big.Int).SetUint64(123),
					affiliateFee: big.NewInt(0),
					affiliateAttributions: []*types.AffiliateAttribution{
						{
							Role:                        types.AffiliateAttribution_ROLE_TAKER,
							ReferrerAddress:             "referrer",
							ReferredVolumeQuoteQuantums: 123,
						},
					},
				},
			},
			&types.BlockStats{
				Fills: []*types.BlockStats_Fill{
					{
						Taker:                         "taker",
						Maker:                         "maker",
						Notional:                      123,
						AffiliateFeeGeneratedQuantums: 0,
						AffiliateAttributions: []*types.AffiliateAttribution{
							{
								Role:                        types.AffiliateAttribution_ROLE_TAKER,
								ReferrerAddress:             "referrer",
								ReferredVolumeQuoteQuantums: 123,
							},
						},
					},
				},
			},
		},
		"multiple fills": {
			[]recordFillArgs{
				{
					taker:        "alice",
					maker:        "bob",
					notional:     new(big.Int).SetUint64(123),
					affiliateFee: big.NewInt(0),
					affiliateAttributions: []*types.AffiliateAttribution{
						{
							Role:                        types.AffiliateAttribution_ROLE_TAKER,
							ReferrerAddress:             "referrer",
							ReferredVolumeQuoteQuantums: 123,
						},
					},
				},
				{
					taker:        "bob",
					maker:        "alice",
					notional:     new(big.Int).SetUint64(321),
					affiliateFee: big.NewInt(0),
					affiliateAttributions: []*types.AffiliateAttribution{
						{
							Role:                        types.AffiliateAttribution_ROLE_TAKER,
							ReferrerAddress:             "referrer",
							ReferredVolumeQuoteQuantums: 321,
						},
					},
				},
			},
			&types.BlockStats{
				Fills: []*types.BlockStats_Fill{
					{
						Taker:                         "alice",
						Maker:                         "bob",
						Notional:                      123,
						AffiliateFeeGeneratedQuantums: 0,
						AffiliateAttributions: []*types.AffiliateAttribution{
							{
								Role:                        types.AffiliateAttribution_ROLE_TAKER,
								ReferrerAddress:             "referrer",
								ReferredVolumeQuoteQuantums: 123,
							},
						},
					},
					{
						Taker:                         "bob",
						Maker:                         "alice",
						Notional:                      321,
						AffiliateFeeGeneratedQuantums: 0,
						AffiliateAttributions: []*types.AffiliateAttribution{
							{
								Role:                        types.AffiliateAttribution_ROLE_TAKER,
								ReferrerAddress:             "referrer",
								ReferredVolumeQuoteQuantums: 321,
							},
						},
					},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.StatsKeeper

			for _, fill := range tc.args {
				k.RecordFill(ctx, fill.taker, fill.maker, fill.notional, fill.affiliateFee, fill.affiliateAttributions)
			}
			require.Equal(t, tc.expectedBlockStats, k.GetBlockStats(ctx))
		})
	}
}

func TestProcessBlockStats(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()

	// Epochs initialize at block height 2
	tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{
		BlockTime: time.Unix(1, 0).UTC(),
	})
	ctx := tApp.AdvanceToBlock(10, testapp.AdvanceToBlockOptions{
		BlockTime: time.Unix(int64(epochstypes.StatsEpochDuration)+1, 0).UTC(),
	})
	k := tApp.App.StatsKeeper

	k.SetBlockStats(ctx, &types.BlockStats{
		Fills: []*types.BlockStats_Fill{
			{
				Taker:    "alice",
				Maker:    "bob",
				Notional: 5,
			},
			{
				Taker:    "bob",
				Maker:    "alice",
				Notional: 10,
			},
		},
	})
	k.ProcessBlockStats(ctx)

	assert.Equal(t, &types.GlobalStats{
		NotionalTraded: 15,
	}, k.GetGlobalStats(ctx))
	assert.Equal(t, &types.UserStats{
		TakerNotional: 5,
		MakerNotional: 10,
	}, k.GetUserStats(ctx, "alice"))
	assert.Equal(t, &types.UserStats{
		TakerNotional: 10,
		MakerNotional: 5,
	}, k.GetUserStats(ctx, "bob"))
	assert.Equal(t, &types.EpochStats{
		EpochEndTime: time.Unix(7200, 0).UTC(),
		Stats: []*types.EpochStats_UserWithStats{
			{
				User: "alice",
				Stats: &types.UserStats{
					TakerNotional: 5,
					MakerNotional: 10,
				},
			},
			{
				User: "bob",
				Stats: &types.UserStats{
					TakerNotional: 10,
					MakerNotional: 5,
				},
			},
		},
	}, k.GetEpochStatsOrNil(ctx, 1))

	k.SetBlockStats(ctx, &types.BlockStats{
		Fills: []*types.BlockStats_Fill{
			{
				Taker:    "bob",
				Maker:    "alice",
				Notional: 10,
			},
		},
	})
	k.ProcessBlockStats(ctx)
	assert.Equal(t, &types.GlobalStats{
		NotionalTraded: 25,
	}, k.GetGlobalStats(ctx))
	assert.Equal(t, &types.UserStats{
		TakerNotional: 5,
		MakerNotional: 20,
	}, k.GetUserStats(ctx, "alice"))
	assert.Equal(t, &types.UserStats{
		TakerNotional: 20,
		MakerNotional: 5,
	}, k.GetUserStats(ctx, "bob"))
	assert.Equal(t, &types.EpochStats{
		EpochEndTime: time.Unix(7200, 0).UTC(),
		Stats: []*types.EpochStats_UserWithStats{
			{
				User: "alice",
				Stats: &types.UserStats{
					TakerNotional: 5,
					MakerNotional: 20,
				},
			},
			{
				User: "bob",
				Stats: &types.UserStats{
					TakerNotional: 20,
					MakerNotional: 5,
				},
			},
		},
	}, k.GetEpochStatsOrNil(ctx, 1))

	// Test affiliate revenue attribution
	k.SetBlockStats(ctx, &types.BlockStats{
		Fills: []*types.BlockStats_Fill{
			{
				Taker:                         "taker",
				Maker:                         "maker",
				Notional:                      100_000_000_000,
				AffiliateFeeGeneratedQuantums: 50_000_000,
				AffiliateAttributions: []*types.AffiliateAttribution{
					{
						Role:                        types.AffiliateAttribution_ROLE_TAKER,
						ReferrerAddress:             "referrer",
						ReferredVolumeQuoteQuantums: 100_000_000_000,
					},
				},
			},
		},
	})
	k.ProcessBlockStats(ctx)

	// Verify referrer's UserStats has the referred volume
	assert.Equal(t, &types.UserStats{
		Affiliate_30DReferredVolumeQuoteQuantums: 100_000_000_000,
	}, k.GetUserStats(ctx, "referrer"))

	// Verify taker has the affiliate fee generated AND attributed volume
	assert.Equal(t, &types.UserStats{
		TakerNotional:                              100_000_000_000,
		Affiliate_30DRevenueGeneratedQuantums:      50_000_000,
		Affiliate_30DAttributedVolumeQuoteQuantums: 100_000_000_000, // Taker's volume was attributed
	}, k.GetUserStats(ctx, "taker"))

	// Verify maker stats
	assert.Equal(t, &types.UserStats{
		MakerNotional: 100_000_000_000,
	}, k.GetUserStats(ctx, "maker"))

	// Verify global stats includes the new fill
	assert.Equal(t, &types.GlobalStats{
		NotionalTraded: 100_000_000_025,
	}, k.GetGlobalStats(ctx))

	// Verify referrer is in epoch stats with correct referred volume
	epochStats := k.GetEpochStatsOrNil(ctx, 1)
	require.NotNil(t, epochStats)
	var referrerFound bool
	for _, userStats := range epochStats.Stats {
		if userStats.User == "referrer" {
			referrerFound = true
			assert.Equal(t, uint64(100_000_000_000), userStats.Stats.Affiliate_30DReferredVolumeQuoteQuantums)
			break
		}
	}
	require.True(t, referrerFound, "referrer should be in epoch stats")

	// Test multiple fills with same referrer - referred volume should accumulate
	k.SetBlockStats(ctx, &types.BlockStats{
		Fills: []*types.BlockStats_Fill{
			{
				Taker:                         "taker2",
				Maker:                         "maker2",
				Notional:                      50_000_000_000,
				AffiliateFeeGeneratedQuantums: 25_000_000,
				AffiliateAttributions: []*types.AffiliateAttribution{
					{
						Role:                        types.AffiliateAttribution_ROLE_TAKER,
						ReferrerAddress:             "referrer",
						ReferredVolumeQuoteQuantums: 50_000_000_000,
					},
				},
			},
		},
	})
	k.ProcessBlockStats(ctx)

	// Verify referrer's referred volume accumulated
	assert.Equal(t, &types.UserStats{
		Affiliate_30DReferredVolumeQuoteQuantums: 150_000_000_000,
	}, k.GetUserStats(ctx, "referrer"))

	// Verify referrer's epoch stats accumulated
	epochStats = k.GetEpochStatsOrNil(ctx, 1)
	require.NotNil(t, epochStats)
	referrerFound = false
	for _, userStats := range epochStats.Stats {
		if userStats.User == "referrer" {
			referrerFound = true
			assert.Equal(t, uint64(150_000_000_000), userStats.Stats.Affiliate_30DReferredVolumeQuoteQuantums)
			break
		}
	}
	require.True(t, referrerFound, "referrer should be in epoch stats")

	// Test fill with capped attributable volume
	k.SetBlockStats(ctx, &types.BlockStats{
		Fills: []*types.BlockStats_Fill{
			{
				Taker:                         "taker3",
				Maker:                         "maker3",
				Notional:                      100_000_000_000,
				AffiliateFeeGeneratedQuantums: 50_000_000,
				AffiliateAttributions: []*types.AffiliateAttribution{
					{
						Role:                        types.AffiliateAttribution_ROLE_TAKER,
						ReferrerAddress:             "referrer2",
						ReferredVolumeQuoteQuantums: 30_000_000_000,
					},
				},
			},
		},
	})
	k.ProcessBlockStats(ctx)

	// Verify referrer2's referred volume reflects the capped amount
	assert.Equal(t, &types.UserStats{
		Affiliate_30DReferredVolumeQuoteQuantums: 30_000_000_000,
	}, k.GetUserStats(ctx, "referrer2"))

	// Verify referrer2's epoch stats has the capped amount
	epochStats = k.GetEpochStatsOrNil(ctx, 1)
	require.NotNil(t, epochStats)
	var referrer2Found bool
	for _, userStats := range epochStats.Stats {
		if userStats.User == "referrer2" {
			referrer2Found = true
			assert.Equal(t, uint64(30_000_000_000), userStats.Stats.Affiliate_30DReferredVolumeQuoteQuantums)
			break
		}
	}
	require.True(t, referrer2Found, "referrer2 should be in epoch stats")

	// Test fill without affiliate revenue attribution - should not affect referrer stats
	k.SetBlockStats(ctx, &types.BlockStats{
		Fills: []*types.BlockStats_Fill{
			{
				Taker:                         "taker4",
				Maker:                         "maker4",
				Notional:                      50_000_000_000,
				AffiliateFeeGeneratedQuantums: 0,
				AffiliateAttributions:         nil,
			},
		},
	})
	k.ProcessBlockStats(ctx)

	// Verify referrer stats unchanged
	assert.Equal(t, &types.UserStats{
		Affiliate_30DReferredVolumeQuoteQuantums: 150_000_000_000,
	}, k.GetUserStats(ctx, "referrer"))

	// Verify referrer2 stats unchanged
	assert.Equal(t, &types.UserStats{
		Affiliate_30DReferredVolumeQuoteQuantums: 30_000_000_000,
	}, k.GetUserStats(ctx, "referrer2"))

	// Test fill where both taker AND maker have affiliate attributions
	k.SetBlockStats(ctx, &types.BlockStats{
		Fills: []*types.BlockStats_Fill{
			{
				Taker:                         "taker5",
				Maker:                         "maker5",
				Notional:                      80_000_000_000,
				AffiliateFeeGeneratedQuantums: 40_000_000,
				AffiliateAttributions: []*types.AffiliateAttribution{
					{
						Role:                        types.AffiliateAttribution_ROLE_TAKER,
						ReferrerAddress:             "referrer_for_taker",
						ReferredVolumeQuoteQuantums: 80_000_000_000,
					},
					{
						Role:                        types.AffiliateAttribution_ROLE_MAKER,
						ReferrerAddress:             "referrer_for_maker",
						ReferredVolumeQuoteQuantums: 80_000_000_000,
					},
				},
			},
		},
	})
	k.ProcessBlockStats(ctx)

	// Verify taker's referrer received the attributed volume
	assert.Equal(t, &types.UserStats{
		Affiliate_30DReferredVolumeQuoteQuantums: 80_000_000_000,
	}, k.GetUserStats(ctx, "referrer_for_taker"))

	// Verify maker's referrer also received the attributed volume
	assert.Equal(t, &types.UserStats{
		Affiliate_30DReferredVolumeQuoteQuantums: 80_000_000_000,
	}, k.GetUserStats(ctx, "referrer_for_maker"))

	// Verify both referrers are in epoch stats
	epochStats = k.GetEpochStatsOrNil(ctx, 1)
	require.NotNil(t, epochStats)

	var takerReferrerFound, makerReferrerFound bool
	for _, userStats := range epochStats.Stats {
		if userStats.User == "referrer_for_taker" {
			takerReferrerFound = true
			assert.Equal(t, uint64(80_000_000_000), userStats.Stats.Affiliate_30DReferredVolumeQuoteQuantums)
		}
		if userStats.User == "referrer_for_maker" {
			makerReferrerFound = true
			assert.Equal(t, uint64(80_000_000_000), userStats.Stats.Affiliate_30DReferredVolumeQuoteQuantums)
		}
	}
	require.True(t, takerReferrerFound, "taker's referrer should be in epoch stats")
	require.True(t, makerReferrerFound, "maker's referrer should be in epoch stats")

	// Verify taker5 and maker5 stats (they're different addresses)
	taker5Stats := k.GetUserStats(ctx, "taker5")
	assert.Equal(t, &types.UserStats{
		TakerNotional:                              80_000_000_000,
		Affiliate_30DRevenueGeneratedQuantums:      40_000_000,
		Affiliate_30DAttributedVolumeQuoteQuantums: 80_000_000_000, // Taker's volume attributed
	}, taker5Stats)

	maker5Stats := k.GetUserStats(ctx, "maker5")
	assert.Equal(t, &types.UserStats{
		MakerNotional: 80_000_000_000,
		Affiliate_30DAttributedVolumeQuoteQuantums: 80_000_000_000, // Maker's volume attributed
	}, maker5Stats)
}

func TestExpireOldStats(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()

	// Epochs start at block height 2
	ctx := tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{
		BlockTime: time.Unix(int64(1), 0).UTC(),
	})
	windowDuration := tApp.App.StatsKeeper.GetWindowDuration(ctx)
	// 5 epochs are out of the window
	tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{
		BlockTime: time.Unix(0, 0).
			Add(windowDuration).
			Add((time.Duration(5*epochstypes.StatsEpochDuration) + 1) * time.Second).
			UTC(),
	})
	ctx = tApp.AdvanceToBlock(100, testapp.AdvanceToBlockOptions{})
	k := tApp.App.StatsKeeper

	// Create a bunch of EpochStats.
	// Odd epochs don't have stats. 30 epochs total.
	for i := 0; i < 30; i++ {
		k.SetEpochStats(ctx, uint32(i*2), &types.EpochStats{
			EpochEndTime: time.Unix(0, 0).
				Add(time.Duration(i*int(epochstypes.StatsEpochDuration)) * time.Second).
				UTC(),
			Stats: []*types.EpochStats_UserWithStats{
				{
					User: "alice",
					Stats: &types.UserStats{
						TakerNotional:                              1,
						MakerNotional:                              2,
						Affiliate_30DReferredVolumeQuoteQuantums:   10_000_000_000,
						Affiliate_30DRevenueGeneratedQuantums:      100_000_000,
						Affiliate_30DAttributedVolumeQuoteQuantums: 5_000_000_000, // 5k attributed per epoch
					},
				},
				{
					User: "bob",
					Stats: &types.UserStats{
						TakerNotional:                              2,
						MakerNotional:                              1,
						Affiliate_30DReferredVolumeQuoteQuantums:   10_000_000_000,
						Affiliate_30DRevenueGeneratedQuantums:      100_000_000,
						Affiliate_30DAttributedVolumeQuoteQuantums: 8_000_000_000, // 8k attributed per epoch
					},
				},
			},
		})
	}
	k.SetUserStats(ctx, "alice", &types.UserStats{
		TakerNotional:                              30,
		MakerNotional:                              60,
		Affiliate_30DReferredVolumeQuoteQuantums:   300_000_000_000,
		Affiliate_30DRevenueGeneratedQuantums:      3_000_000_000,
		Affiliate_30DAttributedVolumeQuoteQuantums: 150_000_000_000, // 30 epochs * 5k per epoch
	})
	k.SetUserStats(ctx, "bob", &types.UserStats{
		TakerNotional:                              60,
		MakerNotional:                              30,
		Affiliate_30DReferredVolumeQuoteQuantums:   300_000_000_000,
		Affiliate_30DRevenueGeneratedQuantums:      3_000_000_000,
		Affiliate_30DAttributedVolumeQuoteQuantums: 240_000_000_000, // 30 epochs * 8k per epoch
	})
	k.SetGlobalStats(ctx, &types.GlobalStats{
		NotionalTraded: 90,
	})
	k.SetStatsMetadata(ctx, &types.StatsMetadata{
		TrailingEpoch: 0,
	})

	// Prune epochs in batches of 2. For each pair, the second epoch is nil.
	// Epochs 1-10 pruned.
	for i := 0; i < 6; i++ {
		// EpochStats exist before pruning
		require.NotNil(t, k.GetEpochStatsOrNil(ctx, uint32(i*2)))

		k.ExpireOldStats(ctx)
		require.Equal(t, &types.UserStats{
			TakerNotional:                              30 - uint64(i+1),
			MakerNotional:                              60 - 2*uint64(i+1),
			Affiliate_30DReferredVolumeQuoteQuantums:   300_000_000_000 - (uint64(i+1) * 10_000_000_000),
			Affiliate_30DRevenueGeneratedQuantums:      3_000_000_000 - (uint64(i+1) * 100_000_000),
			Affiliate_30DAttributedVolumeQuoteQuantums: 150_000_000_000 - (uint64(i+1) * 5_000_000_000),
			// Decreases by 5k per expired epoch
		}, k.GetUserStats(ctx, "alice"))
		require.Equal(t, &types.UserStats{
			TakerNotional:                              60 - 2*uint64(i+1),
			MakerNotional:                              30 - uint64(i+1),
			Affiliate_30DReferredVolumeQuoteQuantums:   300_000_000_000 - (uint64(i+1) * 10_000_000_000),
			Affiliate_30DRevenueGeneratedQuantums:      3_000_000_000 - (uint64(i+1) * 100_000_000),
			Affiliate_30DAttributedVolumeQuoteQuantums: 240_000_000_000 - (uint64(i+1) * 8_000_000_000),
			// Decreases by 8k per expired epoch
		}, k.GetUserStats(ctx, "bob"))
		require.Equal(t, &types.GlobalStats{
			NotionalTraded: 90 - 3*uint64(i+1),
		}, k.GetGlobalStats(ctx))

		// EpochStats removed
		require.Nil(t, k.GetEpochStatsOrNil(ctx, uint32(i*2)))

		k.ExpireOldStats(ctx)

		// Unchanged after pruning nil epoch
		require.Equal(t, &types.UserStats{
			TakerNotional:                              30 - uint64(i+1),
			MakerNotional:                              60 - 2*uint64(i+1),
			Affiliate_30DReferredVolumeQuoteQuantums:   300_000_000_000 - (uint64(i+1) * 10_000_000_000),
			Affiliate_30DRevenueGeneratedQuantums:      3_000_000_000 - (uint64(i+1) * 100_000_000),
			Affiliate_30DAttributedVolumeQuoteQuantums: 150_000_000_000 - (uint64(i+1) * 5_000_000_000),
		}, k.GetUserStats(ctx, "alice"))
		require.Equal(t, &types.UserStats{
			TakerNotional:                              60 - 2*uint64(i+1),
			MakerNotional:                              30 - uint64(i+1),
			Affiliate_30DReferredVolumeQuoteQuantums:   300_000_000_000 - (uint64(i+1) * 10_000_000_000),
			Affiliate_30DRevenueGeneratedQuantums:      3_000_000_000 - (uint64(i+1) * 100_000_000),
			Affiliate_30DAttributedVolumeQuoteQuantums: 240_000_000_000 - (uint64(i+1) * 8_000_000_000),
		}, k.GetUserStats(ctx, "bob"))
		require.Equal(t, &types.GlobalStats{
			NotionalTraded: 90 - 3*uint64(i+1),
		}, k.GetGlobalStats(ctx))
	}

	// Epoch 12 is within the window so it won't get pruned
	k.ExpireOldStats(ctx)
	k.ExpireOldStats(ctx)
	k.ExpireOldStats(ctx)
	require.NotNil(t, k.GetEpochStatsOrNil(ctx, uint32(12)))
}

// TestAffiliateAttribution_ConsistentlyHighVolumeTrader tests the scenario where a user
// is consistently trading at high volume and hitting the attribution cap.
// This test proves there is NO equilibrium trap - the user can continue getting
// attribution as old stats expire, even while trading continuously.
func TestAffiliateAttribution_ConsistentlyHighVolumeTrader(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()

	// Epochs start at block height 2
	ctx := tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{
		BlockTime: time.Unix(int64(1), 0).UTC(),
	})

	// Advance time so first 5 epochs will be ready to expire
	windowDuration := tApp.App.StatsKeeper.GetWindowDuration(ctx)
	tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{
		BlockTime: time.Unix(0, 0).
			Add(windowDuration).
			Add((time.Duration(5*epochstypes.StatsEpochDuration) + 1) * time.Second).
			UTC(),
	})

	// Now advance to a stable block and set up stats
	ctx = tApp.AdvanceToBlock(100, testapp.AdvanceToBlockOptions{})
	k := tApp.App.StatsKeeper

	// Scenario: User is a consistent high-volume trader
	// - Cap is 100k
	// - User has been trading 200k per epoch for 30 epochs
	// - User has attributed 100k total (at cap)
	// - Old epochs will expire, allowing new attribution

	// Create 30 epochs of history where the user traded 200k each epoch
	// but only 100k total was attributed (due to cap being reached early)
	for i := 0; i < 30; i++ {
		var attributedThisEpoch uint64
		if i < 5 {
			// First 5 epochs: 20k attributed per epoch = 100k total (reaches cap)
			attributedThisEpoch = 20_000_000_000
		} else {
			// Epochs 6-30: 0 attributed (already at cap)
			attributedThisEpoch = 0
		}

		k.SetEpochStats(ctx, uint32(i*2), &types.EpochStats{
			EpochEndTime: time.Unix(0, 0).
				Add(time.Duration(i*int(epochstypes.StatsEpochDuration)) * time.Second).
				UTC(),
			Stats: []*types.EpochStats_UserWithStats{
				{
					User: "highVolumeTrader",
					Stats: &types.UserStats{
						TakerNotional:                              200_000_000_000, // 200k traded
						MakerNotional:                              0,
						Affiliate_30DReferredVolumeQuoteQuantums:   0,
						Affiliate_30DRevenueGeneratedQuantums:      0,
						Affiliate_30DAttributedVolumeQuoteQuantums: attributedThisEpoch, // Only first
						// 5 epochs have attribution
					},
				},
				{
					User: "someAffiliate",
					Stats: &types.UserStats{
						TakerNotional:                            0,
						MakerNotional:                            0,
						Affiliate_30DReferredVolumeQuoteQuantums: attributedThisEpoch, // Same as attributed
						// (this is what the affiliate refers)
						Affiliate_30DRevenueGeneratedQuantums: 0,
					},
				},
			},
		})
	}

	// User's current stats: 6000k traded (30 epochs × 200k), but only 100k attributed
	k.SetUserStats(ctx, "highVolumeTrader", &types.UserStats{
		TakerNotional:                              6000_000_000_000, // 6M total volume
		MakerNotional:                              0,
		Affiliate_30DReferredVolumeQuoteQuantums:   0,
		Affiliate_30DRevenueGeneratedQuantums:      0,
		Affiliate_30DAttributedVolumeQuoteQuantums: 100_000_000_000, // At cap (100k)
	})

	k.SetStatsMetadata(ctx, &types.StatsMetadata{
		TrailingEpoch: 0,
	})

	// Set up affiliate's initial referred volume (also at cap from same history)
	k.SetUserStats(ctx, "someAffiliate", &types.UserStats{
		TakerNotional:                            0,
		MakerNotional:                            0,
		Affiliate_30DReferredVolumeQuoteQuantums: 100_000_000_000, // Also at 100k
		Affiliate_30DRevenueGeneratedQuantums:    0,
	})

	// BEFORE expiration: User is at cap (100k)
	userStats := k.GetUserStats(ctx, "highVolumeTrader")
	require.Equal(t, uint64(100_000_000_000), userStats.Affiliate_30DAttributedVolumeQuoteQuantums,
		"User should start at cap")

	affiliateStats := k.GetUserStats(ctx, "someAffiliate")
	require.Equal(t, uint64(100_000_000_000), affiliateStats.Affiliate_30DReferredVolumeQuoteQuantums,
		"Affiliate should start with 100k referred volume")

	// NOW INTERWEAVE: Expire old stats while processing new blocks with attribution
	// This simulates the realistic scenario where a high-volume trader continues trading
	// as old attributed volume expires, keeping their attributed volume at/near the cap

	// Expire first epoch (20k attributed removed) + Process new block (20k attributed added)
	k.ExpireOldStats(ctx) // Removes 20k from epoch 0
	k.SetBlockStats(ctx, &types.BlockStats{
		Fills: []*types.BlockStats_Fill{
			{
				Taker:    "highVolumeTrader",
				Maker:    "someMaker",
				Notional: 20_000_000_000,
				AffiliateAttributions: []*types.AffiliateAttribution{
					{
						Role:                        types.AffiliateAttribution_ROLE_TAKER,
						ReferrerAddress:             "someAffiliate",
						ReferredVolumeQuoteQuantums: 20_000_000_000,
					},
				},
			},
		},
	})
	k.ProcessBlockStats(ctx)

	userStats = k.GetUserStats(ctx, "highVolumeTrader")
	require.Equal(t, uint64(100_000_000_000), userStats.Affiliate_30DAttributedVolumeQuoteQuantums,
		"Attributed volume stays at cap: 100k - 20k (expired) + 20k (new) = 100k")

	affiliateStats = k.GetUserStats(ctx, "someAffiliate")
	require.Equal(t, uint64(100_000_000_000), affiliateStats.Affiliate_30DReferredVolumeQuoteQuantums,
		"Affiliate referred volume stays at maximum: 100k")

	// Skip nil epoch
	k.ExpireOldStats(ctx)

	// Expire second epoch (20k removed) + Process new block (20k added)
	k.ExpireOldStats(ctx) // Removes 20k from epoch 2
	k.SetBlockStats(ctx, &types.BlockStats{
		Fills: []*types.BlockStats_Fill{
			{
				Taker:    "highVolumeTrader",
				Maker:    "someMaker",
				Notional: 20_000_000_000,
				AffiliateAttributions: []*types.AffiliateAttribution{
					{
						Role:                        types.AffiliateAttribution_ROLE_TAKER,
						ReferrerAddress:             "someAffiliate",
						ReferredVolumeQuoteQuantums: 20_000_000_000,
					},
				},
			},
		},
	})
	k.ProcessBlockStats(ctx)

	userStats = k.GetUserStats(ctx, "highVolumeTrader")
	require.Equal(t, uint64(100_000_000_000), userStats.Affiliate_30DAttributedVolumeQuoteQuantums,
		"Attributed volume stays at cap: still 100k after second rotation")

	affiliateStats = k.GetUserStats(ctx, "someAffiliate")
	require.Equal(t, uint64(100_000_000_000), affiliateStats.Affiliate_30DReferredVolumeQuoteQuantums,
		"Affiliate referred volume stays at maximum: 100k")

	// Skip nil epoch
	k.ExpireOldStats(ctx)

	// Expire third epoch (20k removed) + Process new block (20k added)
	k.ExpireOldStats(ctx)
	k.SetBlockStats(ctx, &types.BlockStats{
		Fills: []*types.BlockStats_Fill{
			{
				Taker:    "highVolumeTrader",
				Maker:    "someMaker",
				Notional: 20_000_000_000,
				AffiliateAttributions: []*types.AffiliateAttribution{
					{
						Role:                        types.AffiliateAttribution_ROLE_TAKER,
						ReferrerAddress:             "someAffiliate",
						ReferredVolumeQuoteQuantums: 20_000_000_000,
					},
				},
			},
		},
	})
	k.ProcessBlockStats(ctx)

	userStats = k.GetUserStats(ctx, "highVolumeTrader")
	require.Equal(t, uint64(100_000_000_000), userStats.Affiliate_30DAttributedVolumeQuoteQuantums,
		"Attributed volume stays at cap: still 100k after third rotation")

	affiliateStats = k.GetUserStats(ctx, "someAffiliate")
	require.Equal(t, uint64(100_000_000_000), affiliateStats.Affiliate_30DReferredVolumeQuoteQuantums,
		"Affiliate referred volume stays at maximum: 100k")

	// Skip nil epoch
	k.ExpireOldStats(ctx)

	// Expire fourth epoch (20k removed) + Process new block (20k added)
	k.ExpireOldStats(ctx)
	k.SetBlockStats(ctx, &types.BlockStats{
		Fills: []*types.BlockStats_Fill{
			{
				Taker:    "highVolumeTrader",
				Maker:    "someMaker",
				Notional: 20_000_000_000,
				AffiliateAttributions: []*types.AffiliateAttribution{
					{
						Role:                        types.AffiliateAttribution_ROLE_TAKER,
						ReferrerAddress:             "someAffiliate",
						ReferredVolumeQuoteQuantums: 20_000_000_000,
					},
				},
			},
		},
	})
	k.ProcessBlockStats(ctx)

	userStats = k.GetUserStats(ctx, "highVolumeTrader")
	require.Equal(t, uint64(100_000_000_000), userStats.Affiliate_30DAttributedVolumeQuoteQuantums,
		"Attributed volume stays at cap: still 100k after fourth rotation")

	affiliateStats = k.GetUserStats(ctx, "someAffiliate")
	require.Equal(t, uint64(100_000_000_000), affiliateStats.Affiliate_30DReferredVolumeQuoteQuantums,
		"Affiliate referred volume stays at maximum: 100k")

	// Skip nil epoch
	k.ExpireOldStats(ctx)

	// Expire fifth epoch (20k removed) + Process new block (20k added)
	k.ExpireOldStats(ctx)
	k.SetBlockStats(ctx, &types.BlockStats{
		Fills: []*types.BlockStats_Fill{
			{
				Taker:    "highVolumeTrader",
				Maker:    "someMaker",
				Notional: 20_000_000_000,
				AffiliateAttributions: []*types.AffiliateAttribution{
					{
						Role:                        types.AffiliateAttribution_ROLE_TAKER,
						ReferrerAddress:             "someAffiliate",
						ReferredVolumeQuoteQuantums: 20_000_000_000,
					},
				},
			},
		},
	})
	k.ProcessBlockStats(ctx)

	userStats = k.GetUserStats(ctx, "highVolumeTrader")
	require.Equal(t, uint64(100_000_000_000), userStats.Affiliate_30DAttributedVolumeQuoteQuantums,
		"Attributed volume STAYS AT MAXIMUM: 100k throughout the entire rotation!")

	affiliateStats = k.GetUserStats(ctx, "someAffiliate")
	require.Equal(t, uint64(100_000_000_000), affiliateStats.Affiliate_30DReferredVolumeQuoteQuantums,
		"Affiliate referred volume STAYS AT MAXIMUM: 100k throughout the entire rotation!")
}

func TestGetStakedBaseTokens(t *testing.T) {
	testCases := []struct {
		name                 string
		userShares           uint32
		validatorTotalTokens uint32
		validatorTotalShares uint32
	}{
		{
			name:                 "1 share = 1 base token",
			userShares:           100,
			validatorTotalTokens: 1000,
			validatorTotalShares: 1000,
		},
		{
			name:                 "1 share = 1.5 base tokens",
			userShares:           100,
			validatorTotalTokens: 1500,
			validatorTotalShares: 1000,
		},
		{
			name:                 "1 share = 2 base tokens",
			userShares:           100,
			validatorTotalTokens: 2000,
			validatorTotalShares: 1000,
		},
		{
			name:                 "1 share = 0.5 base tokens",
			userShares:           100,
			validatorTotalTokens: 500,
			validatorTotalShares: 1000,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			statsKeeper := tApp.App.StatsKeeper
			stakingKeeper := tApp.App.StakingKeeper

			expMultiplier, _ := lib.BigPow10(-lib.BaseDenomExponent)

			// Set validator's tokens and shares
			validator, err := stakingKeeper.GetValidator(ctx, constants.AliceValAddress)
			require.NoError(t, err)

			validator.Tokens = math.NewIntFromBigInt(
				new(big.Int).Mul(lib.BigU(tc.validatorTotalTokens), expMultiplier),
			)
			validator.DelegatorShares = math.LegacyNewDecFromBigInt(
				new(big.Int).Mul(lib.BigU(tc.validatorTotalShares), expMultiplier),
			)
			err = stakingKeeper.SetValidator(ctx, validator)
			require.NoError(t, err)

			// Create user delegation
			userSharesBigInt := new(big.Int).Mul(lib.BigU(tc.userShares), expMultiplier)
			delegation := stakingtypes.NewDelegation(
				constants.AliceAccAddress.String(),
				constants.AliceValAddress.String(),
				math.LegacyNewDecFromBigInt(userSharesBigInt),
			)
			err = stakingKeeper.SetDelegation(ctx, delegation)
			require.NoError(t, err)

			// User should have `userShares * (validator tokens / validator shares)` number of base tokens
			expectedTokens := new(big.Int).Mul(userSharesBigInt, lib.BigU(tc.validatorTotalTokens))
			expectedTokens.Div(expectedTokens, lib.BigU(tc.validatorTotalShares))

			actualTokens := statsKeeper.GetStakedBaseTokens(ctx, constants.AliceAccAddress.String())
			require.Equal(t, expectedTokens, actualTokens)
		})
	}
}

func TestGetStakedBaseTokens_Cache_Hit(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	statsKeeper := tApp.App.StatsKeeper
	expMultiplier, _ := lib.BigPow10(-lib.BaseDenomExponent)
	coinsToStakeQuantums := new(big.Int).Mul(
		lib.BigI(100),
		expMultiplier,
	)
	statsKeeper.UnsafeSetCachedStakedBaseTokens(ctx, constants.AliceAccAddress.String(), &types.CachedStakedBaseTokens{
		StakedBaseTokens: dtypes.NewIntFromBigInt(coinsToStakeQuantums),
		CachedAt:         ctx.BlockTime().Unix(),
	})

	receivedCoins := statsKeeper.GetStakedBaseTokens(ctx, constants.AliceAccAddress.String())
	require.Equal(t, coinsToStakeQuantums, receivedCoins)
}

func TestGetStakedBaseTokens_Cache_Miss(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	statsKeeper := tApp.App.StatsKeeper
	stakingKeeper := tApp.App.StakingKeeper

	expMultiplier, _ := lib.BigPow10(-lib.BaseDenomExponent)
	expiredWholeCoinsToStake := 100
	latestWholeCoinsToStake := 200
	expiredCoinsToStakeQuantums := new(big.Int).Mul(
		lib.BigI(expiredWholeCoinsToStake),
		expMultiplier,
	)
	latestCoinsToStakeQuantums := new(big.Int).Mul(
		lib.BigI(latestWholeCoinsToStake),
		expMultiplier,
	)

	// set expired delegation
	statsKeeper.UnsafeSetCachedStakedBaseTokens(ctx, constants.AliceAccAddress.String(), &types.CachedStakedBaseTokens{
		StakedBaseTokens: dtypes.NewIntFromBigInt(expiredCoinsToStakeQuantums),
		CachedAt:         ctx.BlockTime().Unix(),
	})

	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(time.Duration(epochstypes.StatsEpochDuration+1) * time.Second))

	delegation := stakingtypes.NewDelegation(
		constants.AliceAccAddress.String(), constants.AliceValAddress.String(),
		math.LegacyNewDecFromBigInt(latestCoinsToStakeQuantums))
	err := stakingKeeper.SetDelegation(ctx, delegation)
	require.NoError(t, err)

	receivedCoins := statsKeeper.GetStakedBaseTokens(ctx, constants.AliceAccAddress.String())
	require.Equal(t, latestCoinsToStakeQuantums, receivedCoins)
}

// TestExpireOldStats_UnderflowProtection tests that affiliate fields are properly
// clamped to 0 when expiring epochs would cause underflow due to corrupted/inconsistent data.
func TestExpireOldStats_UnderflowProtection(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()

	// Epochs start at block height 2
	ctx := tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{
		BlockTime: time.Unix(int64(1), 0).UTC(),
	})
	windowDuration := tApp.App.StatsKeeper.GetWindowDuration(ctx)

	// Advance time so epochs can expire
	tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{
		BlockTime: time.Unix(0, 0).
			Add(windowDuration).
			Add((time.Duration(2*epochstypes.StatsEpochDuration) + 1) * time.Second).
			UTC(),
	})
	ctx = tApp.AdvanceToBlock(100, testapp.AdvanceToBlockOptions{})
	k := tApp.App.StatsKeeper

	// Simulate a scenario where epoch stats have MORE volume than user's current stats
	// This could happen after the v9.5 migration or other data inconsistencies

	// Create epoch 0 with large affiliate volumes
	k.SetEpochStats(ctx, 0, &types.EpochStats{
		EpochEndTime: time.Unix(0, 0).UTC(),
		Stats: []*types.EpochStats_UserWithStats{
			{
				User: "alice",
				Stats: &types.UserStats{
					TakerNotional:                              100,
					MakerNotional:                              200,
					Affiliate_30DReferredVolumeQuoteQuantums:   1_000_000_000_000, // 1M volume in epoch
					Affiliate_30DAttributedVolumeQuoteQuantums: 500_000_000_000,   // 500k attributed
				},
			},
		},
	})

	// Set user stats with LESS than what's in the epoch
	// This simulates corrupted/inconsistent state
	k.SetUserStats(ctx, "alice", &types.UserStats{
		TakerNotional:                              50,              // Less than epoch
		MakerNotional:                              100,             // Less than epoch
		Affiliate_30DReferredVolumeQuoteQuantums:   100_000_000_000, // Only 100k, but epoch has 1M
		Affiliate_30DAttributedVolumeQuoteQuantums: 50_000_000_000,  // Only 50k, but epoch has 500k
	})

	k.SetGlobalStats(ctx, &types.GlobalStats{
		NotionalTraded: 150,
	})

	k.SetStatsMetadata(ctx, &types.StatsMetadata{
		TrailingEpoch: 0,
	})

	// Expire the epoch - this should NOT cause underflow
	k.ExpireOldStats(ctx)

	// Verify that affiliate fields are clamped to 0, not wrapped around to huge numbers
	aliceStats := k.GetUserStats(ctx, "alice")

	// TakerNotional and MakerNotional can go negative (they wrap around for uint64)
	// But we're testing that the affiliate fields with underflow protection work correctly

	// These fields should be clamped to 0, not underflow
	require.Equal(t, uint64(0), aliceStats.Affiliate_30DReferredVolumeQuoteQuantums,
		"Referred volume should be clamped to 0, not underflow")
	require.Equal(t, uint64(0), aliceStats.Affiliate_30DAttributedVolumeQuoteQuantums,
		"Attributed volume should be clamped to 0, not underflow")

	// Verify the values aren't huge wrapped-around numbers
	// If underflow occurred, these would be close to uint64 max (18446744073709551615)
	require.Less(t, aliceStats.Affiliate_30DReferredVolumeQuoteQuantums, uint64(1000),
		"Referred volume should not be a wrapped-around huge number")
	require.Less(t, aliceStats.Affiliate_30DAttributedVolumeQuoteQuantums, uint64(1000),
		"Attributed volume should not be a wrapped-around huge number")
}

// TestExpireOldStats_UnderflowProtection_EdgeCase tests the exact boundary case
func TestExpireOldStats_UnderflowProtection_EdgeCase(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()

	ctx := tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{
		BlockTime: time.Unix(int64(1), 0).UTC(),
	})
	windowDuration := tApp.App.StatsKeeper.GetWindowDuration(ctx)

	tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{
		BlockTime: time.Unix(0, 0).
			Add(windowDuration).
			Add((time.Duration(2*epochstypes.StatsEpochDuration) + 1) * time.Second).
			UTC(),
	})
	ctx = tApp.AdvanceToBlock(100, testapp.AdvanceToBlockOptions{})
	k := tApp.App.StatsKeeper

	// Test exact boundary: user has exactly the same amount as epoch
	k.SetEpochStats(ctx, 0, &types.EpochStats{
		EpochEndTime: time.Unix(0, 0).UTC(),
		Stats: []*types.EpochStats_UserWithStats{
			{
				User: "bob",
				Stats: &types.UserStats{
					Affiliate_30DReferredVolumeQuoteQuantums:   1000,
					Affiliate_30DAttributedVolumeQuoteQuantums: 500,
				},
			},
		},
	})

	k.SetUserStats(ctx, "bob", &types.UserStats{
		Affiliate_30DReferredVolumeQuoteQuantums:   1000, // Exactly the same
		Affiliate_30DAttributedVolumeQuoteQuantums: 500,  // Exactly the same
	})

	k.SetStatsMetadata(ctx, &types.StatsMetadata{
		TrailingEpoch: 0,
	})

	k.ExpireOldStats(ctx)

	bobStats := k.GetUserStats(ctx, "bob")

	// Should be exactly 0 after subtracting equal amounts
	require.Equal(t, uint64(0), bobStats.Affiliate_30DReferredVolumeQuoteQuantums,
		"Should be exactly 0 after subtracting equal amounts")
	require.Equal(t, uint64(0), bobStats.Affiliate_30DAttributedVolumeQuoteQuantums,
		"Should be exactly 0 after subtracting equal amounts")
}

// TestExpireOldStats_UnderflowProtection_MultipleUsers tests that underflow
// protection works correctly when multiple users have inconsistent data
func TestExpireOldStats_UnderflowProtection_MultipleUsers(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()

	ctx := tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{
		BlockTime: time.Unix(int64(1), 0).UTC(),
	})
	windowDuration := tApp.App.StatsKeeper.GetWindowDuration(ctx)

	tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{
		BlockTime: time.Unix(0, 0).
			Add(windowDuration).
			Add((time.Duration(2*epochstypes.StatsEpochDuration) + 1) * time.Second).
			UTC(),
	})
	ctx = tApp.AdvanceToBlock(100, testapp.AdvanceToBlockOptions{})
	k := tApp.App.StatsKeeper

	// Create epoch with multiple users having different underflow scenarios
	k.SetEpochStats(ctx, 0, &types.EpochStats{
		EpochEndTime: time.Unix(0, 0).UTC(),
		Stats: []*types.EpochStats_UserWithStats{
			{
				User: "alice",
				Stats: &types.UserStats{
					Affiliate_30DReferredVolumeQuoteQuantums:   5_000_000_000,
					Affiliate_30DAttributedVolumeQuoteQuantums: 3_000_000_000,
				},
			},
			{
				User: "bob",
				Stats: &types.UserStats{
					Affiliate_30DReferredVolumeQuoteQuantums:   10_000_000_000,
					Affiliate_30DAttributedVolumeQuoteQuantums: 7_000_000_000,
				},
			},
			{
				User: "carl",
				Stats: &types.UserStats{
					Affiliate_30DReferredVolumeQuoteQuantums:   2_000_000_000,
					Affiliate_30DAttributedVolumeQuoteQuantums: 1_000_000_000,
				},
			},
		},
	})

	// Alice: has more than epoch (normal case)
	k.SetUserStats(ctx, "alice", &types.UserStats{
		Affiliate_30DReferredVolumeQuoteQuantums:   10_000_000_000,
		Affiliate_30DAttributedVolumeQuoteQuantums: 8_000_000_000,
	})

	// Bob: has less than epoch (should clamp to 0)
	k.SetUserStats(ctx, "bob", &types.UserStats{
		Affiliate_30DReferredVolumeQuoteQuantums:   5_000_000_000,
		Affiliate_30DAttributedVolumeQuoteQuantums: 3_000_000_000,
	})

	// Carl: has exactly the same (should become 0)
	k.SetUserStats(ctx, "carl", &types.UserStats{
		Affiliate_30DReferredVolumeQuoteQuantums:   2_000_000_000,
		Affiliate_30DAttributedVolumeQuoteQuantums: 1_000_000_000,
	})

	k.SetStatsMetadata(ctx, &types.StatsMetadata{
		TrailingEpoch: 0,
	})

	k.ExpireOldStats(ctx)

	// Verify Alice: normal subtraction
	aliceStats := k.GetUserStats(ctx, "alice")
	require.Equal(t, uint64(5_000_000_000), aliceStats.Affiliate_30DReferredVolumeQuoteQuantums,
		"Alice should have 5B remaining (10B - 5B)")
	require.Equal(t, uint64(5_000_000_000), aliceStats.Affiliate_30DAttributedVolumeQuoteQuantums,
		"Alice should have 5B remaining (8B - 3B)")

	// Verify Bob: clamped to 0
	bobStats := k.GetUserStats(ctx, "bob")
	require.Equal(t, uint64(0), bobStats.Affiliate_30DReferredVolumeQuoteQuantums,
		"Bob should be clamped to 0 (had 5B, tried to subtract 10B)")
	require.Equal(t, uint64(0), bobStats.Affiliate_30DAttributedVolumeQuoteQuantums,
		"Bob should be clamped to 0 (had 3B, tried to subtract 7B)")

	// Verify Carl: exactly 0
	carlStats := k.GetUserStats(ctx, "carl")
	require.Equal(t, uint64(0), carlStats.Affiliate_30DReferredVolumeQuoteQuantums,
		"Carl should be exactly 0 (had 2B, subtracted 2B)")
	require.Equal(t, uint64(0), carlStats.Affiliate_30DAttributedVolumeQuoteQuantums,
		"Carl should be exactly 0 (had 1B, subtracted 1B)")
}
