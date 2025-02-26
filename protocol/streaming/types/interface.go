package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

type FullNodeStreamingManager interface {
	Enabled() bool
	Stop()

	// Subscribe to streams
	Subscribe(
		clobPairIds []uint32,
		subaccountIds []*satypes.SubaccountId,
		marketIds []uint32,
		filterOrdersBySubaccountId bool,
		srv OutgoingMessageSender,
	) (
		err error,
	)

	// L3+ Orderbook updates.
	InitializeNewStreams(
		getOrderbookSnapshot func(clobPairId clobtypes.ClobPairId) *clobtypes.OffchainUpdates,
		subaccountSnapshots map[satypes.SubaccountId]*satypes.StreamSubaccountUpdate,
		priceSnapshots map[uint32]*pricestypes.StreamPriceUpdate,
		blockHeight uint32,
		execMode sdk.ExecMode,
	)
	GetSubaccountSnapshotsForInitStreams(
		getSubaccountSnapshot func(subaccountId satypes.SubaccountId) *satypes.StreamSubaccountUpdate,
	) map[satypes.SubaccountId]*satypes.StreamSubaccountUpdate
	GetPriceSnapshotsForInitStreams(
		getPriceSnapshot func(marketId uint32) *pricestypes.StreamPriceUpdate,
	) map[uint32]*pricestypes.StreamPriceUpdate
	SendOrderbookUpdates(
		offchainUpdates *clobtypes.OffchainUpdates,
		ctx sdk.Context,
	)
	SendOrderbookFillUpdate(
		orderbookFill clobtypes.StreamOrderbookFill,
		ctx sdk.Context,
		perpetualIdToClobPairId map[uint32][]clobtypes.ClobPairId,
	)
	SendTakerOrderStatus(
		takerOrder clobtypes.StreamTakerOrder,
		ctx sdk.Context,
	)
	SendSubaccountUpdate(
		ctx sdk.Context,
		subaccountUpdate satypes.StreamSubaccountUpdate,
	)
	SendPriceUpdate(
		ctx sdk.Context,
		priceUpdate pricestypes.StreamPriceUpdate,
	)
	GetStagedFinalizeBlockEvents(
		ctx sdk.Context,
	) []clobtypes.StagedFinalizeBlockEvent
	TracksSubaccountId(id satypes.SubaccountId) bool
	TracksMarketId(marketId uint32) bool
	StreamBatchUpdatesAfterFinalizeBlock(
		ctx sdk.Context,
		orderBookUpdatesToSyncLocalOpsQueue *clobtypes.OffchainUpdates,
		perpetualIdToClobPairId map[uint32][]clobtypes.ClobPairId,
	)
}

type OutgoingMessageSender interface {
	Send(*clobtypes.StreamOrderbookUpdatesResponse) error
}
