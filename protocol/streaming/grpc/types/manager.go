package types

import (
	clobtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GrpcStreamingManager interface {
	Enabled() bool

	// L3+ Orderbook updates.
	Subscribe(
		req clobtypes.StreamOrderbookUpdatesRequest,
		srv clobtypes.Query_StreamOrderbookUpdatesServer,
	) (
		err error,
	)
	GetUninitializedClobPairIds() []uint32
	SendOrderbookUpdates(
		offchainUpdates *clobtypes.OffchainUpdates,
		snapshot bool,
		blockHeight uint32,
		execMode sdk.ExecMode,
	)
}
