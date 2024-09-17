package types

import errorsmod "cosmossdk.io/errors"

var (
	// Add x/listing specific errors here
	ErrMarketNotFound = errorsmod.Register(
		ModuleName,
		1,
		"market not found",
	)

	ErrReferencePriceZero = errorsmod.Register(
		ModuleName,
		2,
		"reference price is zero",
	)

	ErrMarketsHardCapReached = errorsmod.Register(
		ModuleName,
		3,
		"listed markets hard cap reached",
	)

	ErrInvalidDepositAmount = errorsmod.Register(
		ModuleName,
		4,
		"invalid vault deposit amount",
	)

	ErrInvalidNumBlocksToLockShares = errorsmod.Register(
		ModuleName,
		5,
		"invalid number of blocks to lock shares",
	)

	ErrInvalidMarketMapTickerMetadata = errorsmod.Register(
		ModuleName,
		6,
		"invalid market map ticker metadata",
	)
)
