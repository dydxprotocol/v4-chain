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
)
