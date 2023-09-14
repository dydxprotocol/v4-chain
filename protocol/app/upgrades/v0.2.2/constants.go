package v0_2_2

import (
	store "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/upgrades"
)

const (
	UpgradeName   = "v0.2.2"
	UpgradeHeight = 500
)

var (
	Fork = upgrades.Fork{
		UpgradeName:   UpgradeName,
		UpgradeHeight: UpgradeHeight,
		UpgradeInfo:   "",
	}

	Upgrade = upgrades.Upgrade{
		UpgradeName:   UpgradeName,
		StoreUpgrades: store.StoreUpgrades{},
	}
)
