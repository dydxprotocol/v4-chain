package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

type FullNodeStreamingManager interface {
	Enabled() bool
	Stop()

	// Subscribe to orderbook streams
	SubscribeToOrderbookStream(
		clobPairIds []uint32,
		srv OutgoingOrderbookMessageSender,
	) (
		err error,
	)

	// Subscribe to subaccount streams
	SubscribeToSubaccountStream(
		subaccountIds []*satypes.SubaccountId,
		srv OutgoingSubaccountMessageSender,
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
	SendOrderbookFillUpdates(
		orderbookFills []clobtypes.StreamOrderbookFill,
		blockHeight uint32,
		execMode sdk.ExecMode,
		perpetualIdToClobPairId map[uint32][]clobtypes.ClobPairId,
	)
}

type OutgoingOrderbookMessageSender interface {
	Send(*clobtypes.StreamOrderbookUpdatesResponse) error
}

type OutgoingSubaccountMessageSender interface {
	Send(*satypes.StreamSubaccountUpdatesResponse) error
}
