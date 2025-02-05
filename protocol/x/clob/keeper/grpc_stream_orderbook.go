package keeper

import (
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

func (k Keeper) StreamOrderbookUpdates(
	req *types.StreamOrderbookUpdatesRequest,
	stream types.Query_StreamOrderbookUpdatesServer,
) error {
	err := k.GetFullNodeStreamingManager().Subscribe(
		req.GetClobPairId(),
		req.GetSubaccountIds(),
		req.GetMarketIds(),
		req.GetFilterOrdersBySubaccountId(),
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
