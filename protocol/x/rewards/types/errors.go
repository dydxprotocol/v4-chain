package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/rewards module sentinel errors
var (
	ErrInvalidTreasuryAccount  = sdkerrors.Register(ModuleName, 1001, "invalid treasury account")
	ErrInvalidFeeMultiplierPpm = sdkerrors.Register(ModuleName, 1002, "invalid FeeMultiplierPpm")
)
