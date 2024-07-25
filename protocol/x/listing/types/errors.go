package types

import errorsmod "cosmossdk.io/errors"

var (
	// Add x/listing specific errors here
	ErrReferencePriceZero = errorsmod.Register(
		ModuleName,
		1,
		"reference price is zero",
	)
)
