package ante

import (
	"bytes"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authante "github.com/cosmos/cosmos-sdk/x/auth/ante"

	txsigning "cosmossdk.io/x/tx/signing"

	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
)

// AuthenticatorDecorator is responsible for processing authentication logic
// before transaction execution.
type AuthenticatorDecorator struct {
	accountPlusKeeper *keeper.Keeper
	accountKeeper     authante.AccountKeeper
	sigModeHandler    *txsigning.HandlerMap
	cdc               codec.Codec
}

// NewAuthenticatorDecorator creates a new instance of AuthenticatorDecorator with the provided parameters.
func NewAuthenticatorDecorator(
	cdc codec.Codec,
	accountPlusKeeper *keeper.Keeper,
	accountKeeper authante.AccountKeeper,
	sigModeHandler *txsigning.HandlerMap,
) AuthenticatorDecorator {
	return AuthenticatorDecorator{
		accountPlusKeeper: accountPlusKeeper,
		accountKeeper:     accountKeeper,
		sigModeHandler:    sigModeHandler,
		cdc:               cdc,
	}
}

// AnteHandle is the authenticator ante handler responsible for processing authentication
// logic before transaction execution.
func (ad AuthenticatorDecorator) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (newCtx sdk.Context, err error) {
	defer metrics.ModuleMeasureSince(
		types.ModuleName,
		metrics.AuthenticatorDecoratorAnteHandleLatency,
		time.Now(),
	)

	// Make sure smart account is active.
	if active := ad.accountPlusKeeper.GetIsSmartAccountActive(ctx); !active {
		return ctx, types.ErrSmartAccountNotActive
	}

	// Authenticators don't support manually setting the fee payer
	err = ad.ValidateAuthenticatorFeePayer(tx)
	if err != nil {
		return ctx, err
	}

	msgs := tx.GetMsgs()
	if len(msgs) == 0 {
		return ctx, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "no messages in transaction")
	}

	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, errorsmod.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	// The fee payer is the first signer of the transaction. This should have been enforced by the
	// ValidateAuthenticatorFeePayer.
	signers, _, err := ad.cdc.GetMsgV1Signers(msgs[0])
	if err != nil {
		return ctx, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "failed to get signers")
	}
	feePayer := sdk.AccAddress(signers[0])
	feeGranter := feeTx.FeeGranter()
	fee := feeTx.GetFee()

	selectedAuthenticators, err := ad.GetSelectedAuthenticators(tx, len(msgs))
	if err != nil {
		return ctx, err
	}

	// tracks are used to make sure that we only write to the store after every message is successful
	var tracks []func() error

	// Authenticate the accounts of all messages
	for msgIndex, msg := range msgs {
		signers, _, err := ad.cdc.GetMsgV1Signers(msg)
		if err != nil {
			return ctx, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "failed to get signers")
		}
		// Enforce only one signer per message
		if len(signers) != 1 {
			return ctx, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "messages must have exactly one signer")
		}

		// Get the account corresponding to the only signer of this message.
		account := sdk.AccAddress(signers[0])

		// Get the currently selected authenticator
		selectedAuthenticatorId := selectedAuthenticators[msgIndex]
		selectedAuthenticator, err := ad.accountPlusKeeper.GetInitializedAuthenticatorForAccount(
			ctx,
			account,
			selectedAuthenticatorId,
		)
		if err != nil {
			return ctx,
				errorsmod.Wrapf(
					err,
					"failed to get initialized authenticator "+
						"(account = %s, authenticator id = %d, msg index = %d, msg type url = %s)",
					account,
					selectedAuthenticatorId,
					msgIndex,
					sdk.MsgTypeURL(msg),
				)
		}

		// Generate the authentication request data
		authenticationRequest, err := lib.GenerateAuthenticationRequest(
			ctx,
			ad.cdc,
			ad.accountKeeper,
			ad.sigModeHandler,
			account,
			feePayer,
			feeGranter,
			fee,
			msg,
			tx,
			msgIndex,
			simulate,
		)
		if err != nil {
			return ctx,
				errorsmod.Wrapf(
					err,
					"failed to generate authentication data "+
						"(account = %s, authenticator id = %d, msg index = %d, msg type url = %s)",
					account,
					selectedAuthenticator.Id,
					msgIndex,
					sdk.MsgTypeURL(msg),
				)
		}

		authenticator := selectedAuthenticator.Authenticator
		stringId := strconv.FormatUint(selectedAuthenticator.Id, 10)
		authenticationRequest.AuthenticatorId = stringId

		// Consume the authenticator's static gas
		ctx.GasMeter().ConsumeGas(authenticator.StaticGas(), "authenticator static gas")

		// Authenticate should never modify state. That's what track is for
		neverWriteCtx, _ := ctx.CacheContext()
		authErr := authenticator.Authenticate(neverWriteCtx, authenticationRequest)

		// If authentication is successful, continue
		if authErr == nil {
			// Append the track closure to be called after every message is authenticated
			// Note: pre-initialize type URL to avoid closure issues from passing a msg
			// loop variable inside the closure.
			currentMsgTypeURL := sdk.MsgTypeURL(msg)
			tracks = append(tracks, func() error {
				err := authenticator.Track(ctx, authenticationRequest)
				if err != nil {
					// track should not fail in normal circumstances,
					// since it is intended to update track state before execution.
					// If it does fail, we log the error.
					metrics.IncrCounter(metrics.AuthenticatorTrackFailed, 1)
					ad.accountPlusKeeper.Logger(ctx).Error(
						"track failed",
						"account", account,
						"feePayer", feePayer,
						"msg", currentMsgTypeURL,
						"authenticatorId", stringId,
						"error", err,
					)

					return errorsmod.Wrapf(
						err,
						"track failed (account = %s, authenticator id = %s, authenticator type, %s, msg index = %d)",
						account,
						stringId,
						authenticator.Type(),
						msgIndex,
					)
				}
				return nil
			})
		}

		// If authentication failed, return an error
		if authErr != nil {
			return ctx, errorsmod.Wrapf(
				authErr,
				"authentication failed for message %d, authenticator id %d, type %s",
				msgIndex,
				selectedAuthenticator.Id,
				selectedAuthenticator.Authenticator.Type(),
			)
		}
	}

	// If the transaction has been authenticated, we call Track(...) on every message
	// to notify its authenticator so that it can handle any state updates.
	for _, track := range tracks {
		if err := track(); err != nil {
			return ctx, err
		}
	}

	return next(ctx, tx, simulate)
}

// ValidateAuthenticatorFeePayer enforces that the tx fee payer has not been set manually
// to an account different to the signer of the first message. This is a requirement
// for the authenticator module.
// The only user of a manually set fee payer is with fee grants, which are not
// available on dydx.
func (ad AuthenticatorDecorator) ValidateAuthenticatorFeePayer(tx sdk.Tx) error {
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return errorsmod.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	// The fee payer by default is the first signer of the transaction
	feePayer := feeTx.FeePayer()

	msgs := tx.GetMsgs()
	if len(msgs) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrTxDecode, "Tx must contain at least one message")
	}
	signers, _, err := ad.cdc.GetMsgV1Signers(msgs[0])
	if err != nil {
		return errorsmod.Wrap(sdkerrors.ErrUnauthorized, "failed to get signers")
	}
	if len(signers) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrTxDecode, "Tx message must contain at least one signer")
	}

	if !bytes.Equal(feePayer, signers[0]) {
		return errorsmod.Wrap(sdkerrors.ErrUnauthorized, "fee payer must be the first signer")
	}
	return nil
}

// GetSelectedAuthenticators retrieves the selected authenticators for the provided transaction extension
// and matches them with the number of messages in the transaction.
// If no selected authenticators are found in the extension, the function initializes the list with -1 values.
// It returns an array of selected authenticators or an error if the number of selected authenticators does not match
// the number of messages in the transaction.
func (ad AuthenticatorDecorator) GetSelectedAuthenticators(
	tx sdk.Tx,
	msgCount int,
) ([]uint64, error) {
	extTx, ok := tx.(authante.HasExtensionOptionsTx)
	if !ok {
		return nil, errorsmod.Wrap(sdkerrors.ErrTxDecode, "Tx must be a HasExtensionOptionsTx to use Authenticators")
	}

	// Get the selected authenticator options from the transaction.
	txOptions := lib.GetAuthenticatorExtension(extTx.GetNonCriticalExtensionOptions(), ad.cdc)
	if txOptions == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest,
			"Cannot get AuthenticatorTxOptions from tx")
	}
	// Retrieve the selected authenticators from the extension.
	selectedAuthenticators := txOptions.GetSelectedAuthenticators()

	if len(selectedAuthenticators) != msgCount {
		// Return an error if the number of selected authenticators does not match the number of messages.
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrInvalidRequest,
			"Mismatch between the number of selected authenticators and messages, "+
				"msg count %d, got %d selected authenticators",
			msgCount,
			len(selectedAuthenticators),
		)
	}

	return selectedAuthenticators, nil
}
