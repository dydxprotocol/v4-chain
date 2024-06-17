package ante

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	txsigning "cosmossdk.io/x/tx/signing"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdkante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	gometrics "github.com/hashicorp/go-metrics"
	"google.golang.org/protobuf/types/known/anypb"
)

// This is the maximum acceptable sequence number.
// Any transactions with a sequence number higher than this will be treated as a nonce timestamp in epoch milliseconds.
// This allows for over 1 trillion sequence numbers. Nonce timestamps must be later than sometime in the year 2004.
const sequence_cutoff = uint64(1) << 40

// This is the maximum acceptable difference between the sequence number and the current block time in milliseconds.
const nonce_timestamp_tolerance_ms = uint64(60_000)

// SigVerificationDecorator verifies all signatures for a tx and return an error if any are invalid. Note,
// the SigVerificationDecorator will not check signatures on ReCheck.
//
// CONTRACT: Pubkeys are set in context for all signers before this decorator runs
// CONTRACT: Tx must implement SigVerifiableTx interface
type SigVerificationDecorator struct {
	ak              sdkante.AccountKeeper
	signModeHandler *txsigning.HandlerMap
}

func NewSigVerificationDecorator(
	ak sdkante.AccountKeeper,
	signModeHandler *txsigning.HandlerMap,
) SigVerificationDecorator {
	return SigVerificationDecorator{
		ak:              ak,
		signModeHandler: signModeHandler,
	}
}

func (svd SigVerificationDecorator) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (newCtx sdk.Context, err error) {
	sigTx, ok := tx.(authsigning.Tx)
	if !ok {
		return ctx, errorsmod.Wrap(sdkerrors.ErrTxDecode, "invalid transaction type")
	}

	// stdSigs contains the sequence number, account number, and signatures.
	// When simulating, this would just be a 0-length slice.
	sigs, err := sigTx.GetSignaturesV2()
	if err != nil {
		return ctx, err
	}

	signers, err := sigTx.GetSigners()
	if err != nil {
		return ctx, err
	}

	// check that signer length and signature length are the same
	if len(sigs) != len(signers) {
		err := errorsmod.Wrapf(
			sdkerrors.ErrUnauthorized,
			"invalid number of signer;  expected: %d, got %d",
			len(signers),
			len(sigs),
		)
		return ctx, err
	}

	// Sequence number validation can be skipped if the given transaction consists of
	// only messages that use `GoodTilBlock` for replay protection.
	skipSequenceValidation := ShouldSkipSequenceValidation(tx.GetMsgs())

	for i, sig := range sigs {
		signer := signers[i]
		acc, err := sdkante.GetSignerAcc(ctx, svd.ak, signer)
		if err != nil {
			return ctx, err
		}

		// retrieve pubkey
		pubKey := acc.GetPubKey()
		if !simulate && pubKey == nil {
			return ctx, errorsmod.Wrap(sdkerrors.ErrInvalidPubKey, "pubkey on account is not set")
		}

		// retrieve signer data
		genesis := ctx.BlockHeight() == 0
		chainID := ctx.ChainID()
		var accNum uint64
		if !genesis {
			accNum = acc.GetAccountNumber()
		}

		// Sequence validation and storage must be performed.
		// The sequence number is used as a nonce for replay protection.
		// There are two different modes depending on whether the provided sequence number is
		// less-than or greater-than the sequence cutoff:
		//  - Less-than: it is treated as an incrementally increasing number.
		//  - Greater-than: it is treated as a nonce timestamp in epoch milliseconds.
		if !skipSequenceValidation {
			sequenceIsSequential := sig.Sequence <= sequence_cutoff
			if sequenceIsSequential {
				// Treat sequence number as an incrementally increasing number.
				if sig.Sequence != acc.GetSequence() {
					labels := make([]gometrics.Label, 0)
					if len(tx.GetMsgs()) > 0 {
						labels = append(
							labels,
							metrics.GetLabelForStringValue(metrics.MessageType, fmt.Sprintf("%T", tx.GetMsgs()[0])),
						)
					}
					telemetry.IncrCounterWithLabels(
						[]string{metrics.SequenceNumber, metrics.Invalid, metrics.Count},
						1,
						labels,
					)
					return ctx, errorsmod.Wrapf(
						sdkerrors.ErrWrongSequence,
						"account sequence mismatch, expected %d, got %d", acc.GetSequence(), sig.Sequence,
					)
				}
				// Increment the sequence number.
				if err := acc.SetSequence(acc.GetSequence() + 1); err != nil {
					panic(err)
				}
				svd.ak.SetAccount(ctx, acc)
			} else {
				// Treat sequence number as a nonce timestamp.
				midSeq := uint64(ctx.BlockTime().Unix())
				if sig.Sequence > midSeq+nonce_timestamp_tolerance_ms {
					return ctx, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "timestamp too far in the future")
				}
				if sig.Sequence < midSeq-nonce_timestamp_tolerance_ms {
					return ctx, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "timestamp too far in the past")
				}
				err := svd.ak.AddNonce(ctx, signer, sig.Sequence)
				if err != nil {
					return ctx, err
				}
			}
		}

		// no need to verify signatures on recheck tx
		if !simulate && !ctx.IsReCheckTx() {
			anyPk, _ := codectypes.NewAnyWithValue(pubKey)

			signerData := txsigning.SignerData{
				Address:       acc.GetAddress().String(),
				ChainID:       chainID,
				AccountNumber: accNum,
				Sequence:      acc.GetSequence(),
				PubKey: &anypb.Any{
					TypeUrl: anyPk.TypeUrl,
					Value:   anyPk.Value,
				},
			}
			adaptableTx, ok := tx.(authsigning.V2AdaptableTx)
			if !ok {
				return ctx, fmt.Errorf("expected tx to implement V2AdaptableTx, got %T", tx)
			}
			txData := adaptableTx.GetSigningTxData()
			err = authsigning.VerifySignature(ctx, pubKey, signerData, sig.Data, svd.signModeHandler, txData)
			if err != nil {
				var errMsg string
				if sdkante.OnlyLegacyAminoSigners(sig.Data) {
					// If all signers are using SIGN_MODE_LEGACY_AMINO, we rely on VerifySignature to check account sequence number,
					// and therefore communicate sequence number as a potential cause of error.
					errMsg = fmt.Sprintf(
						"signature verification failed; please verify account number (%d), sequence (%d)"+
							" and chain-id (%s)",
						accNum,
						acc.GetSequence(),
						chainID,
					)
				} else {
					errMsg = fmt.Sprintf(
						"signature verification failed; please verify account number (%d) and chain-id (%s): (%s)",
						accNum,
						chainID,
						err.Error(),
					)
				}
				return ctx, errorsmod.Wrap(sdkerrors.ErrUnauthorized, errMsg)
			}
		}
	}

	return next(ctx, tx, simulate)
}
