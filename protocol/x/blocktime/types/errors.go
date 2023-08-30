package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
