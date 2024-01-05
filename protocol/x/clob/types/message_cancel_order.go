package types

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
)

const TypeMsgCancelOrder = "cancel_order"

var _ sdk.Msg = &MsgCancelOrder{}

// NewMsgCancelOrderShortTerm constructs a MsgCancelOrder from an `OrderId` and a `GoodTilBlock`.
// `OrderId` must be for a Short-Term order.
func NewMsgCancelOrderShortTerm(orderId OrderId, goodTilBlock uint32) *MsgCancelOrder {
	orderId.MustBeShortTermOrder()
	return &MsgCancelOrder{
		OrderId:      orderId,
		GoodTilOneof: &MsgCancelOrder_GoodTilBlock{GoodTilBlock: goodTilBlock},
	}
}

// NewMsgCancelOrderStateful constructs a MsgCancelOrder from an `OrderId` and a `GoodTillBlockTime`.
// `OrderId` must be for a Stateful Order. Long term and conditional orderIds are acceptable.
func NewMsgCancelOrderStateful(orderId OrderId, goodTilBlockTime uint32) *MsgCancelOrder {
	orderId.MustBeStatefulOrder()
	return &MsgCancelOrder{
		OrderId:      orderId,
		GoodTilOneof: &MsgCancelOrder_GoodTilBlockTime{GoodTilBlockTime: goodTilBlockTime},
	}
}

func (msg *MsgCancelOrder) ValidateBasic() (err error) {
	orderId := msg.GetOrderId()

	defer func() {
		if err != nil {
			telemetry.IncrCounterWithLabels(
				[]string{ModuleName, metrics.CancelOrder, metrics.ValidateBasic, metrics.Error, metrics.Count},
				1,
				msg.OrderId.GetOrderIdLabels(),
			)
		}
	}()

	if err := orderId.Validate(); err != nil {
		return err
	}

	if orderId.IsStatefulOrder() {
		if msg.GetGoodTilBlockTime() == 0 {
			return errorsmod.Wrapf(
				ErrInvalidStatefulOrderGoodTilBlockTime,
				"stateful cancellation goodTilBlockTime cannot be 0, %+v",
				orderId,
			)
		}
	} else {
		if msg.GetGoodTilBlock() == 0 {
			return errorsmod.Wrapf(
				ErrInvalidOrderGoodTilBlock,
				"cancellation goodTilBlock cannot be 0, orderId %+v",
				orderId,
			)
		}
	}
	return nil
}

// MustGetUnixGoodTilBlockTime returns an instance of `Time` that represents the cancel's
// `GoodTilBlockTime`. This function panics when the order is a short-term order or
// when its `GoodTilBlockTime` is zero.
func (msg *MsgCancelOrder) MustGetUnixGoodTilBlockTime() time.Time {
	orderId := msg.GetOrderId()
	orderId.MustBeStatefulOrder()
	goodTilBlockTime := msg.GetGoodTilBlockTime()
	if goodTilBlockTime == 0 {
		panic(
			fmt.Errorf(
				"MustGetUnixGoodTilBlockTime: cancel (%+v) goodTilBlockTime is zero",
				msg,
			),
		)
	}
	return time.Unix(int64(goodTilBlockTime), 0)
}
