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
)
