package types

import authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

var (
	TreasuryAddress = authtypes.NewModuleAddress(TreasuryAccountName)
)
