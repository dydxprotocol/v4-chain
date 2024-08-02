package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/streaming/grpc/client"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

type FullNodeStreamingManager interface {
	Enabled() bool
	Stop()

	// Subscribe to streams
	Subscribe(
		clobPairIds []uint32,
		srv OutgoingMessageSender,
	) (
		err error,
	)

	// L3+ Orderbook updates.
	InitializeNewStreams(
		getOrderbookSnapshot func(clobPairId clobtypes.ClobPairId) *clobtypes.OffchainUpdates,
		blockHeight uint32,
		execMode sdk.ExecMode,
	)
	SendOrderbookUpdates(
		offchainUpdates *clobtypes.OffchainUpdates,
		blockHeight uint32,
		execMode sdk.ExecMode,
	)
	SubscribeTestClient(
		client *client.GrpcClient,
	)
	SendOrderbookFillUpdates(
		orderbookFills []clobtypes.StreamOrderbookFill,
		blockHeight uint32,
		execMode sdk.ExecMode,
		perpetualIdToClobPairId map[uint32][]clobtypes.ClobPairId,
	)
}

type OutgoingMessageSender interface {
	Send(*clobtypes.StreamOrderbookUpdatesResponse) error
}
