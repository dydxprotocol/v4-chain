package keeper

import (
	"fmt"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/cosmos/gogoproto/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// GetNegativeTncSubaccountSeenAtBlock gets the last block height a negative TNC subaccount was
// seen in state for the given collateral pool address and a boolean for whether it exists in state.
func (k Keeper) GetNegativeTncSubaccountSeenAtBlock(
	ctx sdk.Context,
	collateralPoolAddr sdk.AccAddress,
) (uint32, bool) {
	store := prefix.NewStore(
		ctx.KVStore(k.storeKey),
		[]byte(types.NegativeTncSubaccountForCollateralPoolSeenAtBlockKeyPrefix),
	)
	return k.getNegativeTncSubaccountSeenAtBlock(store, collateralPoolAddr)
}

// getNegativeTncSubaccountSeenAtBlock is a helper function that takes a store and returns the last
// block height a negative TNC subaccount was seen in state for the given collateral pool address
// and a boolean for whether it exists in state.
func (k Keeper) getNegativeTncSubaccountSeenAtBlock(
	store storetypes.KVStore,
	collateralPoolAddr sdk.AccAddress,
) (uint32, bool) {
	b := store.Get(
		collateralPoolAddr,
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
	collateralPoolAddr sdk.AccAddress,
	blockHeight uint32,
) {
	store := prefix.NewStore(
		ctx.KVStore(k.storeKey),
		[]byte(types.NegativeTncSubaccountForCollateralPoolSeenAtBlockKeyPrefix),
	)

	// Panic if the stored block height value exists and is greater than the new block height value.
	currentValue, exists := k.getNegativeTncSubaccountSeenAtBlock(store, collateralPoolAddr)
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
		collateralPoolAddr.Bytes(),
		k.cdc.MustMarshal(&blockHeightValue),
	)
}
