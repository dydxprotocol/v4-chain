package types

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// MemClob is an interface that encapsulates all reads and writes to the
// CLOB's in-memory data structures.
type MemClob interface {
	SetClobKeeper(
		keeper MemClobKeeper,
	)
	CancelOrder(
		ctx sdk.Context,
		msgCancelOrder *MsgCancelOrder,
	) (offchainUpdates *OffchainUpdates, err error)
	CreateOrderbook(
		ctx sdk.Context,
		clobPair ClobPair,
	)
	CountSubaccountShortTermOrders(
		ctx sdk.Context,
		subaccountId satypes.SubaccountId,
	) uint32
	GetOperationsToReplay(
		ctx sdk.Context,
	) (
		[]InternalOperation,
		map[OrderHash][]byte,
	)

	GetOperationsRaw(
		ctx sdk.Context,
	) (
		operationsQueue []OperationRaw,
	)
	GetOrder(
		ctx sdk.Context,
		orderId OrderId,
	) (Order, bool)
	GetCancelOrder(
		ctx sdk.Context,
		orderId OrderId,
	) (uint32, bool)
	GetOrderFilledAmount(
		ctx sdk.Context,
		orderId OrderId,
	) satypes.BaseQuantums
	GetOrderRemainingAmount(
		ctx sdk.Context,
		order Order,
	) (
		remainingAmount satypes.BaseQuantums,
		hasRemainingAmount bool,
	)
	GetSubaccountOrders(
		ctx sdk.Context,
		clobPairId ClobPairId,
		subaccountId satypes.SubaccountId,
		side Order_Side,
	) ([]Order, error)
	PlaceOrder(
		ctx sdk.Context,
		order Order,
	) (satypes.BaseQuantums, OrderStatus, *OffchainUpdates, error)
	PlacePerpetualLiquidation(
		ctx sdk.Context,
		liquidationOrder LiquidationOrder,
	) (
		orderSizeOptimisticallyFilledFromMatchingQuantums satypes.BaseQuantums,
		orderStatus OrderStatus,
		offchainUpdates *OffchainUpdates,
		err error,
	)
	DeleverageSubaccount(
		ctx sdk.Context,
		subaccountId satypes.SubaccountId,
		perpetualId uint32,
		deltaQuantums *big.Int,
		isFinalSettlement bool,
	) (
		quantumsDeleveraged *big.Int,
		err error,
	)
	RemoveOrderIfFilled(
		ctx sdk.Context,
		orderId OrderId,
	)
	GetPricePremium(
		ctx sdk.Context,
		clobPair ClobPair,
		params perptypes.GetPricePremiumParams,
	) (
		premiumPpm int32,
		err error,
	)
	RemoveAndClearOperationsQueue(
		ctx sdk.Context,
		localValidatorOperationsQueue []InternalOperation,
	)
	PurgeInvalidMemclobState(
		ctx sdk.Context,
		fullyFilledOrderIds []OrderId,
		expiredStatefulOrderIds []OrderId,
		canceledStatefulOrderIds []OrderId,
		removedStatefulOrderIds []OrderId,
		existingOffchainUpdates *OffchainUpdates,
	) (offchainUpdates *OffchainUpdates)
	ReplayOperations(
		ctx sdk.Context,
		localOperations []InternalOperation,
		shortTermOrderTxBytes map[OrderHash][]byte,
		existingOffchainUpdates *OffchainUpdates,
	) (offchainUpdates *OffchainUpdates)
	SetMemclobGauges(
		ctx sdk.Context,
	)
	GetMidPrice(
		ctx sdk.Context,
		clobPairId ClobPairId,
	) (
		midPrice Subticks,
		bestBid Order,
		bestAsk Order,
		exists bool,
	)
	InsertZeroFillDeleveragingIntoOperationsQueue(
		ctx sdk.Context,
		subaccountId satypes.SubaccountId,
		perpetualId uint32,
	)
}
