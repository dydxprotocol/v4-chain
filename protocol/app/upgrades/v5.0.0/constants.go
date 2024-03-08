package v_5_0_0

import (
	store "cosmossdk.io/store/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/upgrades"

	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

const (
	UpgradeName = "v5.0.0"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName: UpgradeName,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{
			vaulttypes.StoreKey,
		},
	},
}
