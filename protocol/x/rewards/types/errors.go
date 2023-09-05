package types

import moderrors "cosmossdk.io/errors"

// DONTCOVER

// x/rewards module sentinel errors
var (
	ErrInvalidTreasuryAccount  = moderrors.Register(ModuleName, 1001, "invalid treasury account")
	ErrInvalidFeeMultiplierPpm = moderrors.Register(ModuleName, 1002, "invalid FeeMultiplierPpm")
)
