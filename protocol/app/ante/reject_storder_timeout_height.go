package ante

import (
	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type (
	// RejecSTOrderTimeoutHeightDecorator defines an AnteHandler decorator that checks for a
	// tx height timeout and rejects the transaction if a tx height timeout is defined and contains
	// short-term orders.
	RejectSTOrderTimeoutHeightDecorator struct{}

	// TxWithTimeoutHeight defines the interface a tx must implement in order for
	// RejectSTOrderTimeoutHeightDecorator to process the tx.
	TxWithTimeoutHeight interface {
		sdk.Tx

		GetTimeoutHeight() uint64
	}
)

// RejecSTOrderTimeoutHeightDecorator defines an AnteHandler decorator that checks for a
// tx height timeout and rejects the transaction if a tx height timeout is defined and contains
// short-term orders.
func NewRejectSTOrderTimeoutHeightDecorator() RejectSTOrderTimeoutHeightDecorator {
	return RejectSTOrderTimeoutHeightDecorator{}
}

// AnteHandle implements an AnteHandler decorator for the RejectTimeoutHeightDecorator
// type where the tx is checked to see if it contains a non-zero height timeout and is a ST order.
// If a height timeout is provided and a STOrder is in the tx, reject the transaction.
func (txh RejectSTOrderTimeoutHeightDecorator) AnteHandle(
	ctx sdk.Context, tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (sdk.Context, error) {
	if !ShouldSkipSequenceValidation(tx.GetMsgs()) {
		return next(ctx, tx, simulate)
	}

	timeoutTx, ok := tx.(TxWithTimeoutHeight)
	// TimeoutHeightDecorator will handle error-ing on any transactions that don't implement the interface.
	if !ok {
		return next(ctx, tx, simulate)
	}

	timeoutHeight := timeoutTx.GetTimeoutHeight()
	if timeoutHeight > 0 {
		return ctx, errorsmod.Wrap(
			sdkerrors.ErrInvalidRequest,
			"a short term clob message may not have a non-zero timeout height, use goodTilBlock instead",
		)
	}

	return next(ctx, tx, simulate)
}
