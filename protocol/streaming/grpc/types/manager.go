package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
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
	Stop()
	SendOrderbookFillUpdates(
		orderbookFills []clobtypes.StreamOrderbookFill,
		blockHeight uint32,
		execMode sdk.ExecMode,
	)
}
