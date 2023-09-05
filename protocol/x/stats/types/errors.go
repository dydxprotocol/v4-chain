package types

import moderrors "cosmossdk.io/errors"

// DONTCOVER

var (
	ErrNonpositiveDuration = moderrors.Register(
		ModuleName,
		400,
		"Duration is nonpositive",
	)
)
