package ante

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// ClobDecorator is an AnteDecorator which is responsible for:
//   - adding short term order placements and cancelations to the in-memory orderbook (`CheckTx` only).
//   - adding stateful order placements and cancelations to state (`CheckTx` and `RecheckTx` only).
//
// This AnteDecorator also enforces that any Transaction which contains a `MsgPlaceOrder`,
// or a `MsgCancelOrder`, must consist only of a single message.
//
// This AnteDecorator is a no-op if:
//   - No messages in the transaction are `MsgPlaceOrder` or `MsgCancelOrder`.
//   - This AnteDecorator is called during `DeliverTx`.
//
// This AnteDecorator returns an error if:
//   - The transaction contains multiple messages, and one of them is a `MsgPlaceOrder`
//     or `MsgCancelOrder` message.
//   - The underlying `PlaceStatefulOrder`, `PlaceShortTermOrder`, `CancelStatefulOrder`, or `CancelShortTermOrder`
//     methods on the keeper return errors.
type ClobDecorator struct {
	clobKeeper types.ClobKeeper
}

func NewClobDecorator(clobKeeper types.ClobKeeper) ClobDecorator {
	return ClobDecorator{
		clobKeeper,
	}
}

func (cd ClobDecorator) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (sdk.Context, error) {
	// No need to process during `DeliverTx` or simulation, call next `AnteHandler`.
	if lib.IsDeliverTxMode(ctx) || simulate {
		return next(ctx, tx, simulate)
	}

	return ctx, sdkerrors.ErrInvalidRequest
}

// IsSingleClobMsgTx returns `true` if the supplied `tx` consist of a single clob message
// (`MsgPlaceOrder` or `MsgCancelOrder`). If `msgs` consist of multiple clob messages,
// or a mix of on-chain and clob messages, an error is returned.
func IsSingleClobMsgTx(ctx sdk.Context, tx sdk.Tx) (bool, error) {
	msgs := tx.GetMsgs()
	var hasMessage = false

	for _, msg := range msgs {
		switch msg.(type) {
		case *types.MsgCancelOrder, *types.MsgPlaceOrder:
			hasMessage = true
		}

		if hasMessage {
			break
		}
	}

	if !hasMessage {
		return false, nil
	}

	numMsgs := len(msgs)
	if numMsgs > 1 {
		return false, errorsmod.Wrap(
			sdkerrors.ErrInvalidRequest,
			"a transaction containing MsgCancelOrder or MsgPlaceOrder may not contain more than one message",
		)
	}

	return true, nil
}

// IsShortTermClobMsgTx returns `true` if the supplied `tx` consist of a single clob message
// (`MsgPlaceOrder` or `MsgCancelOrder`) which references a Short-Term Order. If `msgs` consist of multiple
// clob messages, or a mix of on-chain and clob messages, an error is returned.
func IsShortTermClobMsgTx(ctx sdk.Context, tx sdk.Tx) (bool, error) {
	msgs := tx.GetMsgs()

	var isShortTermOrder = false

	for _, msg := range msgs {
		switch msg := msg.(type) {
		case *types.MsgCancelOrder:
			{
				if msg.OrderId.IsShortTermOrder() {
					isShortTermOrder = true
				}
			}
		case *types.MsgPlaceOrder:
			{
				if msg.Order.OrderId.IsShortTermOrder() {
					isShortTermOrder = true
				}
			}
		}

		if isShortTermOrder {
			break
		}
	}

	if !isShortTermOrder {
		return false, nil
	}

	numMsgs := len(msgs)
	if numMsgs > 1 {
		return false, errorsmod.Wrap(
			sdkerrors.ErrInvalidRequest,
			"a transaction containing MsgCancelOrder or MsgPlaceOrder may not contain more than one message",
		)
	}

	return true, nil
}
