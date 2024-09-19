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
		subaccountSnapshots map[satypes.SubaccountId]*satypes.StreamSubaccountUpdate,
		blockHeight uint32,
		execMode sdk.ExecMode,
	)
	GetSubaccountSnapshotsForInitStreams(
		getSubaccountSnapshot func(subaccountId satypes.SubaccountId) *satypes.StreamSubaccountUpdate,
	) map[satypes.SubaccountId]*satypes.StreamSubaccountUpdate
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
	SendTakerOrderStatus(
		takerOrder clobtypes.StreamTakerOrder,
		blockHeight uint32,
		execMode sdk.ExecMode,
	)
	SendFinalizedSubaccountUpdates(
		subaccountUpdates []satypes.StreamSubaccountUpdate,
		blockHeight uint32,
		execMode sdk.ExecMode,
	)
	StageFinalizeBlockFill(
		ctx sdk.Context,
		fill clobtypes.StreamOrderbookFill,
	)
	StageFinalizeBlockSubaccountUpdate(
		ctx sdk.Context,
		subaccountUpdate satypes.StreamSubaccountUpdate,
	)
	GetStagedFinalizeBlockEvents(
		ctx sdk.Context,
	) []clobtypes.StagedFinalizeBlockEvent
	TracksSubaccountId(id satypes.SubaccountId) bool
	StreamBatchUpdatesAfterFinalizeBlock(
		ctx sdk.Context,
		orderBookUpdatesToSyncLocalOpsQueue *clobtypes.OffchainUpdates,
		perpetualIdToClobPairId map[uint32][]clobtypes.ClobPairId,
	)
}

type OutgoingMessageSender interface {
	Send(*clobtypes.StreamOrderbookUpdatesResponse) error
}
