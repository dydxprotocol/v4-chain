// Package protectedaccounts defines module accounts whose balances back
// protocol invariants (staked funds, subaccount collateral) and therefore must
// never be configurable as the source of outbound transfers (gov-directed
// sends, vest entries, reward treasuries).
package protectedaccounts

import (
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

var protectedModuleNames = map[string]bool{
	stakingtypes.BondedPoolName:    true,
	stakingtypes.NotBondedPoolName: true,
	satypes.ModuleName:             true,
}

func IsProtectedModuleName(moduleName string) bool {
	return protectedModuleNames[moduleName]
}
