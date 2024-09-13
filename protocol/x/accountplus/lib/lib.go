package lib

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
)

// HasSelectedAuthenticatorTxExtensionSpecified checks to see if the transaction has the correct
// extension, it returns false if we continue to the authenticator flow.
func HasSelectedAuthenticatorTxExtensionSpecified(
	tx sdk.Tx,
	accountPlusKeeper *keeper.Keeper,
) (bool, types.AuthenticatorTxOptions) {
	extTx, ok := tx.(authante.HasExtensionOptionsTx)
	if !ok {
		return false, nil
	}

	// Get the selected authenticator options from the transaction.
	txOptions := accountPlusKeeper.GetAuthenticatorExtension(extTx.GetNonCriticalExtensionOptions())

	// Check if authenticator transaction options are present and there is at least 1 selected.
	if txOptions == nil || len(txOptions.GetSelectedAuthenticators()) < 1 {
		return false, nil
	}

	return true, txOptions
}
