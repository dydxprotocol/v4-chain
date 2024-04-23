package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/streaming/grpc/client"
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
	SubscribeTestClient(
		client *client.GrpcClient,
	)
	GetUninitializedClobPairIds() []uint32
	SendOrderbookUpdates(
		ctx sdk.Context,
		offchainUpdates *clobtypes.OffchainUpdates,
		snapshot bool,
		blockHeight uint32,
		execMode sdk.ExecMode,
	)
	SendOrderbookMatchFillUpdates(
		ctx sdk.Context,
		matches []clobtypes.OrderBookMatchFill,
		blockHeight uint32,
		execMode sdk.ExecMode,
	)
}
