package types

import (
	"math/big"

	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/zyedidia/generic/list"
)

// Subticks is used to represent the price level that an order will be placed on the orderbook.
type Subticks uint64

func (st Subticks) ToBigInt() *big.Int {
	return new(big.Int).SetUint64(st.ToUint64())
}

func (st Subticks) ToBigRat() *big.Rat {
	return new(big.Rat).SetUint64(st.ToUint64())
}

func (st Subticks) ToUint64() uint64 {
	return uint64(st)
}

// SubticksPerTick is used to represent how many subticks are in one tick.
// That is, the subticks of any valid order must be a multiple of this value.
type SubticksPerTick uint32

type ClobPairId uint32

func (cp ClobPairId) ToUint32() uint32 {
	return uint32(cp)
}

// ClobOrder represents an order that is resting on the CLOB.
type ClobOrder struct {
	// The order that is resting on the CLOB.
	Order Order
	// The signature on the transaction containing the `MsgPlaceOrder` message,
	// from the user who placed this order [PENDING ABCI++].
	Signature []byte
}

// LevelOrder represents the queue position of an order that is within a
// specific price level of the CLOB.
type LevelOrder = list.Node[ClobOrder]

// Level represents a price level on the CLOB.
type Level struct {
	// LevelOrders represents a doubly-linked list of `ClobOrder`s sorted in chronical
	// order (ascending). Note that this should always be non-`nil`, since the
	// `Level` should not exist if there are no elements in the linked list.
	LevelOrders list.List[ClobOrder]
}

// PendingOpenOrder is a utility struct used for representing an order a subaccount will open. This is
// used for collateralization checks, to specifically verify that the number of quantums in this order
// can be opened for this subaccount.
// Only used for representing maker orders in add-to-orderbook collat check.
// TODO(CLOB-849) Remove this struct.
type PendingOpenOrder struct {
	// The amount of base quantums that is remaining for this order.
	RemainingQuantums satypes.BaseQuantums
	// True if this is a buy order, false if it's a sell order.
	IsBuy bool
	// The price that this order would be matched at.
	Subticks Subticks
	// The ID of the CLOB this order would be placed on.
	ClobPairId ClobPairId
	// The builder code parameters for this order.
	BuilderCodeParameters *BuilderCodeParameters
}

// AddOrderToOrderbookCollateralizationCheckFn defines a function interface that can be used for verifying
// one or more subaccounts are properly collateralized if their respective order(s) are matched. Returns the result of
// the collateralization check.
type AddOrderToOrderbookCollateralizationCheckFn func(
	subaccountMatchedOrders map[satypes.SubaccountId][]PendingOpenOrder,
) (
	success bool,
	successPerSubaccountUpdate map[satypes.SubaccountId]satypes.UpdateResult,
)

// GetStatePositionFn defines a function interface that can be used for getting the position size
// of an order in state. It is used for determining whether reduce-only orders need to be resized
// or canceled.
type GetStatePositionFn func(
	subaccountId satypes.SubaccountId,
	clobPairId ClobPairId,
) (
	positionSizeQuantums *big.Int,
)

// TakerOrderStatus is a utility struct used for representing the status, remaining size, and optimistically filled
// size of a taker order after attempting to match it on the orderbook.
type TakerOrderStatus struct {
	// The state of the taker order after attempting to match it against the orderbook.
	OrderStatus OrderStatus
	// The amount of remaining (non-matched) base quantums of this taker order.
	RemainingQuantums satypes.BaseQuantums
	// The amount of base quantums that were optimistically filled (from this current matching cycle) of this taker
	// order. Note that if any quantums of this order were optimistically filled or filled in state before the current
	// matching cycle, this value will not include them.
	OrderOptimisticallyFilledQuantums satypes.BaseQuantums
}

// ToStreamingTakerOrderStatus converts the TakerOrderStatus to a StreamTakerOrderStatus
// to be emitted by full node streaming.
func (tos *TakerOrderStatus) ToStreamingTakerOrderStatus() *StreamTakerOrderStatus {
	return &StreamTakerOrderStatus{
		OrderStatus:                  uint32(tos.OrderStatus),
		RemainingQuantums:            tos.RemainingQuantums.ToUint64(),
		OptimisticallyFilledQuantums: tos.OrderOptimisticallyFilledQuantums.ToUint64(),
	}
}

// OrderStatus represents the status of an order after attempting to place it on the orderbook.
type OrderStatus uint

const (
	// Success indicates the order was successfully matched and / or added to the orderbook.
	Success OrderStatus = iota
	// Undercollateralized indicates the order failed collateralization checks while matching or
	// when placed on the orderbook, and was therefore canceled.
	Undercollateralized
	// InternalError indicates the order caused an internal error during collateralization checks
	// while matching or when placed on the orderbook, and was therefore canceled.
	InternalError
	// ImmediateOrCancelWouldRestOnBook indicates this is an IOC order that would have been placed
	// on the orderbook as resting liquidity, and was therefore canceled.
	ImmediateOrCancelWouldRestOnBook
	// ReduceOnlyResized indicates the reduce-only order was resized since it would have changed
	// the user's position side.
	ReduceOnlyResized
	// LiquidationRequiresDeleveraging indicates that there wasn't enough liquidity to liquidate
	// the subaccount profitably on the orderbook and the order was subsequently not fully matched
	// because the insurance fund did not have enough funds to cover the losses from performing
	// the liquidation.
	LiquidationRequiresDeleveraging
	// LiquidationExceededSubaccountMaxNotionalLiquidated indicates that the liquidation order
	// could not be matched because it exceeded the max notional liquidated in this block.
	LiquidationExceededSubaccountMaxNotionalLiquidated
	// LiquidationExceededSubaccountMaxInsuranceLost indicates that the liquidation order could not
	// be matched because it exceeded the maximum funds lost for the insurance fund in this block.
	LiquidationExceededSubaccountMaxInsuranceLost
	// ViolatesIsolatedSubaccountConstraints indicates that matching the order would lead to the
	// subaccount violating constraints for isolated perpetuals, where the subaccount would end up
	// with either multiple positions in isolated perpetuals or both an isolated and a cross perpetual
	// position.
	ViolatesIsolatedSubaccountConstraints
	// PostOnlyWouldCrossMakerOrder indicates that matching the post only taker order would cross the
	// orderbook, and was therefore canceled.
	PostOnlyWouldCrossMakerOrder
)

// String returns a string representation of this `OrderStatus` enum.
func (os OrderStatus) String() string {
	switch os {
	case Success:
		return "Success"
	case Undercollateralized:
		return "Undercollateralized"
	case InternalError:
		return "InternalError"
	case ImmediateOrCancelWouldRestOnBook:
		return "ImmediateOrCancelWouldRestOnBook"
	case ReduceOnlyResized:
		return "ReduceOnlyResized"
	case LiquidationRequiresDeleveraging:
		return "LiquidationRequiresDeleveraging"
	case LiquidationExceededSubaccountMaxNotionalLiquidated:
		return "LiquidationExceededSubaccountMaxNotionalLiquidated"
	case LiquidationExceededSubaccountMaxInsuranceLost:
		return "LiquidationExceededSubaccountMaxInsuranceLost"
	case ViolatesIsolatedSubaccountConstraints:
		return "ViolatesIsolatedSubaccountConstraints"
	default:
		return "Unknown"
	}
}

// IsSuccess returns `true` if this `OrderStatus` enum is `Success`, else returns `false`.
func (os OrderStatus) IsSuccess() bool {
	return os == Success
}

// FillType represents the type of the fill.
type FillType uint

const (
	Trade FillType = iota
	PerpetualLiquidate
	PerpetualDeleverage
)

// MatchableOrder is an interface that a matchable order must conform to. This interface is used
// to generalize matching between standard orders and liquidations.
type MatchableOrder interface {
	// GetSubaccountID returns the `SubaccountId` of the subaccount that placed the order.
	// In the case of a `LiquidationOrder`, it refers to the subaccount that is being liquidated.
	GetSubaccountId() satypes.SubaccountId
	// GetClobPairId returns the CLOB pair ID that this order should be matched against.
	GetClobPairId() ClobPairId
	// IsBuy returns true if this is a buy order, false if it's a sell order.
	IsBuy() bool
	// IsLiquidation returns true if this is a liquidation order, false if not.
	IsLiquidation() bool
	// MustGetOrder returns the underlying order if this is not a liquidation order. Panics if called
	// for a liquidation order.
	MustGetOrder() Order
	// MustGetLiquidationOrder returns the underlying liquidation order if this is not a regular order.
	// Panics if called for a regular order.
	MustGetLiquidationOrder() LiquidationOrder
	// MustGetLiquidatedPerpetualId returns the perpetual ID if this is a liquidation order. Panics
	// if called for a non-liquidation order.
	MustGetLiquidatedPerpetualId() uint32
	// GetBaseQuantums returns the base quantums of this order.
	GetBaseQuantums() satypes.BaseQuantums
	// GetOrderSubticks returns the subticks of this order.
	GetOrderSubticks() Subticks
	// GetOrderHash returns the hash of this order.
	// If this is a liquidation it returns the hash of the `PerpetualLiquidationInfo`.
	// Else, it returns the hash of the `Order` proto.
	GetOrderHash() OrderHash
	// IsReduceOnly returns whether this is a reduce-only order.
	// This always returns false for liquidation orders.
	IsReduceOnly() bool
	// GetOrderRouterAddress returns the order router address for this order.
	GetOrderRouterAddress() string
}
