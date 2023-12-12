package keeper

import (
	"fmt"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	sdklog "cosmossdk.io/log"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/ratelimit/types"
)

var (
	// TempTVLPlacerholder is a placeholder value for TVL.
	// TODO(CORE-836): Remove this after `GetBaseline` is fully implemented.
	TempTVLPlacerholder = big.NewInt(20_000_000_000_000)
)

type (
	Keeper struct {
		cdc      codec.BinaryCodec
		storeKey storetypes.StoreKey

		// TODO(CORE-824): Implement `x/ratelimit` keeper

		// the addresses capable of executing a MsgUpdateParams message.
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

// GetBaseline returns the current capacity baseline for the given limiter.
// `baseline` formula:
//
//	baseline = max(baseline_minimum, baseline_tvl_ppm * current_tvl)
func (k Keeper) GetBaseline(
	ctx sdk.Context,
	denom string,
	limiter types.Limiter,
) *big.Int {
	// Get the current TVL.
	// TODO(CORE-836): Query bank Keeper to get current supply of the token.
	currentTVL := TempTVLPlacerholder

	return lib.BigMax(
		limiter.BaselineMinimum.BigInt(),
		lib.BigIntMulPpm(
			currentTVL,
			limiter.BaselineTvlPpm,
		),
	)
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

	// Initialize the capacity list with the current baseline.
	newCapacityList := make([]dtypes.SerializableInt, len(limitParams.Limiters))
	for i, limiter := range limitParams.Limiters {
		newCapacityList[i] = dtypes.NewIntFromBigInt(
			k.GetBaseline(ctx, limitParams.Denom, limiter),
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
	return ctx.Logger().With(sdklog.ModuleKey, fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) InitializeForGenesis(ctx sdk.Context) {
}
