package types

import errorsmod "cosmossdk.io/errors"

var (
	ErrInvalidAddress = errorsmod.Register(
		ModuleName,
		1,
		"invalid address",
	)

	ErrInvalidRevenueSharePpm = errorsmod.Register(
		ModuleName,
		2,
		"invalid revenue share ppm",
	)
)
