package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DispatchMessagesForBlock executes all delayed messages scheduled for the given block height and deletes
// the messages. If there are no delayed messages scheduled for this block, this function does nothing. It is
// expected that this function is called at the end of every block.
func (k Keeper) DispatchMessagesForBlock(ctx sdk.Context) {
	blockMessageIds, found := k.GetBlockMessageIds(ctx, ctx.BlockHeight())

	// If there are no delayed messages scheduled for this block, return.
	if !found {
		return
	}

	// Execute all delayed messages scheduled for this block and delete them from the store.
	for _, id := range blockMessageIds.Ids {
		delayedMsg, found := k.GetMessage(ctx, id)
		if !found {
			k.Logger(ctx).Error("delayed message %v not found", id)
			continue
		}

		var msg sdk.Msg
		err := k.DecodeMessage(delayedMsg.Msg, &msg)

		if err != nil {
			k.Logger(ctx).Error("failed to decode delayed message with id %v: %v", id, err)
			continue
		}

		handler := k.router.Handler(msg)
		_, err = handler(ctx, msg)

		if err != nil {
			k.Logger(ctx).Error("failed to execute delayed message with id %v: %v", id, err)
		}
	}

	for _, id := range blockMessageIds.Ids {
		if err := k.DeleteMessage(ctx, id); err != nil {
			k.Logger(ctx).Error("failed to delete delayed message: %w", err)
		}
	}
}
