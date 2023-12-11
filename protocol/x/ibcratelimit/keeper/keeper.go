package keeper

import (
	"fmt"

	sdklog "cosmossdk.io/log"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/ibcratelimit/types"
)

type (
	Keeper struct {
		cdc      codec.BinaryCodec
		storeKey storetypes.StoreKey

		// TODO(CORE-824): Implement `x/ibcratelimit` keeper

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

// // Returns whether the withdrawal IBC transfer can proceed without getting rate limited.
// // It compares input amount vs. minimum capacity among all LimitParams for the denom.
// func (k Keeper) CanWithdraw(
// 	ctx sdk.Context,
// 	denom string,
// 	amount *sdkmath.Int,
// ) bool {
// 	denomCapacity := k.GetDenomCapacity(ctx, denom)
// 	for _, capacity := range denomCapacity.CapacityList {
// 		// Remaining capacity is the minimum of all LimitParams.
// 		if capacity.BigInt().Cmp(amount.BigInt()) <= 0 {
// 			return false
// 		}
// 	}
// 	return true
// }

// ProcessDeposit processes a inbound IBC transfer,
// by updating the relevant capacity lists for the denom.
func (k Keeper) ProcessDeposit(
	ctx sdk.Context,
	denom string,
	amount *sdk.Int,
) error {
	denomCapacity := k.GetDenomCapacity(ctx, denom)
	for _, capacity := range denomCapacity.CapacityList {

	}
	return true
}

// SetLimitParams sets `LimitParams` for the given denom.
func (k Keeper) SetLimitParams(
	ctx sdk.Context,
	limitParams types.LimitParams,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.DenomCapacityKeyPrefix))

	key := []byte(limitParams.Denom)
	// If there's no capacity entry to set, delete the key.
	if len(limitParams.Limiters) == 0 {
		if store.Has(key) {
			store.Delete(key)
		}
	} else {
		b := k.cdc.MustMarshal(&limitParams)
		store.Set(key, b)
	}
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
