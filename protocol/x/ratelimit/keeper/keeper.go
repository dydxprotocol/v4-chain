package keeper

import (
	"fmt"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/ratelimit/types"
	ratelimitutil "github.com/dydxprotocol/v4-chain/protocol/x/ratelimit/util"
	gometrics "github.com/hashicorp/go-metrics"
)

type (
	Keeper struct {
		cdc             codec.BinaryCodec
		storeKey        storetypes.StoreKey
		bankKeeper      types.BankKeeper
		blockTimeKeeper types.BlockTimeKeeper

		// TODO(CORE-824): Implement `x/ratelimit` keeper

		// the addresses capable of executing a MsgUpdateParams message.
		authorities map[string]struct{}
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	bankKeeper types.BankKeeper,
	blockTimeKeeper types.BlockTimeKeeper,
	authorities []string,
) *Keeper {
	return &Keeper{
		cdc:             cdc,
		storeKey:        storeKey,
		bankKeeper:      bankKeeper,
		blockTimeKeeper: blockTimeKeeper,
		authorities:     lib.UniqueSliceToSet(authorities),
	}
}

// ProcessWithdrawal processes an outbound IBC transfer,
// by updating the capacity lists for the denom.
func (k Keeper) ProcessWithdrawal(
	ctx sdk.Context,
	denom string,
	amount *big.Int,
) error {
	denomCapacity := k.GetDenomCapacity(ctx, denom)

	newCapacityList := make([]dtypes.SerializableInt,
		0,
		len(denomCapacity.CapacityList),
	)

	for i, capacity := range denomCapacity.CapacityList {
		// Check that the withdrawal amount does not exceed each capacity.
		if capacity.BigInt().Cmp(amount) < 0 {
			return errorsmod.Wrapf(
				types.ErrWithdrawalExceedsCapacity,
				"denom = %v, capacity(index: %v) = %v, amount = %v",
				denom,
				i,
				capacity.BigInt(),
				amount,
			)
		}

		// Debit each capacity in the list by the amount of withdrawal.
		newCapacityList[i] = dtypes.NewIntFromBigInt(
			new(big.Int).Sub(
				capacity.BigInt(),
				amount,
			),
		)
	}

	k.SetDenomCapacity(ctx, types.DenomCapacity{
		Denom:        denom,
		CapacityList: newCapacityList,
	})

	return nil
}

// ProcessDeposit processes a inbound IBC transfer,
// by updating the capacity lists for the denom.
func (k Keeper) ProcessDeposit(
	ctx sdk.Context,
	denom string,
	amount *big.Int,
) {
	denomCapacity := k.GetDenomCapacity(ctx, denom)

	newCapacityList := make([]dtypes.SerializableInt,
		0,
		len(denomCapacity.CapacityList),
	)

	// Credit each capacity in the list by the amount of deposit.
	for i, capacity := range denomCapacity.CapacityList {
		newCapacityList[i] = dtypes.NewIntFromBigInt(
			new(big.Int).Add(
				capacity.BigInt(),
				amount,
			),
		)
	}

	k.SetDenomCapacity(ctx, types.DenomCapacity{
		Denom:        denom,
		CapacityList: newCapacityList,
	})
}

// SetLimitParams sets `LimitParams` for the given denom.
// Also overwrites the existing `DenomCapacity` object for the denom with a default `capacity_list` of the
// same length as the `limiters` list. Each `capacity` is initialized to the current baseline.
func (k Keeper) SetLimitParams(
	ctx sdk.Context,
	limitParams types.LimitParams,
) {
	limitParamsStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.LimitParamsKeyPrefix))
	denomCapacityStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.DenomCapacityKeyPrefix))

	denomKey := []byte(limitParams.Denom)
	// If the length of input limit params is zero, then remove both the
	// limit params and the denom capacity in state.
	if len(limitParams.Limiters) == 0 {
		if limitParamsStore.Has(denomKey) {
			limitParamsStore.Delete(denomKey)
		}
		if denomCapacityStore.Has(denomKey) {
			denomCapacityStore.Delete(denomKey)
		}
		return
	}

	currentTvl := k.bankKeeper.GetSupply(ctx, limitParams.Denom)
	// Initialize the capacity list with the current baseline.
	newCapacityList := make([]dtypes.SerializableInt, len(limitParams.Limiters))
	for i, limiter := range limitParams.Limiters {
		newCapacityList[i] = dtypes.NewIntFromBigInt(
			ratelimitutil.GetBaseline(currentTvl.Amount.BigInt(), limiter),
		)
	}
	// Set correspondong `DenomCapacity` in state.
	k.SetDenomCapacity(ctx, types.DenomCapacity{
		Denom:        limitParams.Denom,
		CapacityList: newCapacityList,
	})

	b := k.cdc.MustMarshal(&limitParams)
	limitParamsStore.Set(denomKey, b)
}

// GetLimitParams returns `LimitParams` for the given denom.
func (k Keeper) GetLimitParams(
	ctx sdk.Context,
	denom string,
) (val types.LimitParams) {
	// Check state for the LimitParams.
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.LimitParamsKeyPrefix))
	b := store.Get([]byte(denom))

	// If LimitParams does not exist in state, return a default value.
	if b == nil {
		return types.LimitParams{
			Denom: denom,
		}
	}

	// If LimitParams does exist in state, unmarshall and return the value.
	k.cdc.MustUnmarshal(b, &val)
	return val
}

// SetDenomCapacity sets `DenomCapacity` for the given denom.
func (k Keeper) SetDenomCapacity(
	ctx sdk.Context,
	denomCapacity types.DenomCapacity,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.DenomCapacityKeyPrefix))

	key := []byte(denomCapacity.Denom)
	// If there's no capacity entry to set, delete the key.
	if len(denomCapacity.CapacityList) == 0 {
		if store.Has(key) {
			store.Delete(key)
		}
	} else {
		b := k.cdc.MustMarshal(&denomCapacity)
		store.Set(key, b)
	}
}

// UpdateAllCapacitiesEndBlocker is called during the EndBlocker to update the capacity for all limit params.
func (k Keeper) UpdateAllCapacitiesEndBlocker(
	ctx sdk.Context,
) {
	// Iterate through all the limit params in state.
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.LimitParamsKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var limitParams types.LimitParams
		k.cdc.MustUnmarshal(iterator.Value(), &limitParams)
		k.updateCapacityForLimitParams(ctx, limitParams)
	}
}

// updateCapacityForLimitParams calculates current baseline for a denom and recovers some amount of capacity
// towards baseline.
// Assumes that the `LimitParams` exist in state.
// Detailed math for calculating the updated capacity:
//
//	`baseline = max(baseline_minimum, baseline_tvl_ppm * tvl)`
//	`capacity_diff = max(baseline, capacity-baseline) * (time_since_last_block / period)`
//
// This is basically saying that the capacity returns to the baseline over the course of the `period`.
// Usually in a linear way, but if the `capacity` is more than twice the `baseline`, then in an exponential way.
//
//	`capacity =`
//	    if `abs(capacity - baseline) < capacity_diff` then `capacity = baseline`
//	    else if `capacity < baseline` then `capacity += capacity_diff`
//	    else `capacity -= capacity_diff`
//
// On a high level, `capacity` trends towards `baseline` by `capacity_diff` but does not “cross” it.
func (k Keeper) updateCapacityForLimitParams(
	ctx sdk.Context,
	limitParams types.LimitParams,
) {
	tvl := k.bankKeeper.GetSupply(ctx, limitParams.Denom)

	capacityList := k.GetDenomCapacity(ctx, limitParams.Denom).CapacityList
	if len(capacityList) != len(limitParams.Limiters) {
		// This violates an invariant. Since this is in the `EndBlocker`, we log an error instead of panicking.
		k.Logger(ctx).Error(
			fmt.Sprintf(
				"denom (%s) capacity list length (%v) != limiters length (%v); skipping capacity update",
				limitParams.Denom,
				len(capacityList),
				len(limitParams.Limiters),
			),
		)
		return
	}

	// Convert timestamps to milliseconds for algebraic operations.
	blockTimeMilli := ctx.BlockTime().UnixMilli()
	prevBlockInfo := k.blockTimeKeeper.GetPreviousBlockInfo(ctx)
	prevBlockTimeMilli := prevBlockInfo.Timestamp.UnixMilli()
	timeSinceLastBlockMilli := new(big.Int).Sub(
		big.NewInt(blockTimeMilli),
		big.NewInt(prevBlockTimeMilli),
	)
	if timeSinceLastBlockMilli.Sign() <= 0 {
		// This violates an invariant (current block time > prev block time).
		// Since this is in the `EndBlocker`, we log an error instead of panicking.
		k.Logger(ctx).Error(
			fmt.Sprintf(
				"timeSinceLastBlockMilli (%v) <= 0; skipping capacity update. prevBlockTimeMilli = %v, blockTimeMilli = %v",
				timeSinceLastBlockMilli,
				prevBlockTimeMilli,
				blockTimeMilli,
			),
		)
		return
	}

	// Declare new capacity list to be populated.
	newCapacityList := make([]dtypes.SerializableInt, len(capacityList))

	for i, limiter := range limitParams.Limiters {
		// For each limiter, calculate the current baseline.
		baseline := ratelimitutil.GetBaseline(tvl.Amount.BigInt(), limiter)

		capacityMinusBaseline := new(big.Int).Sub(
			capacityList[i].BigInt(), // array access is safe because of the invariant check above
			baseline,
		)

		// Calculate left operand: `max(baseline, capacity-baseline)`. This equals `baseline` when `capacity <= 2 * baseline`
		operandL := new(big.Rat).SetInt(
			lib.BigMax(
				baseline,
				capacityMinusBaseline,
			),
		)

		// Calculate right operand: `time_since_last_block / period`
		periodMilli := new(big.Int).Mul(
			new(big.Int).SetUint64(uint64(limiter.PeriodSec)),
			big.NewInt(1000),
		)
		operandR := new(big.Rat).Quo(
			new(big.Rat).SetInt(timeSinceLastBlockMilli),
			new(big.Rat).SetInt(periodMilli),
		)

		// Calculate: `capacity_diff = max(baseline, capacity-baseline) * (time_since_last_block / period)`
		// Since both operands > 0, `capacity_diff` is positive or zero (due to rounding).
		capacityDiff := new(big.Rat).Mul(
			operandL,
			operandR,
		)

		bigRatcapacityMinusBaseline := new(big.Rat).SetInt(capacityMinusBaseline)

		if new(big.Rat).Abs(bigRatcapacityMinusBaseline).Cmp(capacityDiff) < 0 {
			// if `abs(capacity - baseline) < capacity_diff` then `capacity = baseline``
			newCapacityList[i] = dtypes.NewIntFromBigInt(baseline)
		} else if capacityList[i].BigInt().Cmp(baseline) < 0 {
			// else if `capacity < baseline` then `capacity += capacity_diff`
			newCapacityList[i] = dtypes.NewIntFromBigInt(
				new(big.Int).Add(
					capacityList[i].BigInt(),
					lib.BigRatRound(capacityDiff, false), // rounds down `capacity_diff`
				),
			)
		} else {
			// else `capacity -= capacity_diff`
			newCapacityList[i] = dtypes.NewIntFromBigInt(
				new(big.Int).Sub(
					capacityList[i].BigInt(),
					lib.BigRatRound(capacityDiff, false), // rounds down `capacity_diff`
				),
			)
		}

		// Emit telemetry for the new capacity.
		telemetry.SetGaugeWithLabels(
			[]string{types.ModuleName, metrics.Capacity},
			metrics.GetMetricValueFromBigInt(newCapacityList[i].BigInt()),
			[]gometrics.Label{
				metrics.GetLabelForStringValue(metrics.RateLimitDenom, limitParams.Denom),
				metrics.GetLabelForIntValue(metrics.LimiterIndex, i),
			},
		)
	}

	k.SetDenomCapacity(ctx, types.DenomCapacity{
		Denom:        limitParams.Denom,
		CapacityList: newCapacityList,
	})
}

// GetDenomCapacity returns `DenomCapacity` for the given denom.
func (k Keeper) GetDenomCapacity(
	ctx sdk.Context,
	denom string,
) (val types.DenomCapacity) {
	// Check state for the DenomCapacity.
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.DenomCapacityKeyPrefix))
	b := store.Get([]byte(denom))

	// If DenomCapacity does not exist in state, return a default value.
	if b == nil {
		return types.DenomCapacity{
			Denom: denom,
		}
	}

	// If DenomCapacity does exist in state, unmarshall and return the value.
	k.cdc.MustUnmarshal(b, &val)
	return val
}

func (k Keeper) HasAuthority(authority string) bool {
	_, ok := k.authorities[authority]
	return ok
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With(log.ModuleKey, fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) InitializeForGenesis(ctx sdk.Context) {
}
