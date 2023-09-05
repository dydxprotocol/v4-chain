package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/rewards module sentinel errors
var (
	ErrInvalidTreasuryAccount  = sdkerrors.Register(ModuleName, 1001, "invalid treasury account")
	ErrInvalidFeeMultiplierPpm = sdkerrors.Register(ModuleName, 1002, "invalid FeeMultiplierPpm")
)
