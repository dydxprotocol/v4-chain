package v0_2_2

import (
	store "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/upgrades"
)

const UpgradeName = "v0.2.2"

var (
	Fork = upgrades.Fork{
		UpgradeName: UpgradeName,
		// Target upgrade time is Sept 15, 2023, 1pm EST,
		// estimated to occur on block 178500, assuming 1.5s block time.
		UpgradeHeight: 178500,
		UpgradeInfo:   "",
	}

	Upgrade = upgrades.Upgrade{
		UpgradeName:   UpgradeName,
		StoreUpgrades: store.StoreUpgrades{},
	}
)
