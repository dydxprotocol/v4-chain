package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/cosmos/gogoproto/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// GetNegativeTncSubaccountSeenAtBlock gets the last block height a negative TNC subaccount was
// seen in state. Note this defaults to 0 if it has not been set.
func (k Keeper) GetNegativeTncSubaccountSeenAtBlock(
	ctx sdk.Context,
) uint32 {
	store := ctx.KVStore(k.storeKey)
	return k.getNegativeTncSubaccountSeenAtBlock(store)
}

// getNegativeTncSubaccountSeenAtBlock is a helper function that takes a store and returns the last
// block height a negative TNC subaccount was seen in state.
func (k Keeper) getNegativeTncSubaccountSeenAtBlock(
	store sdk.KVStore,
) uint32 {
	b := store.Get(
		[]byte(types.NegativeTncSubaccountSeenAtBlockKey),
	)
	blockHeight := gogotypes.UInt32Value{Value: 0}
	if b != nil {
		k.cdc.MustUnmarshal(b, &blockHeight)
	}

	return blockHeight.Value
}

// SetNegativeTncSubaccountSeenAtBlock sets a block number in state where a negative TNC subaccount
// was seen. This function will overwrite previous values at this key.
// This function will panic if the old block height is greater than the new block height.
func (k Keeper) SetNegativeTncSubaccountSeenAtBlock(
	ctx sdk.Context,
	blockHeight uint32,
) {
	store := ctx.KVStore(k.storeKey)

	// Panic if the new block height is less than the current block height.
	currentValue := k.getNegativeTncSubaccountSeenAtBlock(store)
	if blockHeight < currentValue {
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
		[]byte(types.NegativeTncSubaccountSeenAtBlockKey),
		k.cdc.MustMarshal(&blockHeightValue),
	)
}
