package v_9_5

import (
	"context"
	"fmt"
	"sort"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	epochskeeper "github.com/dydxprotocol/v4-chain/protocol/x/epochs/keeper"
	statskeeper "github.com/dydxprotocol/v4-chain/protocol/x/stats/keeper"
	statstypes "github.com/dydxprotocol/v4-chain/protocol/x/stats/types"
)

// migrate30dReferredVolumeToEpochStats migrates all users' 30d referred volume
// to epoch stats in the current epoch.
func migrate30dReferredVolumeToEpochStats(
	ctx sdk.Context,
	statsKeeper statskeeper.Keeper,
	epochsKeeper epochskeeper.Keeper,
) {
	// Get the current stats epoch
	statsEpochInfo := epochsKeeper.MustGetStatsEpochInfo(ctx)
	currentEpoch := statsEpochInfo.CurrentEpoch

	ctx.Logger().Info(fmt.Sprintf(
		"Migrating 30d referred volume to epoch stats for epoch %d",
		currentEpoch,
	))

	// Get or create epoch stats for current epoch
	epochStats := statsKeeper.GetEpochStatsOrNil(ctx, currentEpoch)
	if epochStats == nil {
		epochStats = &statstypes.EpochStats{
			Stats: []*statstypes.EpochStats_UserWithStats{},
		}
	}

	// Create a map for existing epoch stats for quick lookup
	userStatsMap := make(map[string]*statstypes.EpochStats_UserWithStats)
	for _, userWithStats := range epochStats.Stats {
		userStatsMap[userWithStats.User] = userWithStats
	}

	// Get all addresses with referred volume from the global UserStats
	allAddressesWithReferredVolume := statsKeeper.GetAllAddressesWithReferredVolume(ctx)

	migratedCount := 0

	for _, address := range allAddressesWithReferredVolume {
		// Get the global user stats which contains the 30d referred volume
		globalUserStats := statsKeeper.GetUserStats(ctx, address)
		if globalUserStats == nil {
			continue
		}

		referredVolume := globalUserStats.Affiliate_30DReferredVolumeQuoteQuantums

		// Get or create user stats for this epoch
		epochUserStats, exists := userStatsMap[address]
		if !exists {
			// User not in epoch stats yet, create new entry
			epochUserStats = &statstypes.EpochStats_UserWithStats{
				User: address,
				Stats: &statstypes.UserStats{
					Affiliate_30DReferredVolumeQuoteQuantums: referredVolume,
				},
			}
			userStatsMap[address] = epochUserStats
		} else {
			// User already in epoch stats, add the referred volume
			epochUserStats.Stats.Affiliate_30DReferredVolumeQuoteQuantums += referredVolume
		}

		migratedCount++

		ctx.Logger().Info(fmt.Sprintf(
			"Migrated referred volume for address %s: %d",
			address,
			referredVolume,
		))
	}

	// Convert map back to slice and sort for deterministic ordering
	epochStats.Stats = make([]*statstypes.EpochStats_UserWithStats, 0, len(userStatsMap))
	for _, userStats := range userStatsMap {
		epochStats.Stats = append(epochStats.Stats, userStats)
	}

	// Sort by address for deterministic ordering
	sort.Slice(epochStats.Stats, func(i, j int) bool {
		return epochStats.Stats[i].User < epochStats.Stats[j].User
	})

	// Save the updated epoch stats
	statsKeeper.SetEpochStats(ctx, currentEpoch, epochStats)

	ctx.Logger().Info(fmt.Sprintf(
		"Successfully migrated 30d referred volume for %d addresses to epoch %d",
		migratedCount,
		currentEpoch,
	))
}

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	statsKeeper statskeeper.Keeper,
	epochsKeeper epochskeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := lib.UnwrapSDKContext(ctx, "app/upgrades")
		sdkCtx.Logger().Info(fmt.Sprintf("Running %s Upgrade...", UpgradeName))

		// Migrate 30d referred volume to epoch stats
		migrate30dReferredVolumeToEpochStats(sdkCtx, statsKeeper, epochsKeeper)

		sdkCtx.Logger().Info(fmt.Sprintf("Successfully completed %s Upgrade", UpgradeName))

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
