package keeper

import (
	"fmt"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// GetTotalShares gets total shares.
func (k Keeper) GetTotalShares(
	ctx sdk.Context,
) (
	totalShares types.NumShares,
) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get([]byte(types.TotalSharesKey))
	k.cdc.MustUnmarshal(b, &totalShares)
	return totalShares
}

// SetTotalShares sets total shares. Returns error if `totalShares` is negative.
func (k Keeper) SetTotalShares(
	ctx sdk.Context,
	totalShares types.NumShares,
) error {
	if totalShares.NumShares.Sign() < 0 {
		return types.ErrNegativeShares
	}

	b := k.cdc.MustMarshal(&totalShares)
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.TotalSharesKey), b)

	return nil
}

// GetOwnerShares gets owner shares for an owner.
func (k Keeper) GetOwnerShares(
	ctx sdk.Context,
	owner string,
) (val types.NumShares, exists bool) {
	store := k.getOwnerSharesStore(ctx)

	b := store.Get([]byte(owner))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// SetOwnerShares sets owner shares for an owner. Returns error if `ownerShares` is negative.
func (k Keeper) SetOwnerShares(
	ctx sdk.Context,
	owner string,
	ownerShares types.NumShares,
) error {
	if ownerShares.NumShares.Sign() < 0 {
		return types.ErrNegativeShares
	}

	b := k.cdc.MustMarshal(&ownerShares)
	store := k.getOwnerSharesStore(ctx)
	store.Set([]byte(owner), b)

	return nil
}

// getOwnerSharesStore returns the store for owner shares.
func (k Keeper) getOwnerSharesStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.OwnerSharesKeyPrefix))
}

// GetAllOwnerShares gets all owner shares.
func (k Keeper) GetAllOwnerShares(ctx sdk.Context) []types.OwnerShare {
	allOwnerShares := []types.OwnerShare{}
	ownerSharesStore := k.getOwnerSharesStore(ctx)
	ownerSharesIterator := storetypes.KVStorePrefixIterator(ownerSharesStore, []byte{})
	defer ownerSharesIterator.Close()
	for ; ownerSharesIterator.Valid(); ownerSharesIterator.Next() {
		owner := string(ownerSharesIterator.Key())
		var ownerShares types.NumShares
		k.cdc.MustUnmarshal(ownerSharesIterator.Value(), &ownerShares)
		allOwnerShares = append(allOwnerShares, types.OwnerShare{
			Owner:  owner,
			Shares: ownerShares,
		})
	}
	return allOwnerShares
}

// GetLockedShares gets locked shares for an owner.
func (k Keeper) GetLockedShares(
	ctx sdk.Context,
	owner string,
) (val types.LockedShares, exists bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.LockedSharesKeyPrefix))

	b := store.Get([]byte(owner))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// SetLockedShares sets locked shares for an owner.
func (k Keeper) SetLockedShares(
	ctx sdk.Context,
	owner string,
	lockedShares types.LockedShares,
) error {
	if err := lockedShares.Validate(); err != nil {
		return err
	}

	b := k.cdc.MustMarshal(&lockedShares)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.LockedSharesKeyPrefix))
	store.Set([]byte(owner), b)

	return nil
}

// GetAllLockedShares gets all locked shares.
func (k Keeper) GetAllLockedShares(ctx sdk.Context) []types.LockedShares {
	allLockedShares := []types.LockedShares{}
	lockedSharesStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.LockedSharesKeyPrefix))
	lockedSharesIterator := storetypes.KVStorePrefixIterator(lockedSharesStore, []byte{})
	defer lockedSharesIterator.Close()
	for ; lockedSharesIterator.Valid(); lockedSharesIterator.Next() {
		var lockedShares types.LockedShares
		k.cdc.MustUnmarshal(lockedSharesIterator.Value(), &lockedShares)
		allLockedShares = append(allLockedShares, lockedShares)
	}
	return allLockedShares
}

// LockShares locks `sharesToLock` for `ownerAddress` until height `tilBlock`.
func (k Keeper) LockShares(
	ctx sdk.Context,
	ownerAddress string,
	sharesToLock types.NumShares,
	tilBlock uint32,
) error {
	if ownerAddress == "" || sharesToLock.NumShares.Sign() <= 0 || tilBlock <= uint32(ctx.BlockHeight()) {
		return fmt.Errorf(
			`invalid parameters of shares locking:
				owner: %s, sharesToLock: %s, tilBlock: %d, current block height: %d`,
			ownerAddress,
			sharesToLock,
			tilBlock,
			ctx.BlockHeight(),
		)
	}

	lockedShares, exists := k.GetLockedShares(ctx, ownerAddress)
	if !exists {
		// Initialize locked shares.
		lockedShares = types.LockedShares{
			OwnerAddress:      ownerAddress,
			TotalLockedShares: sharesToLock,
			UnlockDetails: []types.UnlockDetail{
				{
					Shares:            sharesToLock,
					UnlockBlockHeight: tilBlock,
				},
			},
		}
	} else {
		// Increment total unlocked shares.
		totalLockedShares := lockedShares.TotalLockedShares.NumShares.BigInt()
		totalLockedShares.Add(totalLockedShares, sharesToLock.NumShares.BigInt())
		lockedShares.TotalLockedShares = types.BigIntToNumShares(totalLockedShares)

		// Add new unlock detail.
		lockedShares.UnlockDetails = append(lockedShares.UnlockDetails, types.UnlockDetail{
			Shares:            sharesToLock,
			UnlockBlockHeight: tilBlock,
		})
	}

	// TODO (TRA-565): delay a MsgUnlockShares.

	err := k.SetLockedShares(ctx, ownerAddress, lockedShares)
	if err != nil {
		return err
	}

	return nil
}
