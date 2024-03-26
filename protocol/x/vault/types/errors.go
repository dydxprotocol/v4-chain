package types

// DONTCOVER

import errorsmod "cosmossdk.io/errors"

var (
	ErrNegativeShares = errorsmod.Register(
		ModuleName,
		1,
		"Shares are negative",
	)
	ErrClobPairNotFound = errorsmod.Register(
		ModuleName,
		2,
		"ClobPair not found",
	)
	ErrMarketParamNotFound = errorsmod.Register(
		ModuleName,
		3,
		"MarketParam not found",
	)
	ErrInvalidSkewFactorPpm = errorsmod.Register(
		ModuleName,
		4,
		"SkewFactorPpm must be strictly greater than 0 and less than 1",
	)
	ErrInvalidOrderSizePpm = errorsmod.Register(
		ModuleName,
		5,
		"OrderSizePpm must be strictly greater than 0",
	)
	ErrInvalidOrderExpirationSeconds = errorsmod.Register(
		ModuleName,
		6,
		"OrderExpirationSeconds must be strictly greater than 0",
	)
	ErrInvalidSpreadMinPpm = errorsmod.Register(
		ModuleName,
		7,
		"SpreadMinPpm must be strictly greater than 0",
	)
)
