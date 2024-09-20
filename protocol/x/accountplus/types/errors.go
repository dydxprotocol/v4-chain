package types

import errorsmod "cosmossdk.io/errors"

var (
	ErrAuthenticatorNotFound = errorsmod.Register(
		ModuleName,
		1,
		"Authenticator is not found",
	)
	ErrInvalidAccountAddress = errorsmod.Register(
		ModuleName,
		2,
		"Invalid account address",
	)
	ErrAuthenticatorDataExceedsMaximumLength = errorsmod.Register(
		ModuleName,
		3,
		"Authenticator data exceeds maximum length",
	)
	ErrInitializingAuthenticator = errorsmod.Register(
		ModuleName,
		4,
		"Error initializing authenticator",
	)
)
