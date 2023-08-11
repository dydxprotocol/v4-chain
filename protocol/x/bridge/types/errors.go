package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/bridge module sentinel errors
var (
	ErrNegativeDuration = sdkerrors.Register(
		ModuleName,
		400,
		"Duration is negative",
	)
	ErrRateOutOfBounds = sdkerrors.Register(
		ModuleName,
		401,
		"Rate is out of bounds",
	)
)
