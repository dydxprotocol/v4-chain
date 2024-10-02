package ante

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdkante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	accountpluskeeper "github.com/dydxprotocol/v4-chain/protocol/x/accountplus/keeper"
	gometrics "github.com/hashicorp/go-metrics"
)

type ReplayProtectionDecorator struct {
	ak  sdkante.AccountKeeper
	akp accountpluskeeper.Keeper
}

func NewReplayProtectionDecorator(
	ak sdkante.AccountKeeper,
	akp accountpluskeeper.Keeper,
) ReplayProtectionDecorator {
	return ReplayProtectionDecorator{
		ak:  ak,
		akp: akp,
	}
}

func (rpd ReplayProtectionDecorator) AnteHandle(
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

	// Check that signer length and signature length are the same.
	// The ordering of the sigs and signers have matching ordering (sigs[i] belongs to signers[i]).
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

	if !skipSequenceValidation {
		// Iterate on sig and signer pairs.
		for i, sig := range sigs {
			acc, err := sdkante.GetSignerAcc(ctx, rpd.ak, signers[i])
			if err != nil {
				return ctx, err
			}

			// Check account sequence number.
			// Skip individual sequence number validation since this transaction use
			// `GoodTilBlock` for replay protection.
			if accountpluskeeper.IsTimestampNonce(sig.Sequence) {
				if err := rpd.akp.ProcessTimestampNonce(ctx, acc, sig.Sequence); err != nil {
					telemetry.IncrCounterWithLabels(
						[]string{metrics.TimestampNonce, metrics.Invalid, metrics.Count},
						1,
						[]gometrics.Label{metrics.GetLabelForIntValue(metrics.ExecMode, int(ctx.ExecMode()))},
					)
					return ctx, errorsmod.Wrap(sdkerrors.ErrWrongSequence, err.Error())
				}
				telemetry.IncrCounterWithLabels(
					[]string{metrics.TimestampNonce, metrics.Valid, metrics.Count},
					1,
					[]gometrics.Label{metrics.GetLabelForIntValue(metrics.ExecMode, int(ctx.ExecMode()))},
				)
			} else {
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
			}
		}
	}

	return next(ctx, tx, simulate)
}
