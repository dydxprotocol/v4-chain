package types

import moderrors "cosmossdk.io/errors"

// DONTCOVER

var (
	ErrNonpositiveDuration = moderrors.Register(
		ModuleName,
		400,
		"Durations must be positive",
	)
	ErrUnorderedDurations = moderrors.Register(
		ModuleName,
		401,
		"Durations must be in ascending order by length",
	)
)
