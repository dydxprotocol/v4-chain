package ante

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	accountpluskeeper "github.com/dydxprotocol/v4-chain/protocol/x/accountplus/keeper"
)

// TODO: combine increment sequence and sequence verification into one decorator
// https://github.com/cosmos/cosmos-sdk/pull/18817

// Forked from "github.com/cosmos/cosmos-sdk/x/auth/ante/sigverify.go"
type IncrementSequenceDecorator struct {
	ak ante.AccountKeeper
}

func NewIncrementSequenceDecorator(ak ante.AccountKeeper) IncrementSequenceDecorator {
	return IncrementSequenceDecorator{
		ak: ak,
	}
}

func (isd IncrementSequenceDecorator) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (sdk.Context, error) {
	sigTx, ok := tx.(authsigning.SigVerifiableTx)
	if !ok {
		return ctx, errorsmod.Wrap(sdkerrors.ErrTxDecode, "invalid transaction type")
	}

	signatures, err := sigTx.GetSignaturesV2()
	if err != nil {
		return sdk.Context{}, err
	}

	for _, signature := range signatures {
		if accountpluskeeper.IsTimestampNonce(signature.Sequence) {
			// Skip increment for this signature
			continue
		}

		acc := isd.ak.GetAccount(ctx, signature.PubKey.Address().Bytes())
		if err := acc.SetSequence(acc.GetSequence() + 1); err != nil {
			panic(err)
		}

		isd.ak.SetAccount(ctx, acc)
	}

	return next(ctx, tx, simulate)
}
