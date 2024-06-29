package app

import (
	errorsmod "cosmossdk.io/errors"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"

	customante "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ante"
	libante "github.com/StreamFinance-Protocol/stream-chain/protocol/lib/ante"
	clobante "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/ante"
	clobtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
	consumerkeeper "github.com/ethos-works/ethos/ethos-chain/x/ccv/consumer/keeper"
)

// HandlerOptions are the options required for constructing an SDK AnteHandler.
// Note: This struct is defined here in order to add `ClobKeeper`. We use
// struct embedding to include the normal cosmos-sdk `HandlerOptions`.
type HandlerOptions struct {
	ante.HandlerOptions
	Codec          codec.Codec
	AuthStoreKey   storetypes.StoreKey
	ClobKeeper     clobtypes.ClobKeeper
	IBCKeeper      ibckeeper.Keeper
	ConsumerKeeper consumerkeeper.Keeper
}

// NewAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & account numbers, deducts fees from the first
// signer, and handles in-memory clob messages.
//
// Note that the contract for the forked version of Cosmos SDK is that during `checkTx` the ante handler
// is responsible for branching and writing the state store. During this time the forked Cosmos SDK has
// a read lock allowing for parallel state reads but no writes. The `AnteHandler` is responsible for ensuring
// the linearization of reads and writes by having locks cover each other. This requires any ante decorators
// that read state that can be mutated during `checkTx` to acquire an appropriate lock. Today that is:
//   - account keeper params / consensus params (and all other state) are only read during `checkTx` and only
//     mutated during `deliverTx` thus no additional locking is needed to linearize reads and writes.
//   - accounts require the per account lock to be acquired since accounts have have pub keys set or the
//     sequence number incremented to linearize reads and writes.
//   - banks / fee state (and all other state) that can be mutated during `checkTx` requires the global
//     lock to be acquired before it is read or written to linearize reads and writes.
//
// During `deliverTx` and simulation the Cosmos SDK is responsible for branching and writing the state store
// so no additional locking is necessary to linearize state reads and writes. Note that simulation only ever occurs
// on a past block and not the current `checkState` so there is no opportunity for it to collide with concurrent
// `checkTx` invocations.
//
// Also note that all the ante decorators that are used return immediately the results of invoking `next` allowing
// us to significantly reduce the stack by saving and passing forward the context to the next ante decorator.
// This allows us to have a method that contains the order in which all the ante decorators are invoked including
// a single place to reason about the locking semantics without needing to look at several ante decorators.
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

	anteDecorators := NewAnteDecoratorChain(options)

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

// NewAnteDecoratorChain returns a list of AnteDecorators in the expected application chain ordering
func NewAnteDecoratorChain(options HandlerOptions) []sdk.AnteDecorator {
	return []sdk.AnteDecorator{
		baseapp.NewLockAndCacheContextAnteDecorator(),
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

		// Note: app-injected messages are not signed on purpose.
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
		// Note: app-injected messages are not signed on purpose.
		libante.NewAppInjectedMsgAnteWrapper(
			ante.NewSetPubKeyDecorator(options.AccountKeeper),
		),

		ante.NewValidateSigCountDecorator(options.AccountKeeper),

		// Note: app-injected messages don't require Gas fees.
		libante.NewAppInjectedMsgAnteWrapper(
			ante.NewSigGasConsumeDecorator(options.AccountKeeper, options.SigGasConsumer),
		),

		// Note: app-injected messages are not signed on purpose.
		libante.NewAppInjectedMsgAnteWrapper(
			customante.NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
		),

		// Note: app-injected messages, and short-term clob txs don't have sequence numbers on purpose.
		libante.NewAppInjectedMsgAnteWrapper(
			clobante.NewShortTermSingleMsgClobTxAnteWrapper(
				ante.NewIncrementSequenceDecorator(options.AccountKeeper),
			),
		),

		clobante.NewRateLimitDecorator(options.ClobKeeper),
		clobante.NewClobDecorator(options.ClobKeeper),
	}
}
