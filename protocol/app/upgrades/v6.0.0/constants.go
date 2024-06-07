package v_6_0_0

import (
	store "cosmossdk.io/store/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/upgrades"
	listingtypes "github.com/dydxprotocol/v4-chain/protocol/x/listing/types"
)

const (
	UpgradeName = "v6.0.0"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName: UpgradeName,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{
			listingtypes.StoreKey,
		},
	},
}
