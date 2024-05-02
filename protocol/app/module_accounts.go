package app

import (
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/config"
	clobmoduletypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	rewardsmoduletypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/rewards/types"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	vestmoduletypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/vest/types"

	"golang.org/x/exp/maps"
)

func init() {
	// SetAddressPrefixes() explicitly in order to set the `dydx` address prefixes.
	config.SetAddressPrefixes()
}

var (
	// Module account permissions. Contains all module accounts on dYdX chain.
	maccPerms = map[string][]string{
		// -------- Native SDK module accounts --------
		authtypes.FeeCollectorName:     nil,
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
		icatypes.ModuleName:            nil,
		// -------- dYdX custom module accounts --------
		// subaccounts module account holds tokens for all subaccounts.
		satypes.ModuleName: nil,
		// clob insurance fund account manages insurance fund for liquidations.
		clobmoduletypes.InsuranceFundName: nil,
		// rewards treasury account distribute funds trading accounts.
		rewardsmoduletypes.TreasuryAccountName: nil,
		// rewards vester account vest rewards tokens into the rewards treasury.
		rewardsmoduletypes.VesterAccountName: nil,
		// community treasury account holds funds for community use.
		vestmoduletypes.CommunityTreasuryAccountName: nil,
		// community vester account vests funds into the community treasury.
		vestmoduletypes.CommunityVesterAccountName: nil,
	}
	// Blocked module accounts which cannot receive external funds.
	// By default, all non-custom modules (except for gov) are blocked. This prevents
	// unexpected violation of invariants (for example, https://github.com/cosmos/cosmos-sdk/issues/4795)
	blockedModuleAccounts = map[string]bool{
		authtypes.FeeCollectorName:     true,
		stakingtypes.BondedPoolName:    true,
		stakingtypes.NotBondedPoolName: true,
		ibctransfertypes.ModuleName:    true,
		icatypes.ModuleName:            true,
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
	// Shallow clone.
	return maps.Clone(maccPerms)
}

// BlockedAddresses returns all the app's blocked account addresses.
func BlockedAddresses() map[string]bool {
	// By default, returns all the app's blocked module account addresses.
	// Other regular addresses can also be added here.
	return moduleAccToAddress(blockedModuleAccounts)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func ModuleAccountAddrs() map[string]bool {
	return moduleAccToAddress(maccPerms)
}
