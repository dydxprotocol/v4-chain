package v_2_0_0

import (
	store "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/upgrades"
)

const (
	UpgradeName = "v2.0.0"

	UpgradeHeight = 1805000 // Estimated to be 4:00 PM ET Nov 23, 2023
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
