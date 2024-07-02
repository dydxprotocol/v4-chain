package v_6_0_0

import (
	store "cosmossdk.io/store/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/upgrades"
	listingtypes "github.com/dydxprotocol/v4-chain/protocol/x/listing/types"
	revsharetypes "github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
	marketmapmoduletypes "github.com/skip-mev/slinky/x/marketmap/types"
)

const (
	UpgradeName = "v6.0.0"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName: UpgradeName,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{
			listingtypes.StoreKey,
			revsharetypes.StoreKey,
			marketmapmoduletypes.StoreKey,
		},
	},
}
