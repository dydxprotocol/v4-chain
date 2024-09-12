package types

import errorsmod "cosmossdk.io/errors"

var (
	ErrAuthenticatorNotFound = errorsmod.Register(
		ModuleName,
		1,
		"Authenticator is not found",
	)
)
