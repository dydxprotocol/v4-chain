package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	accountpluslib "github.com/dydxprotocol/v4-chain/protocol/x/accountplus/lib"
	aptypes "github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// MaybeValidateAuthenticators checks if the transaction has authenticators specified and if so,
// validates them. It returns an error if the authenticators are not valid or removed from state.
func (k Keeper) MaybeValidateAuthenticators(ctx sdk.Context, txBytes []byte) error {
	// Decode the tx from the tx bytes.
	tx, err := k.txDecoder(txBytes)
	if err != nil {
		return err
	}

	// Check if the tx had authenticator specified.
	specified, txOptions := accountpluslib.HasSelectedAuthenticatorTxExtensionSpecified(tx, k.cdc)
	if !specified {
		return nil
	}

	// The tx had authenticators specified.
	// First make sure smart account flow is enabled.
	if active := k.accountPlusKeeper.GetIsSmartAccountActive(ctx); !active {
		return aptypes.ErrSmartAccountNotActive
	}

	// Make sure txn is a SigVerifiableTx and get signers from the tx.
	sigVerifiableTx, ok := tx.(authsigning.SigVerifiableTx)
	if !ok {
		return errorsmod.Wrap(sdkerrors.ErrTxDecode, "invalid transaction type")
	}

	signers, err := sigVerifiableTx.GetSigners()
	if err != nil {
		return err
	}

	if len(signers) != 1 {
		return errorsmod.Wrap(types.ErrTxnHasMultipleSigners, "only one signer is allowed")
	}

	account := sdk.AccAddress(signers[0])

	// Retrieve the selected authenticators from the extension and make sure they are valid, i.e. they
	// are registered and not removed from state.
	//
	// Note that we only verify the existence of the authenticators here without actually
	// runnning them. This is because all current authenticators are stateless and do not read/modify any states.
	selectedAuthenticators := txOptions.GetSelectedAuthenticators()
	for _, authenticatorId := range selectedAuthenticators {
		_, err := k.accountPlusKeeper.GetInitializedAuthenticatorForAccount(
			ctx,
			account,
			authenticatorId,
		)
		if err != nil {
			return errorsmod.Wrapf(
				err,
				"selected authenticator (%s, %d) is not registered or removed from state",
				account.String(),
				authenticatorId,
			)
		}
	}
	return nil
}
