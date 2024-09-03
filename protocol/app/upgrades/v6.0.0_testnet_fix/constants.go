package v_6_0_0_testnet_fix

import (
	store "cosmossdk.io/store/types"

	"github.com/dydxprotocol/v4-chain/protocol/app/upgrades"
)

const (
	// v6_0_0_testnet_fix not intended for prod use.
	UpgradeName = "v6.0.0_testnet_fix"
)

var (
	Upgrade = upgrades.Upgrade{
		UpgradeName:   UpgradeName,
		StoreUpgrades: store.StoreUpgrades{},
	}
)
