package ante

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/lib"
)

// CircuitBreakerDecorator routes transactions through appropriate ante handlers based on
// the existence of `TxExtension`.
type CircuitBreakerDecorator struct {
	cdc                          codec.BinaryCodec
	authenticatorAnteHandlerFlow sdk.AnteHandler
	defaultAnteHandlerFlow       sdk.AnteHandler
}

// NewCircuitBreakerDecorator creates a new instance of CircuitBreakerDecorator with the provided parameters.
func NewCircuitBreakerDecorator(
	cdc codec.BinaryCodec,
	authenticatorAnteHandlerFlow sdk.AnteHandler,
	defaultAnteHandlerFlow sdk.AnteHandler,
) CircuitBreakerDecorator {
	return CircuitBreakerDecorator{
		cdc:                          cdc,
		authenticatorAnteHandlerFlow: authenticatorAnteHandlerFlow,
		defaultAnteHandlerFlow:       defaultAnteHandlerFlow,
	}
}

// AnteHandle checks if a tx is a smart account tx and routes it through the correct series of ante handlers.
//
// Note that whether or not to use the new authenticator flow is determined by the presence of the `TxExtension`.
// This is different from the Osmosis's implementation, which falls back to the original flow if
// smart account is disabled.
// The reason for this is because only minimal validation is done on resting maker orders when they get matched
// and this approach mitigates an issue for maker orders when smart account gets disabled through governance.
func (ad CircuitBreakerDecorator) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (newCtx sdk.Context, err error) {
	// Check that the authenticator flow is active
	if specified, _ := lib.HasSelectedAuthenticatorTxExtensionSpecified(tx, ad.cdc); specified {
		// Return and call the AnteHandle function on all the authenticator decorators.
		return ad.authenticatorAnteHandlerFlow(ctx, tx, simulate)
	}

	// Return and call the AnteHandle function on all the original decorators.
	return ad.defaultAnteHandlerFlow(ctx, tx, simulate)
}
