package types

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"
	"time"

	proto "github.com/cosmos/gogoproto/proto"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	gometrics "github.com/hashicorp/go-metrics"
)

var _ MatchableOrder = &Order{}

func (o *Order) MustBeValidOrderSide() {
	// If `o.Side` is an invalid side, panic.
	if o.Side != Order_SIDE_BUY && o.Side != Order_SIDE_SELL {
		panic(ErrInvalidOrderSide)
	}
}

// MustBeConditionalOrder panics if the order is not a conditional order.
func (o *Order) MustBeConditionalOrder() {
	o.OrderId.MustBeConditionalOrder()
}

// IsBuy returns true if this is a buy order, false if not.
// This function is necessary for the `Order` type to implement the `MatchableOrder` interface.
func (o *Order) IsBuy() bool {
	return o.Side == Order_SIDE_BUY
}

// GetOrderHash returns the SHA256 hash of this order.
// This function is necessary for the `Order` type to implement the `MatchableOrder` interface.
func (o *Order) GetOrderHash() OrderHash {
	orderBytes, err := o.Marshal()
	if err != nil {
		panic(err)
	}
	return sha256.Sum256(orderBytes)
}

// GetOrderTextString returns the JSON representation of this order.
func (o *Order) GetOrderTextString() string {
	return proto.MarshalTextString(o)
}

// GetBaseQuantums returns the quantums of this order.
// This function is necessary for the `Order` type to implement the `MatchableOrder` interface.
func (o *Order) GetBaseQuantums() satypes.BaseQuantums {
	return satypes.BaseQuantums(o.Quantums)
}

// GetBigQuantums returns the quantums of this order. The returned quantums is positive
// if long, negative if short, and zero if `o.Quantums == 0`.
// This function is necessary for the `Order` type to implement the `MatchableOrder` interface.
func (o *Order) GetBigQuantums() *big.Int {
	bigQuantums := new(big.Int).SetUint64(o.Quantums)
	if !o.IsBuy() {
		bigQuantums.Neg(bigQuantums)
	}
	return bigQuantums
}

// GetOrderSubticks returns the subticks of this order.
// This function is necessary for the `Order` type to implement the `MatchableOrder` interface.
func (o *Order) GetOrderSubticks() Subticks {
	return Subticks(o.Subticks)
}

// MustCmpReplacementOrder compares x to y and returns:
// 1 if x > y
// 0 if x = y
// -1 if x < y
// The orders are compared primarily by `GoodTilBlock` for Short-Term orders and `GoodTilBlockTime`
// for stateful orders. If the order expirations are equal, then they are compared by their SHA256 hash.
// Note that this function panics if the order IDs are not equal.
func (x *Order) MustCmpReplacementOrder(y *Order) int {
	if x.OrderId != y.OrderId {
		panic(
			fmt.Sprintf(
				"MustCmpReplacementOrder: order ID (%v) does not equal order ID (%v)",
				x.OrderId,
				y.OrderId,
			),
		)
	}

	var orderXExpiration uint32
	var orderYExpiration uint32

	// If this is a Short-Term order, use the `GoodTilBlock` for comparison.
	// Else this is a stateful order, therefore use `GoodTilBlockTime` for comparison.
	if x.IsShortTermOrder() {
		orderXExpiration = x.GetGoodTilBlock()
		orderYExpiration = y.GetGoodTilBlock()
	} else {
		orderXExpiration = x.GetGoodTilBlockTime()
		orderYExpiration = y.GetGoodTilBlockTime()
	}

	if orderXExpiration > orderYExpiration {
		return 1
	} else if orderXExpiration < orderYExpiration {
		return -1
	}

	// If both orders have the same expiration, use the SHA256 hash for comparison.
	xHash := x.GetOrderHash()
	yHash := y.GetOrderHash()
	return bytes.Compare(xHash[:], yHash[:])
}

// GetSubaccountId returns the subaccount ID that placed this order.
// This function is necessary for the `Order` type to implement the `MatchableOrder` interface.
func (o *Order) GetSubaccountId() satypes.SubaccountId {
	return o.OrderId.SubaccountId
}

// IsLiquidation always returns false since this order is not a liquidation.
// This function is necessary for the `Order` type to implement the `MatchableOrder` interface.
func (o *Order) IsLiquidation() bool {
	return false
}

// MustGetOrder returns the underlying `Order` type.
// This function is necessary for the `Order` type to implement the `MatchableOrder` interface.
func (o *Order) MustGetOrder() Order {
	return *o
}

// MustGetLiquidationOrder always panics since Order is not a Liquidation Order.
// This function is necessary for the `Order` type to implement the `MatchableOrder` interface.
func (o *Order) MustGetLiquidationOrder() LiquidationOrder {
	panic("MustGetLiquidationOrder: Order is not a liquidation order")
}

// MustGetLiquidatedPerpetualId always panics since there is no underlying perpetual ID for a `Order`.
// This function is necessary for the `Order` type to implement the `MatchableOrder` interface.
func (o *Order) MustGetLiquidatedPerpetualId() uint32 {
	panic("MustGetLiquidatedPerpetualId: No liquidated perpetual on an Order type.")
}

// IsReduceOnly returns whether this is a reduce-only order.
func (o *Order) IsReduceOnly() bool {
	return o.ReduceOnly
}

// IsTakeProfitOrder returns whether this is order is a conditional take profit order.
func (o *Order) IsTakeProfitOrder() bool {
	return o.IsConditionalOrder() && o.ConditionType == Order_CONDITION_TYPE_TAKE_PROFIT
}

// IsStopLossOrder returns whether this is order is a conditional stop loss order.
func (o *Order) IsStopLossOrder() bool {
	return o.IsConditionalOrder() && o.ConditionType == Order_CONDITION_TYPE_STOP_LOSS
}

// RequiresImmediateExecution returns whether this order has to be executed immediately.
func (o *Order) RequiresImmediateExecution() bool {
	return o.GetTimeInForce() == Order_TIME_IN_FORCE_IOC || o.GetTimeInForce() == Order_TIME_IN_FORCE_FILL_OR_KILL
}

// IsShortTermOrder returns whether this is a Short-Term order.
func (o *Order) IsShortTermOrder() bool {
	return o.OrderId.IsShortTermOrder()
}

// IsStatefulOrder returns whether this order is a stateful order, which is true for Long-Term
// and conditional orders and false for Short-Term orders.
func (o *Order) IsStatefulOrder() bool {
	return o.OrderId.IsStatefulOrder()
}

// IsConditionalOrder returns whether this order is a conditional order.
func (o *Order) IsConditionalOrder() bool {
	return o.OrderId.IsConditionalOrder()
}

// IsTwapOrder returns whether this order is a TWAP order.
func (o *Order) IsTwapOrder() bool {
	return o.OrderId.IsTwapOrder()
}

// IsTwapSuborder returns whether this order is a TWAP suborder.
func (o *Order) IsTwapSuborder() bool {
	return o.OrderId.IsTwapSuborder()
}

// IsCollateralCheckRequired returns whether this order needs
// to pass collateral checks. This is true for all non-internal
// orders and for generated TWAP suborders.
func (o *Order) IsCollateralCheckRequired(isInternalOrder bool) bool {
	return (!isInternalOrder && !o.IsConditionalOrder()) || (isInternalOrder && o.IsTwapSuborder())
}

// IsPostOnlyOrder returns whether this order is a post only order.
func (o *Order) IsPostOnlyOrder() bool {
	return o.GetTimeInForce() == Order_TIME_IN_FORCE_POST_ONLY
}

// CanTrigger returns if a condition order is eligible to be triggered based on a given
// subticks value. Function will panic if order is not a conditional order.
func (o *Order) CanTrigger(subticks Subticks) bool {
	o.MustBeConditionalOrder()
	orderTriggerSubticks := Subticks(o.ConditionalOrderTriggerSubticks)

	// Take profit buys and stop loss sells trigger when the oracle price goes lower
	// than or equal to the trigger price.
	if o.ConditionType == Order_CONDITION_TYPE_TAKE_PROFIT && o.IsBuy() ||
		o.ConditionType == Order_CONDITION_TYPE_STOP_LOSS && !o.IsBuy() {
		return orderTriggerSubticks >= subticks
	}
	// Take profit sells and stop loss buys trigger when the oracle price goes higher
	// than or equal to the trigger price.
	return orderTriggerSubticks <= subticks
}

// MustGetUnixGoodTilBlockTime returns an instance of `Time` that represents the order's
// `GoodTilBlockTime`. This function panics when the order is a short-term order or
// when its `GoodTilBlockTime` is zero.
func (o *Order) MustGetUnixGoodTilBlockTime() time.Time {
	o.MustBeStatefulOrder()
	goodTilBlockTime := o.GetGoodTilBlockTime()
	if goodTilBlockTime == 0 {
		panic(
			fmt.Errorf(
				"MustGetUnixGoodTilBlockTime: order (%v) goodTilBlockTime is zero",
				o,
			),
		)
	}
	return time.Unix(int64(goodTilBlockTime), 0)
}

// MustBeStatefulOrder panics if the order is not a stateful order, else it does nothing.
func (o *Order) MustBeStatefulOrder() {
	o.OrderId.MustBeStatefulOrder()
}

// GetClobPairId returns the CLOB pair ID of this order.
// This function implements the `MatchableOrder` interface.
func (o *Order) GetClobPairId() ClobPairId {
	return ClobPairId(o.OrderId.GetClobPairId())
}

// GetOrderLabels returns the telemetry labels of this order.
func (o *Order) GetOrderLabels() []gometrics.Label {
	return append(
		[]gometrics.Label{
			metrics.GetLabelForStringValue(metrics.TimeInForce, o.GetTimeInForce().String()),
			metrics.GetLabelForBoolValue(metrics.ReduceOnly, o.IsReduceOnly()),
			metrics.GetLabelForStringValue(metrics.OrderSide, o.GetSide().String()),
		},
		o.OrderId.GetOrderIdLabels()...,
	)
}

func (o *Order) GetTotalLegsTWAPOrder() uint32 {
	if o.IsTwapOrder() {
		return o.TwapParameters.Duration / o.TwapParameters.Interval
	}
	return 0
}

// GetTWAPTriggerKey returns the key for a TWAP trigger order.
func GetTWAPTriggerKey(triggerTime int64, orderId OrderId) []byte {
	// Write trigger time as big-endian uint64
	timeBytes := TriggerTimeToBytes(triggerTime)

	// Get marshaled orderId
	orderIdBytes := orderId.ToStateKey()

	// Combine time and orderId bytes
	return append(timeBytes, orderIdBytes...)
}

func TriggerTimeToBytes(triggerTime int64) []byte {
	timeBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(timeBytes, uint64(triggerTime))
	return timeBytes
}

func TimeFromTriggerKey(triggerKey []byte) int64 {
	return int64(binary.BigEndian.Uint64(triggerKey[0:8]))
}

// IsCompleted returns true if the TWAP order has no remaining legs to execute.
func (t *TwapOrderPlacement) IsCompleted() bool {
	return t.RemainingLegs == 0
}
