package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"

	customante "github.com/dydxprotocol/v4/app/ante"
	libante "github.com/dydxprotocol/v4/lib/ante"
	clobante "github.com/dydxprotocol/v4/x/clob/ante"
	clobtypes "github.com/dydxprotocol/v4/x/clob/types"
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
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "account keeper is required for ante builder")
	}

	if options.BankKeeper == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "bank keeper is required for ante builder")
	}

	if options.ClobKeeper == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "clob keeper is required for ante builder")
	}

	if options.SignModeHandler == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "sign mode handler is required for ante builder")
	}

	anteDecorators := newAnteDecoratorChain(options)
	return sdk.ChainAnteDecorators(anteDecorators...), nil
}

// newAnteDecoratorChain returns a list of AnteDecorators in the expected application chain ordering
func newAnteDecoratorChain(options HandlerOptions) []sdk.AnteDecorator {
	return []sdk.AnteDecorator{
		// Note: app-injected messages, and clob off-chain transactions don't require Gas fees.
		libante.NewAppInjectedMsgAnteWrapper(
			clobante.NewOffChainSingleMsgClobTxAnteWrapper(
				ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
			),
		),
		// Set `FreeInfiniteGasMeter` for app-injected messages, and clob off-chain transactions.
		customante.NewFreeInfiniteGasDecorator(),
		ante.NewExtensionOptionsDecorator(options.ExtensionOptionChecker),
		customante.NewValidateMsgTypeDecorator(),
		libante.NewAppInjectedMsgAnteWrapper(
			ante.NewValidateBasicDecorator(),
		),
		ante.NewTxTimeoutHeightDecorator(),
		ante.NewValidateMemoDecorator(options.AccountKeeper),
		ante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),

		// Note: app-injected messages, and clob off-chain transactions don't require Gas fees.
		libante.NewAppInjectedMsgAnteWrapper(
			clobante.NewOffChainSingleMsgClobTxAnteWrapper(
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
			customante.NewIncrementSequenceDecorator(options.AccountKeeper),
		),
		clobante.NewClobDecorator(options.ClobKeeper),
	}
}
