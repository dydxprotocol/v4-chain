package types

import errorsmod "cosmossdk.io/errors"

var (
	// Add x/listing specific errors here
	ErrReferencePriceZero = errorsmod.Register(
		ModuleName,
		1,
		"reference price is zero",
	)

	ErrMarketsHardCapReached = errorsmod.Register(
		ModuleName,
		2,
		"listed markets hard cap reached",
	)
)
