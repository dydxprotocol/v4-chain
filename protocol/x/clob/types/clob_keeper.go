package types

import (
	"math/big"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"
)

type ClobKeeper interface {
	DeliverTxCancelOrder(ctx sdk.Context, msgCancelOrder *MsgCancelOrder, memclob MemClob) error
	CheckTxCancelOrder(ctx sdk.Context, msgCancelOrder *MsgCancelOrder) error
	CreatePerpetualClobPair(
		ctx sdk.Context,
		perpetualId uint32,
		stepSizeInBaseQuantums satypes.BaseQuantums,
		minOrderInBaseQuantums satypes.BaseQuantums,
		quantumConversionExponent int32,
		subticksPerTick uint32,
		status ClobPair_Status,
		makerFeePpm uint32,
		takerFeePpm uint32,
	) (
		ClobPair,
		error,
	)
	GetAllClobPair(ctx sdk.Context) (list []ClobPair)
	GetClobPair(ctx sdk.Context, id ClobPairId) (val ClobPair, found bool)
	DeliverTxPlaceOrder(
		ctx sdk.Context,
		msg *MsgPlaceOrder,
		performAddToOrderbookCollatCheck bool,
		memclob MemClob,
	) (
		orderSizeOptimisticallyFilledFromMatchingQuantums satypes.BaseQuantums,
		orderStatus OrderStatus,
		err error,
	)
	CheckTxPlaceOrder(ctx sdk.Context, msg *MsgPlaceOrder) (
		orderSizeOptimisticallyFilledFromMatchingQuantums satypes.BaseQuantums,
		orderStatus OrderStatus,
		err error,
	)
	PruneStateFillAmountsForShortTermOrders(
		ctx sdk.Context,
	)
	PruneExpiredSeenPlaceOrders(
		ctx sdk.Context,
		goodTilBlockHeightToPrune uint32,
	)

	RemoveClobPair(ctx sdk.Context, id ClobPairId)
	ProcessProposerOperations(
		ctx sdk.Context,
		operations []Operation,
		addToOrderbookCollatCheckOrderHashesSet map[OrderHash]bool,
	) error
	LiquidationsKeeper
	LiquidationsConfigKeeper
	AddOrderToOrderbookCollatCheck(
		ctx sdk.Context,
		clobPairId ClobPairId,
		subaccountOpenOrders map[satypes.SubaccountId][]PendingOpenOrder,
	) (
		success bool,
		successPerUpdate map[satypes.SubaccountId]satypes.UpdateResult,
	)
	GetStatePosition(
		ctx sdk.Context,
		subaccountId satypes.SubaccountId,
		clobPairId ClobPairId,
	) (
		positionSizeQuantums *big.Int,
	)
	ProcessSingleMatch(
		ctx sdk.Context,
		matchWithOrders MatchWithOrders,
	) (
		success bool,
		takerUpdateResult satypes.UpdateResult,
		makerUpdateResult satypes.UpdateResult,
		offchainUpdates *OffchainUpdates,
		err error,
	)
	SetStatefulOrderPlacement(
		ctx sdk.Context,
		order Order,
		blockHeight uint32,
	)
	GetStatefulOrderPlacement(
		ctx sdk.Context,
		orderId OrderId,
	) (val StatefulOrderPlacement, found bool)
	DeleteStatefulOrderPlacement(
		ctx sdk.Context,
		orderId OrderId,
	)
	DoesStatefulOrderExistInState(
		ctx sdk.Context,
		order Order,
	) bool
	RemoveOrderFillAmount(ctx sdk.Context, orderId OrderId)
	MustAddOrderToStatefulOrdersTimeSlice(
		ctx sdk.Context,
		goodTilBlockTime time.Time,
		orderId OrderId,
	)
	GetStatefulOrdersTimeSlice(ctx sdk.Context, goodTilBlockTime time.Time) (
		orderIds []OrderId,
	)
	MustRemoveStatefulOrder(
		ctx sdk.Context,
		goodTilBlockTime time.Time,
		orderId OrderId,
	)
	RemoveExpiredStatefulOrdersTimeSlices(ctx sdk.Context, blockTime time.Time) (
		expiredOrderIds []OrderId,
	)
	GetProcessProposerMatchesEvents(ctx sdk.Context) ProcessProposerMatchesEvents
	MustSetProcessProposerMatchesEvents(
		ctx sdk.Context,
		processProposerMatchesEvents ProcessProposerMatchesEvents,
	)
	SetBlockTimeForLastCommittedBlock(ctx sdk.Context)
	MustGetBlockTimeForLastCommittedBlock(ctx sdk.Context) (blockTime time.Time)
	GetNumClobPairs(ctx sdk.Context) uint32
}
