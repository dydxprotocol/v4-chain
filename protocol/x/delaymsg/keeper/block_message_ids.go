package keeper

import (
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
)

// newBlockIdStore creates a new prefix store for BlockMessageIds.
func (k Keeper) newBlockMessageIdsStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.BlockMessageIdsPrefix))
}

// GetBlockMessageIds gets the ids of delayed messages to execute at a given block.
func (k Keeper) GetBlockMessageIds(
	ctx sdk.Context,
	blockHeight int64,
) (
	blockMessageIds types.BlockMessageIds,
	found bool,
) {
	store := k.newBlockMessageIdsStore(ctx)
	b := store.Get(lib.Int64ToBytesForState(blockHeight))

	if b == nil {
		return types.BlockMessageIds{}, false
	}

	blockMessageIds = types.BlockMessageIds{}
	k.cdc.MustUnmarshal(b, &blockMessageIds)
	return blockMessageIds, true
}

// addMessageIdToBlock adds a message id to the list of message ids for a block. This method should only
// be called from DelayMessageByBlocks whenever a new message is added. When this restriction is followed and the id is
// only called during DelayedMessage creation, the message ids for a block will always be in ascending order.
func (k Keeper) addMessageIdToBlock(
	ctx sdk.Context,
	id uint32,
	blockHeight int64,
) {
	store := k.newBlockMessageIdsStore(ctx)
	var blockMessageIds types.BlockMessageIds
	key := lib.Int64ToBytesForState(blockHeight)
	if b := store.Get(key); b != nil {
		k.cdc.MustUnmarshal(b, &blockMessageIds)
		blockMessageIds.Ids = append(blockMessageIds.Ids, id)
	} else {
		blockMessageIds = types.BlockMessageIds{
			Ids: []uint32{id},
		}
	}
	store.Set(key, k.cdc.MustMarshal(&blockMessageIds))
}

// deleteMessageIdFromBlock deletes a message id from the list of message ids for a block. This method should
// only be called from DeleteMessage whenever a message is deleted. Message id removal assumes non-duplicate ids and
// respects the original ordering of message ids for a block, so if the message ids were already sorted, they will
// remain sorted.
func (k Keeper) deleteMessageIdFromBlock(
	ctx sdk.Context,
	id uint32,
	blockHeight int64,
) (
	err error,
) {
	blockMessageIds, found := k.GetBlockMessageIds(ctx, blockHeight)
	if !found {
		return sdkerrors.Wrapf(
			types.ErrInvalidInput,
			"block %v not found",
			blockHeight,
		)
	}

	for i, blockMessageId := range blockMessageIds.Ids {
		// Skip ids that don't match.
		if blockMessageId != id {
			continue
		}

		// Reconstruct the id list in the same order, with this id removed. This reconstruction of ids preserves the
		// original ordering of ids.
		blockMessageIds.Ids = append(blockMessageIds.Ids[:i], blockMessageIds.Ids[i+1:]...)

		// If the remaining list of ids is empty, go ahead and delete the BlockMessageIds from the store.
		if len(blockMessageIds.Ids) == 0 {
			k.newBlockMessageIdsStore(ctx).Delete(lib.Int64ToBytesForState(blockHeight))
		} else {
			// Otherwise, update the BlockMessageIds to have the id of this delayed message removed.
			k.newBlockMessageIdsStore(ctx).Set(
				lib.Int64ToBytesForState(blockHeight),
				k.cdc.MustMarshal(&blockMessageIds),
			)
		}
		return nil
	}

	// If we make it here, the message id was not found in the block.
	return sdkerrors.Wrapf(
		types.ErrInvalidInput,
		"message id %v not found in block %v",
		id,
		blockHeight,
	)
}
