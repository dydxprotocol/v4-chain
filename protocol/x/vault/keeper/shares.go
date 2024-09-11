package keeper

import (
	"fmt"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	delaymsgtypes "github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
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

// GetOwnerShareUnlocks gets share unlocks for an owner.
func (k Keeper) GetOwnerShareUnlocks(
	ctx sdk.Context,
	owner string,
) (val types.OwnerShareUnlocks, exists bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.OwnerShareUnlocksKeyPrefix))

	b := store.Get([]byte(owner))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// SetOwnerShareUnlocks sets share unlocks for an owner.
func (k Keeper) SetOwnerShareUnlocks(
	ctx sdk.Context,
	owner string,
	ownerShareUnlocks types.OwnerShareUnlocks,
) error {
	if err := ownerShareUnlocks.Validate(); err != nil {
		return err
	}

	b := k.cdc.MustMarshal(&ownerShareUnlocks)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.OwnerShareUnlocksKeyPrefix))
	store.Set([]byte(owner), b)

	return nil
}

// GetAllOwnerShareUnlocks gets all `OwnerShareUnlocks`.
func (k Keeper) GetAllOwnerShareUnlocks(ctx sdk.Context) []types.OwnerShareUnlocks {
	allOwnerShareUnlocks := []types.OwnerShareUnlocks{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.OwnerShareUnlocksKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var ownerShareUnlocks types.OwnerShareUnlocks
		k.cdc.MustUnmarshal(iterator.Value(), &ownerShareUnlocks)
		allOwnerShareUnlocks = append(allOwnerShareUnlocks, ownerShareUnlocks)
	}
	return allOwnerShareUnlocks
}

// LockShares locks `sharesToLock` for `ownerAddress` until height `tilBlock`.
// Note: cannot lock more than the total number of shares that owner has.
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

	ownerShareUnlocks, exists := k.GetOwnerShareUnlocks(ctx, ownerAddress)
	if !exists {
		// Initialize share unlocks of this owner.
		ownerShareUnlocks = types.OwnerShareUnlocks{
			OwnerAddress: ownerAddress,
			ShareUnlocks: []types.ShareUnlock{
				{
					Shares:            sharesToLock,
					UnlockBlockHeight: tilBlock,
				},
			},
		}
	} else {
		// Add new instance of share unlock.
		ownerShareUnlocks.ShareUnlocks = append(ownerShareUnlocks.ShareUnlocks, types.ShareUnlock{
			Shares:            sharesToLock,
			UnlockBlockHeight: tilBlock,
		})
	}

	// Total locked shares cannot exceed total owner shares.
	ownerShares, exists := k.GetOwnerShares(ctx, ownerAddress)
	if !exists {
		return errorsmod.Wrapf(types.ErrOwnerNotFound, "owner: %s", ownerAddress)
	}
	totalLockedShares := ownerShareUnlocks.GetTotalLockedShares()
	if ownerShareUnlocks.GetTotalLockedShares().Cmp(ownerShares.NumShares.BigInt()) == 1 {
		return errorsmod.Wrapf(
			types.ErrLockedSharesExceedsOwnerShares,
			"owner: %s, ownerShares: %s, sharesToLock: %s, total shares that will be locked: %s",
			ownerAddress,
			ownerShares,
			sharesToLock,
			totalLockedShares,
		)
	}

	// Schedule an unlock at block height `tilBlock`.
	_, err := k.delayMsgKeeper.DelayMessageByBlocks(
		ctx,
		&types.MsgUnlockShares{
			Authority:    delaymsgtypes.ModuleAddress.String(),
			OwnerAddress: ownerAddress,
		},
		tilBlock-uint32(ctx.BlockHeight()),
	)
	if err != nil {
		return err
	}

	err = k.SetOwnerShareUnlocks(ctx, ownerAddress, ownerShareUnlocks)
	if err != nil {
		return err
	}

	return nil
}

// UnlockShares unlocks all shares of an owner that are due to unlock at or before current block height.
func (k Keeper) UnlockShares(
	ctx sdk.Context,
	ownerAddress string,
) (unlockedShares types.NumShares, err error) {
	ownerShareUnlocks, exists := k.GetOwnerShareUnlocks(ctx, ownerAddress)
	if !exists {
		return unlockedShares, types.ErrOwnerNotFound
	}

	// Process all unlocks that are due.
	totalSharesUnlocked := big.NewInt(0)
	currentBlockHeight := uint32(ctx.BlockHeight())
	remainingUnlocks := []types.ShareUnlock{}
	for _, unlock := range ownerShareUnlocks.ShareUnlocks {
		if unlock.UnlockBlockHeight <= currentBlockHeight {
			totalSharesUnlocked.Add(totalSharesUnlocked, unlock.Shares.NumShares.BigInt())
		} else {
			remainingUnlocks = append(remainingUnlocks, unlock)
		}
	}

	// Remove from store if no more unlocks. Update otherwise.
	if len(remainingUnlocks) == 0 {
		ownerShareUnlocksStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.OwnerShareUnlocksKeyPrefix))
		ownerShareUnlocksStore.Delete([]byte(ownerAddress))
	} else {
		err = k.SetOwnerShareUnlocks(
			ctx,
			ownerAddress,
			types.OwnerShareUnlocks{
				OwnerAddress: ownerAddress,
				ShareUnlocks: remainingUnlocks,
			},
		)
		if err != nil {
			return unlockedShares, err
		}
	}

	return types.BigIntToNumShares(totalSharesUnlocked), nil
}
