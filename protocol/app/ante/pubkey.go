package ante

import (
	"encoding/base64"
	"fmt"

	errorsmod "cosmossdk.io/errors"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
)

// NewEmitPubKeyEventsDecorator emits events for each signer's public key.
// CONTRACT: Tx must implement SigVerifiableTx interface
type EmitPubKeyEventsDecorator struct{}

func NewEmitPubKeyEventsDecorator() EmitPubKeyEventsDecorator {
	return EmitPubKeyEventsDecorator{}
}

func (eped EmitPubKeyEventsDecorator) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (sdk.Context, error) {
	sigTx, ok := tx.(authsigning.SigVerifiableTx)
	if !ok {
		return ctx, errorsmod.Wrap(sdkerrors.ErrTxDecode, "invalid tx type")
	}

	signers, err := sigTx.GetSigners()
	if err != nil {
		return ctx, err
	}

	signerStrs := make([]string, len(signers))

	// Also emit the following events, so that txs can be indexed by these
	// indices:
	// - signature (via `tx.signature='<sig_as_base64>'`),
	// - concat(address,"/",sequence) (via `tx.acc_seq='cosmos1abc...def/42'`).
	sigs, err := sigTx.GetSignaturesV2()
	if err != nil {
		return ctx, err
	}

	var events sdk.Events
	for i, sig := range sigs {
		events = append(events, sdk.NewEvent(sdk.EventTypeTx,
			sdk.NewAttribute(sdk.AttributeKeyAccountSequence, fmt.Sprintf("%s/%d", signerStrs[i], sig.Sequence)),
		))

		sigBzs, err := signatureDataToBz(sig.Data)
		if err != nil {
			return ctx, err
		}
		for _, sigBz := range sigBzs {
			events = append(events, sdk.NewEvent(sdk.EventTypeTx,
				sdk.NewAttribute(sdk.AttributeKeySignature, base64.StdEncoding.EncodeToString(sigBz)),
			))
		}
	}

	ctx.EventManager().EmitEvents(events)

	return next(ctx, tx, simulate)
}

// signatureDataToBz converts a SignatureData into raw bytes signature.
// For SingleSignatureData, it returns the signature raw bytes.
// For MultiSignatureData, it returns an array of all individual signatures,
// as well as the aggregated signature.
func signatureDataToBz(data signing.SignatureData) ([][]byte, error) {
	if data == nil {
		return nil, fmt.Errorf("got empty SignatureData")
	}

	switch data := data.(type) {
	case *signing.SingleSignatureData:
		return [][]byte{data.Signature}, nil
	case *signing.MultiSignatureData:
		sigs := [][]byte{}
		var err error

		for _, d := range data.Signatures {
			nestedSigs, err := signatureDataToBz(d)
			if err != nil {
				return nil, err
			}
			sigs = append(sigs, nestedSigs...)
		}

		multiSignature := cryptotypes.MultiSignature{
			Signatures: sigs,
		}
		aggregatedSig, err := multiSignature.Marshal()
		if err != nil {
			return nil, err
		}
		sigs = append(sigs, aggregatedSig)

		return sigs, nil
	default:
		return nil, sdkerrors.ErrInvalidType.Wrapf("unexpected signature data type %T", data)
	}
}
