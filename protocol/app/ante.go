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
<<<<<<< Updated upstream
=======
	return h.AnteHandle, nil
}

// An ante handler that returns the context.
func noOpAnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) {
	return ctx, nil
}

type lockingAnteHandler struct {
	globalLock   sync.Mutex
	authStoreKey storetypes.StoreKey

	setupContextDecorator    ante.SetUpContextDecorator
	freeInfiniteGasDecorator customante.FreeInfiniteGasDecorator
	extensionOptionsChecker  sdk.AnteDecorator
	validateMsgType          customante.ValidateMsgTypeDecorator
	txTimeoutHeight          ante.TxTimeoutHeightDecorator
	validateMemo             ante.ValidateMemoDecorator
	validateBasic            ante.ValidateBasicDecorator
	validateSigCount         ante.ValidateSigCountDecorator
	incrementSequence        ante.IncrementSequenceDecorator
	sigVerification          customante.SigVerificationDecorator
	consumeTxSizeGas         ante.ConsumeTxSizeGasDecorator
	deductFee                ante.DeductFeeDecorator
	setPubKey                ante.SetPubKeyDecorator
	sigGasConsume            ante.SigGasConsumeDecorator
	clobRateLimit            clobante.ClobRateLimitDecorator
	clob                     clobante.ClobDecorator
}

func (h *lockingAnteHandler) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) {
	ctx = log.AddPersistentTagsToLogger(ctx,
		log.Callback, lib.TxMode(ctx),
		log.BlockHeight, ctx.BlockHeight()+1,
	)

	isClob, err := clobante.IsSingleClobMsgTx(tx)
	if err != nil {
		return ctx, err
	} else if isClob {
		return h.clobAnteHandle(ctx, tx, simulate)
	}
	if libante.IsSingleAppInjectedMsg(tx.GetMsgs()) {
		return h.appInjectedMsgAnteHandle(ctx, tx, simulate)
	}

	return h.otherMsgAnteHandle(ctx, tx, simulate)
}

func (h *lockingAnteHandler) clobAnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool) (
	newCtx sdk.Context,
	err error,
) {
	// These ante decorators access state but only state that is mutated during `deliverTx`. The Cosmos SDK
	// is responsible for linearizing the reads and writes during `deliverTx`.
	if ctx, err = h.freeInfiniteGasDecorator.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
		return ctx, err
	}
	if ctx, err = h.extensionOptionsChecker.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
		return ctx, err
	}
	if ctx, err = h.validateMsgType.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
		return ctx, err
	}
	if ctx, err = h.validateBasic.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
		return ctx, err
	}
	if ctx, err = h.txTimeoutHeight.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
		return ctx, err
	}
	if ctx, err = h.validateMemo.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
		return ctx, err
	}

	// During `deliverTx` and simulation the Cosmos SDK is responsible for branching and writing the state store.
	// During `checkTx` we acquire a per account lock to prevent stale reads of state that can be mutated during
	// `checkTx`. Note that these messages are common so we use a row level like lock for each account and branch
	// the state store to support writes in the ante decorators that follow.
	var cacheMs storetypes.CacheMultiStore
	if !simulate && (ctx.IsCheckTx() || ctx.IsReCheckTx()) {
		sigTx, ok := tx.(authsigning.SigVerifiableTx)
		if !ok {
			return ctx, errorsmod.Wrap(sdkerrors.ErrTxDecode, "Tx must be a sigTx")
		}
		var signers [][]byte
		signers, err = sigTx.GetSigners()
		if err != nil {
			return ctx, err
		}

		cacheMs = ctx.MultiStore().(cachemulti.Store).CacheMultiStoreWithLocking(map[storetypes.StoreKey][][]byte{
			h.authStoreKey: signers,
		})
		defer cacheMs.(storetypes.LockingStore).Unlock()
		ctx = ctx.WithMultiStore(cacheMs)
	}

	if ctx, err = h.consumeTxSizeGas.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
		return ctx, err
	}
	if ctx, err = h.setPubKey.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
		return ctx, err
	}
	if ctx, err = h.validateSigCount.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
		return ctx, err
	}
	if ctx, err = h.sigGasConsume.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
		return ctx, err
	}
	if ctx, err = h.sigVerification.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
		return ctx, err
	}

	var isShortTerm bool
	if isShortTerm, err = clobante.IsShortTermClobMsgTx(ctx, tx); err != nil {
		return ctx, err
	}
	if !isShortTerm {
		if ctx, err = h.incrementSequence.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
			return ctx, err
		}
	}

	// We now acquire the global ante handler since the clob decorator is not thread safe and performs
	// several reads and writes across many stores.
	if !simulate && (ctx.IsCheckTx() || ctx.IsReCheckTx()) {
		h.globalLock.Lock()
		defer h.globalLock.Unlock()
	}

	if ctx, err = h.clobRateLimit.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
		return ctx, err
	}

	// During non-simulated `checkTx` we must write the store since we own branching and writing.
	// During `deliverTx` and simulation the Cosmos SDK is responsible for branching and writing.
	if err == nil && !simulate && (ctx.IsCheckTx() || ctx.IsReCheckTx()) {
		cacheMs.Write()
	}

	return ctx, err
}

// appInjectedMsgAnteHandle processes app injected messages through the necessary and sufficient set
// of ante decorators.
//
// Note that app injected messages do not require gas and are unsigned thus we do not need to:
//   - setup context to install a gas meter.
//   - validate basic for the signature check.
//   - set the pub key on the account.
//   - verify the signature.
//   - consume gas.
//   - deduct fees.
//   - increment the sequence number.
//   - rate limit or handle through the clob decorator since this isn't a clob message.
func (h *lockingAnteHandler) appInjectedMsgAnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool) (
	newCtx sdk.Context,
	err error,
) {
	// Note that app injected messages are only sent during `deliverTx` (checked by validateMsgType)
	// and hence we do not require any additional locking beyond what the Cosmos SDK already provides.
	if ctx, err = h.freeInfiniteGasDecorator.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
		return ctx, err
	}
	if ctx, err = h.extensionOptionsChecker.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
		return ctx, err
	}
	if ctx, err = h.validateMsgType.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
		return ctx, err
	}
	if ctx, err = h.txTimeoutHeight.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
		return ctx, err
	}
	if ctx, err = h.validateMemo.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
		return ctx, err
	}
	if ctx, err = h.consumeTxSizeGas.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
		return ctx, err
	}
	if ctx, err = h.validateSigCount.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
		return ctx, err
	}

	return ctx, err
}

// otherMsgAnteHandle processes all non-clob and non-app injected messages through the necessary and sufficient
// set of ante decorators.
//
// Note that these messages will never need to use the clob ante decorators so they are omitted.
func (h *lockingAnteHandler) otherMsgAnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool) (
	newCtx sdk.Context,
	err error,
) {
	// During `deliverTx` we hold an exclusive lock on `app.mtx` and have a context with a branched state store
	// allowing us to not have to perform any further locking or state store branching.
	//
	// For `checkTx`, these ante decorators access state but only state that is mutated during `deliverTx`
	// and hence since we already hold a read lock on `app.mtx` we can be certain that no state writes will occur.
	if ctx, err = h.setupContextDecorator.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
		return ctx, err
	}
	if ctx, err = h.freeInfiniteGasDecorator.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
		return ctx, err
	}
	if ctx, err = h.extensionOptionsChecker.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
		return ctx, err
	}
	if ctx, err = h.validateMsgType.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
		return ctx, err
	}
	if ctx, err = h.validateBasic.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
		return ctx, err
	}
	if ctx, err = h.txTimeoutHeight.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
		return ctx, err
	}
	if ctx, err = h.validateMemo.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
		return ctx, err
	}

	// During `deliverTx` and simulation the Cosmos SDK is responsible for branching and writing the state store.
	// During `checkTx` we acquire a per account lock to prevent stale reads of state that can be mutated during
	// `checkTx`. Note that we acquire row level like locks per account and also acquire the global lock to ensure
	// that we linearize reads and writes for accounts with the clobAnteHandle and the global lock ensures that
	// we linearize reads and writes to other stores since the deduct fees decorator mutates state outside of the
	// account keeper and those stores are currently not safe for concurrent use.
	var cacheMs storetypes.CacheMultiStore
	if !simulate && (ctx.IsCheckTx() || ctx.IsReCheckTx()) {
		sigTx, ok := tx.(authsigning.SigVerifiableTx)
		if !ok {
			return ctx, errorsmod.Wrap(sdkerrors.ErrTxDecode, "Tx must be a sigTx")
		}
		var signers [][]byte
		signers, err = sigTx.GetSigners()
		if err != nil {
			return ctx, err
		}

		cacheMs = ctx.MultiStore().(cachemulti.Store).CacheMultiStoreWithLocking(map[storetypes.StoreKey][][]byte{
			h.authStoreKey: signers,
		})
		defer cacheMs.(storetypes.LockingStore).Unlock()
		ctx = ctx.WithMultiStore(cacheMs)

		h.globalLock.Lock()
		defer h.globalLock.Unlock()
	}

	if ctx, err = h.consumeTxSizeGas.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
		return ctx, err
	}
	if ctx, err = h.deductFee.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
		return ctx, err
	}
	if ctx, err = h.setPubKey.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
		return ctx, err
	}
	if ctx, err = h.validateSigCount.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
		return ctx, err
	}
	if ctx, err = h.sigGasConsume.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
		return ctx, err
	}
	if ctx, err = h.sigVerification.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
		return ctx, err
	}
	if ctx, err = h.incrementSequence.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
		return ctx, err
	}

	// During non-simulated `checkTx` we must write the store since we own branching and writing.
	// During `deliverTx` and simulation the Cosmos SDK is responsible for branching and writing.
	if err == nil && !simulate && (ctx.IsCheckTx() || ctx.IsReCheckTx()) {
		cacheMs.Write()
	}

	return ctx, err
>>>>>>> Stashed changes
}
