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
	ErrInvalidAuthority = errorsmod.Register(
		ModuleName,
		402,
		"Authority is invalid",
	)
	ErrNegativeNextBlockDelay = errorsmod.Register(
		ModuleName,
		403,
		"next_block_delay must be non-negative",
	)
)
