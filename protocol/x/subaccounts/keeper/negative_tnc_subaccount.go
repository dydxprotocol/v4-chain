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

	suffix, err := k.getNegativeTncSubaccountStoreSuffx(ctx, perpetualId)
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

	storeSuffix, err := k.getNegativeTncSubaccountStoreSuffx(ctx, perpetualId)
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

func (k Keeper) getNegativeTncSubaccountStoreSuffx(
	ctx sdk.Context,
	perpetualId uint32,
) (string, error) {
	isIsolated, err := k.perpetualsKeeper.IsIsolatedPerpetual(ctx, perpetualId)
	if err != nil {
		return "", err
	}
	if isIsolated {
		return types.CrossCollateralSuffix, nil
	} else {
		return lib.UintToString(perpetualId), nil
	}
}
