package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AUTHENTICATOR_DATA_MAX_LENGTH is the maximum length of the data field in an authenticator.
//
// This is used as a light-weight spam mitigation measure to prevent users from adding
// arbitrarily complex authenticators that are too resource intensive.
const AUTHENTICATOR_DATA_MAX_LENGTH = 1024

// ValidateBasic performs stateless validation for the `MsgAddAuthenticator` msg.
func (msg *MsgAddAuthenticator) ValidateBasic() (err error) {
	// Validate account sender.
	_, err = sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return ErrInvalidAccountAddress
	}

	// Make sure that the authenticator data does not exceed the maximum length.
	if len(msg.Data) > AUTHENTICATOR_DATA_MAX_LENGTH {
		return ErrAuthenticatorDataExceedsMaximumLength
	}

	return nil
}
