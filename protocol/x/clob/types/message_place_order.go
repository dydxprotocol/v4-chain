package types

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
)

const TypeMsgPlaceOrder = "place_order"

var _ sdk.Msg = &MsgPlaceOrder{}

func NewMsgPlaceOrder(order Order) *MsgPlaceOrder {
	return &MsgPlaceOrder{
		Order: order,
	}
}

func (msg *MsgPlaceOrder) ValidateBasic() (err error) {
	defer func() {
		if err != nil {
			telemetry.IncrCounterWithLabels(
				[]string{ModuleName, metrics.PlaceOrder, metrics.ValidateBasic, metrics.Error, metrics.Count},
				1,
				msg.Order.GetOrderLabels(),
			)
		}
	}()

	// Check for deprecated fields.
	if msg.Order.TimeInForce == Order_TIME_IN_FORCE_FILL_OR_KILL {
		return errorsmod.Wrapf(ErrDeprecatedField, "Fill-or-kill has been deprecated")
	}

	err = msg.Order.OrderId.SubaccountId.Validate()
	if err != nil {
		return err
	}

	// Verify that enum type values are valid.
	if _, exists := Order_Side_name[int32(msg.Order.Side)]; !exists {
		return errorsmod.Wrapf(ErrInvalidOrderSide, "invalid order side (%s)", msg.Order.Side)
	}

	if _, exists := Order_TimeInForce_name[int32(msg.Order.TimeInForce)]; !exists {
		return errorsmod.Wrapf(ErrInvalidTimeInForce, "invalid time in force (%s)", msg.Order.TimeInForce)
	}

	if _, exists := Order_ConditionType_name[int32(msg.Order.ConditionType)]; !exists {
		return errorsmod.Wrapf(ErrInvalidConditionType, "invalid condition type (%s)", msg.Order.ConditionType)
	}

	if msg.Order.Side == Order_SIDE_UNSPECIFIED {
		return errorsmod.Wrapf(ErrInvalidOrderSide, "UNSPECIFIED is not a valid order side")
	}

	if msg.Order.Quantums == uint64(0) {
		return errorsmod.Wrapf(ErrInvalidOrderQuantums, "order size quantums cannot be 0")
	}

	orderId := msg.Order.GetOrderId()
	if orderId.IsShortTermOrder() {
		// This also implicitly verifies that GoodTilBlockTime is not set / is zero for short-term orders.
		if msg.Order.GetGoodTilBlock() == uint32(0) {
			return errorsmod.Wrapf(ErrInvalidOrderGoodTilBlock, "order goodTilBlock cannot be 0")
		}
	} else if orderId.IsStatefulOrder() {
		if msg.Order.GetGoodTilBlockTime() == uint32(0) {
			return errorsmod.Wrapf(
				ErrInvalidStatefulOrderGoodTilBlockTime,
				"stateful order goodTilBlockTime cannot be 0",
			)
		}
	} else {
		return errorsmod.Wrapf(ErrInvalidOrderFlag, "invalid order flag %v", orderId.OrderFlags)
	}

	if orderId.IsLongTermOrder() && msg.Order.RequiresImmediateExecution() {
		return ErrLongTermOrdersCannotRequireImmediateExecution
	}

	if msg.Order.ReduceOnly && !msg.Order.RequiresImmediateExecution() {
		return errorsmod.Wrapf(ErrReduceOnlyDisabled, "reduce only orders must be short term IOC orders")
	}

	if msg.Order.Subticks == uint64(0) && !msg.Order.IsTwapOrder() {
		return errorsmod.Wrapf(ErrInvalidOrderSubticks, "order subticks cannot be 0 for this order type")
	}

	if orderId.IsConditionalOrder() {
		if msg.Order.ConditionType == Order_CONDITION_TYPE_UNSPECIFIED {
			return errorsmod.Wrapf(ErrInvalidConditionType, "condition type cannot be unspecified")
		}

		if msg.Order.ConditionalOrderTriggerSubticks == uint64(0) {
			return errorsmod.Wrapf(ErrInvalidConditionalOrderTriggerSubticks, "conditional order trigger subticks cannot be 0")
		}
	} else {
		if msg.Order.ConditionType != Order_CONDITION_TYPE_UNSPECIFIED {
			return errorsmod.Wrapf(ErrInvalidConditionType, "condition type specified for non-conditional order")
		}

		if msg.Order.ConditionalOrderTriggerSubticks != uint64(0) {
			return errorsmod.Wrapf(
				ErrInvalidConditionalOrderTriggerSubticks,
				"conditional order trigger subticks greater than 0 for non-conditional order",
			)
		}
	}

	if msg.Order.IsTwapOrder() {
		if msg.Order.TwapConfig == nil {
			return errorsmod.Wrapf(
				ErrInvalidPlaceOrder,
				"TWAP order must have a TWAP config",
			)
		}
		if msg.Order.TwapConfig.Interval < 30 || msg.Order.TwapConfig.Interval > 3600 {
			return errorsmod.Wrapf(
				ErrInvalidPlaceOrder,
				"TWAP order interval must be between 30 seconds and 3600 seconds (1 hour)",
			)
		}
		if msg.Order.TwapConfig.Duration < 300 || msg.Order.TwapConfig.Duration > 86400 {
			return errorsmod.Wrapf(
				ErrInvalidPlaceOrder,
				"TWAP order duration must be between 300 seconds (5 minutes) and 86400 seconds (24 hours)",
			)
		}
		if msg.Order.TwapConfig.Duration%msg.Order.TwapConfig.Interval != 0 {
			return errorsmod.Wrapf(
				ErrInvalidPlaceOrder,
				"TWAP order duration must be a multiple of the interval",
			)
		}
		if msg.Order.TwapConfig.SlippagePercent > 5000 {
			return errorsmod.Wrapf(
				ErrInvalidPlaceOrder,
				"TWAP order slippage percent must be between 0 and 5000",
			)
		}
	}

	return nil
}
