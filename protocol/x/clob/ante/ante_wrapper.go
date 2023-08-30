package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SingleMsgClobTxAnteWrapper is a wrapper for Antehandlers that need to be skipped for
// single msg clob txs `MsgPlaceOrder` and `MsgCancelOrder`. These transactions should always have `0` Gas,
// and therefore should never be charged a gas fee.
type SingleMsgClobTxAnteWrapper struct {
	antehandler sdk.AnteDecorator
}

func NewSingleMsgClobTxAnteWrapper(handler sdk.AnteDecorator) SingleMsgClobTxAnteWrapper {
	return SingleMsgClobTxAnteWrapper{
		antehandler: handler,
	}
}

func (antWrapper SingleMsgClobTxAnteWrapper) GetAnteHandler() sdk.AnteDecorator {
	return antWrapper.antehandler
}

func (anteWrapper SingleMsgClobTxAnteWrapper) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (sdk.Context, error) {
	isSingleClobMsgTx, err := IsSingleClobMsgTx(ctx, tx)
	if err != nil {
		return ctx, err
	}

	if isSingleClobMsgTx {
		return next(ctx, tx, simulate)
	}

	return anteWrapper.antehandler.AnteHandle(ctx, tx, simulate, next)
}

// ShortTermSingleMsgClobTxAnteWrapper is a wrapper for Antehandlers that need to be skipped for
// single msg clob txs `MsgPlaceOrder` and `MsgCancelOrder` which reference Short-Term orders.
// For example, these transactions do not require sequence number validation.
type ShortTermSingleMsgClobTxAnteWrapper struct {
	antehandler sdk.AnteDecorator
}

func NewShortTermSingleMsgClobTxAnteWrapper(handler sdk.AnteDecorator) ShortTermSingleMsgClobTxAnteWrapper {
	return ShortTermSingleMsgClobTxAnteWrapper{
		antehandler: handler,
	}
}

func (antWrapper ShortTermSingleMsgClobTxAnteWrapper) GetAnteHandler() sdk.AnteDecorator {
	return antWrapper.antehandler
}

func (anteWrapper ShortTermSingleMsgClobTxAnteWrapper) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (sdk.Context, error) {
	isShortTermClobMsgTx, err := IsShortTermClobMsgTx(ctx, tx)
	if err != nil {
		return ctx, err
	}

	if isShortTermClobMsgTx {
		return next(ctx, tx, simulate)
	}

	return anteWrapper.antehandler.AnteHandle(ctx, tx, simulate, next)
}
