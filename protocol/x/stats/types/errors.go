package types

// DONTCOVER

import errorsmod "cosmossdk.io/errors"

var (
	ErrNonpositiveDuration = errorsmod.Register(
		ModuleName,
		400,
		"Duration is nonpositive",
	)
)
