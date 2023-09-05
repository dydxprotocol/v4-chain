package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

var (
	ErrNonpositiveDuration = sdkerrors.Register(
		ModuleName,
		400,
		"Duration is nonpositive",
	)
)
