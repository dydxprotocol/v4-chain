package app

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"

	customante "github.com/dydxprotocol/v4-chain/protocol/app/ante"
	libante "github.com/dydxprotocol/v4-chain/protocol/lib/ante"
	clobante "github.com/dydxprotocol/v4-chain/protocol/x/clob/ante"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// HandlerOptions are the options required for constructing an SDK AnteHandler.
// Note: This struct is defined here in order to add `ClobKeeper`. We use
// struct embedding to include the normal cosmos-sdk `HandlerOptions`.
type HandlerOptions struct {
	ante.HandlerOptions
	ClobKeeper clobtypes.ClobKeeper
}

// NewAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & account numbers, deducts fees from the first
// signer, and handles in-memory clob messages.
//
// Link to default `AnteHandler` used by cosmos sdk:
// https://github.com/cosmos/cosmos-sdk/blob/3bb27795742dab2451b232bab02b82566d1a0192/x/auth/ante/ante.go#L25
func NewAnteHandler(options HandlerOptions) (sdk.AnteHandler, error) {
	if options.AccountKeeper == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "account keeper is required for ante builder")
	}

	if options.BankKeeper == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "bank keeper is required for ante builder")
	}

	if options.ClobKeeper == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "clob keeper is required for ante builder")
	}

	if options.SignModeHandler == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "sign mode handler is required for ante builder")
	}

	anteDecorators := newAnteDecoratorChain(options)

	// TODO(STAB-24): This change can be reverted to using ChainAnteDecorators again once
	// https://github.com/cosmos/cosmos-sdk/pull/16076 is merged, released, and we pick-up the SDK version containing
	// the change.
	anteHandlers := make([]sdk.AnteHandler, len(anteDecorators)+1)
	// Install the terminator ante handler.
	anteHandlers[len(anteDecorators)] = func(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) {
		return ctx, nil
	}
	for i := 0; i < len(anteDecorators); i++ {
		// Make a copy of the value to ensure that we can hold a reference to it. This avoids the golang common mistake:
		// https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		ii := i
		anteHandlers[ii] = func(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) {
			return anteDecorators[ii].AnteHandle(ctx, tx, simulate, anteHandlers[ii+1])
		}
	}

	return anteHandlers[0], nil
}

// newAnteDecoratorChain returns a list of AnteDecorators in the expected application chain ordering
func newAnteDecoratorChain(options HandlerOptions) []sdk.AnteDecorator {
	return []sdk.AnteDecorator{
		// Note: app-injected messages, and clob transactions don't require Gas fees.
		libante.NewAppInjectedMsgAnteWrapper(
			clobante.NewSingleMsgClobTxAnteWrapper(
				ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
			),
		),
		// Set `FreeInfiniteGasMeter` for app-injected messages, and clob transactions.
		customante.NewFreeInfiniteGasDecorator(),
		ante.NewExtensionOptionsDecorator(options.ExtensionOptionChecker),
		customante.NewValidateMsgTypeDecorator(),
		libante.NewAppInjectedMsgAnteWrapper(
			ante.NewValidateBasicDecorator(),
		),
		ante.NewTxTimeoutHeightDecorator(),
		ante.NewValidateMemoDecorator(options.AccountKeeper),
		ante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),

		// Note: app-injected messages, and clob transactions don't require Gas fees.
		libante.NewAppInjectedMsgAnteWrapper(
			clobante.NewSingleMsgClobTxAnteWrapper(
				ante.NewDeductFeeDecorator(
					options.AccountKeeper,
					options.BankKeeper,
					options.FeegrantKeeper,
					options.TxFeeChecker,
				),
			),
		),

		// SetPubKeyDecorator must be called before all signature verification decorators
		libante.NewAppInjectedMsgAnteWrapper(
			ante.NewSetPubKeyDecorator(options.AccountKeeper),
		),
		ante.NewValidateSigCountDecorator(options.AccountKeeper),
		libante.NewAppInjectedMsgAnteWrapper(
			ante.NewSigGasConsumeDecorator(options.AccountKeeper, options.SigGasConsumer),
		),
		libante.NewAppInjectedMsgAnteWrapper(
			customante.NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
		),
		libante.NewAppInjectedMsgAnteWrapper(
			clobante.NewShortTermSingleMsgClobTxAnteWrapper(
				ante.NewIncrementSequenceDecorator(options.AccountKeeper),
			),
		),
		clobante.NewRateLimitDecorator(options.ClobKeeper),
		clobante.NewClobDecorator(options.ClobKeeper),
	}
}
