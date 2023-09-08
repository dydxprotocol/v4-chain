package types

// DONTCOVER

import errorsmod "cosmossdk.io/errors"

var (
	ErrNonpositiveDuration = errorsmod.Register(
		ModuleName,
		400,
		"Durations must be positive",
	)
	ErrUnorderedDurations = errorsmod.Register(
		ModuleName,
		401,
		"Durations must be in ascending order by length",
	)
)
