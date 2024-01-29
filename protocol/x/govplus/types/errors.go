package types

// DONTCOVER

import errorsmod "cosmossdk.io/errors"

var (
	ErrValidatorAddress = errorsmod.Register(
		ModuleName,
		400,
		"Could not convert validator consensus address from bech32",
	)
	ErrInvalidAuthority = errorsmod.Register(
		ModuleName,
		401,
		"Authority is invalid",
	)
	ErrInvalidSlashFactor = errorsmod.Register(
		ModuleName,
		402,
		"slash_factor must be between 0 and 1 inclusive",
	)
	ErrInvalidTokensAtInfractionHeight = errorsmod.Register(
		ModuleName,
		403,
		"tokens_at_infraction_height must be positive",
	)
)
