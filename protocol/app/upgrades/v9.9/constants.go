package v_9_9

import (
	store "cosmossdk.io/store/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/upgrades"
	canonicalusdckeeper "github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/keeper"
	canonicalusdctypes "github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"
)

const UpgradeName = "v9.9"

// Canonical USDC routes. The upgrade handler enters GRADUAL bound to the
// pinned channels on a chain listed in PinnedRoutes, and elsewhere to the one
// open transfer channel to each chain. It leaves the module DISABLED for a
// governance proposal when a binding does not validate or is ambiguous.
const (
	NobleChainID     = "noble-1"
	InjectiveChainID = "injective-1"
	NoblePacketDenom = "uusdc"
	// Native USDC on Injective, Circle's MultiVM ERC20.
	InjectivePacketDenom  = "erc20:0xa00C59fF5a080D2b954d0c75e46E22a0c371235a"
	MaxPendingSettlements = 10_000
)

type Routes struct {
	Noble     canonicalusdckeeper.RouteBinding
	Injective canonicalusdckeeper.RouteBinding
}

// PinnedRoutes were checked against dydx-mainnet-1 on 2026-11-12: the only
// open transfer channels to noble-1 and injective-1, the latter already
// carrying USDC vouchers.
var PinnedRoutes = map[string]Routes{
	"dydx-mainnet-1": {
		Noble:     canonicalusdckeeper.RouteBinding{Channel: "channel-0", Connection: "connection-0", Client: "07-tendermint-0"},
		Injective: canonicalusdckeeper.RouteBinding{Channel: "channel-62", Connection: "connection-30", Client: "07-tendermint-26"},
	},
}

var Upgrade = upgrades.Upgrade{
	UpgradeName: UpgradeName,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{canonicalusdctypes.StoreKey},
	},
}
