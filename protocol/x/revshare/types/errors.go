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

	ErrMarketMapperRevShareDetailsNotFound = errorsmod.Register(
		ModuleName,
		3,
		"MarketMapperRevShareDetails not found for marketId",
	)
)
