package types

// DONTCOVER

import errorsmod "cosmossdk.io/errors"

var (
	ErrNegativeShares = errorsmod.Register(
		ModuleName,
		1,
		"Shares are negative",
	)
	ErrMarketParamNotFound = errorsmod.Register(
		ModuleName,
		2,
		"MarketParam not found",
	)
	ErrInvalidDepositAmount = errorsmod.Register(
		ModuleName,
		3,
		"Deposit amount is invalid",
	)
	ErrNonPositiveEquity = errorsmod.Register(
		ModuleName,
		4,
		"Equity is non-positive",
	)
	ErrZeroDenominator = errorsmod.Register(
		ModuleName,
		5,
		"Denominator is zero",
	)
	ErrNilFraction = errorsmod.Register(
		ModuleName,
		6,
		"Fraction is nil",
	)
	ErrInvalidOrderSizePctPpm = errorsmod.Register(
		ModuleName,
		7,
		"OrderSizePctPpm must be strictly greater than 0",
	)
	ErrInvalidOrderExpirationSeconds = errorsmod.Register(
		ModuleName,
		8,
		"OrderExpirationSeconds must be strictly greater than 0",
	)
	ErrInvalidSpreadMinPpm = errorsmod.Register(
		ModuleName,
		9,
		"SpreadMinPpm must be strictly greater than 0",
	)
	ErrInvalidLayers = errorsmod.Register(
		ModuleName,
		10,
		"Layers must be less than or equal to MaxUint8",
	)
	ErrZeroSharesToMint = errorsmod.Register(
		ModuleName,
		11,
		"Cannot mint zero shares",
	)
	ErrInvalidActivationThresholdQuoteQuantums = errorsmod.Register(
		ModuleName,
		12,
		"ActivationThresholdQuoteQuantums must be non-negative",
	)
	ErrInvalidOrderSize = errorsmod.Register(
		ModuleName,
		13,
		"OrderSize is invalid",
	)
	ErrInvalidOwner = errorsmod.Register(
		ModuleName,
		14,
		"Owner is invalid",
	)
	ErrMismatchedTotalAndOwnerShares = errorsmod.Register(
		ModuleName,
		15,
		"TotalShares does not match sum of OwnerShares",
	)
	ErrZeroMarketPrice = errorsmod.Register(
		ModuleName,
		16,
		"MarketPrice is zero",
	)
	ErrOrdersAndOrderIdsDiffLen = errorsmod.Register(
		ModuleName,
		17,
		"Orders and OrderIds must have the same length",
	)
)
