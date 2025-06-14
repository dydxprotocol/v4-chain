package app

import (
	"sync"

	sending "github.com/dydxprotocol/v4-chain/protocol/x/sending/types"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/store/cachemulti"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"

	customante "github.com/dydxprotocol/v4-chain/protocol/app/ante"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	libante "github.com/dydxprotocol/v4-chain/protocol/lib/ante"
	"github.com/dydxprotocol/v4-chain/protocol/lib/log"
	accountplusante "github.com/dydxprotocol/v4-chain/protocol/x/accountplus/ante"
	accountpluskeeper "github.com/dydxprotocol/v4-chain/protocol/x/accountplus/keeper"
	clobante "github.com/dydxprotocol/v4-chain/protocol/x/clob/ante"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perpetualstypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

// HandlerOptions are the options required for constructing an SDK AnteHandler.
// Note: This struct is defined here in order to add `ClobKeeper`. We use
// struct embedding to include the normal cosmos-sdk `HandlerOptions`.
type HandlerOptions struct {
	ante.HandlerOptions
	Codec             codec.Codec
	AuthStoreKey      storetypes.StoreKey
	AccountplusKeeper *accountpluskeeper.Keeper
	ClobKeeper        clobtypes.ClobKeeper
	PerpetualsKeeper  perpetualstypes.PerpetualsKeeper
	PricesKeeper      pricestypes.PricesKeeper
	MarketMapKeeper   customante.MarketMapKeeper
	SendingKeeper     sending.SendingKeeper
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

	if options.AccountplusKeeper == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "accountplus keeper is required for ante builder")
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

	if options.PerpetualsKeeper == nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrLogic, "perpetuals keeper is required for ante builder")
	}

	if options.PricesKeeper == nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrLogic, "prices keeper is required for ante builder")
	}

	if options.MarketMapKeeper == nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrLogic, "market map keeper is required for ante builder")
	}

	if options.SendingKeeper == nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrLogic, "sending keeper is required for ante builder")
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
		replayProtection: customante.NewReplayProtectionDecorator(
			options.AccountKeeper,
			*options.AccountplusKeeper,
		),
		sigVerification: accountplusante.NewCircuitBreakerDecorator(
			options.Codec,
			sdk.ChainAnteDecorators(
				customante.NewEmitPubKeyEventsDecorator(),
				accountplusante.NewAuthenticatorDecorator(
					options.Codec,
					options.AccountplusKeeper,
					options.AccountKeeper,
					options.SignModeHandler,
				),
			),
			sdk.ChainAnteDecorators(
				ante.NewSetPubKeyDecorator(options.AccountKeeper),
				ante.NewSigGasConsumeDecorator(options.AccountKeeper, options.SigGasConsumer),
				customante.NewSigVerificationDecorator(
					options.AccountKeeper,
					options.SignModeHandler,
				),
			),
		),
		consumeTxSizeGas: ante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		deductFee: ante.NewDeductFeeDecorator(
			options.AccountKeeper,
			options.BankKeeper,
			options.FeegrantKeeper,
			options.TxFeeChecker,
		),
		clobRateLimit: clobante.NewRateLimitDecorator(options.ClobKeeper),
		clob:          clobante.NewClobDecorator(options.ClobKeeper, options.SendingKeeper),
		marketUpdates: customante.NewValidateMarketUpdateDecorator(
			options.PerpetualsKeeper, options.PricesKeeper, options.MarketMapKeeper,
		),
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
	replayProtection         customante.ReplayProtectionDecorator
	sigVerification          accountplusante.CircuitBreakerDecorator
	consumeTxSizeGas         ante.ConsumeTxSizeGasDecorator
	deductFee                ante.DeductFeeDecorator
	clobRateLimit            clobante.ClobRateLimitDecorator
	clob                     clobante.ClobDecorator
	marketUpdates            customante.ValidateMarketUpdateDecorator
}

func (h *lockingAnteHandler) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) {
	ctx = log.AddPersistentTagsToLogger(ctx,
		log.Callback, lib.TxMode(ctx),
		log.BlockHeight, ctx.BlockHeight()+1,
	)

	if clobante.HasClobMsg(tx) {
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
	if ctx, err = h.validateSigCount.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
		return ctx, err
	}
	if ctx, err = h.replayProtection.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
		return ctx, err
	}
	if ctx, err = h.sigVerification.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
		return ctx, err
	}

	var isShortTerm bool
	if isShortTerm, err = clobante.IsShortTermClobMsgTx(ctx, tx); err != nil {
		return ctx, err
	}

	var isTimestampNonce bool
	if isTimestampNonce, err = accountplusante.IsTimestampNonceTx(ctx, tx); err != nil {
		return ctx, err
	}

	if !isShortTerm && !isTimestampNonce {
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
	if ctx, err = h.marketUpdates.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
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
	if ctx, err = h.validateSigCount.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
		return ctx, err
	}
	if ctx, err = h.replayProtection.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
		return ctx, err
	}
	if ctx, err = h.sigVerification.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
		return ctx, err
	}

	var isTimestampNonce bool
	if isTimestampNonce, err = accountplusante.IsTimestampNonceTx(ctx, tx); err != nil {
		return ctx, err
	}

	if !isTimestampNonce {
		if ctx, err = h.incrementSequence.AnteHandle(ctx, tx, simulate, noOpAnteHandle); err != nil {
			return ctx, err
		}
	}

	// During non-simulated `checkTx` we must write the store since we own branching and writing.
	// During `deliverTx` and simulation the Cosmos SDK is responsible for branching and writing.
	if err == nil && !simulate && (ctx.IsCheckTx() || ctx.IsReCheckTx()) {
		cacheMs.Write()
	}

	return ctx, err
}
