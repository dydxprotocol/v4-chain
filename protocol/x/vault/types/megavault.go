package types

import (
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

var (
	MegavaultMainAddress = authtypes.NewModuleAddress(MegavaultAccountName)
	// MegavaultMainSubaccount is subaccount 0 of the module account derived from string "megavault".
	MegavaultMainSubaccount = satypes.SubaccountId{
		Owner:  MegavaultMainAddress.String(),
		Number: 0,
	}
)
