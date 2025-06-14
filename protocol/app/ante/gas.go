package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/ante/types"
	libante "github.com/dydxprotocol/v4-chain/protocol/lib/ante"
	clobante "github.com/dydxprotocol/v4-chain/protocol/x/clob/ante"
)

// FreeInfiniteGasDecorator is an AnteHandler that sets `GasMeter` to
// `FreeInfiniteGasMeter` for off-chain single clob msg transactions, and app-injected transactions.
// These transactions should not use any gas, and the sender should not be charged any gas.
// Using this meter means gas will never be consumed for these transactions.
// Also note that not explicitly setting a `gasMeter` means that the `gasMeter` from the previous transaction
// or from `BeginBlock` will be used. Not doing this could result in consensus failure as demonstrated in #869.
// Cosmos SDK expects an explicit call to `WithGasMeter` at the beginning of the AnteHandler chain.
type FreeInfiniteGasDecorator struct {
}

func NewFreeInfiniteGasDecorator() FreeInfiniteGasDecorator {
	return FreeInfiniteGasDecorator{}
}

func (dec FreeInfiniteGasDecorator) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (newCtx sdk.Context, err error) {
	hasClobMsg := clobante.HasClobMsg(tx)
	if err != nil {
		return ctx, err
	}

	// If this is a clob msg tx, or a single app-injected msg tx, then set the gas meter to
	// FreeInfiniteGasMeter.
	// ValidateClobMsgTx will enforce that at most 1 transfer msg is allowed in a clob msg tx,
	// which makes it safe to set the gas meter to FreeInfiniteGasMeter
	if hasClobMsg || libante.IsSingleAppInjectedMsg(tx.GetMsgs()) {
		newCtx = ctx.WithGasMeter(types.NewFreeInfiniteGasMeter())
		return next(newCtx, tx, simulate)
	}

	return next(ctx, tx, simulate)
}
