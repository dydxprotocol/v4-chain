package app

import (
	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/store/cachemulti"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"sync"

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
	Codec        codec.Codec
	AuthStoreKey storetypes.StoreKey
	ClobKeeper   clobtypes.ClobKeeper
}

// NewAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & account numbers, deducts fees from the first
// signer, and handles in-memory clob messages.
//
// Note that the contract for the forked version of Cosmos SDK is that during `checkTx` the ante handler
// is responsible for branching and writing the state store. During `deliverTx` and simulation the Cosmos SDK
// is responsible for branching and writing the state store.
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

	if options.Codec == nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrLogic, "codec is required for ante builder")
	}

	if options.AuthStoreKey == nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrLogic, "auth store key is required for ante builder")
	}

	h := &lockingAnteHandler{
		authStoreKey:             options.AuthStoreKey,
		setupContextDecorator:    ante.NewSetUpContextDecorator(),
		freeInfiniteGasDecorator: customante.NewFreeInfiniteGasDecorator(),
		extensionOptionsChecker:  ante.NewExtensionOptionsDecorator(options.ExtensionOptionChecker),
		validateMsgType:          customante.NewValidateMsgTypeDecorator(),
		txTimeoutHeight:          ante.NewTxTimeoutHeightDecorator(),
		validateMemo:             ante.NewValidateMemoDecorator(options.AccountKeeper),
		validateBasic:            ante.NewValidateBasicDecorator(),
		validateSigCount:         ante.NewValidateSigCountDecorator(options.AccountKeeper),
		incrementSequence:        ante.NewIncrementSequenceDecorator(options.AccountKeeper),
		sigVerification:          customante.NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
		consumeTxSizeGas:         ante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		deductFee: ante.NewDeductFeeDecorator(
			options.AccountKeeper,
			options.BankKeeper,
			options.FeegrantKeeper,
			options.TxFeeChecker,
		),
		setPubKey:     ante.NewSetPubKeyDecorator(options.AccountKeeper),
		sigGasConsume: ante.NewSigGasConsumeDecorator(options.AccountKeeper, options.SigGasConsumer),
		clobRateLimit: clobante.NewRateLimitDecorator(options.ClobKeeper),
		clob:          clobante.NewClobDecorator(options.ClobKeeper),
	}
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

func (h *lockingAnteHandler) clobAnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, err error) {
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

		accountStoreKeys := make([][]byte, len(signers))
		for i, signer := range signers {
			var encodedSigner []byte
			encodedSigner, err = collections.EncodeKeyWithPrefix(authtypes.AddressStoreKeyPrefix, sdk.AccAddressKey, signer)
			if err != nil {
				return ctx, err
			}
			accountStoreKeys[i] = encodedSigner
		}

		cacheMs = ctx.MultiStore().(cachemulti.Store).CacheMultiStoreWithLocking(map[storetypes.StoreKey][][]byte{
			h.authStoreKey: accountStoreKeys,
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
	if ctx, err = h.clob.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
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
func (h *lockingAnteHandler) appInjectedMsgAnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, err error) {
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
	if ctx, err = h.txTimeoutHeight.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
		return ctx, err
	}
	if ctx, err = h.validateMemo.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
		return ctx, err
	}

	// During `deliverTx` and simulation the Cosmos SDK is responsible for branching and writing the state store.
	// During `checkTx` we acquire a lock to prevent stale reads of state that can be mutated during `checkTx`,
	// specifically account information. Note that these messages are rare so we can use a simple solution of using
	// a global lock for simplicity and maintenance here instead of a per account lock like we do in the
	// `clobAnteHandler` and that we do not have to branch the state store since the ante decorators that follow do
	// not mutate state.
	if !simulate && (ctx.IsCheckTx() || ctx.IsReCheckTx()) {
		h.globalLock.Lock()
		defer h.globalLock.Unlock()
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
func (h *lockingAnteHandler) otherMsgAnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, err error) {
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

	// During `deliverTx` we are guaranteed to hold an exclusive lock. During `checkTx` we
	// acquire a lock to prevent stale reads of state that can be mutated during `checkTx`, specifically
	// account information. Note that these messages are rare so we can use a simple solution of using a global lock
	// for simplicity and maintenance here instead of a per account lock like we do in the `clobAnteHandler` and
	// that we must branch the state store since the ante decorators that follow mutate account information.
	var cacheMs storetypes.CacheMultiStore
	if ctx.IsCheckTx() || ctx.IsReCheckTx() {
		h.globalLock.Lock()
		defer h.globalLock.Unlock()
		cacheMs = ctx.MultiStore().(cachemulti.Store).CacheMultiStore()
		ctx = ctx.WithMultiStore(cacheMs)
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
	if err == nil && (ctx.IsCheckTx() || ctx.IsReCheckTx()) {
		cacheMs.Write()
	}

	return ctx, err
}
