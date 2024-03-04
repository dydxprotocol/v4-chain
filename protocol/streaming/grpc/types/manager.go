package types

import (
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

type GrpcStreamingManager interface {
	Enabled() bool

	// L3+ Orderbook updates.
	Subscribe(
		req clobtypes.StreamOrderbookUpdatesRequest,
		srv clobtypes.Query_StreamOrderbookUpdatesServer,
	) (
		finished chan bool,
		err error,
	)
	SendOrderbookUpdates(*clobtypes.OffchainUpdates)
}
