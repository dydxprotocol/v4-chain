package keeper

import (
	"fmt"
	"math/big"
	"time"

	errorsmod "cosmossdk.io/errors"
	cosmoslog "cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/log"
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

		// the addresses capable of executing MsgSetLimitParams message.
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
) (err error) {
	if err := limitParams.Validate(); err != nil {
		return err
	}

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

	return nil
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

	// Emit telemetry for the new capacity list.
	for i, capacity := range denomCapacity.CapacityList {
		telemetry.SetGaugeWithLabels(
			[]string{types.ModuleName, metrics.Capacity},
			metrics.GetMetricValueFromBigInt(capacity.BigInt()),
			[]gometrics.Label{
				metrics.GetLabelForStringValue(metrics.RateLimitDenom, denomCapacity.Denom),
				metrics.GetLabelForIntValue(metrics.LimiterIndex, i),
			},
		)
	}
}

// UpdateAllCapacitiesEndBlocker is called during the EndBlocker to update the capacity for all limit params.
func (k Keeper) UpdateAllCapacitiesEndBlocker(
	ctx sdk.Context,
) {
	timeSinceLastBlock := k.blockTimeKeeper.GetTimeSinceLastBlock(ctx)

	if timeSinceLastBlock < 0 {
		// This violates an invariant (current block time > prev block time).
		// Since this is in the `EndBlocker`, we log an error instead of panicking.
		log.ErrorLog(
			ctx,
			fmt.Sprintf(
				"timeSinceLastBlock (%v) <= 0; skipping UpdateAllCapacitiesEndBlocker",
				timeSinceLastBlock,
			),
		)
		return
	}

	// Iterate through all the limit params in state.
	limitParams := k.GetAllLimitParams(ctx)
	for _, limitParams := range limitParams {
		k.updateCapacityForLimitParams(ctx, limitParams, timeSinceLastBlock)
	}
}

// updateCapacityForLimitParams calculates current baseline for a denom and recovers some amount of capacity
// towards baseline.
// Assumes that the `LimitParams` exist in state.
func (k Keeper) updateCapacityForLimitParams(
	ctx sdk.Context,
	limitParams types.LimitParams,
	timeSinceLastBlock time.Duration,
) {
	tvl := k.bankKeeper.GetSupply(ctx, limitParams.Denom)

	capacityList := k.GetDenomCapacity(ctx, limitParams.Denom).CapacityList

	newCapacityList, err := ratelimitutil.CalculateNewCapacityList(
		tvl.Amount.BigInt(),
		limitParams,
		capacityList,
		timeSinceLastBlock,
	)

	if err != nil {
		log.ErrorLogWithError(
			ctx,
			fmt.Sprintf(
				"error calculating new capacity list for denom %v. Skipping update.",
				limitParams.Denom,
			),
			err,
		)
		return
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

// GetAllLimitParams returns `LimitParams` stored in state
func (k Keeper) GetAllLimitParams(ctx sdk.Context) (list []types.LimitParams) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.LimitParamsKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.LimitParams
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) HasAuthority(authority string) bool {
	_, ok := k.authorities[authority]
	return ok
}

func (k Keeper) Logger(ctx sdk.Context) cosmoslog.Logger {
	return ctx.Logger().With(cosmoslog.ModuleKey, fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) InitializeForGenesis(ctx sdk.Context) {}
