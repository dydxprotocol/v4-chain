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
	ErrInvalidRevShareConfig = errorsmod.Register(
		ModuleName,
		4,
		"invalid unconditional revshare config",
	)

	ErrRevShareSafetyViolation = errorsmod.Register(
		ModuleName,
		5,
		"rev shares greater than or equal to 100%",
	)
	ErrTotalFeesSharedExceedsNetFees = errorsmod.Register(
		ModuleName,
		6,
		"total fees shared exceeds net fees",
	)
	ErrAffiliateFeesSharedExceedsNetFees = errorsmod.Register(
		ModuleName,
		7,
		"affiliate fees shared exceeds net fees",
	)
)
