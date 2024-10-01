package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MaybeValidateAuthenticators checks if the transaction has authenticators specified and if so,
// validates them. It returns an error if the authenticators are not valid or removed from state.
func (k Keeper) MaybeValidateAuthenticators(ctx sdk.Context, txBytes []byte) error {
	// Decode the tx from the tx bytes.
	tx, err := k.txDecoder(txBytes)
	if err != nil {
		return err
	}

	// Perform a light-weight validation of the authenticators via the accountplus module.
	//
	// Note that alternatively we could have been calling the ante handler directly on this transaction,
	// but there are some deadlock issues that are non-trivial to resolve.
	return k.accountPlusKeeper.MaybeValidateAuthenticators(ctx, tx)
}
