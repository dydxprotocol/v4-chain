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
	originalAnteHandlerFlow      sdk.AnteHandler
}

// NewCircuitBreakerDecorator creates a new instance of CircuitBreakerDecorator with the provided parameters.
func NewCircuitBreakerDecorator(
	cdc codec.BinaryCodec,
	auth sdk.AnteHandler,
	classic sdk.AnteHandler,
) CircuitBreakerDecorator {
	return CircuitBreakerDecorator{
		cdc:                          cdc,
		authenticatorAnteHandlerFlow: auth,
		originalAnteHandlerFlow:      classic,
	}
}

// AnteHandle checks if a tx is a smart account tx and routes it through the correct series of ante handlers.
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
	return ad.originalAnteHandlerFlow(ctx, tx, simulate)
}
