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
	ErrInvalidOrderSizePctPpm = errorsmod.Register(
		ModuleName,
		8,
		"OrderSizePctPpm must be strictly greater than 0",
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
	ErrZeroSharesToMint = errorsmod.Register(
		ModuleName,
		12,
		"Cannot mint zero shares",
	)
	ErrInvalidActivationThresholdQuoteQuantums = errorsmod.Register(
		ModuleName,
		13,
		"ActivationThresholdQuoteQuantums must be non-negative",
	)
	ErrInvalidOrderSize = errorsmod.Register(
		ModuleName,
		14,
		"OrderSize is invalid",
	)
	ErrInvalidOwner = errorsmod.Register(
		ModuleName,
		15,
		"Owner is invalid",
	)
	ErrMismatchedTotalAndOwnerShares = errorsmod.Register(
		ModuleName,
		16,
		"TotalShares does not match sum of OwnerShares",
	)
	ErrInvalidWithdrawalAmount = errorsmod.Register(
		ModuleName,
		17,
		"Withdrawal amount is invalid",
	)
	ErrVaultNotFound = errorsmod.Register(
		ModuleName,
		18,
		"Vault not found",
	)
	ErrOwnerShareNotFound = errorsmod.Register(
		ModuleName,
		19,
		"Owner share not found",
	)
	ErrNotEnoughOwnerShare = errorsmod.Register(
		ModuleName,
		20,
		"Owner share is less than or equal to zero",
	)
)
