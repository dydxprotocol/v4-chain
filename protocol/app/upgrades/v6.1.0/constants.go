package v_6_1_0

import (
	store "cosmossdk.io/store/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/upgrades"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
)

const (
	UpgradeName = "v6.1.0"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName: UpgradeName,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{
			wasmtypes.StoreKey,
		},
	},
}
