package keeper

import (
	"fmt"

	storetypes "cosmossdk.io/store/types"

	errorsmod "cosmossdk.io/errors"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/log"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/epochs/types"
	gometrics "github.com/hashicorp/go-metrics"
)

func (k Keeper) getEpochInfoStore(
	ctx sdk.Context,
) prefix.Store {
	return prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.EpochInfoKeyPrefix))
}

func (k Keeper) setEpochInfo(ctx sdk.Context, epochInfo types.EpochInfo) {
	store := k.getEpochInfoStore(ctx)
	b := k.cdc.MustMarshal(&epochInfo)
	store.Set([]byte(epochInfo.Name), b)
}

// MaybeStartNextEpoch initializes and/or ticks the next epoch.
// First, initializes `EpochInfo` if all below conditions are met:
// - `EpochInfo.IsInitialized` is false
// - `BlockHeight` >= 2 (for accurate current block time)
// - `BlockTime` >= `EpochInfo.NextTick`
// If `EpochInfo.IsInitialized` is true, checks if current `BlockTime` has reached `NextTick`
// of the epoch, and if so starts a new epoch by updating `NextTick`, `CurrentEpoch` and
// `NextEpochStartBlock`.
func (k Keeper) MaybeStartNextEpoch(ctx sdk.Context, id types.EpochInfoName) (nextEpochStarted bool, err error) {
	epoch, found := k.GetEpochInfo(ctx, id)
	if !found {
		return false, errorsmod.Wrapf(types.ErrEpochInfoNotFound, "EpochInfo Id not found (%s)", id)
	}

	blockTime := uint32(ctx.BlockTime().Unix())

	if !epoch.IsInitialized {
		// Require `blockHeight >= 2` since both genesis and the first block use
		// genesis file time as blockTime, so we need blockTime from the second block
		// (which represents current time) to initialize NextTick.
		shouldInitialize := ctx.BlockHeight() >= 2 && blockTime >= epoch.NextTick

		if !shouldInitialize {
			// `EpochInfo` not ready for initialization, don't tick
			return false, nil
		}

		// Initialize `EpochInfo`.
		epoch.IsInitialized = true

		// Set `NextTick` to the smallest value `x` greater than
		// the current block time such that `(x - NextTick) % duration = 0`.
		if epoch.FastForwardNextTick {
			// `durationMultiplier` is equal to the number of `duration`s between
			// genesis `NextTick` and the nearest time in future after current block time,
			// rounded up.
			durationMultiplier := (blockTime-epoch.NextTick)/epoch.Duration + 1
			epoch.NextTick = epoch.NextTick + epoch.Duration*durationMultiplier
		}
		k.setEpochInfo(ctx, epoch)
	}

	if blockTime < epoch.NextTick {
		// NextTick not reached yet.
		return false, nil
	}

	// Starts next epoch.
	currentTick := epoch.NextTick

	epoch.NextTick = epoch.NextTick + epoch.Duration
	epoch.CurrentEpoch++
	epoch.CurrentEpochStartBlock = lib.MustConvertIntegerToUint32(ctx.BlockHeight())
	k.setEpochInfo(ctx, epoch)

	log.InfoLog(
		ctx,
		fmt.Sprintf(
			"Starting new epoch for [%s], current block time = %d, new epoch info = %+v",
			epoch.Name,
			ctx.BlockTime().Unix(),
			epoch,
		),
	)

	ctx.EventManager().EmitEvent(
		types.NewEpochEvent(ctx, epoch, currentTick),
	)

	// Stat latest epoch number.
	telemetry.SetGaugeWithLabels(
		[]string{types.ModuleName, types.AttributeKeyEpochNumber},
		float32(epoch.CurrentEpoch),
		[]gometrics.Label{
			metrics.GetLabelForStringValue(types.AttributeKeyEpochInfoName, epoch.Name),
		},
	)

	return true, nil
}

// CreateEpochInfo creates a new EpochInfo.
// Return an error if the epoch fails validation, if the epoch Id already exists.
func (k Keeper) CreateEpochInfo(ctx sdk.Context, epochInfo types.EpochInfo) error {
	// Perform stateless validation on the provided `EpochInfo`.
	err := epochInfo.Validate()
	if err != nil {
		return err
	}

	// Check if identifier already exists
	if _, found := k.GetEpochInfo(ctx, epochInfo.GetEpochInfoName()); found {
		return errorsmod.Wrapf(types.ErrEpochInfoAlreadyExists, "epochInfo.Name already exists (%s)", epochInfo.Name)
	}

	k.setEpochInfo(ctx, epochInfo)
	log.InfoLog(
		ctx,
		fmt.Sprintf(
			"Created new epoch info (current block time = %v): %+v",
			ctx.BlockTime().Unix(),
			epochInfo,
		),
	)
	return nil
}

// GetEpochInfo returns an epochInfo from its id
func (k Keeper) GetEpochInfo(
	ctx sdk.Context,
	id types.EpochInfoName,
) (val types.EpochInfo, found bool) {
	store := k.getEpochInfoStore(ctx)

	b := store.Get([]byte(id))

	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// GetAllEpochInfo returns all epochInfos
func (k Keeper) GetAllEpochInfo(ctx sdk.Context) (list []types.EpochInfo) {
	store := k.getEpochInfoStore(ctx)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.EpochInfo
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// NumBlocksSinceEpochStart returns the number of blocks since the epoch started.
func (k Keeper) NumBlocksSinceEpochStart(
	ctx sdk.Context,
	id types.EpochInfoName,
) (
	uint32,
	error,
) {
	epoch, found := k.GetEpochInfo(ctx, id)
	if !found {
		return 0, errorsmod.Wrapf(types.ErrEpochInfoNotFound, "EpochInfo Id not found (%s)", id)
	}

	return lib.MustConvertIntegerToUint32(ctx.BlockHeight() - int64(epoch.CurrentEpochStartBlock)), nil
}

func (k Keeper) MustGetFundingTickEpochInfo(
	ctx sdk.Context,
) types.EpochInfo {
	return k.mustGetEpochInfo(ctx, types.FundingTickEpochInfoName)
}

func (k Keeper) MustGetFundingSampleEpochInfo(
	ctx sdk.Context,
) types.EpochInfo {
	return k.mustGetEpochInfo(ctx, types.FundingSampleEpochInfoName)
}

func (k Keeper) MustGetStatsEpochInfo(
	ctx sdk.Context,
) types.EpochInfo {
	return k.mustGetEpochInfo(ctx, types.StatsEpochInfoName)
}

func (k Keeper) mustGetEpochInfo(
	ctx sdk.Context,
	epochInfoName types.EpochInfoName,
) types.EpochInfo {
	epochInfo, found := k.GetEpochInfo(
		ctx,
		epochInfoName,
	)
	if !found {
		panic(errorsmod.Wrapf(
			types.ErrEpochInfoNotFound,
			"name: %s",
			epochInfoName,
		))
	}
	return epochInfo
}
