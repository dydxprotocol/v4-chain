package keeper

import (
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/cometbft/cometbft/libs/log"

	sdklog "cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/stats/types"
)

type (
	Keeper struct {
		cdc               codec.BinaryCodec
		epochsKeeper      types.EpochsKeeper
		storeKey          storetypes.StoreKey
		transientStoreKey storetypes.StoreKey
		authorities       map[string]struct{}
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	epochsKeeper types.EpochsKeeper,
	storeKey storetypes.StoreKey,
	transientStoreKey storetypes.StoreKey,
	authorities []string,
) *Keeper {
	return &Keeper{
		cdc:               cdc,
		epochsKeeper:      epochsKeeper,
		storeKey:          storeKey,
		transientStoreKey: transientStoreKey,
		authorities:       lib.UniqueSliceToSet(authorities),
	}
}

func (k Keeper) HasAuthority(authority string) bool {
	_, ok := k.authorities[authority]
	return ok
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With(sdklog.ModuleKey, fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) InitializeForGenesis(ctx sdk.Context) {}

func (k Keeper) GetBlockStats(ctx sdk.Context) *types.BlockStats {
	store := ctx.TransientStore(k.transientStoreKey)
	bytes := store.Get([]byte(types.BlockStatsKey))

	if bytes == nil {
		return &types.BlockStats{}
	}

	var blockStats types.BlockStats
	k.cdc.MustUnmarshal(bytes, &blockStats)
	return &blockStats
}

func (k Keeper) SetBlockStats(ctx sdk.Context, blockStats *types.BlockStats) {
	store := ctx.TransientStore(k.transientStoreKey)
	b := k.cdc.MustMarshal(blockStats)
	store.Set([]byte(types.BlockStatsKey), b)
}

// Record a match in BlockStats, which is stored in the transient store
func (k Keeper) RecordFill(ctx sdk.Context, takerAddress string, makerAddress string, notional *big.Int) {
	blockStats := k.GetBlockStats(ctx)
	blockStats.Fills = append(
		blockStats.Fills,
		&types.BlockStats_Fill{
			Taker:    takerAddress,
			Maker:    makerAddress,
			Notional: notional.Uint64(),
		},
	)
	k.SetBlockStats(ctx, blockStats)
}

func (k Keeper) GetStatsMetadata(ctx sdk.Context) *types.StatsMetadata {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get([]byte(types.StatsMetadataKey))

	if bytes == nil {
		return &types.StatsMetadata{}
	}

	var metadata types.StatsMetadata
	k.cdc.MustUnmarshal(bytes, &metadata)
	return &metadata
}

func (k Keeper) SetStatsMetadata(ctx sdk.Context, metadata *types.StatsMetadata) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(metadata)
	store.Set([]byte(types.StatsMetadataKey), b)
}

// GetEpochStatsOrNil returns the EpochStats for an epoch. This function returns nil
// if epoch stats aren't found.
func (k Keeper) GetEpochStatsOrNil(ctx sdk.Context, epoch uint32) *types.EpochStats {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.EpochStatsKeyPrefix))
	bytes := store.Get(lib.Uint32ToKey(epoch))

	if bytes == nil {
		return nil
	}

	var epochStats types.EpochStats
	k.cdc.MustUnmarshal(bytes, &epochStats)
	return &epochStats
}

func (k Keeper) SetEpochStats(ctx sdk.Context, epoch uint32, epochStats *types.EpochStats) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.EpochStatsKeyPrefix))
	b := k.cdc.MustMarshal(epochStats)
	store.Set(lib.Uint32ToKey(epoch), b)
}

func (k Keeper) deleteEpochStats(ctx sdk.Context, epoch uint32) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.EpochStatsKeyPrefix))
	store.Delete(lib.Uint32ToKey(epoch))
}

func (k Keeper) GetUserStats(ctx sdk.Context, address string) *types.UserStats {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.UserStatsKeyPrefix))
	bytes := store.Get([]byte(address))

	if bytes == nil {
		return &types.UserStats{}
	}

	var userStats types.UserStats
	k.cdc.MustUnmarshal(bytes, &userStats)
	return &userStats
}

func (k Keeper) SetUserStats(ctx sdk.Context, address string, userStats *types.UserStats) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.UserStatsKeyPrefix))
	b := k.cdc.MustMarshal(userStats)
	store.Set([]byte(address), b)
}

func (k Keeper) GetGlobalStats(ctx sdk.Context) *types.GlobalStats {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get([]byte(types.GlobalStatsKey))

	if bytes == nil {
		return &types.GlobalStats{}
	}

	var globalStats types.GlobalStats
	k.cdc.MustUnmarshal(bytes, &globalStats)
	return &globalStats
}

func (k Keeper) SetGlobalStats(ctx sdk.Context, globalStats *types.GlobalStats) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(globalStats)
	store.Set([]byte(types.GlobalStatsKey), b)
}

// ProcessBlockStats persists the info from this block's BlockStats this epoch's stats.
// It also appropriately increments the overall stats globally and for each user
func (k Keeper) ProcessBlockStats(ctx sdk.Context) {
	epochInfo := k.epochsKeeper.MustGetStatsEpochInfo(ctx)
	blockStats := k.GetBlockStats(ctx)

	if len(blockStats.Fills) == 0 {
		return
	}

	epochStats := k.GetEpochStatsOrNil(ctx, epochInfo.CurrentEpoch)
	if epochStats == nil {
		epochStats = &types.EpochStats{
			Stats: []*types.EpochStats_UserWithStats{},
		}
	}
	// We expect entries in the list to already be unique
	userStatsMap := map[string]*types.EpochStats_UserWithStats{}
	for _, userWithStats := range epochStats.Stats {
		userStatsMap[userWithStats.User] = userWithStats
	}

	// NB: These unsigned ints can technically overflow and wrap around, but the trading volume
	// required to do so is unrealistic.
	for _, fill := range blockStats.Fills {
		userStats := k.GetUserStats(ctx, fill.Taker)
		userStats.TakerNotional += fill.Notional
		k.SetUserStats(ctx, fill.Taker, userStats)

		userStats = k.GetUserStats(ctx, fill.Maker)
		userStats.MakerNotional += fill.Notional
		k.SetUserStats(ctx, fill.Maker, userStats)

		if _, ok := userStatsMap[fill.Taker]; !ok {
			userStatsMap[fill.Taker] = &types.EpochStats_UserWithStats{
				User:  fill.Taker,
				Stats: &types.UserStats{},
			}
		}
		if _, ok := userStatsMap[fill.Maker]; !ok {
			userStatsMap[fill.Maker] = &types.EpochStats_UserWithStats{
				User:  fill.Maker,
				Stats: &types.UserStats{},
			}
		}
		userStatsMap[fill.Taker].Stats.TakerNotional += fill.Notional
		userStatsMap[fill.Maker].Stats.MakerNotional += fill.Notional

		globalStats := k.GetGlobalStats(ctx)
		globalStats.NotionalTraded += fill.Notional
		k.SetGlobalStats(ctx, globalStats)
	}

	keys := make([]string, 0, len(userStatsMap))
	for k := range userStatsMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	epochStats.Stats = make([]*types.EpochStats_UserWithStats, 0, len(userStatsMap))
	for _, k := range keys {
		epochStats.Stats = append(epochStats.Stats, userStatsMap[k])
	}
	epochStats.EpochEndTime = time.Unix(int64(epochInfo.NextTick), 0).UTC()
	k.SetEpochStats(ctx, epochInfo.CurrentEpoch, epochStats)
}

// ExpireOldStats expiration of stats when they fall out of the window.
// TrailingEpoch is next epoch that can potentially fall out of the window.
// Attempt to expire the next epoch. TrailingEpoch will be advanced at most once.
func (k Keeper) ExpireOldStats(ctx sdk.Context) {
	currentEpoch := k.epochsKeeper.MustGetStatsEpochInfo(ctx).CurrentEpoch
	metadata := k.GetStatsMetadata(ctx)

	// Current epoch can't be expired.
	if metadata.TrailingEpoch == currentEpoch {
		return
	}

	epochStats := k.GetEpochStatsOrNil(ctx, metadata.TrailingEpoch)
	// Empty epoch falls out of window
	if epochStats == nil {
		metadata.TrailingEpoch += 1
		k.SetStatsMetadata(ctx, metadata)
		return
	}

	// Epoch not ready to fall out of window
	if !epochStats.EpochEndTime.Before(ctx.BlockTime().Add(-k.GetWindowDuration(ctx))) {
		return
	}

	globalStats := k.GetGlobalStats(ctx)
	for _, removedStats := range epochStats.Stats {
		stats := k.GetUserStats(ctx, removedStats.User)
		stats.TakerNotional -= removedStats.Stats.TakerNotional
		stats.MakerNotional -= removedStats.Stats.MakerNotional
		k.SetUserStats(ctx, removedStats.User, stats)

		// Just remove TakerNotional to avoid double counting
		globalStats.NotionalTraded -= removedStats.Stats.TakerNotional
	}
	k.SetGlobalStats(ctx, globalStats)
	k.deleteEpochStats(ctx, metadata.TrailingEpoch)
	metadata.TrailingEpoch += 1
	k.SetStatsMetadata(ctx, metadata)
}
