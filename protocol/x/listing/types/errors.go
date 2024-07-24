package types

import errorsmod "cosmossdk.io/errors"

var (
	// Add x/listing specific errors here
	ErrReferencePriceZero = errorsmod.Register(
		ModuleName,
		1,
		"reference price is zero",
	)

	ErrReferencePriceOutOfRange = errorsmod.Register(
		ModuleName,
		2,
		"reference price is out of range of int32",
	)
)
