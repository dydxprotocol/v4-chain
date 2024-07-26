package types

import errorsmod "cosmossdk.io/errors"

var (
	// Add x/listing specific errors here
	ErrMarketNotFound = errorsmod.Register(ModuleName, 1, "market not found")
)
