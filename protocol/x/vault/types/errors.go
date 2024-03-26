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
	ErrInvalidOrderSizePpm = errorsmod.Register(
		ModuleName,
		4,
		"OrderSizePpm must be strictly greater than 0",
	)
	ErrInvalidOrderExpirationSeconds = errorsmod.Register(
		ModuleName,
		5,
		"OrderExpirationSeconds must be strictly greater than 0",
	)
	ErrInvalidSpreadMinPpm = errorsmod.Register(
		ModuleName,
		6,
		"SpreadMinPpm must be strictly greater than 0",
	)
)
