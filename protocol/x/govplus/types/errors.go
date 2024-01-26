package types

// DONTCOVER

import errorsmod "cosmossdk.io/errors"

var (
	ErrValidatorAddress = errorsmod.Register(
		ModuleName,
		400,
		"Could not convert validator consensus address from bech32",
	)
)
