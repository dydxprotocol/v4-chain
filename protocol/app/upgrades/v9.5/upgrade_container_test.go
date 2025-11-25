package v_9_5_test

import (
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"

	v_9_5 "github.com/dydxprotocol/v4-chain/protocol/app/upgrades/v9.5"
	"github.com/dydxprotocol/v4-chain/protocol/testing/containertest"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	epochstypes "github.com/dydxprotocol/v4-chain/protocol/x/epochs/types"
	statstypes "github.com/dydxprotocol/v4-chain/protocol/x/stats/types"
)
func TestMigrate30dReferredVolumeToEpochStats(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	tApp.InitChain()

	statsKeeper := tApp.App.StatsKeeper
	epochsKeeper := tApp.App.EpochsKeeper

	// Advance to the next epoch so we have a current epoch with some activity
	ctx := tApp.AdvanceToBlock(10, testapp.AdvanceToBlockOptions{})

	// Get current epoch directly from keeper instead of using containertest.Query
	statsEpochInfo := epochsKeeper.MustGetStatsEpochInfo(ctx)
	currentEpoch := statsEpochInfo.CurrentEpoch

	// Create some initial epoch stats (simulating existing trading activity)
	initialEpochStats := &statstypes.EpochStats{
		Stats: []*statstypes.EpochStats_UserWithStats{
			{
				User: constants.AliceAccAddress.String(),
				Stats: &statstypes.UserStats{
					TakerNotional: 50,
					MakerNotional: 75,
					// No referred volume set yet
				},
			},
			{
				User: constants.BobAccAddress.String(),
				Stats: &statstypes.UserStats{
					TakerNotional: 100,
					MakerNotional: 150,
					// No referred volume set yet
				},
			},
		},
	}
	statsKeeper.SetEpochStats(ctx, currentEpoch, initialEpochStats)

	// Set up global user stats with referred volume (this represents the 30d cumulative volume)
	testUsers := []struct {
		address        string
		referredVolume uint64
	}{
		{constants.AliceAccAddress.String(), 1_000_000_000}, // 1k volume
		{constants.BobAccAddress.String(), 5_000_000_000},   // 5k volume
		{constants.CarlAccAddress.String(), 10_000_000_000}, // 10k volume
		{constants.DaveAccAddress.String(), 0},              // no referred volume
	}

	for _, user := range testUsers {
		userStats := &statstypes.UserStats{
			TakerNotional:                            100,
			MakerNotional:                            200,
			Affiliate_30DReferredVolumeQuoteQuantums: user.referredVolume,
		}
		statsKeeper.SetUserStats(ctx, user.address, userStats)
	}

	// Verify initial state - epoch stats should not have referred volume
	preUpgradeEpochStats := statsKeeper.GetEpochStatsOrNil(ctx, currentEpoch)
	require.NotNil(t, preUpgradeEpochStats)
	for _, userStats := range preUpgradeEpochStats.Stats {
		require.Equal(t, uint64(0), userStats.Stats.Affiliate_30DReferredVolumeQuoteQuantums,
			"Referred volume should be 0 before migration for user %s", userStats.User)
	}

	// Run the migration function directly
	v_9_5.Migrate30dReferredVolumeToEpochStats(ctx, statsKeeper, epochsKeeper)

	// Verify migration results
	postUpgradeEpochStats := statsKeeper.GetEpochStatsOrNil(ctx, currentEpoch)
	require.NotNil(t, postUpgradeEpochStats)

	// Create a map for easier verification

func TestMigrate30dReferredVolumeToEpochStats_EmptyEpochStats(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	tApp.InitChain()

	statsKeeper := tApp.App.StatsKeeper
	epochsKeeper := tApp.App.EpochsKeeper

	// Advance to next epoch
	ctx := tApp.AdvanceToBlock(10, testapp.AdvanceToBlockOptions{})

	// Get current epoch directly from keeper instead of using containertest.Query
	statsEpochInfo := epochsKeeper.MustGetStatsEpochInfo(ctx)
	currentEpoch := statsEpochInfo.CurrentEpoch

	// Setup: Create user stats with referred volume but no epoch stats
	userStats := &statstypes.UserStats{
		TakerNotional:                            100,
		MakerNotional:                            200,
		Affiliate_30DReferredVolumeQuoteQuantums: 1_000_000_000,
	}
	statsKeeper.SetUserStats(ctx, constants.AliceAccAddress.String(), userStats)

	// Verify no epoch stats exist initially
	preUpgradeEpochStats := statsKeeper.GetEpochStatsOrNil(ctx, currentEpoch)
	require.Nil(t, preUpgradeEpochStats)

	// Run the migration function directly
	v_9_5.Migrate30dReferredVolumeToEpochStats(ctx, statsKeeper, epochsKeeper)

	// Verify Alice was added to epoch stats even though she wasn't trading
	postUpgradeEpochStats := statsKeeper.GetEpochStatsOrNil(ctx, currentEpoch)
	require.NotNil(t, postUpgradeEpochStats)
	require.Len(t, postUpgradeEpochStats.Stats, 1,
		"Alice should be added to epoch stats because she has referred volume")

	aliceStats := postUpgradeEpochStats.Stats[0]
func TestStateUpgrade(t *testing.T) {
	testnet, err := containertest.NewTestnetWithPreupgradeGenesis()
	require.NoError(t, err, "failed to create testnet - is docker daemon running?")
	err = testnet.Start()
	require.NoError(t, err)
	defer testnet.MustCleanUp()
	node := testnet.Nodes["alice"]
	nodeAddress := constants.AliceAccAddress.String()

	preUpgradeSetups(node, t)
	preUpgradeChecks(node, t)

	err = containertest.UpgradeTestnet(nodeAddress, t, node, v_9_5.UpgradeName)
	require.NoError(t, err)

	postUpgradeChecks(node, t)
}

func preUpgradeSetups(node *containertest.Node, t *testing.T) {
	// Set up user stats with referred volume before upgrade
	// This simulates users having 30d referred volume in global stats
}

func preUpgradeChecks(node *containertest.Node, t *testing.T) {
}

func postUpgradeChecks(node *containertest.Node, t *testing.T) {
}
