package types

import (
	"math/big"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4/indexer/indexer_manager"
	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"
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
		matchWithOrders MatchWithOrders,
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
	CheckTxCancelOrder(
		ctx sdk.Context,
		msgCancelOrder *MsgCancelOrder,
	) error
	DoesStatefulOrderExistInState(
		ctx sdk.Context,
		order Order,
	) bool
	GetStatefulOrderPlacement(
		ctx sdk.Context,
		orderId OrderId,
	) (val StatefulOrderPlacement, found bool)
	SetStatefulOrderPlacement(
		ctx sdk.Context,
		order Order,
		blockHeight uint32,
	)
	MustAddOrderToStatefulOrdersTimeSlice(
		ctx sdk.Context,
		goodTilBlockTime time.Time,
		orderId OrderId,
	)
	GetIndexerEventManager() indexer_manager.IndexerEventManager
}
