package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AppInjectedMsgAnteWrapper is a wrapper for AnteHandlers that need to be skipped for
// "App-injected message" tx due to the fact that the "App-injected message" txs do
// not have certain features of regular txs like signatures.
type AppInjectedMsgAnteWrapper struct {
	antehandler sdk.AnteDecorator
}

func NewAppInjectedMsgAnteWrapper(handler sdk.AnteDecorator) AppInjectedMsgAnteWrapper {
	return AppInjectedMsgAnteWrapper{
		antehandler: handler,
	}
}

func (imaw AppInjectedMsgAnteWrapper) GetAnteHandler() sdk.AnteDecorator {
	return imaw.antehandler
}

func (imaw AppInjectedMsgAnteWrapper) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (sdk.Context, error) {
	// "App-injected message" tx must skip certain AnteHandlers.
	// For example, "App-injected message" tx is not signed on purpose and if a tx is not signed,
	// signature related functions like `GetSigners` will fail due to signer field being empty.
	// Therefore, we skip the AnteHandlers that require the tx to be signed.
	if IsSingleAppInjectedMsg(tx.GetMsgs()) {
		return next(ctx, tx, simulate)
	}

	return imaw.antehandler.AnteHandle(ctx, tx, simulate, next)
}
