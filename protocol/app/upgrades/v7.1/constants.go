package v_7_1

import (
	store "cosmossdk.io/store/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/upgrades"
	affiliatetypes "github.com/dydxprotocol/v4-chain/protocol/x/affiliates/types"
)

const (
	UpgradeName = "v7.1"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName: UpgradeName,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{
			affiliatetypes.StoreKey,
		},
	},
}
