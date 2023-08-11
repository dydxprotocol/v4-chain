package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrNonpositiveDuration = sdkerrors.Register(
		ModuleName,
		400,
		"Duration is nonpositive",
	)
)
