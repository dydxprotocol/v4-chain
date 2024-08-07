package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

type FullNodeStreamingManager interface {
	Enabled() bool
	Stop()

	// Subscribe to streams
	Subscribe(
		clobPairIds []uint32,
		subaccountIds []*satypes.SubaccountId,
		srv OutgoingMessageSender,
	) (
		err error,
	)

	// L3+ Orderbook updates.
	InitializeNewStreams(
		getOrderbookSnapshot func(clobPairId clobtypes.ClobPairId) *clobtypes.OffchainUpdates,
		getSubaccountSnapshot func(subaccountId satypes.SubaccountId) *satypes.StreamSubaccountUpdate,
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
	SendSubaccountUpdates(
		subaccountUpdates []satypes.StreamSubaccountUpdate,
		blockHeight uint32,
		execMode sdk.ExecMode,
	)
}

type OutgoingMessageSender interface {
	Send(*clobtypes.StreamOrderbookUpdatesResponse) error
}
