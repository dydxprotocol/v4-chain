package v_9_8

import (
	store "cosmossdk.io/store/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/upgrades"
	canonicalusdctypes "github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"
)

const UpgradeName = "v9.8"

var Upgrade = upgrades.Upgrade{
	UpgradeName: UpgradeName,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{canonicalusdctypes.StoreKey},
	},
}
