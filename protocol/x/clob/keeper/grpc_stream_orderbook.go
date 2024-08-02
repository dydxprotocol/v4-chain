package keeper

import (
	"fmt"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"io"
)

func (k Keeper) StreamOrderbookUpdates(stream types.Query_StreamOrderbookUpdatesServer) error {
	// A channel to handle incoming requests
	reqChan := make(chan *types.StreamOrderbookUpdatesRequest)

	// Goroutine to handle receiving requests from the client
	go func() {
		defer close(reqChan)
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				// Client closed the stream
				return
			}
			if err != nil {
				panic(fmt.Sprintf("Error receiving request in orderbook GRPC stream: %v", err))
			}
			reqChan <- req
		}
	}()

	// Context to handle stream lifecycle
	ctx := stream.Context()

	for {
		select {
		case req := <-reqChan:
			// Handle the incoming request and subscribe
			err := k.GetFullNodeStreamingManager().Subscribe(
				req.GetClobPairId(),
				stream,
			)
			if err != nil {
				return err
			}

		case <-ctx.Done():
			// Stream context is done, exit the function
			return nil
		}
	}
}
