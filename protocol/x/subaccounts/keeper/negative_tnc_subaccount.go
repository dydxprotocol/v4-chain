package keeper

import (
	"fmt"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/cosmos/gogoproto/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// GetNegativeTncSubaccountSeenAtBlock gets the last block height a negative TNC subaccount was
// seen in state for the given collateral pool address and a boolean for whether it exists in state.
func (k Keeper) GetNegativeTncSubaccountSeenAtBlock(
	ctx sdk.Context,
	perpetualId uint32,
) (uint32, bool, error) {
	store := prefix.NewStore(
		ctx.KVStore(k.storeKey),
		[]byte(types.NegativeTncSubaccountForCollateralPoolSeenAtBlockKeyPrefix),
	)

	suffix, err := k.getNegativeTncSubaccountStoreSuffix(ctx, perpetualId)
	if err != nil {
		return 0, false, err
	}

	blockHeight, exists := k.getNegativeTncSubaccountSeenAtBlock(store, suffix)
	return blockHeight, exists, nil
}

// Internal helper method to read the store using a store suffix.
func (k Keeper) getNegativeTncSubaccountSeenAtBlockWithSuffix(
	ctx sdk.Context,
	storeSuffix string,
) (uint32, bool) {
	store := prefix.NewStore(
		ctx.KVStore(k.storeKey),
		[]byte(types.NegativeTncSubaccountForCollateralPoolSeenAtBlockKeyPrefix),
	)

	return k.getNegativeTncSubaccountSeenAtBlock(store, storeSuffix)
}

// getNegativeTncSubaccountSeenAtBlock is a helper function that takes a store and returns the last
// block height a negative TNC subaccount was seen in state for the given collateral pool address
// and a boolean for whether it exists in state.
func (k Keeper) getNegativeTncSubaccountSeenAtBlock(
	store storetypes.KVStore,
	storeSuffix string,
) (uint32, bool) {
	b := store.Get(
		[]byte(storeSuffix),
	)
	blockHeight := gogotypes.UInt32Value{Value: 0}
	exists := false
	if b != nil {
		k.cdc.MustUnmarshal(b, &blockHeight)
		exists = true
	}

	return blockHeight.Value, exists
}

// SetNegativeTncSubaccountSeenAtBlock sets a block number in state where a negative TNC subaccount
// was seen for a specific collateral pool. This function will overwrite previous values at this key.
// This function will panic if the old block height is greater than the new block height.
func (k Keeper) SetNegativeTncSubaccountSeenAtBlock(
	ctx sdk.Context,
	perpetualId uint32,
	blockHeight uint32,
) error {
	store := prefix.NewStore(
		ctx.KVStore(k.storeKey),
		[]byte(types.NegativeTncSubaccountForCollateralPoolSeenAtBlockKeyPrefix),
	)

	storeSuffix, err := k.getNegativeTncSubaccountStoreSuffix(ctx, perpetualId)
	if err != nil {
		return err
	}

	// Panic if the stored block height value exists and is greater than the new block height value.
	currentValue, exists := k.getNegativeTncSubaccountSeenAtBlock(store, storeSuffix)
	if exists && blockHeight < currentValue {
		panic(
			fmt.Sprintf(
				"SetNegativeTncSubaccountSeenAtBlock: new block height (%d) is less than the current block height (%d)",
				blockHeight,
				currentValue,
			),
		)
	}

	blockHeightValue := gogotypes.UInt32Value{Value: blockHeight}
	store.Set(
		[]byte(storeSuffix),
		k.cdc.MustMarshal(&blockHeightValue),
	)

	return nil
}

func (k Keeper) getNegativeTncSubaccountStoreSuffix(
	ctx sdk.Context,
	perpetualId uint32,
) (string, error) {
	isIsolated, err := k.perpetualsKeeper.IsIsolatedPerpetual(ctx, perpetualId)
	if err != nil {
		return "", err
	}
	if isIsolated {
		return lib.UintToString(perpetualId), nil
	} else {
		return types.CrossCollateralSuffix, nil
	}
}

// getNegativeTncSubaccountStoresuffixes gets a slice of negative tnc subaccount store suffixes for
// the subaccounts in the slice of `settledUpdate`s passed in.
// The slice will be de-duplicated and will contain unique store suffixes.
func (k Keeper) getNegativeTncSubaccountStoresuffixes(
	ctx sdk.Context,
	settledUpdates []types.SettledUpdate,
) (
	suffixes []string,
	err error,
) {
	suffixesMap := make(map[string]bool)
	suffixes = make([]string, 0)
	for _, u := range settledUpdates {
		var suffix string
		if len(u.SettledSubaccount.PerpetualPositions) == 0 {
			suffix = types.CrossCollateralSuffix
		} else {
			suffix, err = k.getNegativeTncSubaccountStoreSuffix(ctx, u.SettledSubaccount.PerpetualPositions[0].PerpetualId)
			if err != nil {
				return nil, err
			}
		}
		if _, exists := suffixesMap[suffix]; !exists {
			suffixes = append(suffixes, suffix)
			suffixesMap[suffix] = true
		}
	}
	return suffixes, nil
}

// getLastBlockNegativeSubaccountSeen gets the last block where a subaccount with negative total net
// collateral was seen for subaccounts in a slice of settled updates.
func (k Keeper) getLastBlockNegativeSubaccountSeen(
	ctx sdk.Context,
	settledUpdates []types.SettledUpdate,
) (
	lastBlockNegativeSubaccountSeen uint32,
	negativeSubaccountExists bool,
	err error,
) {
	negativeTncSubaccountStoresuffixes, err := k.getNegativeTncSubaccountStoresuffixes(
		ctx,
		settledUpdates,
	)
	if err != nil {
		return uint32(0), false, err
	}
	lastBlockNegativeSubaccountSeen = uint32(0)
	negativeSubaccountExists = false
	for _, storeSuffix := range negativeTncSubaccountStoresuffixes {
		blockHeight, exists := k.getNegativeTncSubaccountSeenAtBlockWithSuffix(ctx, storeSuffix)
		if exists && blockHeight > lastBlockNegativeSubaccountSeen {
			lastBlockNegativeSubaccountSeen = blockHeight
			negativeSubaccountExists = true
		}
	}
	return lastBlockNegativeSubaccountSeen, negativeSubaccountExists, nil
}

// LegacyGetNegativeTncSubaccountSeenAtBlock gets the last block height a negative TNC subaccount was
// seen in state and a boolean for whether it exists in state.
// Deprecated: This is the legacy implementation and meant to be used for the v5.0.0 state migration.
// Use `GetNegativeTncSubaccountSeenAtBlock` instead.
func (k Keeper) LegacyGetNegativeTncSubaccountSeenAtBlock(
	ctx sdk.Context,
) (uint32, bool) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(
		// Key used in v4.0.0.
		[]byte("NegSA:"),
	)
	blockHeight := gogotypes.UInt32Value{Value: 0}
	exists := false
	if b != nil {
		k.cdc.MustUnmarshal(b, &blockHeight)
		exists = true
	}

	return blockHeight.Value, exists
}
