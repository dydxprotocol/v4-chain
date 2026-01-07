package v_9_6

import (
	store "cosmossdk.io/store/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/upgrades"
)

const (
	UpgradeName = "v9.6"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:   UpgradeName,
	StoreUpgrades: store.StoreUpgrades{},
}
