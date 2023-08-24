package keeper

import (
	"fmt"
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
			panic(fmt.Errorf("delayed message %v not found", id))
		}

		var msg sdk.Msg
		err := k.DecodeMessage(delayedMsg.Msg, &msg)

		if err != nil {
			panic(fmt.Errorf("Failed to decode delayed message: %w", err))
		}

		handler := k.router.Handler(msg)
		_, err = handler(ctx, msg)

		if err != nil {
			panic(fmt.Errorf("Failed to execute delayed message: %w", err))
		}

		err = k.DeleteMessage(ctx, id)
		if err != nil {
			panic(fmt.Errorf("Failed to delete delayed message: %w", err))
		}
	}
}
