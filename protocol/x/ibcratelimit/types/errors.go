package types

// DONTCOVER

import errorsmod "cosmossdk.io/errors"

// x/ibcratelimit module sentinel errors
var (
	ErrInvalidAuthority = errorsmod.Register(
		ModuleName,
		1001,
		"Authority is invalid",
	)
)
