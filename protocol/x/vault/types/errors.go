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
	ErrInvalidDepositAmount = errorsmod.Register(
		ModuleName,
		4,
		"Deposit amount is invalid",
	)
	ErrNonPositiveEquity = errorsmod.Register(
		ModuleName,
		5,
		"Equity is non-positive",
	)
	ErrZeroDenominator = errorsmod.Register(
		ModuleName,
		6,
		"Denominator is zero",
	)
	ErrNilFraction = errorsmod.Register(
		ModuleName,
		7,
		"Fraction is nil",
	)
	ErrInvalidOrderSizePpm = errorsmod.Register(
		ModuleName,
		8,
		"OrderSizePpm must be strictly greater than 0",
	)
	ErrInvalidOrderExpirationSeconds = errorsmod.Register(
		ModuleName,
		9,
		"OrderExpirationSeconds must be strictly greater than 0",
	)
	ErrInvalidSpreadMinPpm = errorsmod.Register(
		ModuleName,
		10,
		"SpreadMinPpm must be strictly greater than 0",
	)
	ErrInvalidLayers = errorsmod.Register(
		ModuleName,
		11,
		"Layers must be less than or equal to MaxUint8",
	)
)
