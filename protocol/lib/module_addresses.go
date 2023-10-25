package lib

import (
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

var (
	// GovModuleAddress is the module address for the gov module.
	GovModuleAddress = authtypes.NewModuleAddress(govtypes.ModuleName)
)
