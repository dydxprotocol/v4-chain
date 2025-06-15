package types

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
)

const TypeMsgPlaceOrder = "place_order"

const (
	MaxBuilderCodeFeePpm uint32 = 10_000
)

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

	if msg.Order.Subticks == uint64(0) {
		return errorsmod.Wrapf(ErrInvalidOrderSubticks, "order subticks cannot be 0")
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

	if msg.Order.BuilderCodeParameters != nil {
		if _, err := sdk.AccAddressFromBech32(msg.Order.BuilderCodeParameters.BuilderAddress); err != nil {
			return errorsmod.Wrapf(
				ErrInvalidBuilderCode,
				"builder code address '%s' must be a valid bech32 address, but got error '%v'",
				msg.Order.BuilderCodeParameters.BuilderAddress,
				err.Error(),
			)
		}
		if msg.Order.BuilderCodeParameters.FeePpm <= 0 || msg.Order.BuilderCodeParameters.FeePpm > MaxBuilderCodeFeePpm {
			return errorsmod.Wrapf(
				ErrInvalidBuilderCode,
				"builder code fee ppm '%d' must be in the range (0, %d]",
				msg.Order.BuilderCodeParameters.FeePpm,
				MaxBuilderCodeFeePpm,
			)
		}
	}

	return nil
}
