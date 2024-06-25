package v_cosmwasm_0

import (
	store "cosmossdk.io/store/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/upgrades"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

const (
	UpgradeName = "v6.6.6"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName: UpgradeName,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{
			vaulttypes.StoreKey,
			wasmtypes.StoreKey,
		},
	},
}
