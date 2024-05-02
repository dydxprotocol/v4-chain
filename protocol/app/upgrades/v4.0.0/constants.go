package v_4_0_0

import (
	store "cosmossdk.io/store/types"
	circuittypes "cosmossdk.io/x/circuit/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/upgrades"

	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/types"
)

const (
	UpgradeName = "v4.0.0"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName: UpgradeName,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{

			// Add circuittypes as per 0.47 to 0.50 upgrade handler
			// https://github.com/cosmos/cosmos-sdk/blob/b7d9d4c8a9b6b8b61716d2023982d29bdc9839a6/simapp/upgrades.go#L21
			circuittypes.ModuleName,

			// Add new ICA stores that are needed by ICA host types as of v8.
			icacontrollertypes.StoreKey,

			// Add authz module to allow granting arbitrary privileges from one account to another acocunt.
			authzkeeper.StoreKey,
		},
	},
}
