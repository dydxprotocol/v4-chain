package types

import (
	"math/big"
	"time"

	"cosmossdk.io/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

type ClobKeeper interface {
	LiquidationsKeeper
	LiquidationsConfigKeeper

	IsInitialized() bool
	Initialize(ctx sdk.Context)

	AddOrderToOrderbookSubaccountUpdatesCheck(
		ctx sdk.Context,
		clobPairId ClobPairId,
		subaccountOpenOrders map[satypes.SubaccountId][]PendingOpenOrder,
	) (
		success bool,
		successPerUpdate map[satypes.SubaccountId]satypes.UpdateResult,
	)
	BatchCancelShortTermOrder(
		ctx sdk.Context,
		msg *MsgBatchCancel,
	) (success []uint32, failure []uint32, err error)
	CancelShortTermOrder(ctx sdk.Context, msg *MsgCancelOrder) error
	CancelStatefulOrder(ctx sdk.Context, msg *MsgCancelOrder) error
	CreatePerpetualClobPair(
		ctx sdk.Context,
		clobPairId uint32,
		perpetualId uint32,
		stepSizeInBaseQuantums satypes.BaseQuantums,
		quantumConversionExponent int32,
		subticksPerTick uint32,
		status ClobPair_Status,
	) (
		ClobPair,
		error,
	)
	HandleMsgCancelOrder(
		ctx sdk.Context,
		msg *MsgCancelOrder,
	) (err error)
	HandleMsgPlaceOrder(
		ctx sdk.Context,
		msg *MsgPlaceOrder,
		isInternalOrder bool,
	) (err error)
	GetAllClobPairs(ctx sdk.Context) (list []ClobPair)
	GetClobPair(ctx sdk.Context, id ClobPairId) (val ClobPair, found bool)
	HasAuthority(authority string) bool
	PlaceShortTermOrder(ctx sdk.Context, msg *MsgPlaceOrder) (
		orderSizeOptimisticallyFilledFromMatchingQuantums satypes.BaseQuantums,
		orderStatus OrderStatus,
		err error,
	)
	PlaceStatefulOrder(
		ctx sdk.Context,
		msg *MsgPlaceOrder,
		isInternalOrder bool,
	) error

	PruneStateFillAmountsForShortTermOrders(
		ctx sdk.Context,
	)

	RemoveClobPair(ctx sdk.Context, id ClobPairId)
	ProcessProposerOperations(
		ctx sdk.Context,
		operations []OperationRaw,
	) error
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
	RateLimitBatchCancel(ctx sdk.Context, order *MsgBatchCancel) error
	InitializeBlockRateLimit(ctx sdk.Context, config BlockRateLimitConfiguration) error
	GetBlockRateLimitConfiguration(
		ctx sdk.Context,
	) (config BlockRateLimitConfiguration)
	InitializeEquityTierLimit(ctx sdk.Context, config EquityTierLimitConfiguration) error
	Logger(ctx sdk.Context) log.Logger
	UpdateClobPair(
		ctx sdk.Context,
		clobPair ClobPair,
	) error
	UpdateLiquidationsConfig(ctx sdk.Context, config LiquidationsConfig) error
	// Gprc streaming
	InitializeNewGrpcStreams(ctx sdk.Context)
	SendOrderbookUpdates(
		ctx sdk.Context,
		offchainUpdates *OffchainUpdates,
		snapshot bool,
	)
	MigratePruneableOrders(ctx sdk.Context)
}
