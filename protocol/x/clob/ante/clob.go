package ante

import (
	"fmt"

	"github.com/cometbft/cometbft/crypto/tmhash"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/dydxprotocol/v4/x/clob/types"
)

// ClobDecorator is an AnteDecorator which is responsible for adding order placements
// and cancelations to the in-memory orderbook during `CheckTx`.
// This AnteDecorator also enforces that any Transaction which contains a `MsgPlaceOrder`,
// or a `MsgCancelOrder`, must consist only of a single message.
//
// This AnteDecorator is a no-op if:
// - No messages in the transaction are `MsgPlaceOrder` or `MsgCancelOrder`.
// - This AnteDecorator is called during `DeliverTx` or `ReCheckTx`.
//
// This AnteDecorator returns an error if:
//   - The transaction contains multiple messages, and one of them is a `MsgPlaceOrder`
//     or `MsgCancelOrder` message.
//   - The underlying `PlaceOrder` or `CancelOrder` methods on the keeper return errors.
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
	// No need to process on `ReCheckTx`, or `DeliverTx`, call next `AnteHandler`.
	if !ctx.IsCheckTx() || simulate || ctx.IsReCheckTx() {
		return next(ctx, tx, simulate)
	}

	isClobOffChainMessage, err := IsOffChainSingleClobMsgTx(ctx, tx)
	if err != nil {
		return ctx, err
	}

	if !isClobOffChainMessage {
		return next(ctx, tx, simulate)
	}

	msgs := tx.GetMsgs()
	var msg = msgs[0]

	switch msg := msg.(type) {
	case *types.MsgCancelOrder:
		// Note that `msg.ValidateBasic` is called before the AnteHandlers.
		// This guarantees that `MsgCancelOrder` has undergone stateless validation.
		err := cd.clobKeeper.CheckTxCancelOrder(ctx, msg)
		txBytes := ctx.TxBytes()
		ctx.Logger().Info("Received new order cancelation",
			"tx",
			fmt.Sprintf("%X", tmhash.Sum(txBytes)),
			"msg",
			msg,
			"err",
			err,
			"block",
			ctx.BlockHeight(),
		)

		if err != nil {
			return ctx, err
		}
	case *types.MsgPlaceOrder:
		// Note that `msg.ValidateBasic` is called before all AnteHandlers.
		// This guarantees that `MsgPlaceOrder` has undergone stateless validation.
		orderSizeOptimisticallyFilledFromMatchingQuantums, status, err := cd.clobKeeper.CheckTxPlaceOrder(ctx, msg)
		txBytes := ctx.TxBytes()
		ctx.Logger().Info("Received new order",
			"tx",
			fmt.Sprintf("%X", tmhash.Sum(txBytes)),
			"orderHash",
			fmt.Sprintf("%X", msg.Order.GetOrderHash()),
			"msg",
			msg,
			"status",
			status,
			"orderSizeOptimisticallyFilledFromMatchingQuantums",
			orderSizeOptimisticallyFilledFromMatchingQuantums,
			"err",
			err,
			"block",
			ctx.BlockHeight(),
		)

		if err != nil {
			return ctx, err
		}
	}

	return next(ctx, tx, simulate)
}

// IsOffChainSingleClobMsgTx returns `true` if the supplied `tx` consist of a single off-chain clob message
// (such as a `MsgPlaceOrder` or `MsgCancelOrder`). If `msgs` consist of multiple off-chain clob messages,
// or a mix of on-chain and off-chain clob messages, an error is returned.
func IsOffChainSingleClobMsgTx(ctx sdk.Context, tx sdk.Tx) (bool, error) {
	msgs := tx.GetMsgs()
	var hasOffChainMessage = false

	for _, msg := range msgs {
		switch msg.(type) {
		case *types.MsgCancelOrder, *types.MsgPlaceOrder:
			hasOffChainMessage = true
		}

		if hasOffChainMessage {
			break
		}
	}

	if !hasOffChainMessage {
		return false, nil
	}

	numMsgs := len(msgs)
	if numMsgs > 1 {
		return false, sdkerrors.Wrap(
			sdkerrors.ErrInvalidRequest,
			"a transaction containing MsgCancelOrder or MsgPlaceOrder may not contain more than one message",
		)
	}

	return true, nil
}
