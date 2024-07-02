package v_cosmwasm_0

import (
	store "cosmossdk.io/store/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/upgrades"
	listingtypes "github.com/dydxprotocol/v4-chain/protocol/x/listing/types"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	revsharetypes "github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
)

const (
	UpgradeName = "v6.cosmwasm.0"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName: UpgradeName,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{
			listingtypes.StoreKey,
			wasmtypes.StoreKey,
			revsharetypes.StoreKey,
		},
	},
}
