package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/abci"
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
)

// DispatchMessagesForBlock executes all delayed messages scheduled for the given block height and deletes
// the messages. If there are no delayed messages scheduled for this block, this function does nothing. It is
// expected that this function is called at the end of every block.
func DispatchMessagesForBlock(k types.DelayMsgKeeper, ctx sdk.Context) {
	blockMessageIds, found := k.GetBlockMessageIds(ctx, ctx.BlockHeight())

	// If there are no delayed messages scheduled for this block, return.
	if !found {
		return
	}

	// Maintain a list of events emitted by all delayed messages executed in this block.
	// As each delayed message is executed in a cached context, their emitted events need
	// to be explicitly propagated to the current context.
	// Note: events in EndBlocker can be found in `end_block_events` in response from
	// `/block_results` endpoint.
	var events sdk.Events

	// Execute all delayed messages scheduled for this block and delete them from the store.
	for _, id := range blockMessageIds.Ids {
		delayedMsg, found := k.GetMessage(ctx, id)
		if !found {
			k.Logger(ctx).Error("delayed message %v not found", id)
			continue
		}

		msg, err := delayedMsg.GetMessage()
		if err != nil {
			k.Logger(ctx).Error("failed to decode delayed message with id %v: %v", id, err)
			continue
		}

		if err = abci.RunCached(ctx, func(ctx sdk.Context) error {
			handler := k.Router().Handler(msg)
			res, err := handler(ctx, msg)
			if err != nil {
				return err
			}
			// Append events emitted in cached context to `events`.
			events = append(events, res.GetEvents()...)
			return nil
		}); err != nil {
			k.Logger(ctx).Error("failed to execute delayed message with id %v: %v", id, err)
		}
	}

	// Propagate events emitted in cached contexts to current context.
	ctx.EventManager().EmitEvents(events)

	for _, id := range blockMessageIds.Ids {
		if err := k.DeleteMessage(ctx, id); err != nil {
			k.Logger(ctx).Error("failed to delete delayed message: %w", err)
		}
	}
}
