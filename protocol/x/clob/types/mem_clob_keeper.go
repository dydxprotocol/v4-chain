package types

import (
	"math/big"
	"time"

	"github.com/cometbft/cometbft/libs/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

type MemClobKeeper interface {
	GetOrderFillAmount(
		ctx sdk.Context,
		orderId OrderId,
	) (
		exists bool,
		fillAmount satypes.BaseQuantums,
		prunableBlockHeight uint32,
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
	AddOrderToOrderbookCollatCheck(
		ctx sdk.Context,
		clobPairId ClobPairId,
		subaccountOpenOrders map[satypes.SubaccountId][]PendingOpenOrder,
	) (
		success bool,
		successPerUpdate map[satypes.SubaccountId]satypes.UpdateResult,
	)
	CanDeleverageSubaccount(
		ctx sdk.Context,
		subaccountId satypes.SubaccountId,
	) (bool, error)
	GetStatePosition(
		ctx sdk.Context,
		subaccountId satypes.SubaccountId,
		clobPairId ClobPairId,
	) (
		positionSizeQuantums *big.Int,
	)
	ReplayPlaceOrder(
		ctx sdk.Context,
		msg *MsgPlaceOrder,
	) (
		orderSizeOptimisticallyFilledFromMatchingQuantums satypes.BaseQuantums,
		orderStatus OrderStatus,
		offchainUpdates *OffchainUpdates,
		err error,
	)
	CancelShortTermOrder(
		ctx sdk.Context,
		msgCancelOrder *MsgCancelOrder,
	) error
	GetLongTermOrderPlacement(
		ctx sdk.Context,
		orderId OrderId,
	) (val LongTermOrderPlacement, found bool)
	SetLongTermOrderPlacement(
		ctx sdk.Context,
		order Order,
		blockHeight uint32,
	)
	MustAddOrderToStatefulOrdersTimeSlice(
		ctx sdk.Context,
		goodTilBlockTime time.Time,
		orderId OrderId,
	)
	OffsetSubaccountPerpetualPosition(
		ctx sdk.Context,
		liquidatedSubaccountId satypes.SubaccountId,
		perpetualId uint32,
		deltaQuantumsTotal *big.Int,
	) (
		fills []MatchPerpetualDeleveraging_Fill,
		deltaQuantumsRemaining *big.Int,
	)
	GetIndexerEventManager() indexer_manager.IndexerEventManager
	IsLiquidatable(
		ctx sdk.Context,
		subaccountId satypes.SubaccountId,
	) (
		bool,
		error,
	)
	ValidateSubaccountEquityTierLimitForNewOrder(
		ctx sdk.Context,
		order Order,
	) error
	Logger(
		ctx sdk.Context,
	) log.Logger
}
