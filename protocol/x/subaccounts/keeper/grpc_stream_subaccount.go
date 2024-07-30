package keeper

import (
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

func (k Keeper) StreamSubaccountUpdates(
	req *types.StreamSubaccountUpdatesRequest,
	stream types.Query_StreamSubaccountUpdatesServer,
) error {
	err := k.GetFullNodeStreamingManager().SubscribeToSubaccountStream(
		req.GetSubaccountIds(),
		stream,
	)
	if err != nil {
		return err
	}

	// Keep this scope alive because once this scope exits - the stream is closed
	ctx := stream.Context()
	<-ctx.Done()
	return nil
}
