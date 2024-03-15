package types

// DONTCOVER

import errorsmod "cosmossdk.io/errors"

var (
	ErrNegativeShares = errorsmod.Register(
		ModuleName,
		1,
		"Shares are negative",
	)
)
