package types

import (
	"math/big"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

type ClobKeeper interface {
	AddOrderToOrderbookCollatCheck(
		ctx sdk.Context,
		clobPairId ClobPairId,
		subaccountOpenOrders map[satypes.SubaccountId][]PendingOpenOrder,
	) (
		success bool,
		successPerUpdate map[satypes.SubaccountId]satypes.UpdateResult,
	)
	CancelShortTermOrder(ctx sdk.Context, msg *MsgCancelOrder) error
	CancelStatefulOrder(ctx sdk.Context, msg *MsgCancelOrder) error
	CreatePerpetualClobPair(
		ctx sdk.Context,
		perpetualId uint32,
		minOrderInBaseQuantums satypes.BaseQuantums,
		stepSizeInBaseQuantums satypes.BaseQuantums,
		quantumConversionExponent int32,
		subticksPerTick uint32,
		status ClobPair_Status,
	) (
		ClobPair,
		error,
	)
	GetAllClobPair(ctx sdk.Context) (list []ClobPair)
	GetClobPair(ctx sdk.Context, id ClobPairId) (val ClobPair, found bool)
	PlaceShortTermOrder(ctx sdk.Context, msg *MsgPlaceOrder) (
		orderSizeOptimisticallyFilledFromMatchingQuantums satypes.BaseQuantums,
		orderStatus OrderStatus,
		err error,
	)
	PlaceStatefulOrder(ctx sdk.Context, msg *MsgPlaceOrder) error
	PruneStateFillAmountsForShortTermOrders(
		ctx sdk.Context,
	)

	RemoveClobPair(ctx sdk.Context, id ClobPairId)
	ProcessProposerOperations(
		ctx sdk.Context,
		operations []OperationRaw,
	) error
	LiquidationsKeeper
	LiquidationsConfigKeeper
	GetStatePosition(
		ctx sdk.Context,
		subaccountId satypes.SubaccountId,
		clobPairId ClobPairId,
	) (
		positionSizeQuantums *big.Int,
	)
	ProcessSingleMatch(
		ctx sdk.Context,
		matchWithOrders *MatchWithOrders,
	) (
		success bool,
		takerUpdateResult satypes.UpdateResult,
		makerUpdateResult satypes.UpdateResult,
		offchainUpdates *OffchainUpdates,
		err error,
	)
	SetLongTermOrderPlacement(
		ctx sdk.Context,
		order Order,
		blockHeight uint32,
	)
	GetLongTermOrderPlacement(
		ctx sdk.Context,
		orderId OrderId,
	) (val LongTermOrderPlacement, found bool)
	DeleteLongTermOrderPlacement(
		ctx sdk.Context,
		orderId OrderId,
	)
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
	GetNumClobPairs(ctx sdk.Context) uint32
	PerformOrderCancellationStatefulValidation(
		ctx sdk.Context,
		msgCancelOrder *MsgCancelOrder,
		blockHeight uint32,
	) error
	PerformStatefulOrderValidation(
		ctx sdk.Context,
		order *Order,
		blockHeight uint32,
		isPreexistingStatefulOrder bool,
	) error
	GetIndexerEventManager() indexer_manager.IndexerEventManager
	RateLimitCancelOrder(ctx sdk.Context, order *MsgCancelOrder) error
	RateLimitPlaceOrder(ctx sdk.Context, order *MsgPlaceOrder) error
}
