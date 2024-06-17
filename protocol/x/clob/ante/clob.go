package ante

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/cometbft/cometbft/crypto/tmhash"
	cometbftlog "github.com/cometbft/cometbft/libs/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/log"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
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

	// Ensure that if this is a clob message then that there is only one.
	// If it isn't a clob message then pass to the next AnteHandler.
	isSingleClobMsgTx, err := IsSingleClobMsgTx(tx)
	if err != nil {
		return ctx, err
	}

	if !isSingleClobMsgTx {
		return next(ctx, tx, simulate)
	}

	// Disable order placement and cancelation processing if the clob keeper is not initialized.
	if !cd.clobKeeper.IsInitialized() {
		return ctx, errorsmod.Wrap(
			types.ErrClobNotInitialized,
			"clob keeper is not initialized. Please wait for the next block.",
		)
	}

	msgs := tx.GetMsgs()
	var msg = msgs[0]

	switch msg := msg.(type) {
	case *types.MsgCancelOrder:
		if msg.OrderId.IsStatefulOrder() {
			err = cd.clobKeeper.CancelStatefulOrder(ctx, msg)
		} else {
			// No need to process short term order cancelations on `ReCheckTx`.
			if ctx.IsReCheckTx() {
				return next(ctx, tx, simulate)
			}

			// Note that `msg.ValidateBasic` is called before the AnteHandlers.
			// This guarantees that `MsgCancelOrder` has undergone stateless validation.
			err = cd.clobKeeper.CancelShortTermOrder(ctx, msg)
		}

		log.DebugLog(ctx, "Received new order cancellation",
			log.Tx, cometbftlog.NewLazySprintf("%X", tmhash.Sum(ctx.TxBytes())),
			log.Error, err,
		)

	case *types.MsgPlaceOrder:
		if msg.Order.OrderId.IsStatefulOrder() {
			err = cd.clobKeeper.PlaceStatefulOrder(ctx, msg, false)

			log.DebugLog(ctx, "Received new stateful order",
				log.Tx, cometbftlog.NewLazySprintf("%X", tmhash.Sum(ctx.TxBytes())),
				log.OrderHash, cometbftlog.NewLazySprintf("%X", msg.Order.GetOrderHash()),
				log.Error, err,
			)
		} else {
			// No need to process short term orders on `ReCheckTx`.
			if ctx.IsReCheckTx() {
				return next(ctx, tx, simulate)
			}

			var orderSizeOptimisticallyFilledFromMatchingQuantums satypes.BaseQuantums
			var status types.OrderStatus
			// Note that `msg.ValidateBasic` is called before all AnteHandlers.
			// This guarantees that `MsgPlaceOrder` has undergone stateless validation.
			orderSizeOptimisticallyFilledFromMatchingQuantums, status, err = cd.clobKeeper.PlaceShortTermOrder(
				ctx,
				msg,
			)

			log.DebugLog(ctx, "Received new short term order",
				log.Tx, cometbftlog.NewLazySprintf("%X", tmhash.Sum(ctx.TxBytes())),
				log.OrderHash, cometbftlog.NewLazySprintf("%X", msg.Order.GetOrderHash()),
				log.OrderStatus, status,
				log.OrderSizeOptimisticallyFilledFromMatchingQuantums, orderSizeOptimisticallyFilledFromMatchingQuantums,
				log.Error, err,
			)
		}
	case *types.MsgBatchCancel:
		// MsgBatchCancel currently only processes short-term cancels right now.
		// No need to process short term orders on `ReCheckTx`.
		if ctx.IsReCheckTx() {
			return next(ctx, tx, simulate)
		}

		success, failures, err := cd.clobKeeper.BatchCancelShortTermOrder(
			ctx,
			msg,
		)
		// If there are no successful cancellations and no validation errors,
		// return an error indicating no cancels have succeeded.
		if len(success) == 0 && err == nil {
			err = errorsmod.Wrapf(
				types.ErrBatchCancelFailed,
				"No successful cancellations. Failures: %+v",
				failures,
			)
		}

		log.DebugLog(
			ctx,
			"Received new batch cancellation",
			log.Tx, cometbftlog.NewLazySprintf("%X", tmhash.Sum(ctx.TxBytes())),
			log.Error, err,
		)
	}
	if err != nil {
		return ctx, err
	}

	return next(ctx, tx, simulate)
}

// IsSingleClobMsgTx returns `true` if the supplied `tx` consist of a single clob message
// (`MsgPlaceOrder` or `MsgCancelOrder` or `MsgBatchCancel`). If `msgs` consist of multiple
// clob messages, or a mix of on-chain and clob messages, an error is returned.
func IsSingleClobMsgTx(tx sdk.Tx) (bool, error) {
	msgs := tx.GetMsgs()
	var hasMessage = false

	for _, msg := range msgs {
		switch msg.(type) {
		case *types.MsgCancelOrder, *types.MsgPlaceOrder, *types.MsgBatchCancel, *types.MsgXOperate:
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
			"a transaction containing MsgCancelOrder or MsgPlaceOrder or MsgBatchCancel may not contain more than one message",
		)
	}

	return true, nil
}
