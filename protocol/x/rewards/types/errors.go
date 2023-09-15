package types

// DONTCOVER

import errorsmod "cosmossdk.io/errors"

// x/rewards module sentinel errors
var (
	ErrInvalidTreasuryAccount  = errorsmod.Register(ModuleName, 1001, "invalid treasury account")
	ErrInvalidFeeMultiplierPpm = errorsmod.Register(ModuleName, 1002, "invalid FeeMultiplierPpm")
)
