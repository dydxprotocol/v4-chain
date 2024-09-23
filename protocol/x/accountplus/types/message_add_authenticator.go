package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const AUTHENTICATOR_DATA_MAX_LENGTH = 256

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
