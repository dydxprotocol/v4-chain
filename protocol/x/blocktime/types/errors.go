package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

var (
	ErrNonpositiveDuration = sdkerrors.Register(
		ModuleName,
		400,
		"Durations must be positive",
	)
	ErrUnorderedDurations = sdkerrors.Register(
		ModuleName,
		401,
		"Durations must be in ascending order by length",
	)
)
