package app

import (
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/config"
	"github.com/dydxprotocol/v4-chain/protocol/lib/maps"
	bridgemoduletypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	clobmoduletypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	rewardsmoduletypes "github.com/dydxprotocol/v4-chain/protocol/x/rewards/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

func init() {
	// This package does not contain the `app/config` package in its import chain, and therefore needs to call
	// SetAddressPrefixes() explicitly in order to set the `dydx` address prefixes.
	config.SetAddressPrefixes()
}

var (
	// Module account permissions. Contains all module accounts on dYdX chain.
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:             nil,
		bridgemoduletypes.ModuleName:           {authtypes.Minter},
		distrtypes.ModuleName:                  nil,
		stakingtypes.BondedPoolName:            {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName:         {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:                    {authtypes.Burner},
		ibctransfertypes.ModuleName:            {authtypes.Minter, authtypes.Burner},
		satypes.ModuleName:                     nil,
		clobmoduletypes.InsuranceFundName:      nil,
		rewardsmoduletypes.TreasuryAccountName: nil,
		rewardsmoduletypes.VesterAccountName:   nil,
	}
	// Blocked module accounts which cannot receive external funds.
	blockedModuleAccounts = map[string]bool{
		authtypes.FeeCollectorName:     true,
		bridgemoduletypes.ModuleName:   true,
		distrtypes.ModuleName:          true,
		stakingtypes.BondedPoolName:    true,
		stakingtypes.NotBondedPoolName: true,
		ibctransfertypes.ModuleName:    true,
	}
	// Module accounts which are not blocked. This includes:
	// - governance module account (needed for https://github.com/cosmos/cosmos-sdk/pull/12852)
	// - dYdX custom module accounts
	whitelistedModuleAccounts = map[string]bool{
		govtypes.ModuleName:                    true,
		satypes.ModuleName:                     true,
		clobmoduletypes.InsuranceFundName:      true,
		rewardsmoduletypes.TreasuryAccountName: true,
		rewardsmoduletypes.VesterAccountName:   true,
	}
)

func moduleAccToAddress[V any](accs map[string]V) map[string]bool {
	addrs := make(map[string]bool)
	for acc := range accs {
		addrs[authtypes.NewModuleAddress(acc).String()] = true
	}
	return addrs
}

// GetMaccPerms returns a copy of the module account permissions
func GetMaccPerms() map[string][]string {
	return maps.Copy(maccPerms)
}

// BlockedAddresses returns all the app's blocked account addresses.
func BlockedAddresses() map[string]bool {
	// By default, returns all the app's blocked module account addresses.
	// Other regular addresses can also be added here.
	return moduleAccToAddress(blockedModuleAccounts)
}

// WhitelistedModuleAddresses returns all the app's unblocked module account addresses.
func WhitelistedModuleAddresses() map[string]bool {
	return moduleAccToAddress(whitelistedModuleAccounts)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func ModuleAccountAddrs() map[string]bool {
	return moduleAccToAddress(maccPerms)
}
