package keeper

import (
	"fmt"
	"time"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
)

type (
	Keeper struct {
		cdc         codec.BinaryCodec
		storeKey    storetypes.StoreKey
		authorities map[string]struct{}
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	authorities []string,
) *Keeper {
	return &Keeper{
		cdc:         cdc,
		storeKey:    storeKey,
		authorities: lib.UniqueSliceToSet(authorities),
	}
}

func (k Keeper) HasAuthority(authority string) bool {
	_, ok := k.authorities[authority]
	return ok
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With(log.ModuleKey, fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) InitializeForGenesis(ctx sdk.Context) {}

func (k Keeper) GetAllDowntimeInfo(ctx sdk.Context) *types.AllDowntimeInfo {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get([]byte(types.AllDowntimeInfoKey))

	if bytes == nil {
		return &types.AllDowntimeInfo{}
	}

	var info types.AllDowntimeInfo
	k.cdc.MustUnmarshal(bytes, &info)
	return &info
}

// SetAllDowntimeInfo sets AllDowntimeInfo in state. Durations in AllDowntimeInfo must match
// the durations in DowntimeParams. If not, behavior of this module is undefined.
func (k Keeper) SetAllDowntimeInfo(ctx sdk.Context, info *types.AllDowntimeInfo) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(info)
	store.Set([]byte(types.AllDowntimeInfoKey), b)
}

func (k Keeper) GetPreviousBlockInfo(ctx sdk.Context) types.BlockInfo {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get([]byte(types.PreviousBlockInfoKey))

	if bytes == nil {
		return types.BlockInfo{}
	}

	var info types.BlockInfo
	k.cdc.MustUnmarshal(bytes, &info)
	return info
}

// GetTimeSinceLastBlock returns the time delta between the current block time and the last block time.
func (k Keeper) GetTimeSinceLastBlock(ctx sdk.Context) time.Duration {
	prevBlockInfo := k.GetPreviousBlockInfo(ctx)
	return ctx.BlockTime().Sub(prevBlockInfo.Timestamp)
}

func (k Keeper) SetPreviousBlockInfo(ctx sdk.Context, info *types.BlockInfo) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(info)
	store.Set([]byte(types.PreviousBlockInfoKey), b)
}

// UpdateAllDowntimeInfo updates AllDowntimeInfo by considering the downtime between the current block and
// the previous block and updating the DowntimeInfo for each observed duration.
func (k Keeper) UpdateAllDowntimeInfo(ctx sdk.Context) {
	previousBlockInfo := k.GetPreviousBlockInfo(ctx)
	delta := ctx.BlockTime().Sub(previousBlockInfo.Timestamp)
	// Report block time in milliseconds.
	telemetry.SetGauge(
		float32(delta.Milliseconds()),
		types.ModuleName,
		metrics.BlockTimeMs,
	)

	metrics.AddSampleWithLabels(
		metrics.BlockTimeDistribution,
		float32(delta.Milliseconds()),
		metrics.GetLabelForStringValue(
			metrics.Proposer,
			sdk.ConsAddress(ctx.BlockHeader().ProposerAddress).String(),
		),
	)

	allInfo := k.GetAllDowntimeInfo(ctx)

	for _, info := range allInfo.Infos {
		if delta >= info.Duration {
			info.BlockInfo = types.BlockInfo{
				Height:    uint32(ctx.BlockHeight()),
				Timestamp: ctx.BlockTime(),
			}
		} else {
			break
		}
	}
	k.SetAllDowntimeInfo(ctx, allInfo)
}

// GetDowntimeInfoFor gets the DowntimeInfo for a specific duration. If the exact duration is not observed, it
// returns the DowntimeInfo for the largest duration that is smaller than the input duration. If the input
// duration is smaller than all observed durations, then return a DowntimeInfo with duration 0 and the current
// block's info.
func (k Keeper) GetDowntimeInfoFor(ctx sdk.Context, duration time.Duration) types.AllDowntimeInfo_DowntimeInfo {
	allInfo := k.GetAllDowntimeInfo(ctx)
	ret := types.AllDowntimeInfo_DowntimeInfo{
		Duration: 0,
		BlockInfo: types.BlockInfo{
			Height:    uint32(ctx.BlockHeight()),
			Timestamp: ctx.BlockTime(),
		},
	}
	for _, info := range allInfo.Infos {
		if duration >= info.Duration {
			ret = *info
		} else {
			break
		}
	}
	return ret
}
