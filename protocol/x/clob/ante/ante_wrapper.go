package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// OffChainSingleMsgClobTxAnteWrapper is a wrapper for Antehandlers that need to be skipped for
// off-chain single msg clob txs. These transactions should always have `0` Gas, and therefore
// should never be charged a gas fee.
type OffChainSingleMsgClobTxAnteWrapper struct {
	antehandler sdk.AnteDecorator
}

func NewOffChainSingleMsgClobTxAnteWrapper(handler sdk.AnteDecorator) OffChainSingleMsgClobTxAnteWrapper {
	return OffChainSingleMsgClobTxAnteWrapper{
		antehandler: handler,
	}
}

func (antWrapper OffChainSingleMsgClobTxAnteWrapper) GetAnteHandler() sdk.AnteDecorator {
	return antWrapper.antehandler
}

func (anteWrapper OffChainSingleMsgClobTxAnteWrapper) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (sdk.Context, error) {
	isOffChainSingleClobMsgTx, err := IsOffChainSingleClobMsgTx(ctx, tx)
	if err != nil {
		return ctx, err
	}

	if isOffChainSingleClobMsgTx {
		return next(ctx, tx, simulate)
	}

	return anteWrapper.antehandler.AnteHandle(ctx, tx, simulate, next)
}
