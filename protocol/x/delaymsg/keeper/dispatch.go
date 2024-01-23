package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/abci"
	"github.com/dydxprotocol/v4-chain/protocol/lib/log"
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
)

// DispatchMessagesForBlock executes all delayed messages scheduled for the given block height and deletes
// the messages. If there are no delayed messages scheduled for this block, this function does nothing. It is
// expected that this function is called at the end of every block.
func DispatchMessagesForBlock(k types.DelayMsgKeeper, ctx sdk.Context) {
	blockMessageIds, found := k.GetBlockMessageIds(ctx, lib.MustConvertIntegerToUint32(ctx.BlockHeight()))

	// If there are no delayed messages scheduled for this block, return.
	if !found {
		return
	}

	// Maintain a list of events emitted by all delayed messages executed in this block.
	// As message handlers create new event managers, such emitted events need to be
	// explicitly propagated to the current context.
	// Note: events in EndBlocker can be found in `end_block_events` in response from
	// `/block_results` endpoint.
	var events sdk.Events

	// Execute all delayed messages scheduled for this block.
	for _, id := range blockMessageIds.Ids {
		delayedMsg, found := k.GetMessage(ctx, id)
		if !found {
			log.ErrorLog(
				ctx,
				"delayed message not found",
				types.IdLogKey, id,
			)
			continue
		}

		msg, err := delayedMsg.GetMessage()
		if err != nil {
			log.ErrorLogWithError(
				ctx,
				"failed to decode delayed message",
				err,
				types.IdLogKey, id,
			)
			continue
		}

		// Execute the message in a cached context.
		if err = abci.RunCached(ctx, func(ctx sdk.Context) error {
			handler := k.Router().Handler(msg)
			res, err := handler(ctx, msg)
			if err != nil {
				return err
			}
			// Append events emitted in message handler to `events`.
			events = append(events, res.GetEvents()...)
			return nil
		}); err != nil {
			log.ErrorLogWithError(
				ctx,
				"failed to execute delayed message",
				err,
				types.IdLogKey, id,
				types.MessageContentLogKey, msg.String(),
				types.MessageTypeUrlLogKey, delayedMsg.Msg.TypeUrl,
			)
		}
	}

	// Propagate events emitted in message handlers to current context.
	ctx.EventManager().EmitEvents(events)

	// Delete executed messages.
	for _, id := range blockMessageIds.Ids {
		if err := k.DeleteMessage(ctx, id); err != nil {
			log.ErrorLogWithError(
				ctx,
				"failed to delete delayed message",
				err,
				types.IdLogKey, id,
			)
		}
	}
}
