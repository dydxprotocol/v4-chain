package keeper

import (
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

func (k Keeper) StreamOrderbookUpdates(
	req *types.StreamOrderbookUpdatesRequest,
	stream types.Query_StreamOrderbookUpdatesServer,
) error {
	finished, err := k.GetGrpcStreamingManager().Subscribe(*req, stream)
	if err != nil {
		return err
	}

	// Keep this scope alive because once this scope exits - the stream is closed
	ctx := stream.Context()
	for {
		select {
		case <-finished:
			return nil
		case <-ctx.Done():
			return nil
		}
	}
}
