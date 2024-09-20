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
	ErrInvalidQuoteQuantums = errorsmod.Register(
		ModuleName,
		4,
		"QuoteQuantums must be positive and less than 2^64",
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
	ErrZeroMarketPrice = errorsmod.Register(
		ModuleName,
		17,
		"MarketPrice is zero",
	)
	ErrOrdersAndOrderIdsDiffLen = errorsmod.Register(
		ModuleName,
		18,
		"Orders and OrderIds must have the same length",
	)
	ErrUnspecifiedVaultStatus = errorsmod.Register(
		ModuleName,
		19,
		"VaultStatus is unspecified",
	)
	ErrVaultParamsNotFound = errorsmod.Register(
		ModuleName,
		20,
		"VaultParams not found",
	)
	ErrEmptyOwnerAddress = errorsmod.Register(
		ModuleName,
		21,
		"Empty owner address",
	)
	ErrOwnerNotFound = errorsmod.Register(
		ModuleName,
		22,
		"Owner not found",
	)
	ErrLockedSharesExceedsOwnerShares = errorsmod.Register(
		ModuleName,
		23,
		"Locked shares exceeds owner shares",
	)
	ErrEmptyOperator = errorsmod.Register(
		ModuleName,
		24,
		"Empty operator address",
	)
	ErrInvalidSharesToWithdraw = errorsmod.Register(
		ModuleName,
		25,
		"Shares to withdraw must be positive and less than or equal to total shares",
	)
	ErrInvalidAuthority = errorsmod.Register(
		ModuleName,
		26,
		"Authority must be a module authority or operator",
	)
	ErrInsufficientWithdrawableShares = errorsmod.Register(
		ModuleName,
		27,
		"Insufficient withdrawable shares",
	)
	ErrInsufficientRedeemedQuoteQuantums = errorsmod.Register(
		ModuleName,
		28,
		"Insufficient redeemed quote quantums",
	)
	ErrDeactivatePositiveEquityVault = errorsmod.Register(
		ModuleName,
		29,
		"Cannot deactivate vaults with positive equity",
	)
	ErrNonPositiveShares = errorsmod.Register(
		ModuleName,
		30,
		"Shares must be positive",
	)
	ErrInvalidSkewFactor = errorsmod.Register(
		ModuleName,
		31,
		"Skew factor times order_size_pct must be less than 2 to avoid skewing over the spread",
	)
)
